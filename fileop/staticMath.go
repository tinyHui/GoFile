package fileop

import (
	"math"
	"sync"
)

func calcMeanAndStd(data chan uint, wg *sync.WaitGroup) (float64, float64) {
	defer wg.Done()

	// we use Welford's method to do online calculation
	var n int64 = 0
	var mean float64 = 0
	var M2 float64 = 0
	var delta float64 = 0

	for x := range data {
		n += 1
		delta = float64(x) - mean
		mean += delta / float64(n)
		M2 += delta * (float64(x) - mean)
	}

	var std float64
	if n != 1 {
		std = math.Sqrt(M2 / float64(n-1))
	}
	return mean, std
}
