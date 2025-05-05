package indicators

import (
	"fmt"
	"urjith.dev/algobattle/pkg/models"
)

// EMA represents an Exponential Moving Average indicator
type EMA struct {
	Smoothing    int
	PeriodLength int
}

// Name returns the name of the indicator
func (ema *EMA) Name() string {
	return fmt.Sprintf("EMA %d %d", ema.Smoothing, ema.PeriodLength)
}

// Apply applies the EMA indicator to the given rows
func (ema *EMA) Apply(rows []*models.Row, getTarget func(index int) float64, setValue func(index int, value float64), getIndicator func(index int, indicator string) float64) {
	name := ema.Name()

	// Smoothing factor
	sf := float64(ema.Smoothing) / float64(ema.PeriodLength+1)

	sum := 0.0

	for i := 0; i < len(rows); i++ {
		if i < ema.PeriodLength {
			sum += getTarget(i)
			setValue(i, sum/float64(i+1))
		} else {
			setValue(i, getTarget(i)*sf+getIndicator(i-1, name)*(1-sf))
		}
	}
}
