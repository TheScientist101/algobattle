package main

import (
	"cloud.google.com/go/firestore"
	"time"
)

type Transaction struct {
	Time      time.Time              `json:"time" firestore:"time"`
	NumShares float64                `json:"numShares" firestore:"numShares"`
	UnitCost  float64                `json:"unitCost" firestore:"unitCost"`
	Ticker    string                 `json:"ticker" firestore:"ticker"`
	Action    string                 `json:"action" firestore:"action"`
	Bot       *firestore.DocumentRef `json:"-" firestore:"bot"`
}
