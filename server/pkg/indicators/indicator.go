package indicators

import (
	"urjith.dev/algobattle/pkg/models"
)

// Indicator is an interface for stock indicators like EMA and MACD
type Indicator interface {
	// Name returns the name of the indicator
	Name() string

	// Apply applies the indicator to the given rows
	Apply(rows []*models.Row, getTarget func(index int) float64, setValue func(index int, value float64), getIndicator func(index int, indicator string) float64)
}

// CalculateIndicators calculates all indicators for the given history
func CalculateIndicators(history *models.History, indicators []Indicator) {
	for ticker, meta := range history.Tickers {
		startIndex, _ := history.GetClosestRowBefore(meta.Start)
		endIndex, _ := history.GetClosestRowBefore(meta.End)

		if startIndex == -1 || endIndex == -1 {
			continue
		}

		getTarget := func(index int) float64 {
			if _, ok := history.Rows[index+startIndex].Data.Load(ticker); !ok {
				return -1
			}

			data, _ := history.Rows[index+startIndex].Data.Load(ticker)
			return data.AdjClose
		}

		getIndicator := func(index int, indicator string) float64 {
			if _, ok := history.Rows[index+startIndex].Data.Load(ticker); !ok {
				return -1
			}

			data, _ := history.Rows[index+startIndex].Data.Load(ticker)
			return data.Indicators[indicator]
		}

		for _, indicator := range indicators {
			name := indicator.Name()

			setValue := func(index int, value float64) {
				data, ok := history.Rows[index+startIndex].Data.Load(ticker)

				if !ok {
					return
				}

				if data.Indicators == nil {
					data.Indicators = make(map[string]float64)
				}

				data.Indicators[name] = value
			}

			indicator.Apply(history.Rows[startIndex:endIndex+1], getTarget, setValue, getIndicator)
		}
	}
}
