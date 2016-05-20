package models

import (
)
import (
	"math"
	"fmt"
)

type UpdateReliabilityStrategy interface {
	Name() string
	Calculate(*UpdateBehavior) float64

}

type MeanAbsoluteDeviation struct {

}

func (m MeanAbsoluteDeviation) Name() string{
	return "Mean Absolute Deviation"
}

func (m MeanAbsoluteDeviation) Calculate(updateStats *UpdateBehavior) float64 {
	if len(updateStats.UpdateIntervalsInSeconds) < 2 {
		return 0.0
	}
	deviationSum := float64(0)
	for _,interval := range updateStats.UpdateIntervalsInSeconds {
		deviation := math.Abs(interval - updateStats.AverageUpdateIntervalInSeconds)
		deviationSum +=  deviation
	}
	 meanAbsoluteDeviation := deviationSum / float64(len(updateStats.UpdateIntervalsInSeconds))

	//Reset Array to current average + deviation to avoid memory leak
	if len(updateStats.UpdateIntervalsInSeconds) == 1000 {
		updateStats.UpdateIntervalsInSeconds = make([]float64,0,1000)
		updateStats.UpdateIntervalsInSeconds[0] = updateStats.AverageUpdateIntervalInSeconds + updateStats.UpdateReliability[m.Name()]
	}
	return meanAbsoluteDeviation
}
