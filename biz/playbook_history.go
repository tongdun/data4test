package biz

import (
	"fmt"
)

func RunHistoryPlaybook(id, mode string) (err error) {
	playbook, _ := GetHistoryPlaybook(id)
	playbookInfo, productInfo, err := GetPlRunInfo("history", id)
	privateParameter := productInfo.GetPrivateParameter()
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

	if err != nil {
		return
	}
	isPass := 0
	var lastFile string
	for k := range runApis {
		playbook.Tag = tag + k
		subResult, historyApi, err1 := playbook.RunPlaybookContent(playbook.Product, privateParameter, "history")
		if err1 != nil {
			Logger.Error("%s", err1)
		}
		playbook.HistoryApis = append(playbook.HistoryApis, historyApi)
		if subResult == "fail" {
			isPass++
			err = WritePlaybookHistoryResult(id, subResult, historyApi, mode, productInfo.EnvType, err1)
			if err != nil {
				Logger.Error("%s", err)
			}
			return

		}
		if err1 != nil {
			err = err1
			err = WritePlaybookHistoryResult(id, subResult, historyApi, mode, productInfo.EnvType, err1)
			if err != nil {
				Logger.Error("%s", err)
			}
			return
		}

		if isPass != 0 {
			lastFile = historyApi
		}

	}

	var result string
	if isPass != 0 {
		result = "fail"
	} else {
		result = "pass"
		lastFile = ""
	}

	if playbookInfo.SceneType == 2 {
		result, err = CompareResult(playbook.Apis, mode)
	}
	err = WritePlaybookHistoryResult(id, result, lastFile, mode, productInfo.EnvType, err)
	if err != nil {
		Logger.Error("%s", err)
		return
	}

	if result != "pass" {
		err = fmt.Errorf("test %v", result)
	}

	if err != nil {
		Logger.Error("%s", err)
	}
	return
}
