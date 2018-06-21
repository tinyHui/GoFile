package fileop

import (
	"sync"
)

type StaticResult struct {
	TotalFileNumber uint64
	CharNumAvg      uint64
	CharNumStd      float64
	WordLengthAvg   uint64
	WordLengthStd   float64
	TotalBytes      uint64
}

type FileStatic interface {
	GetStaticResult(rootDir string) StaticResult
}

type fileStatic struct {
}

func NewFileStatic() *fileStatic {
	return &fileStatic{}
}

func (s *fileStatic) GetStaticResult(rootDir string) StaticResult {
	files := make(chan osFile, 50)

	var wg sync.WaitGroup

	wg.Add(1)
	go dirWalk(rootDir, files, &wg)

	result := &StaticResult{
		TotalFileNumber: 0,
		CharNumAvg:      0,
		CharNumStd:      0,
		WordLengthAvg:   0,
		WordLengthStd:   0,
		TotalBytes:      0,
	}

	go func(files *chan osFile, wg *sync.WaitGroup, result *StaticResult) {
		countFileNumChan := make(chan osFile, 50)
		calcTotalBytesChan := make(chan osFile, 50)
		readContentChan := make(chan osFile, 50)

		go func() {
			for file := range *files {
				countFileNumChan <- file
				calcTotalBytesChan <- file
				readContentChan <- file
			}
			close(countFileNumChan)
			close(calcTotalBytesChan)
			close(readContentChan)
		}()

		wg.Add(1)
		go countFileNumber(countFileNumChan, wg, result)

		wg.Add(1)
		go calcTotalBytes(calcTotalBytesChan, wg, result)

		go func(readContentChan *chan osFile, wg *sync.WaitGroup, result *StaticResult) {
			contentChan := make(chan string, 10)
			alphnumCountChan := make(chan uint, 100)
			wordLengthChan := make(chan uint, 100)

			go readContent(*readContentChan, contentChan)

			wg.Add(1)
			go gatherContent(
				contentChan, wg,
				&gatherEngine{alphnumCountChan, getAlphanumCount},
				&gatherEngine{wordLengthChan, wordLengthCount},
			)

			go func(result *StaticResult) {
				wg.Add(1)
				mean, std := calcMeanAndStd(alphnumCountChan, wg)
				result.CharNumAvg = uint64(mean)
				result.CharNumStd = std
			}(result)
			go func(result *StaticResult) {
				wg.Add(1)
				mean, std := calcMeanAndStd(wordLengthChan, wg)
				result.WordLengthAvg = uint64(mean)
				result.WordLengthStd = std
			}(result)

			wg.Wait()
		}(&readContentChan, wg, result)

	}(&files, &wg, result)

	wg.Wait()

	return *result
}
