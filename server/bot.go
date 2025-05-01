package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/gin-gonic/gin"
	"github.com/olahol/melody"
)

type Portfolio struct {
	AccountValue           float64                  `json:"accountValue" firestore:"accountValue"`
	HistoricalAccountValue []*AccountValueHistory   `json:"historicalAccountValue" firestore:"historicalAccountValue"`
	Cash                   float64                  `json:"cash" firestore:"cash"`
	Holdings               map[string]*Holding      `json:"holdings" firestore:"holdings"`
	Transactions           []*Transaction           `json:"transactions" firestore:"-"`
	TransactionReferences  []*firestore.DocumentRef `json:"-" firestore:"transactions"`
}

type AccountValueHistory struct {
	Date  time.Time `json:"date" firestore:"date"`
	Value float64   `json:"value" firestore:"value"`
}

func NewPortfolio(startingCash float64) *Portfolio {
	return &Portfolio{
		Cash:                  startingCash,
		Holdings:              make(map[string]*Holding),
		Transactions:          make([]*Transaction, 0),
		TransactionReferences: make([]*firestore.DocumentRef, 0),
	}
}

func (p *Portfolio) Buy(transaction *Transaction) error {
	switch {
	case p.Cash < transaction.NumShares*transaction.UnitCost:
		return fmt.Errorf("not enough cash to buy %f shares of %s", transaction.NumShares, transaction.Ticker)
	case transaction.NumShares < 0:
		return fmt.Errorf("cannot buy negative number of shares")
	}

	if p.Holdings == nil {
		p.Holdings = make(map[string]*Holding)
	}

	p.Cash -= transaction.NumShares * transaction.UnitCost
	if holding, ok := p.Holdings[transaction.Ticker]; !ok {
		p.Holdings[transaction.Ticker] = &Holding{
			NumShares:     transaction.NumShares,
			PurchaseValue: transaction.UnitCost,
		}
	} else {
		holding.NumShares += transaction.NumShares
		holding.PurchaseValue = (holding.PurchaseValue*holding.NumShares + transaction.NumShares*transaction.UnitCost) / (holding.NumShares + transaction.NumShares)
	}

	return nil
}

func (p *Portfolio) Execute(transaction *Transaction) error {
	switch transaction.Action {
	case "buy":
		return p.Buy(transaction)
	case "sell":
		return p.Sell(transaction)
	default:
		return fmt.Errorf("invalid transaction action: %s", transaction.Action)
	}
}

func (p *Portfolio) Sell(transaction *Transaction) error {
	switch {
	case p.Holdings[transaction.Ticker].NumShares < transaction.NumShares:
		return fmt.Errorf("not enough shares to sell %f shares of %s", transaction.NumShares, transaction.Ticker)
	case transaction.NumShares < 0:
		return fmt.Errorf("cannot sell negative number of shares")
	}

	p.Cash += transaction.NumShares * transaction.UnitCost
	p.Holdings[transaction.Ticker].NumShares -= transaction.NumShares
	p.Holdings[transaction.Ticker].PurchaseValue = transaction.UnitCost

	return nil
}

type Holding struct {
	NumShares     float64 `json:"numShares" firestore:"numShares"`
	PurchaseValue float64 `json:"purchaseValue" firestore:"purchaseValue"`
}

type DataPacket struct {
	Type    string `json:"type"`
	Payload any    `json:"payload"`
}

func (dp *DataPacket) JSON() []byte {
	b, err := json.Marshal(dp)
	if err != nil {
		panic(err)
	}

	return b
}

type ResultData struct {
	Message string `json:"payload"`
	Success bool   `json:"success"`
}

func NewResultPacket(message string, success bool) *DataPacket {
	return &DataPacket{
		Type:    "result",
		Payload: &ResultData{message, success},
	}
}

type AuthData struct {
	Key string `json:"key"`
}

type TransactionRequestData struct {
	Action    string  `json:"action"`
	NumShares float64 `json:"numShares"`
	Ticker    string  `json:"ticker"`
}

type BotWorker struct {
	db           *firestore.Client
	tiingo       *Tiingo
	latestPrices map[string]float64
}

func (bw *BotWorker) calculateAccountValue(doc *firestore.DocumentSnapshot) {
	portfolio := &Portfolio{}
	doc.DataTo(portfolio)
	log.Printf("calculating portfolio: %v\n", doc.Ref.ID)

	// Noise Generator
	//portfolio.HistoricalAccountValue = make([]*AccountValueHistory, 0)
	//currTime := time.Date(2025, 3, 23, 0, 0, 0, 0, time.UTC)
	//var lastValue float64 = 100
	//for currTime.Before(time.Now()) {
	//	portfolio.HistoricalAccountValue = append(portfolio.HistoricalAccountValue, &AccountValueHistory{currTime, lastValue})
	//	currTime = currTime.Add(time.Hour * 24)
	//	lastValue += (rand.Float64() - 0.44) * 50
	//}
	//
	//portfolio.AccountValue = portfolio.HistoricalAccountValue[len(portfolio.HistoricalAccountValue)-1].Value
	//
	//_, err = doc.Ref.Update(context.Background(), []firestore.Update{
	//	{Path: "accountValue", Value: portfolio.AccountValue},
	//	{Path: "historicalAccountValue", Value: portfolio.HistoricalAccountValue},
	//})
	//if err != nil {
	//	log.Println(err)
	//}

	portfolio.AccountValue = portfolio.Cash

	for ticker, holding := range portfolio.Holdings {
		price, ok := bw.latestPrices[ticker]
		if !ok {
			bw.tiingo.AddTickers(ticker)
			log.Printf("failed to find ticker data for \"%s\" while calculating portfolio: %v\n", ticker, doc.Ref.ID)
			return
		}

		portfolio.AccountValue += holding.NumShares * price
	}

	if len(portfolio.HistoricalAccountValue) == 0 {
		portfolio.HistoricalAccountValue = make([]*AccountValueHistory, 0)
		portfolio.HistoricalAccountValue = append(portfolio.HistoricalAccountValue, &AccountValueHistory{
			Date:  time.Now(),
			Value: portfolio.AccountValue,
		})
	} else if portfolio.HistoricalAccountValue[len(portfolio.HistoricalAccountValue)-1].Date.Add(time.Hour * 24).Before(time.Now()) {
		portfolio.HistoricalAccountValue = append(portfolio.HistoricalAccountValue, &AccountValueHistory{
			Date:  time.Now(),
			Value: portfolio.AccountValue,
		})
	} else {
		portfolio.HistoricalAccountValue[len(portfolio.HistoricalAccountValue)-1].Value = portfolio.AccountValue
		portfolio.HistoricalAccountValue[len(portfolio.HistoricalAccountValue)-1].Date = time.Now()
	}

	log.Printf("updated portfolio: %v\n", doc.Ref.ID)
	_, err := doc.Ref.Update(context.Background(), []firestore.Update{
		{Path: "accountValue", Value: portfolio.AccountValue},
		{Path: "historicalAccountValue", Value: portfolio.HistoricalAccountValue},
	})
	if err != nil {
		log.Println(err)
	}
}

func NewBotWorker(db *firestore.Client, tiingo *Tiingo) *BotWorker {
	bw := &BotWorker{
		db:           db,
		tiingo:       tiingo,
		latestPrices: make(map[string]float64),
	}

	dataDownloader := time.NewTicker(time.Minute * 5)
	go func() {
		for {
			select {
			case <-dataDownloader.C:
				if time.Now().In(time.UTC).Hour() < 14 || time.Now().In(time.UTC).Hour() > 21 {
					log.Println("skipping data download because it is not in the trading hours")
					continue
				}
				bw.tiingo.DownloadAllTickers()

				for ticker := range tiingo.tickers.All() {
					_, row := tiingo.DailyCache.GetClosestRowBefore(tiingo.DailyCache.Tickers[ticker].End)
					data, ok := row.Data.Load(ticker)
					if !ok {
						log.Printf("error retrieving data for ticker %s\n", ticker)
					} else {
						bw.latestPrices[ticker] = data.Close
					}
				}

				log.Printf("updated prices: %v\n", bw.latestPrices)
			}
		}
	}()

	// TODO: Change this to a webhook
	accountValuer := time.NewTicker(time.Second * 10)
	go func() {
		for {
			select {
			case <-accountValuer.C:
				docs, err := bw.db.Collection("bots").Documents(context.Background()).GetAll()
				if err != nil {
					log.Printf("error retrieving bots: %v\n", err)
					continue
				}

				for _, doc := range docs {
					go bw.calculateAccountValue(doc)
				}
			}
		}
	}()

	return bw
}

func (bw *BotWorker) authenticate(s *melody.Session, data any) {
	auth, ok := data.(AuthData)
	if !ok {
		s.CloseWithMsg(NewResultPacket("error parsing ws packet", false).JSON())
		return
	}

	s.Set("apiKey", auth.Key)

	bot, err := bw.db.Collection("bots").Where("apiKey", "==", auth.Key).Documents(context.Background()).Next()
	if err != nil || bot == nil {
		s.CloseWithMsg(NewResultPacket("error finding bot with specified api key", false).JSON())
		return
	}

	portfolio := &Portfolio{}
	bot.DataTo(portfolio)

	s.Set("bot", portfolio)
	s.Set("db_ref", bot.Ref)
}

//func (bw *BotWorker) transact(s *melody.Session, data any) {
//	transaction, ok := data.(TransactionRequestData)
//	if !ok {
//		s.Write(NewResultPacket("error parsing ws packet", false).JSON())
//		return
//	}
//
//	bot, ok := s.Get("bot")
//	if !ok {
//		s.CloseWithMsg(NewResultPacket("error: not authenticated", false).JSON())
//		return
//	}
//
//	price, ok := bw.latestPrices[transaction.Ticker]
//	if !ok {
//		s.Write(NewResultPacket("error: ticker data not available, make sure to subscribe and receive a ticker data update first", false).JSON())
//		return
//	}
//
//	portfolio := bot.(*Portfolio)
//	portfolio.Buy(&Transaction{
//		Time:      time.Now(),
//		NumShares: transaction.NumShares,
//		UnitCost:  price,
//		Ticker:    transaction.Ticker,
//	})
//}

// TODO: Scope this out
// TradingStream Initial Packet Must Be Of Type "auth"
//func (bw *BotWorker) TradingStream(s *melody.Session, msg []byte) {
//	packet := &DataPacket{}
//	err := json.Unmarshal(msg, packet)
//	if err != nil {
//		s.CloseWithMsg(NewResultPacket("error parsing ws packet", false).JSON())
//		return
//	}
//
//	if _, ok := s.Get("apiKey"); packet.Type != "auth" && !ok {
//		s.CloseWithMsg(NewResultPacket("unauthenticated", false).JSON())
//		return
//	}
//
//	switch packet.Type {
//	case "auth":
//		bw.authenticate(s, packet.Payload)
//	case "transact":
//		bw.transact(s, packet.Payload)
//	case "add_subscription":
//	}
//}
//
// TODO: Scope this out
//func (bw *BotWorker) StockDataBroadcast(m *melody.Melody) {
//	c := time.NewTicker(time.Second * 30)
//	go func() {
//		for {
//			select {
//			case <-c.C:
//				bw.tiingo.LoadCaches(false)
//
//				jsonData, err := json.Marshal(bw.tiingo.DailyCache.Pack())
//				if err != nil {
//					log.Println(err)
//					return
//				}
//
//				m.Broadcast(jsonData)
//			}
//		}
//	}()
//}

func (bw *BotWorker) AuthHandler(c *gin.Context) {
	apikey := c.GetHeader("Authorization")
	bot, err := bw.db.Collection("bots").Where("apiKey", "==", apikey).Documents(context.Background()).Next()
	if err != nil || bot == nil {
		c.AbortWithStatusJSON(401, NewResultPacket("error finding bot with specified api key", false))
		return
	}

	portfolio := &Portfolio{}
	bot.DataTo(portfolio)

	c.Set("db_ref", bot.Ref)
	c.Set("bot", portfolio)
}

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
		{Path: "cash", Value: botUntyped.(*Portfolio).Cash},
		{Path: "holdings", Value: botUntyped.(*Portfolio).Holdings},
		{Path: "transactions", Value: botUntyped.(*Portfolio).TransactionReferences},
	})
}

func (bw *BotWorker) AddTicker(c *gin.Context) {
	tickers, ok := c.GetQueryArray("ticker")
	if !ok {
		c.AbortWithStatusJSON(400, NewResultPacket("error parsing ticker query", false))
		return
	}

	bw.tiingo.AddTickers(tickers...)
	c.JSON(200, NewResultPacket(fmt.Sprintf("successfully added tickers: %v", tickers), true))

	if time.Now().In(time.UTC).Hour() < 14 || time.Now().In(time.UTC).Hour() > 21 {
		for _, ticker := range tickers {
			err := bw.tiingo.historicalDaily(ticker)
			if err != nil {
				log.Printf("error downloading historical data for ticker %s: %v\n", ticker, err)
				continue
			}
		}

		err := bw.tiingo.SaveCaches()
		if err != nil {
			log.Printf("error saving historical data for tickers: %v\n", err)
		}

		for ticker := range bw.tiingo.tickers.All() {
			_, row := bw.tiingo.DailyCache.GetClosestRowBefore(bw.tiingo.DailyCache.Tickers[ticker].End)
			data, ok := row.Data.Load(ticker)
			if !ok {
				log.Printf("error retrieving data for ticker %s\n", ticker)
			} else {
				bw.latestPrices[ticker] = data.Close
			}
		}
	}
}

func (bw *BotWorker) GetStockData(c *gin.Context) {
	c.JSON(200, &DataPacket{"stock_data", bw.tiingo.DailyCache.Pack()})
}

func (bw *BotWorker) MakeTransaction(c *gin.Context) {
	bot, ok := c.Get("bot")
	if !ok {
		c.AbortWithStatusJSON(401, NewResultPacket("error: not authenticated", false))
		return
	}

	portfolio, ok := bot.(*Portfolio)
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

	transaction := &Transaction{
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

func (bw *BotWorker) GetPortfolio(c *gin.Context) {
	bot, ok := c.Get("bot")
	if !ok {
		c.AbortWithStatusJSON(401, NewResultPacket("error: not authenticated", false))
		return
	}

	portfolio, ok := bot.(*Portfolio)
	if !ok {
		c.AbortWithStatusJSON(500, NewResultPacket("error: failed to retrieve portfolio information", false))
		return
	}

	portfolio.Transactions = make([]*Transaction, 0, len(portfolio.TransactionReferences))
	for _, ref := range portfolio.TransactionReferences {
		doc, err := ref.Get(context.Background())
		if err != nil {
			c.AbortWithStatusJSON(500, NewResultPacket("error: failed to retrieve transaction information", false))
			return
		}

		transaction := &Transaction{}
		doc.DataTo(transaction)
		portfolio.Transactions = append(portfolio.Transactions, transaction)
	}

	c.JSON(200, &DataPacket{"portfolio", portfolio})
}
