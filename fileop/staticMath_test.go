package fileop

import (
	"math"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCalcMeanAndStd(t *testing.T) {
	t.Run("should get mean and standard deviation", func(t *testing.T) {
		data := []uint{100, 53, 32, 123, 1, 213, 13, 321}
		wg := new(sync.WaitGroup)

		dataChan := make(chan uint)

		go func() {
			for _, d := range data {
				dataChan <- d
			}
			close(dataChan)
		}()

		wg.Add(1)
		meanR, stdR := calcMeanAndStd(dataChan, wg)

		wg.Wait()

		assert.Equal(t, mean(data), meanR)
		assert.Equal(t, std(data), stdR)
	})

	t.Run("should get mean and standard deviation when only get one element", func(t *testing.T) {
		data := []uint{3}
		wg := new(sync.WaitGroup)

		dataChan := make(chan uint)

		go func() {
			for _, d := range data {
				dataChan <- d
			}
			close(dataChan)
		}()

		wg.Add(1)
		meanR, stdR := calcMeanAndStd(dataChan, wg)

		wg.Wait()

		assert.Equal(t, mean(data), meanR)
		assert.Equal(t, 0, int(stdR))
	})
}

func sum(numbers []uint) (total float64) {
	for _, x := range numbers {
		total += float64(x)
	}
	return total
}

func mean(data []uint) float64 {
	return sum(data) / float64(len(data))
}

func std(data []uint) float64 {
	total := 0.0

	meanR := mean(data)

	for _, number := range data {
		total += math.Pow(float64(number)-meanR, 2)
	}
	variance := total / float64(len(data)-1)
	return math.Sqrt(variance)
}
