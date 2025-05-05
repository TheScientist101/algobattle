// Package models defines the data structures used throughout the AlgoBattle application.
// It includes models for portfolios, transactions, stock data, and related entities.
package models

import (
	"cloud.google.com/go/firestore"
	"fmt"
	"time"
)

// Portfolio represents a user's portfolio of stocks.
// It tracks the current account value, cash balance, stock holdings,
// and transaction history.
type Portfolio struct {
	// AccountValue is the total value of the portfolio (cash + holdings)
	AccountValue float64 `json:"accountValue" firestore:"accountValue"`

	// HistoricalAccountValue tracks the portfolio value over time
	HistoricalAccountValue []*AccountValueHistory `json:"historicalAccountValue" firestore:"historicalAccountValue"`

	// Cash is the available cash balance
	Cash float64 `json:"cash" firestore:"cash"`

	// Holdings maps ticker symbols to stock holdings
	Holdings map[string]*Holding `json:"holdings" firestore:"holdings"`

	// Transactions is the list of transactions (not stored in Firestore)
	Transactions []*Transaction `json:"transactions" firestore:"-"`

	// TransactionReferences stores references to transaction documents in Firestore
	TransactionReferences []*firestore.DocumentRef `json:"-" firestore:"transactions"`
}

// AccountValueHistory represents a historical account value at a specific date.
// This is used to track portfolio performance over time.
type AccountValueHistory struct {
	Date  time.Time `json:"date" firestore:"date"`   // The date of the valuation
	Value float64   `json:"value" firestore:"value"` // The total portfolio value on that date
}

// Holding represents a stock holding in a portfolio.
// It tracks the number of shares and their average purchase value.
type Holding struct {
	NumShares     float64 `json:"numShares" firestore:"numShares"`         // Number of shares held
	PurchaseValue float64 `json:"purchaseValue" firestore:"purchaseValue"` // Average purchase price per share
}

// NewPortfolio creates a new portfolio with the given starting cash.
// It initializes all the necessary maps and slices for a new portfolio.
func NewPortfolio(startingCash float64) *Portfolio {
	return &Portfolio{
		Cash:                  startingCash,
		Holdings:              make(map[string]*Holding),
		Transactions:          make([]*Transaction, 0),
		TransactionReferences: make([]*firestore.DocumentRef, 0),
	}
}

// Buy adds a stock purchase to the portfolio.
// It validates the transaction, updates the cash balance, and adds or updates
// the holding in the portfolio. The purchase value is recalculated as a weighted
// average when adding to an existing position.
func (p *Portfolio) Buy(transaction *Transaction) error {
	// Validate the transaction
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

// Sell removes shares from a stock holding in the portfolio.
// It validates the transaction, updates the cash balance, and reduces
// the number of shares in the holding.
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

// Execute executes a transaction (buy or sell) on the portfolio.
// It routes the transaction to the appropriate handler based on the action.
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
