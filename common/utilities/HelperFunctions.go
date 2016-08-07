package utilities

import (
	"math"
	"github.com/tkrex/IDS/common/models"
)

func Min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func Max(x, y int) int {
	if x > y {
		return x
	}
	return y
}

func RoundUp(input float64, places int) (newVal float64) {
	var round float64
	pow := math.Pow(10, float64(places))
	digit := pow * input
	round = math.Ceil(digit)
	newVal = round / pow
	return
}



// Returns the first index of the target string `t`, or
// -1 if no match is found.
func Index(array []*models.RealWorldDomain, element *models.RealWorldDomain) int {
	for i, v := range array {
		if v.Name == element.Name {
			return i
		}
	}
	return -1
}

// Returns `true` if the target string t is in the
// slice.
func IncludIncludee(vs []*models.RealWorldDomain, t *models.RealWorldDomain) bool {
	return Index(vs, t) >= 0
}

