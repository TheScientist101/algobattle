package models

import (
	"encoding/json"
	"fmt"
	"slices"
	"time"

	"github.com/puzpuzpuz/xsync/v3"
)

// TickerPeriod represents stock data for a specific ticker and time period
type TickerPeriod struct {
	Open        float64            `json:"open"`
	High        float64            `json:"high"`
	Low         float64            `json:"low"`
	Close       float64            `json:"close"`
	Volume      int64              `json:"volume"`
	AdjClose    float64            `json:"adjClose"`
	AdjHigh     float64            `json:"adjHigh"`
	AdjLow      float64            `json:"adjLow"`
	AdjOpen     float64            `json:"adjOpen"`
	AdjVolume   int64              `json:"adjVolume"`
	DivCash     float64            `json:"divCash"`
	SplitFactor float64            `json:"splitFactor"`
	Indicators  map[string]float64 `json:"indicators,omitempty"`
}

// PackedPeriod represents stock data from the API
type PackedPeriod struct {
	Date        time.Time `json:"date"`
	Open        float64   `json:"open"`
	High        float64   `json:"high"`
	Low         float64   `json:"low"`
	Close       float64   `json:"close"`
	Volume      int64     `json:"volume"`
	AdjClose    float64   `json:"adjClose"`
	AdjHigh     float64   `json:"adjHigh"`
	AdjLow      float64   `json:"adjLow"`
	AdjOpen     float64   `json:"adjOpen"`
	AdjVolume   int64     `json:"adjVolume"`
	DivCash     float64   `json:"divCash"`
	SplitFactor float64   `json:"splitFactor"`
}

// TickerMeta contains metadata about a ticker's data range
type TickerMeta struct {
	Start time.Time `json:"dataStart"`
	End   time.Time `json:"dataEnd"`
}

// Row represents stock data for all tickers at a specific date
type Row struct {
	Date time.Time                           `json:"date"`
	Data *xsync.MapOf[string, *TickerPeriod] `json:"data"`
}

// Compare compares two rows by date
func (r Row) Compare(other *Row) int {
	return r.Date.Compare(other.Date)
}

// PackedRow is a serializable version of Row
type PackedRow struct {
	Date time.Time                `json:"date"`
	Data map[string]*TickerPeriod `json:"data"`
}

// UnmarshalJSON implements the json.Unmarshaler interface
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

// Pack converts a Row to a PackedRow for serialization
func (r *Row) Pack() *PackedRow {
	packedRow := &PackedRow{
		Date: r.Date,
		Data: xsync.ToPlainMapOf(r.Data),
	}

	return packedRow
}

// Unpack converts a PackedRow to a Row
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

// History contains stock data for multiple tickers over time
type History struct {
	Tickers map[string]TickerMeta `json:"tickers"`
	Rows    []*Row                `json:"rows"`
}

// PackedHistory is a serializable version of History
type PackedHistory struct {
	Tickers map[string]TickerMeta `json:"tickers"`
	Rows    []*PackedRow          `json:"rows"`
}

// NewHistory creates a new History
func NewHistory() *History {
	history := &History{
		make(map[string]TickerMeta),
		make([]*Row, 0, 365*5),
	}

	return history
}

// Pack converts a History to a PackedHistory for serialization
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

// Unpack converts a PackedHistory to a History
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

// GetClosestRowBefore finds the row closest to but before the given date
func (h *History) GetClosestRowBefore(date time.Time) (index int, row *Row) {
	if len(h.Rows) == 0 {
		return -1, nil
	}

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

// AddData adds stock data for a ticker to the history
func (h *History) AddData(periods []PackedPeriod, ticker string) {
	if len(periods) == 0 {
		return
	}

	h.Tickers[ticker] = TickerMeta{
		periods[0].Date,
		periods[len(periods)-1].Date,
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
			make(map[string]float64),
		})
	}
}