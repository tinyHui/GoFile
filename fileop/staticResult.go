package fileop

import (
	"io/ioutil"
	"regexp"
	"strings"
	"sync"
	"unicode"
)

func countFileNumber(files chan osFile, wg *sync.WaitGroup, result *StaticResult) {
	defer wg.Done()

	for file := range files {
		if file.f != nil {
			result.TotalFileNumber += 1
		}
	}
}

func calcTotalBytes(files chan osFile, wg *sync.WaitGroup, result *StaticResult) {
	defer wg.Done()

	for file := range files {
		fs := file.stats

		result.TotalBytes += uint64(fs.Size())
	}
}

type gatherEngine struct {
	channel  chan uint
	function func(string) interface{}
}

func readContent(files chan osFile, contentChan chan string) {
	defer close(contentChan)

	for file := range files {
		buf, _ := ioutil.ReadFile(file.f.Name())
		contentChan <- string(buf)
	}
}

func gatherContent(contentChan chan string, wg *sync.WaitGroup, gatherOps ...*gatherEngine) {
	defer wg.Done()

	gatherWaiter := new(sync.WaitGroup)

	for content := range contentChan {
		for _, gatherOp := range gatherOps {
			gatherWaiter.Add(1)

			go func(content string, gatherFun *gatherEngine, gatherWait *sync.WaitGroup) {
				defer gatherWaiter.Done()

				vT := gatherFun.function(content)

				v, ok := vT.(uint)
				if ok {
					gatherFun.channel <- v
					return
				}

				vList, ok := vT.([]uint)
				if ok {
					for _, v := range vList {
						gatherFun.channel <- v
					}
					return
				}

			}(content, gatherOp, gatherWaiter)
		}
	}

	gatherWaiter.Wait()

	for _, gatherOp := range gatherOps {
		close(gatherOp.channel)
	}
}

func getAlphanumCount(content string) interface{} {
	count := uint(0)
	for _, c := range content {
		if unicode.IsNumber(c) || unicode.IsLetter(c) {
			count += 1
		}
	}
	return count
}

var IsLetter = regexp.MustCompile(`^[a-zA-Z']+$`).MatchString

func wordLengthCount(content string) interface{} {
	var result []uint

	split := func(r rune) bool {
		return r == ' ' || r == '\''
	}

	for _, runeCombine := range strings.FieldsFunc(content, split) {
		if IsLetter(runeCombine) {
			runeCombine = strings.Replace(runeCombine, "'", "", -1)
			result = append(result, uint(len(runeCombine)))
		}
	}
	return result
}
