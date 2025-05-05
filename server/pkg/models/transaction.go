package models

import (
	"cloud.google.com/go/firestore"
	"time"
)

// Transaction represents a buy or sell transaction for a stock
type Transaction struct {
	Time      time.Time              `json:"time" firestore:"time"`
	NumShares float64                `json:"numShares" firestore:"numShares"`
	UnitCost  float64                `json:"unitCost" firestore:"unitCost"`
	Ticker    string                 `json:"ticker" firestore:"ticker"`
	Action    string                 `json:"action" firestore:"action"`
	Bot       *firestore.DocumentRef `json:"-" firestore:"bot"`
}