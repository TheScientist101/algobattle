package main

import (
	"encoding/json"
	"fmt"
	"log"
	"slices"
	"sync"
	"time"

	"github.com/puzpuzpuz/xsync/v3"
)

type Indicator interface {
	Name() string
	Apply(rows []*Row, getTarget func(index int) float64, setValue func(index int, value float64), getIndicator func(index int, indicator string) float64)
}

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

func NewHistory() *History {
	history := &History{
		make(map[string]TickerMeta),
		make([]*Row, 0, 365*5),
	}

	return history
}

func (h *History) addData(periods []PackedPeriod, ticker string) {
	if len(periods) == 0 {
		log.Println("no periods for", ticker)
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

type TickerMeta struct {
	Start time.Time `json:"dataStart"`
	End   time.Time `json:"dataEnd"`
}

type Row struct {
	Date time.Time                           `json:"date"`
	Data *xsync.MapOf[string, *TickerPeriod] `json:"data"`
}

func (r Row) Compare(other *Row) int {
	return r.Date.Compare(other.Date)
}

type PackedRow struct {
	Date time.Time                `json:"date"`
	Data map[string]*TickerPeriod `json:"data"`
}

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

func (r *Row) Pack() *PackedRow {
	packedRow := &PackedRow{
		Date: r.Date,
		Data: xsync.ToPlainMapOf(r.Data),
	}

	return packedRow
}

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

// History Contains TickerPeriod info plus Indicators, each history represents one ticker
// TODO: one array of rows sorted by Date, each row is a pair of a Date and Period, insertion sort
type History struct {
	Tickers map[string]TickerMeta `json:"tickers"`
	Rows    []*Row                `json:"rows"`
}

type PackedHistory struct {
	Tickers map[string]TickerMeta `json:"tickers"`
	Rows    []*PackedRow          `json:"rows"`
}

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

type EMA struct {
	smoothing    int
	periodLength int
}

func (ema *EMA) Name() string {
	return fmt.Sprintf("EMA %d %d", ema.smoothing, ema.periodLength)
}

// TODO: Signal line?
type MACD struct {
	shortPeriod int
	longPeriod  int
}

func (macd *MACD) Name() string {
	return fmt.Sprintf("MACD %d %d", macd.shortPeriod, macd.longPeriod)
}

func (ema *EMA) Apply(rows []*Row, getTarget func(index int) float64, setValue func(index int, value float64), getIndicator func(index int, indicator string) float64) {
	name := ema.Name()

	// Smoothing factor
	sf := float64(ema.smoothing) / float64(ema.periodLength+1)

	sum := 0.0

	for i := range rows {

		if i < ema.periodLength {
			sum += getTarget(i)
			setValue(i, sum/float64(i+1))
		} else {
			setValue(i, getTarget(i)*sf+getIndicator(i-1, name)*(1-sf))
		}
	}
}

func (h *History) CalculateIndicators(indicators []Indicator) {
	var wg sync.WaitGroup

	wg.Add(len(h.Tickers))

	for ticker, meta := range h.Tickers {
		startIndex, _ := h.GetClosestRowBefore(meta.Start)
		endIndex, _ := h.GetClosestRowBefore(meta.End)

		if startIndex == -1 || endIndex == -1 {
			log.Println("Skipping ticker", ticker, "data does not exist")
			wg.Done()
			continue
		}

		getTarget := func(index int) float64 {
			if _, ok := h.Rows[index+startIndex].Data.Load(ticker); !ok {
				return -1
			}

			data, _ := h.Rows[index+startIndex].Data.Load(ticker)
			return data.AdjClose
		}

		getIndicator := func(index int, indicator string) float64 {
			if _, ok := h.Rows[index+startIndex].Data.Load(ticker); !ok {
				return -1
			}

			data, _ := h.Rows[index+startIndex].Data.Load(ticker)
			return data.Indicators[indicator]
		}

		go func() {
			defer wg.Done()

			for _, indicator := range indicators {
				name := indicator.Name()

				setValue := func(index int, value float64) {
					data, ok := h.Rows[index+startIndex].Data.Load(ticker)

					if !ok {
						return
					}

					if data.Indicators == nil {
						data.Indicators = make(map[string]float64)
					}

					data.Indicators[name] = value
				}

				indicator.Apply(h.Rows[startIndex:endIndex+1], getTarget, setValue, getIndicator)
			}
		}()
	}

	wg.Wait()
}

func (macd *MACD) Apply(rows []*Row, getTarget func(index int) float64, setValue func(index int, value float64), _ func(index int, indicator string) float64) {
	if macd.shortPeriod >= macd.longPeriod {
		panic("MACD shortPeriod should be less than longPeriod")
	}

	shortEMAIndicator := &EMA{2, macd.shortPeriod}

	longEMAIndicator := &EMA{2, macd.longPeriod}

	shortEMAs := make([]float64, len(rows))
	longEMAs := make([]float64, len(rows))

	// TODO: Requirements system?
	shortEMAIndicator.Apply(rows, getTarget, func(index int, value float64) {
		shortEMAs[index] = value
	}, func(index int, _ string) float64 {
		return shortEMAs[index]
	})

	longEMAIndicator.Apply(rows, getTarget, func(index int, value float64) {
		longEMAs[index] = value
	}, func(index int, _ string) float64 {
		return longEMAs[index]
	})

	for i := range rows {
		if i < macd.longPeriod {
			continue
		}

		setValue(i, shortEMAs[i]-longEMAs[i])
	}
}
