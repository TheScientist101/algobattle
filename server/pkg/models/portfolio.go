package models

import (
	"cloud.google.com/go/firestore"
	"fmt"
	"time"
)

// Portfolio represents a user's portfolio of stocks
type Portfolio struct {
	AccountValue           float64                  `json:"accountValue" firestore:"accountValue"`
	HistoricalAccountValue []*AccountValueHistory   `json:"historicalAccountValue" firestore:"historicalAccountValue"`
	Cash                   float64                  `json:"cash" firestore:"cash"`
	Holdings               map[string]*Holding      `json:"holdings" firestore:"holdings"`
	Transactions           []*Transaction           `json:"transactions" firestore:"-"`
	TransactionReferences  []*firestore.DocumentRef `json:"-" firestore:"transactions"`
}

// AccountValueHistory represents a historical account value at a specific date
type AccountValueHistory struct {
	Date  time.Time `json:"date" firestore:"date"`
	Value float64   `json:"value" firestore:"value"`
}

// Holding represents a stock holding in a portfolio
type Holding struct {
	NumShares     float64 `json:"numShares" firestore:"numShares"`
	PurchaseValue float64 `json:"purchaseValue" firestore:"purchaseValue"`
}

// NewPortfolio creates a new portfolio with the given starting cash
func NewPortfolio(startingCash float64) *Portfolio {
	return &Portfolio{
		Cash:                  startingCash,
		Holdings:              make(map[string]*Holding),
		Transactions:          make([]*Transaction, 0),
		TransactionReferences: make([]*firestore.DocumentRef, 0),
	}
}

// Buy adds a stock purchase to the portfolio
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

// Sell removes a stock sale from the portfolio
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

// Execute executes a transaction (buy or sell)
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