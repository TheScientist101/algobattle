package indicators

import (
	"fmt"
	"urjith.dev/algobattle/pkg/models"
)

// MACD represents a Moving Average Convergence Divergence indicator
type MACD struct {
	ShortPeriod int
	LongPeriod  int
}

// Name returns the name of the indicator
func (macd *MACD) Name() string {
	return fmt.Sprintf("MACD %d %d", macd.ShortPeriod, macd.LongPeriod)
}

// Apply applies the MACD indicator to the given rows
func (macd *MACD) Apply(rows []*models.Row, getTarget func(index int) float64, setValue func(index int, value float64), _ func(index int, indicator string) float64) {
	if macd.ShortPeriod >= macd.LongPeriod {
		panic("MACD shortPeriod should be less than longPeriod")
	}

	shortEMAIndicator := &EMA{2, macd.ShortPeriod}
	longEMAIndicator := &EMA{2, macd.LongPeriod}

	shortEMAs := make([]float64, len(rows))
	longEMAs := make([]float64, len(rows))

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
		if i < macd.LongPeriod {
			continue
		}

		setValue(i, shortEMAs[i]-longEMAs[i])
	}
}