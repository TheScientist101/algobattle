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

// AuthHandler authenticates a request
func (bw *BotWorker) AuthHandler(c *gin.Context) {
	apikey := c.GetHeader("Authorization")
	bot, err := bw.db.Collection("bots").Where("apiKey", "==", apikey).Documents(context.Background()).Next()
	if err != nil || bot == nil {
		c.AbortWithStatusJSON(401, NewResultPacket("error finding bot with specified api key", false))
		return
	}

	portfolio := &models.Portfolio{}
	bot.DataTo(portfolio)

	c.Set("db_ref", bot.Ref)
	c.Set("bot", portfolio)
}

// SavePortfolio saves a portfolio to the database
func (bw *BotWorker) SavePortfolio(c *gin.Context) {
	refUntyped, ok := c.Get("db_ref")
	if !ok {
		c.AbortWithStatusJSON(401, NewResultPacket("error: not authenticated", false))
		return
	}

	botUntyped, ok := c.Get("bot")
	if !ok {
		c.AbortWithStatusJSON(401, NewResultPacket("error: failed to save portfolio information", false))
		return
	}

	ref := refUntyped.(*firestore.DocumentRef)
	ref.Update(context.Background(), []firestore.Update{
		{Path: "cash", Value: botUntyped.(*models.Portfolio).Cash},
		{Path: "holdings", Value: botUntyped.(*models.Portfolio).Holdings},
		{Path: "transactions", Value: botUntyped.(*models.Portfolio).TransactionReferences},
	})
}

// AddTicker adds a ticker to the watchlist
func (bw *BotWorker) AddTicker(c *gin.Context) {
	tickers, ok := c.GetQueryArray("ticker")
	if !ok {
		c.AbortWithStatusJSON(400, NewResultPacket("error parsing ticker query", false))
		return
	}

	err := bw.addTickers(tickers...)
	if err != nil {
		log.Printf("error while adding ticker: %v\n", err)
		c.AbortWithStatusJSON(500, NewResultPacket("failed to add at least one ticker", false))
		return
	}

	c.JSON(200, NewResultPacket(fmt.Sprintf("successfully added tickers: %v", tickers), true))
}

// addTickers adds tickers to the watchlist
func (bw *BotWorker) addTickers(tickers ...string) error {
	bw.tiingo.AddTickers(tickers...)
	bw.updateCurrPrices()
	return bw.tiingo.DownloadMissingTickers()
}

// GetDailyStockData returns the daily stock data
func (bw *BotWorker) GetDailyStockData(c *gin.Context) {
	c.JSON(200, &DataPacket{"daily_stock_data", bw.tiingo.DailyCache.Pack()})
}

// MakeTransaction executes a transaction
func (bw *BotWorker) MakeTransaction(c *gin.Context) {
	bot, ok := c.Get("bot")
	if !ok {
		c.AbortWithStatusJSON(401, NewResultPacket("error: not authenticated", false))
		return
	}

	portfolio, ok := bot.(*models.Portfolio)
	if !ok {
		c.AbortWithStatusJSON(500, NewResultPacket("error: failed to retrieve portfolio information", false))
		return
	}

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

	cost, ok := bw.latestPrices[request.Ticker]
	if !ok {
		c.AbortWithStatusJSON(500, NewResultPacket("error: ticker data not available, make sure to subscribe and receive a ticker data update first", false))
		return
	}

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

	transaction := &models.Transaction{
		Time:      time.Now(),
		NumShares: request.NumShares,
		UnitCost:  cost,
		Ticker:    request.Ticker,
		Action:    request.Action,
		Bot:       ref,
	}

	err = portfolio.Execute(transaction)
	if err != nil {
		c.AbortWithStatusJSON(401, NewResultPacket(err.Error(), false))
		return
	}

	doc, _, err := bw.db.Collection("transactions").Add(context.Background(), transaction)
	if err != nil {
		c.AbortWithStatusJSON(500, NewResultPacket("error: failed to save transaction", false))
		return
	}

	portfolio.TransactionReferences = append(portfolio.TransactionReferences, doc)
	c.JSON(200, NewResultPacket("successfully executed transaction", true))
}

// GetPortfolio returns the portfolio
func (bw *BotWorker) GetPortfolio(c *gin.Context) {
	bot, ok := c.Get("bot")
	if !ok {
		c.AbortWithStatusJSON(401, NewResultPacket("error: not authenticated", false))
		return
	}

	portfolio, ok := bot.(*models.Portfolio)
	if !ok {
		c.AbortWithStatusJSON(500, NewResultPacket("error: failed to retrieve portfolio information", false))
		return
	}

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

	c.JSON(200, &DataPacket{"portfolio", portfolio})
}

// GetLiveStockData returns the live stock data
func (bw *BotWorker) GetLiveStockData(c *gin.Context) {
	c.JSON(200, &DataPacket{"live_stock_data", bw.latestPrices})
}

// updateCurrPrices updates the current prices
func (bw *BotWorker) updateCurrPrices() {
	bw.latestPrices = bw.tiingo.FetchCurrPrices()
	log.Printf("updated prices: %v\n", bw.latestPrices)
}