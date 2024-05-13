package biz

import (
	"fmt"
	"sync"
)

func RunHistoryPlaybook(id, mode string) (err error) {
	playbook, _ := GetHistoryPlaybook(id)
	playbookInfo, productList, err := GetPlRunInfo("history", id)
	if err != nil {
		return
	}

	productInfo := productList[0]

	var runApis []string
	var tag int
	if len(playbook.LastFile) != 0 {
		index := GetSliceIndex(playbook.Apis, playbook.LastFile)
		if index != -1 {
			runApis = playbook.Apis[index:]
			tag = index
		} else {
			runApis = playbook.Apis
		}
	}

	// 数据置为最初状态
	if mode != "continue" {
		runApis = playbook.Apis
		tag = 0
	}

	envType, _ := GetEnvTypeByName(playbook.Product)
	isFail := 0
	var lastFile string
	var result string

	switch playbook.SceneType {
	case 1, 2:
		for k := range runApis {
			playbook.Tag = tag + k
			subResult, historyApi, errTmp := playbook.RunPlaybookContent(envType, "history")
			if errTmp != nil {
				Logger.Error("%s", errTmp)
				if err != nil {
					err = errTmp
				} else {
					err = fmt.Errorf("%s; %s", err, errTmp)
				}
			}
			playbook.HistoryApis = append(playbook.HistoryApis, historyApi)
			if subResult == "fail" {
				errTmp2 := WritePlaybookHistoryResult(id, subResult, historyApi, mode, envType, errTmp)
				if errTmp2 != nil {
					Logger.Error("%s", errTmp2)
					if err != nil {
						err = errTmp2
					} else {
						err = fmt.Errorf("%s; %s", err, errTmp2)
					}
				}
				return
			}
		}
	case 3:
		for k := range runApis {
			playbook.Tag = tag + k
			subResult, historyApi, errTmp := playbook.RunPlaybookContent(envType, "history")
			if errTmp != nil {
				Logger.Error("%s", errTmp)
				if err != nil {
					err = fmt.Errorf("%s执行失败: %s", playbook.Apis[playbook.Tag], errTmp)
				} else {
					err = fmt.Errorf("%s; %s执行失败：%s", err, playbook.Apis[playbook.Tag], errTmp)
				}
			}
			playbook.HistoryApis = append(playbook.HistoryApis, historyApi)
			if subResult == "fail" {
				isFail++
				errTmp2 := WritePlaybookHistoryResult(id, subResult, historyApi, mode, envType, errTmp)
				if errTmp2 != nil {
					Logger.Error("%s", errTmp2)
					if err != nil {
						err = errTmp2
					} else {
						err = fmt.Errorf("%s; %s", err, errTmp2)
					}
				}
			}
		}
	case 4, 5:
		wg := sync.WaitGroup{}
		for k := range runApis {
			wg.Add(1)
			go func(inPlaybook Playbook, id string, startIndex, index, envType int, errIn error) {
				inPlaybook.Tag = startIndex + index
				subResult, historyApi, errTmp := inPlaybook.RunPlaybookContent(envType, "history")
				if errTmp != nil {
					Logger.Error("%s", errTmp)
					if err != nil {
						err = fmt.Errorf("%s执行失败: %s", inPlaybook.Apis[playbook.Tag], errTmp)
					} else {
						err = fmt.Errorf("%s; %s执行失败：%s", err, inPlaybook.Apis[playbook.Tag], errTmp)
					}
				}

				inPlaybook.HistoryApis = append(inPlaybook.HistoryApis, historyApi)
				if subResult == "fail" {
					isFail++
					errTmp2 := WritePlaybookHistoryResult(id, subResult, historyApi, mode, envType, errTmp)
					if errTmp2 != nil {
						Logger.Error("%s", errTmp2)
						if errIn != nil {
							errIn = errTmp2
						} else {
							errIn = fmt.Errorf("%s; %s", errIn, errTmp2)
						}
					}
				}
				if errIn != nil {
					Logger.Error("%s", errIn)
				}
			}(playbook, id, tag, k, envType, err)
		}
		wg.Wait()
	}

	if err == nil && isFail == 0 {
		result = "pass"
		lastFile = " "
	} else {
		result = "fail"
	}

	if err == nil && isFail == 0 {
		result = "pass"
		lastFile = " "
	} else {
		result = "fail"
	}

	var errTmp error
	if playbookInfo.SceneType == 2 || playbookInfo.SceneType == 4 {
		Logger.Debug("开始比较")
		result, errTmp = CompareResult(playbook.HistoryApis, mode)
		if errTmp != nil {
			Logger.Error("%s", errTmp)
		}
	}
	Logger.Debug("lastFile: %s", lastFile)
	errTmp = WritePlaybookHistoryResult(id, result, lastFile, mode, productInfo.EnvType, err)

	if errTmp != nil {
		Logger.Error("%s", errTmp)
		return
	}

	if result != "pass" {
		err = fmt.Errorf("测试 %v", result)
	}

	return
}
