// Package models defines the data structures used throughout the AlgoBattle application.
// It includes models for portfolios, transactions, stock data, and related entities.
package models

import (
	"encoding/json"
	"fmt"
	"slices"
	"time"

	"github.com/puzpuzpuz/xsync/v3"
)

// TickerPeriod represents stock data for a specific ticker and time period.
// It contains all the OHLC (Open, High, Low, Close) data as well as volume
// and adjustment factors for dividends and splits. It can also store calculated
// technical indicators.
type TickerPeriod struct {
	Open        float64            `json:"open"`                 // Opening price for the period
	High        float64            `json:"high"`                 // Highest price during the period
	Low         float64            `json:"low"`                  // Lowest price during the period
	Close       float64            `json:"close"`                // Closing price for the period
	Volume      int64              `json:"volume"`               // Trading volume for the period
	AdjClose    float64            `json:"adjClose"`             // Adjusted closing price (for dividends/splits)
	AdjHigh     float64            `json:"adjHigh"`              // Adjusted high price
	AdjLow      float64            `json:"adjLow"`               // Adjusted low price
	AdjOpen     float64            `json:"adjOpen"`              // Adjusted opening price
	AdjVolume   int64              `json:"adjVolume"`            // Adjusted volume
	DivCash     float64            `json:"divCash"`              // Cash dividend amount
	SplitFactor float64            `json:"splitFactor"`          // Stock split factor
	Indicators  map[string]float64 `json:"indicators,omitempty"` // Calculated technical indicators
}

// PackedPeriod represents stock data as received from the API.
// It includes a date field and all the price and volume data for that date.
type PackedPeriod struct {
	Date        time.Time `json:"date"`        // The date of this data point
	Open        float64   `json:"open"`        // Opening price
	High        float64   `json:"high"`        // Highest price
	Low         float64   `json:"low"`         // Lowest price
	Close       float64   `json:"close"`       // Closing price
	Volume      int64     `json:"volume"`      // Trading volume
	AdjClose    float64   `json:"adjClose"`    // Adjusted closing price
	AdjHigh     float64   `json:"adjHigh"`     // Adjusted high price
	AdjLow      float64   `json:"adjLow"`      // Adjusted low price
	AdjOpen     float64   `json:"adjOpen"`     // Adjusted opening price
	AdjVolume   int64     `json:"adjVolume"`   // Adjusted volume
	DivCash     float64   `json:"divCash"`     // Cash dividend amount
	SplitFactor float64   `json:"splitFactor"` // Stock split factor
}

// TickerMeta contains metadata about a ticker's data range.
// It tracks the start and end dates of available data for a ticker.
type TickerMeta struct {
	Start time.Time `json:"dataStart"` // First date with available data
	End   time.Time `json:"dataEnd"`   // Last date with available data
}

// Row represents stock data for all tickers at a specific date.
// It uses a thread-safe map to store ticker data for concurrent access.
type Row struct {
	Date time.Time                           `json:"date"` // The date of this data row
	Data *xsync.MapOf[string, *TickerPeriod] `json:"data"` // Map of ticker symbols to their data
}

// Compare compares two rows by date for sorting purposes.
// Returns a negative value if r is earlier than other, zero if they're equal,
// and a positive value if r is later than other.
func (r Row) Compare(other *Row) int {
	return r.Date.Compare(other.Date)
}

// PackedRow is a serializable version of Row that uses a standard map
// instead of xsync.MapOf for JSON serialization.
type PackedRow struct {
	Date time.Time                `json:"date"` // The date of this data row
	Data map[string]*TickerPeriod `json:"data"` // Map of ticker symbols to their data
}

// UnmarshalJSON implements the json.Unmarshaler interface for Row.
// It first unmarshals into a PackedRow and then converts to a Row with a thread-safe map.
func (r *Row) UnmarshalJSON(bytes []byte) error {
	temp := &PackedRow{}

	if err := json.Unmarshal(bytes, &temp); err != nil {
		return fmt.Errorf("failed to unmarshal JSON into MapOf: %v", err)
	}

	r.Date = temp.Date
	r.Data = xsync.NewMapOf[string, *TickerPeriod]()

	for key, value := range temp.Data {
		r.Data.Store(key, value)
	}

	return nil
}

// Pack converts a Row to a PackedRow for serialization.
// This converts the thread-safe map to a regular map for JSON encoding.
func (r *Row) Pack() *PackedRow {
	packedRow := &PackedRow{
		Date: r.Date,
		Data: xsync.ToPlainMapOf(r.Data),
	}

	return packedRow
}

// Unpack converts a PackedRow to a Row.
// This converts the regular map to a thread-safe map for concurrent access.
func (pr *PackedRow) Unpack() *Row {
	row := &Row{
		Date: pr.Date,
		Data: xsync.NewMapOf[string, *TickerPeriod](),
	}

	for key, value := range pr.Data {
		row.Data.Store(key, value)
	}

	return row
}

// History contains stock data for multiple tickers over time.
// It stores metadata about available tickers and a chronological series of data rows.
type History struct {
	Tickers map[string]TickerMeta `json:"tickers"` // Metadata for each ticker
	Rows    []*Row                `json:"rows"`    // Chronological rows of stock data
}

// PackedHistory is a serializable version of History.
// It uses PackedRows instead of Rows for JSON serialization.
type PackedHistory struct {
	Tickers map[string]TickerMeta `json:"tickers"` // Metadata for each ticker
	Rows    []*PackedRow          `json:"rows"`    // Chronological rows of stock data
}

// NewHistory creates a new History instance with initialized maps and slices.
// The rows slice is pre-allocated with capacity for 5 years of daily data.
func NewHistory() *History {
	history := &History{
		make(map[string]TickerMeta),
		make([]*Row, 0, 365*5), // Pre-allocate 5 years of daily data
	}

	return history
}

// Pack converts a History to a PackedHistory for serialization.
// This method converts all Rows to PackedRows for JSON encoding.
func (h *History) Pack() *PackedHistory {
	packedHistory := &PackedHistory{
		Tickers: h.Tickers,
		Rows:    make([]*PackedRow, len(h.Rows), len(h.Rows)),
	}

	for i := range h.Rows {
		packedHistory.Rows[i] = h.Rows[i].Pack()
	}

	return packedHistory
}

// Unpack converts a PackedHistory to a History.
// This method converts all PackedRows to Rows for thread-safe access.
func (ph *PackedHistory) Unpack() *History {
	history := &History{
		Tickers: ph.Tickers,
		Rows:    make([]*Row, len(ph.Rows), len(ph.Rows)),
	}

	for i := range ph.Rows {
		history.Rows[i] = ph.Rows[i].Unpack()
	}

	return history
}

// GetClosestRowBefore finds the row closest to but before the given date.
// It uses binary search to efficiently find the row in the sorted array.
// Returns the index and row if found, or (-1, nil) if not found or history is empty.
func (h *History) GetClosestRowBefore(date time.Time) (index int, row *Row) {
	if len(h.Rows) == 0 {
		return -1, nil
	}

	// Binary search implementation
	left, right := 0, len(h.Rows)-1
	closest := (*Row)(nil)
	target := date.Unix()

	for left <= right {
		mid := left + (right-left)/2

		if h.Rows[mid].Date.Unix() == target {
			return mid, h.Rows[mid]
		} else if h.Rows[mid].Date.Unix() < target {
			closest = h.Rows[mid]
			left = mid + 1
		} else {
			right = mid - 1
		}
	}

	if closest == nil {
		return -1, nil
	}

	return right, closest
}

// AddData adds stock data for a ticker to the history.
// It updates the ticker metadata and inserts the data points in chronological order.
// If a row already exists for a date, the ticker data is added to that row.
func (h *History) AddData(periods []PackedPeriod, ticker string) {
	if len(periods) == 0 {
		return
	}

	h.Tickers[ticker] = TickerMeta{
		periods[0].Date,              // Start date
		periods[len(periods)-1].Date, // End date
	}

	i, _ := h.GetClosestRowBefore(periods[0].Date)

	for _, p := range periods {
		if i == -1 {
			h.Rows = slices.Insert(h.Rows, 0, &Row{p.Date, xsync.NewMapOf[string, *TickerPeriod]()})
			i++
		}

		for len(h.Rows) > i && h.Rows[i].Date.Before(p.Date) {
			i++
		}

		if i == len(h.Rows) {
			h.Rows = slices.Insert(h.Rows, i, &Row{p.Date, xsync.NewMapOf[string, *TickerPeriod]()})
		}

		h.Rows[i].Data.Store(ticker, &TickerPeriod{
			p.Open,
			p.High,
			p.Low,
			p.Close,
			p.Volume,
			p.AdjClose,
			p.AdjHigh,
			p.AdjLow,
			p.AdjOpen,
			p.AdjVolume,
			p.DivCash,
			p.SplitFactor,
			make(map[string]float64), // Initialize empty indicators map
		})
	}
}
