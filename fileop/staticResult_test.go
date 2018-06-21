package fileop

import (
	"os"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCountFileNumber(t *testing.T) {
	t.Run("should count files in channel", func(t *testing.T) {
		files := make(chan osFile)

		go func() {
			for i := 0; i < 100; i++ {
				files <- osFile{
					f: &os.File{},
				}
			}
			close(files)
		}()

		var wg sync.WaitGroup

		result := &StaticResult{
			TotalFileNumber: 0,
			CharNumAvg:      0,
			CharNumStd:      0,
			WordLengthAvg:   0,
			WordLengthStd:   0,
			TotalBytes:      0,
		}

		wg.Add(1)
		go countFileNumber(files, &wg, result)

		wg.Wait()

		assert.Equal(t, 100, int(result.TotalFileNumber))
	})
}

func TestCalcTotalBytes(t *testing.T) {
	t.Run("should calculate total bytes", func(t *testing.T) {
		files := make(chan osFile)
		go func() {
			for i := 0; i < 100; i++ {
				m := new(mockFileInfo)
				m.On("Size").Return(100)

				files <- osFile{
					stats: m,
				}
			}
			close(files)
		}()

		var wg sync.WaitGroup

		result := &StaticResult{
			TotalFileNumber: 0,
			CharNumAvg:      0,
			CharNumStd:      0,
			WordLengthAvg:   0,
			WordLengthStd:   0,
			TotalBytes:      0,
		}

		wg.Add(1)
		go calcTotalBytes(files, &wg, result)

		wg.Wait()

		assert.Equal(t, uint64(100*100), result.TotalBytes)
	})
}

func TestGatherContent(t *testing.T) {
	t.Run("should gather uint content from channel", func(t *testing.T) {
		contentChan := make(chan string)

		go func() {
			for i := 0; i < 100; i++ {
				contentChan <- "random content"
			}
			close(contentChan)
		}()

		funcCallCount := uint32(0)
		followFunc := func(content string) interface{} {
			atomic.AddUint32(&funcCallCount, 1)
			return uint(0)
		}
		followChan := make(chan uint)

		gatherOp := &gatherEngine{
			channel:  followChan,
			function: followFunc,
		}

		var wg sync.WaitGroup
		wg.Add(1)
		go gatherContent(contentChan, &wg, gatherOp)

		followElemCount := 0
		for range followChan {
			followElemCount += 1
		}

		wg.Wait()

		assert.Equal(t, 100, int(funcCallCount))
		assert.Equal(t, 100, int(followElemCount))
	})

	t.Run("should gather uint from list content from channel", func(t *testing.T) {
		contentChan := make(chan string)

		go func() {
			for i := 0; i < 100; i++ {
				contentChan <- "random content"
			}
			close(contentChan)
		}()

		funcCallCount := uint32(0)
		followFunc := func(content string) interface{} {
			atomic.AddUint32(&funcCallCount, 1)
			return []uint{0, 0, 0}
		}
		followChan := make(chan uint)

		gatherOp := &gatherEngine{
			channel:  followChan,
			function: followFunc,
		}

		var wg sync.WaitGroup
		wg.Add(1)
		go gatherContent(contentChan, &wg, gatherOp)

		followElemCount := 0
		for range followChan {
			followElemCount += 1
		}

		wg.Wait()

		assert.Equal(t, 100, int(funcCallCount))
		assert.Equal(t, 300, int(followElemCount))
	})

	t.Run("should use all gather operations", func(t *testing.T) {
		contentChan := make(chan string)

		go func() {
			for i := 0; i < 100; i++ {
				contentChan <- "random content"
			}
			close(contentChan)
		}()

		func1CallCount := uint32(0)
		func2CallCount := uint32(0)
		followFunc1 := func(content string) interface{} {
			atomic.AddUint32(&func1CallCount, 1)
			return uint(1)
		}
		followFunc2 := func(content string) interface{} {
			atomic.AddUint32(&func2CallCount, 1)
			return []uint{2, 2, 2}
		}
		followChan1 := make(chan uint)
		followChan2 := make(chan uint)

		gatherOp1 := &gatherEngine{
			channel:  followChan1,
			function: followFunc1,
		}
		gatherOp2 := &gatherEngine{
			channel:  followChan2,
			function: followFunc2,
		}

		var wg sync.WaitGroup
		wg.Add(3)
		go gatherContent(contentChan, &wg, gatherOp1, gatherOp2)

		follow1ElemCount := 0
		follow2ElemCount := 0
		go func(followChan *chan uint, wg *sync.WaitGroup) {
			defer wg.Done()
			for range *followChan {
				follow1ElemCount += 1
			}
		}(&followChan1, &wg)
		go func(followChan *chan uint, wg *sync.WaitGroup) {
			defer wg.Done()
			for range *followChan {
				follow2ElemCount += 1
			}
		}(&followChan2, &wg)

		wg.Wait()

		assert.Equal(t, 100, int(func1CallCount))
		assert.Equal(t, 100, int(func2CallCount))
		assert.Equal(t, 100, int(follow1ElemCount))
		assert.Equal(t, 300, int(follow2ElemCount))
	})
}

func TestGetAlphanumCount(t *testing.T) {
	t.Run("should count alpha numeric characters", func(t *testing.T) {
		content := "abc 123 '-=1 `~]\\ 😀h😀app🤣y"

		count := getAlphanumCount(content).(uint)

		assert.Equal(t, 12, int(count))
	})
}

func TestWordLengthCount(t *testing.T) {
	t.Run("should count word length", func(t *testing.T) {
		content := " abc I'm   good  them've 123 '-=1 `~]\\ 😀h😀app🤣y"

		counts := wordLengthCount(content).([]uint)

		assert.Equal(t, 6, len(counts))
		assert.Equal(t, 3, int(counts[0]))
		assert.Equal(t, 1, int(counts[1]))
		assert.Equal(t, 1, int(counts[2]))
		assert.Equal(t, 4, int(counts[3]))
		assert.Equal(t, 4, int(counts[4]))
		assert.Equal(t, 2, int(counts[5]))
	})
}
