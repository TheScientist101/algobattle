// Package bot implements the core business logic for the AlgoBattle trading platform.
// It handles portfolio management, stock data retrieval, and transaction processing.
// The package also provides HTTP handlers for the REST API endpoints.
package bot

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/gin-gonic/gin"
	"urjith.dev/algobattle/pkg/models"
	"urjith.dev/algobattle/pkg/services"
)

// DataPacket represents a data packet sent over WebSocket
type DataPacket struct {
	Type    string `json:"type"`
	Payload any    `json:"payload"`
}

// JSON converts the DataPacket to JSON
func (dp *DataPacket) JSON() []byte {
	b, err := json.Marshal(dp)
	if err != nil {
		panic(err)
	}

	return b
}

// ResultData represents a result message
type ResultData struct {
	Message string `json:"payload"`
	Success bool   `json:"success"`
}

// NewResultPacket creates a new result packet
func NewResultPacket(message string, success bool) *DataPacket {
	return &DataPacket{
		Type:    "result",
		Payload: &ResultData{message, success},
	}
}

// AuthData represents authentication data
type AuthData struct {
	Key string `json:"key"`
}

// TransactionRequestData represents a transaction request
type TransactionRequestData struct {
	Action    string  `json:"action"`
	NumShares float64 `json:"numShares"`
	Ticker    string  `json:"ticker"`
}

// BotWorker manages bots and their portfolios
type BotWorker struct {
	db           *firestore.Client
	tiingo       *services.Tiingo
	latestPrices map[string]float64
}

// NewBotWorker creates a new BotWorker
func NewBotWorker(db *firestore.Client, tiingo *services.Tiingo) *BotWorker {
	bw := &BotWorker{
		db:           db,
		tiingo:       tiingo,
		latestPrices: make(map[string]float64),
	}

	liveDownloader := time.NewTicker(time.Minute * 5)
	dailyDownloader := time.NewTicker(time.Hour * 24)
	accountValuer := make(chan bool)
	go func() {
		for ; true; <-liveDownloader.C {
			if time.Now().In(time.UTC).Hour() < 14 || time.Now().In(time.UTC).Hour() > 21 {
				log.Println("skipping data download because it is not in the trading hours")
				continue
			}

			bw.updateCurrPrices()
			accountValuer <- true
		}
	}()

	go func() {
		for ; true; <-dailyDownloader.C {
			err := bw.tiingo.DownloadAllTickers()
			if err != nil {
				log.Printf("error downloading daily stock data: %v\n", err)
			}
		}
	}()

	// TODO: Change this to a webhook
	go func() {
		for ; true; <-accountValuer {
			docs, err := bw.db.Collection("bots").Documents(context.Background()).GetAll()
			if err != nil {
				log.Printf("error retrieving bots: %v\n", err)
				continue
			}

			for _, doc := range docs {
				go bw.calculateAccountValue(doc)
			}
		}
	}()

	return bw
}

// calculateAccountValue calculates the account value for a portfolio
func (bw *BotWorker) calculateAccountValue(doc *firestore.DocumentSnapshot) {
	portfolio := &models.Portfolio{}
	doc.DataTo(portfolio)
	log.Printf("calculating portfolio: %v\n", doc.Ref.ID)

	oldAccountValue := portfolio.AccountValue
	historyChanged := false

	portfolio.AccountValue = portfolio.Cash

	for ticker, holding := range portfolio.Holdings {
		price, ok := bw.latestPrices[ticker]
		if !ok {
			log.Printf("failed to find ticker data for \"%s\" while calculating portfolio: %v\nadding %s to watchlist...\n", ticker, doc.Ref.ID, ticker)
			err := bw.addTickers(ticker)

			if err != nil {
				log.Printf("error while adding ticker: %v\n", err)
			}
			return
		}

		portfolio.AccountValue += holding.NumShares * price
	}

	if len(portfolio.HistoricalAccountValue) == 0 {
		portfolio.HistoricalAccountValue = make([]*models.AccountValueHistory, 0)
		portfolio.HistoricalAccountValue = append(portfolio.HistoricalAccountValue, &models.AccountValueHistory{
			Date:  time.Now(),
			Value: portfolio.AccountValue,
		})
		historyChanged = true
	} else if portfolio.HistoricalAccountValue[len(portfolio.HistoricalAccountValue)-1].Date.YearDay() != time.Now().YearDay() {
		portfolio.HistoricalAccountValue = append(portfolio.HistoricalAccountValue, &models.AccountValueHistory{
			Date:  time.Now(),
			Value: portfolio.AccountValue,
		})
		historyChanged = true
	} else if portfolio.HistoricalAccountValue[len(portfolio.HistoricalAccountValue)-1].Value != portfolio.AccountValue {
		portfolio.HistoricalAccountValue[len(portfolio.HistoricalAccountValue)-1].Value = portfolio.AccountValue
		portfolio.HistoricalAccountValue[len(portfolio.HistoricalAccountValue)-1].Date = time.Now()
		historyChanged = true
	}

	if !historyChanged && oldAccountValue == portfolio.AccountValue {
		log.Printf("no change in account value for portfolio: %v\n", doc.Ref.ID)
		return
	}

	log.Printf("updated portfolio: %v\nlatest account value: %v\n", doc.Ref.ID, portfolio.AccountValue)
	_, err := doc.Ref.Update(context.Background(), []firestore.Update{
		{Path: "accountValue", Value: portfolio.AccountValue},
		{Path: "historicalAccountValue", Value: portfolio.HistoricalAccountValue},
	})
	if err != nil {
		log.Println(err)
	}
}

// AuthHandler authenticates a request using the API key in the Authorization header.
// It loads the user's portfolio and sets it in the context for downstream handlers.
// This middleware should be applied to all routes that require authentication.
func (bw *BotWorker) AuthHandler(c *gin.Context) {
	// Get API key from Authorization header
	apikey := c.GetHeader("Authorization")

	// Find the bot with the matching API key
	bot, err := bw.db.Collection("bots").Where("apiKey", "==", apikey).Documents(context.Background()).Next()
	if err != nil || bot == nil {
		c.AbortWithStatusJSON(401, NewResultPacket("error finding bot with specified api key", false))
		return
	}

	// Load the portfolio data
	portfolio := &models.Portfolio{}
	bot.DataTo(portfolio)

	// Set the database reference and portfolio in the context
	c.Set("db_ref", bot.Ref)
	c.Set("bot", portfolio)
}

// SavePortfolio saves the updated portfolio to the database.
// This middleware should be applied after handlers that modify the portfolio.
// @Summary Save portfolio changes
// @Description Saves any changes made to the portfolio during the request
// @Tags portfolio
// @Accept json
// @Produce json
// @Success 200 {object} ResultData "Portfolio saved"
// @Failure 401 {object} ResultData "Not authenticated"
// @Router /transact [post]
func (bw *BotWorker) SavePortfolio(c *gin.Context) {
	// Get the database reference from the context
	refUntyped, ok := c.Get("db_ref")
	if !ok {
		c.AbortWithStatusJSON(401, NewResultPacket("error: not authenticated", false))
		return
	}

	// Get the portfolio from the context
	botUntyped, ok := c.Get("bot")
	if !ok {
		c.AbortWithStatusJSON(401, NewResultPacket("error: failed to save portfolio information", false))
		return
	}

	// Update the portfolio in the database
	ref := refUntyped.(*firestore.DocumentRef)
	ref.Update(context.Background(), []firestore.Update{
		{Path: "cash", Value: botUntyped.(*models.Portfolio).Cash},
		{Path: "holdings", Value: botUntyped.(*models.Portfolio).Holdings},
		{Path: "transactions", Value: botUntyped.(*models.Portfolio).TransactionReferences},
	})
}

// AddTicker adds one or more tickers to the watchlist for monitoring.
// @Summary Add ticker to watchlist
// @Description Adds one or more stock tickers to the watchlist for price monitoring and data collection
// @Tags stocks
// @Accept json
// @Produce json
// @Param ticker query []string true "Ticker symbols to add (can specify multiple)"
// @Success 200 {object} ResultData "Tickers added successfully"
// @Failure 400 {object} ResultData "Invalid request"
// @Failure 500 {object} ResultData "Server error"
// @Router /add_ticker [get]
func (bw *BotWorker) AddTicker(c *gin.Context) {
	// Get ticker symbols from query parameters
	tickers, ok := c.GetQueryArray("ticker")
	if !ok {
		c.AbortWithStatusJSON(400, NewResultPacket("error parsing ticker query", false))
		return
	}

	// Add tickers to the watchlist and download their data
	err := bw.addTickers(tickers...)
	if err != nil {
		log.Printf("error while adding ticker: %v\n", err)
		c.AbortWithStatusJSON(500, NewResultPacket("failed to add at least one ticker", false))
		return
	}

	// Return success response
	c.JSON(200, NewResultPacket(fmt.Sprintf("successfully added tickers: %v", tickers), true))
}

// addTickers adds tickers to the watchlist
func (bw *BotWorker) addTickers(tickers ...string) error {
	bw.tiingo.AddTickers(tickers...)
	bw.updateCurrPrices()
	return bw.tiingo.DownloadMissingTickers()
}

// GetDailyStockData returns historical daily stock data for all watched tickers.
// @Summary Get historical stock data
// @Description Retrieves daily historical stock data for all tickers in the watchlist
// @Tags stocks
// @Accept json
// @Produce json
// @Success 200 {object} DataPacket "Historical daily stock data"
// @Failure 401 {object} ResultData "Not authenticated"
// @Router /daily_stock_data [get]
func (bw *BotWorker) GetDailyStockData(c *gin.Context) {
	// Pack and return the daily cache as JSON
	c.JSON(200, &DataPacket{"daily_stock_data", bw.tiingo.DailyCache.Pack()})
}

// MakeTransaction executes a buy or sell transaction for a stock.
// @Summary Execute a stock transaction
// @Description Processes a buy or sell transaction for a specified ticker and number of shares
// @Tags transactions
// @Accept json
// @Produce json
// @Param transaction body TransactionRequestData true "Transaction details"
// @Success 200 {object} ResultData "Transaction successful"
// @Failure 401 {object} ResultData "Not authenticated or insufficient funds/shares"
// @Failure 500 {object} ResultData "Server error"
// @Router /transact [post]
func (bw *BotWorker) MakeTransaction(c *gin.Context) {
	// Get the bot from the context (set by AuthHandler)
	bot, ok := c.Get("bot")
	if !ok {
		c.AbortWithStatusJSON(401, NewResultPacket("error: not authenticated", false))
		return
	}

	// Type assertion to get the portfolio
	portfolio, ok := bot.(*models.Portfolio)
	if !ok {
		c.AbortWithStatusJSON(500, NewResultPacket("error: failed to retrieve portfolio information", false))
		return
	}

	// Read and parse the request body
	body, err := c.GetRawData()
	if err != nil {
		c.AbortWithStatusJSON(500, NewResultPacket("error: failed to retrieve request body", false))
		return
	}

	request := &TransactionRequestData{}
	err = json.Unmarshal(body, request)
	if err != nil {
		c.AbortWithStatusJSON(500, NewResultPacket("error: failed to parse request body", false))
		return
	}

	// Get the current price for the ticker
	cost, ok := bw.latestPrices[request.Ticker]
	if !ok {
		c.AbortWithStatusJSON(500, NewResultPacket("error: ticker data not available, make sure to subscribe and receive a ticker data update first", false))
		return
	}

	// Get the database reference
	refUntyped, ok := c.Get("db_ref")
	if !ok {
		c.AbortWithStatusJSON(500, NewResultPacket("error: failed to retrieve portfolio database reference", false))
		return
	}

	ref, ok := refUntyped.(*firestore.DocumentRef)
	if !ok {
		c.AbortWithStatusJSON(500, NewResultPacket("error: failed to retrieve portfolio database reference", false))
		return
	}

	// Create the transaction object
	transaction := &models.Transaction{
		Time:      time.Now(),
		NumShares: request.NumShares,
		UnitCost:  cost,
		Ticker:    request.Ticker,
		Action:    request.Action,
		Bot:       ref,
	}

	// Execute the transaction on the portfolio
	err = portfolio.Execute(transaction)
	if err != nil {
		c.AbortWithStatusJSON(401, NewResultPacket(err.Error(), false))
		return
	}

	// Save the transaction to the database
	doc, _, err := bw.db.Collection("transactions").Add(context.Background(), transaction)
	if err != nil {
		c.AbortWithStatusJSON(500, NewResultPacket("error: failed to save transaction", false))
		return
	}

	// Add the transaction reference to the portfolio
	portfolio.TransactionReferences = append(portfolio.TransactionReferences, doc)
	c.JSON(200, NewResultPacket("successfully executed transaction", true))
}

// GetPortfolio returns the user's portfolio with all holdings and transactions.
// @Summary Get user portfolio
// @Description Retrieves the authenticated user's portfolio including cash balance, holdings, and transaction history
// @Tags portfolio
// @Accept json
// @Produce json
// @Success 200 {object} DataPacket "Portfolio data"
// @Failure 401 {object} ResultData "Not authenticated"
// @Failure 500 {object} ResultData "Server error"
// @Router /portfolio [get]
func (bw *BotWorker) GetPortfolio(c *gin.Context) {
	// Get the bot from the context (set by AuthHandler)
	bot, ok := c.Get("bot")
	if !ok {
		c.AbortWithStatusJSON(401, NewResultPacket("error: not authenticated", false))
		return
	}

	// Type assertion to get the portfolio
	portfolio, ok := bot.(*models.Portfolio)
	if !ok {
		c.AbortWithStatusJSON(500, NewResultPacket("error: failed to retrieve portfolio information", false))
		return
	}

	// Load all transactions from references
	portfolio.Transactions = make([]*models.Transaction, 0, len(portfolio.TransactionReferences))
	for _, ref := range portfolio.TransactionReferences {
		doc, err := ref.Get(context.Background())
		if err != nil {
			c.AbortWithStatusJSON(500, NewResultPacket("error: failed to retrieve transaction information", false))
			return
		}

		transaction := &models.Transaction{}
		doc.DataTo(transaction)
		portfolio.Transactions = append(portfolio.Transactions, transaction)
	}

	// Return the portfolio as JSON
	c.JSON(200, &DataPacket{"portfolio", portfolio})
}

// GetLiveStockData returns the current stock prices for all watched tickers.
// @Summary Get live stock prices
// @Description Retrieves the latest stock prices for all tickers in the watchlist
// @Tags stocks
// @Accept json
// @Produce json
// @Success 200 {object} DataPacket "Live stock price data"
// @Failure 401 {object} ResultData "Not authenticated"
// @Router /live_stock_data [get]
func (bw *BotWorker) GetLiveStockData(c *gin.Context) {
	// Return the latest prices as JSON
	c.JSON(200, &DataPacket{"live_stock_data", bw.latestPrices})
}

// updateCurrPrices updates the current prices
func (bw *BotWorker) updateCurrPrices() {
	bw.latestPrices = bw.tiingo.FetchCurrPrices()
	log.Printf("updated prices: %v\n", bw.latestPrices)
}
