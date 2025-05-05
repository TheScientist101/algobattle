// Package models defines the data structures used throughout the AlgoBattle application.
// It includes models for portfolios, transactions, stock data, and related entities.
package models

import (
	"cloud.google.com/go/firestore"
	"time"
)

// Transaction represents a buy or sell transaction for a stock.
// It records all details of the transaction including time, shares, cost,
// ticker symbol, action type (buy/sell), and a reference to the bot that executed it.
type Transaction struct {
	Time      time.Time              `json:"time" firestore:"time"`           // When the transaction occurred
	NumShares float64                `json:"numShares" firestore:"numShares"` // Number of shares bought or sold
	UnitCost  float64                `json:"unitCost" firestore:"unitCost"`   // Price per share at transaction time
	Ticker    string                 `json:"ticker" firestore:"ticker"`       // Stock ticker symbol
	Action    string                 `json:"action" firestore:"action"`       // "buy" or "sell"
	Bot       *firestore.DocumentRef `json:"-" firestore:"bot"`               // Reference to the bot that executed the transaction
}
