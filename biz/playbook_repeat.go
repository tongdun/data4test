package biz

import (
	"fmt"
	"sync"
)

type playbookRepeatResult struct {
	result, lastFile string
	err              error
}

// Workers write separate result slots; reduction happens after Wait, so a later
// successful worker cannot clear an earlier failure (including a 401).
func runConcurrentPlaybookRepeats(playbook Playbook, count, limit int, run func(Playbook) (string, string, error)) (result, lastFile string, err error) {
	results := make([]playbookRepeatResult, count)
	for start := 0; start < count; start += limit {
		var wg sync.WaitGroup
		end := start + limit
		if end > count {
			end = count
		}
		for index := start; index < end; index++ {
			wg.Add(1)
			go func(index int) {
				defer wg.Done()
				defer func() {
					if recovered := recover(); recovered != nil {
						results[index].result = "fail"
						results[index].err = fmt.Errorf(T("error.scene_repeat_panic"), recovered)
					}
				}()
				item := playbook
				item.Apis = append([]string(nil), playbook.Apis...)
				item.HistoryApis = append([]string(nil), playbook.HistoryApis...)
				results[index].result, results[index].lastFile, results[index].err = run(item)
			}(index)
		}
		wg.Wait()
	}
	for _, item := range results {
		if result != "fail" {
			result = item.result
		}
		if item.result == "fail" || item.err != nil {
			result = "fail"
		}
		if item.lastFile != "" {
			lastFile = item.lastFile
		}
		err = combineRunErrors(err, item.err)
	}
	return
}
