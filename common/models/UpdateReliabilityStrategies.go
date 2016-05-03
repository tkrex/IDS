package models

import (
)
import "math"

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
	deviationSum := 0.0
	for interval := range updateStats.UpdateIntervalsInSeconds {
		deviationSum += math.Abs(float64(interval) - updateStats.AverageUpdateIntervalInSeconds)
	}

	 meanAbsoluteDeviation := float64(deviationSum) / float64(len(updateStats.UpdateIntervalsInSeconds))

	//Reset Array to current average + deviation to avoid memory leak
	if len(updateStats.UpdateIntervalsInSeconds) == 1000 {
		updateStats.UpdateIntervalsInSeconds = make([]float64,0,1000)
		updateStats.UpdateIntervalsInSeconds[0] = updateStats.AverageUpdateIntervalInSeconds + updateStats.UpdateReliability[m.Name()]
	}
	return meanAbsoluteDeviation
}
