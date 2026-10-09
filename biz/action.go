package biz

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
	//"github.com/derekgr/hivething"
	"github.com/tealeg/xlsx"
	"gopkg.in/yaml.v2"
)

func WriteDataInFile(filePath, fileType, splitTag string, fields []string) (err error) {
	switch fileType {
	case "csv", "txt":
		var content string
		for _, value := range fields {
			if len(content) == 0 {
				content = fmt.Sprintf("%v", value)
			} else {
				content = fmt.Sprintf("%s%s%v", content, splitTag, value)
			}
		}
		WriteDataInCommonFile(filePath, content)
	case "xls":
		WriteDataInXls(filePath, fields)
	}

	return
}

func WriteDataInCommonFile(filePath, content string) (err error) {
	fileHandle, _ := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	defer fileHandle.Close()

	write := bufio.NewWriter(fileHandle)
	contentWithLineFeed := fmt.Sprintf("%s\n", content)
	_, _ = write.WriteString(contentWithLineFeed)
	write.Flush()

	return
}

func WriteDataInXls(filePath string, fields []string) (err error) {
	sheetName := "Sheet1"
	_, err = os.Stat(filePath)
	if err == nil {
		file, errTmp := xlsx.OpenFile(filePath)
		if errTmp != nil {
			err = errTmp
			return
		}
		sheet := file.Sheet[sheetName]
		row := sheet.AddRow()
		for _, item := range fields {
			cell := row.AddCell()
			cell.Value = item
		}
		file.Save(filePath)
	} else {
		file := xlsx.NewFile()
		sheet, _ := file.AddSheet(sheetName)
		row := sheet.AddRow()
		for _, item := range fields {
			cell := row.AddCell()
			cell.Value = item
		}
		file.Save(filePath)
	}
	return
}

// GetTrigger 返回动作的触发时机：before/after/only。
// 未设置时按类型取默认值：sleep 默认 after，其余默认 before。
func (a SceneAction) GetTrigger() string {
	switch a.Trigger {
	case "before", "after", "only":
		return a.Trigger
	default:
		if a.Type == "sleep" {
			return "after"
		}
		return "before"
	}
}

// IsOnlyGenerateMode 文件内任一 action 设置 trigger: only 时，整个执行跳过 HTTP 请求。
func (df DataFile) IsOnlyGenerateMode() bool {
	for _, a := range df.Action {
		if a.Trigger == "only" {
			return true
		}
	}
	return false
}

func (df DataFile) SetSleepAction() (err error) {
	return df.setSleepAction("")
}

func (df DataFile) SetSleepActionByTrigger(trigger string) (err error) {
	return df.setSleepAction(trigger)
}

func (df DataFile) setSleepAction(trigger string) (err error) {
	if len(df.Action) == 0 {
		return
	}
	for _, item := range df.Action {
		if item.Type != "sleep" {
			continue
		}
		if trigger != "" && item.GetTrigger() != trigger {
			continue
		}
		valueType := fmt.Sprintf("%T", item.Value)
		var sleepSecond int
		if valueType == "string" {
			sleepSecondStr := item.Value.(string)
			sleepSecond, _ = strconv.Atoi(sleepSecondStr)
		} else {
			sleepSecond = item.Value.(int)
		}
		//Logger.Debug("开始Sleep")
		time.Sleep(time.Duration(sleepSecond) * time.Second)
		//Logger.Debug("结束Sleep")
	}

	return
}

func (df DataFile) ChangeOutputValue(outputRaw map[string][]interface{}) (outputMap map[string][]interface{}, err error) {
	if len(df.Action) > 0 {
		for _, item := range df.Action {
			if item.Type == "change_output" {
				valueType := fmt.Sprintf("%T", item.Value)
				if valueType == "string" {
					vlaueStr := item.Value.(string)
					tmpList := strings.Split(vlaueStr, ":")
					if len(tmpList) < 3 {
						err = fmt.Errorf(T("error.change_output_invalid"))
						Logger.Error("%v", err)
						return
					}
					var valueKey, old, new, newValue string
					var changeNum int
					valueKey = tmpList[0]
					old = tmpList[1]
					new = tmpList[2]
					if len(tmpList) >= 4 {
						changeNumStr := tmpList[3]
						changeNum, err = strconv.Atoi(changeNumStr)
					}
					if _, ok := outputRaw[valueKey]; ok {
						var newList []interface{}
						for _, subValue := range outputRaw[valueKey] {
							tmpStr := Interface2Str(subValue)
							if changeNum > 0 {
								newValue = strings.Replace(tmpStr, old, new, changeNum)
							} else {
								newValue = strings.Replace(tmpStr, old, new, -1)
							}
							newList = append(newList, newValue)
						}
						outputRaw[valueKey] = newList
					}
				}
			}
		}
	}

	return outputRaw, err
}

func (df DataFile) RecordDataOrderByKey(bodys []map[string]interface{}) (err error) {
	if len(bodys) == 0 {
		return
	}

	var isRecordCSV, isRecordXLS, isRecordTxt bool
	var csvValue, xlsValue, txtValue interface{}

	if len(df.Action) > 0 {
		for _, item := range df.Action {
			if item.Type == "record_csv" {
				isRecordCSV = true
				csvValue = item.Value
			} else if item.Type == "record_xls" || item.Type == "record_excel" || item.Type == "record_xlsx" {
				isRecordXLS = true
				xlsValue = item.Value
			} else if item.Type == "record_txt" {
				isRecordTxt = true
				txtValue = item.Value
			}
		}
	}

	if isRecordCSV {
		df.RecordTargetFile("csv", csvValue)
	}

	if isRecordXLS {
		df.RecordTargetFile("xls", xlsValue)
	}

	if isRecordTxt {
		df.RecordTargetFile("txt", txtValue)
	}

	return
}

func (df DataFile) ModifyFileWithData(bodys []map[string]interface{}) (err error) {
	if len(bodys) == 0 {
		return
	}

	var isRecordFile bool
	var fileValue interface{}

	if len(df.Action) > 0 {
		for _, item := range df.Action {
			if item.Type == "modify_file" {
				isRecordFile = true
				fileValue = item.Value
				break // 如有多文件模板需要修改，到时再变更代码支持
			}
		}
	}

	if isRecordFile {
		var templateName, targetName string
		tmpValue := Interface2Str(fileValue)
		if len(tmpValue) == 0 {
			err = fmt.Errorf(T("error.modify_file_undefined"))
			Logger.Error("%s", err)
			return
		}

		tmps := strings.Split(tmpValue, ":") // name.xml:name_{uniValueVarName}.xml

		tmpBody := make(map[string]interface{})
		tmpBody = bodys[0]

		if len(tmps) >= 2 {
			templateName = tmps[0]
			targetName = tmps[1]
			comReg := regexp.MustCompile(`\{(.+)\}`)
			comMatch := comReg.FindAllSubmatch([]byte(targetName), -1)
			if len(comMatch) > 0 {
				for i := range comMatch {
					var ret string
					dataName := string(comMatch[i][1])
					rawStrDef := string(comMatch[i][0])
					if _, ok := tmpBody[dataName]; ok {
						ret = Interface2Str(tmpBody[dataName])
						targetName = strings.Replace(targetName, rawStrDef, ret, -1)
					} else {
						err = fmt.Errorf(T("error.variable_not_found"), dataName)
						Logger.Error("%s", err)
						return
					}
				}
			} else {
				Logger.Warning(T("warning.no_data_to_replace"), targetName)
			}
		} else {
			err = fmt.Errorf(T("error.modify_file_incomplete"))
			Logger.Error("%s", err)
			return
		}

		templateFilePath := fmt.Sprintf("%s/%s", UploadBasePath, templateName)
		targetFilePath := fmt.Sprintf("%s/%s", DownloadBasePath, targetName)

		content, errTmp := ioutil.ReadFile(templateFilePath)
		if errTmp != nil {
			err = fmt.Errorf("Error: %s, filePath: %s", errTmp, templateFilePath)
			Logger.Error("%s", err)
			return
		}

		strByte := []byte(content)
		newStr := string(content)
		// 匹配字符串
		comReg := regexp.MustCompile(`\{(.+)\}`)
		comMatch := comReg.FindAllSubmatch(strByte, -1)
		if len(comMatch) > 0 {
			for i := range comMatch {
				var ret string
				dataName := string(comMatch[i][1])
				rawStrDef := string(comMatch[i][0])
				if _, ok := tmpBody[dataName]; ok {
					ret = Interface2Str(tmpBody[dataName])
				} else {
					err = fmt.Errorf(T("error.variable_not_found"), dataName)
					Logger.Error("%s", err)
					return
				}

				if len(ret) > 0 {
					newStr = strings.Replace(newStr, rawStrDef, ret, -1)
				}
			}
			err = ioutil.WriteFile(targetFilePath, []byte(newStr), 0644)
			if err != nil {
				Logger.Error("%s", err)
			}
		} else {
			Logger.Warning(T("warning.no_placeholder_data"), templateFilePath)
		}
	}
	return
}

func (df DataFile) CreateDataOrderByKey(lang, dataFileName string, depOutVars map[string][]interface{}) (err error) { // 单线程速度太慢，待重构
	var isCreateCSV, isCreateXLS, isCreateHiveSQL, isCreateTxt bool
	var csvValue, xlsValue, hiveSQLValue, txtValue interface{}

	if len(df.Action) > 0 {
		for _, item := range df.Action {
			if item.Type == "create_csv" {
				isCreateCSV = true
				csvValue = item.Value
			} else if item.Type == "create_excel" || item.Type == "create_xls" || item.Type == "create_xlsx" {
				isCreateXLS = true
				xlsValue = item.Value
			} else if item.Type == "create_hive_table_sql" {
				isCreateHiveSQL = true
				hiveSQLValue = item.Value
			} else if item.Type == "create_txt" {
				isCreateTxt = true
				txtValue = item.Value
			}
		}
	}

	if isCreateCSV {
		CreateTargetFile(lang, dataFileName, "csv", csvValue, depOutVars)
	}

	if isCreateXLS {
		CreateTargetFile(lang, dataFileName, "xls", xlsValue, depOutVars)
	}

	if isCreateHiveSQL {
		df.RecordTargetFile("sql", hiveSQLValue)
	}

	if isCreateTxt {
		CreateTargetFile(lang, dataFileName, "txt", txtValue, depOutVars)
	}

	return
}

func CreateTargetFile(lang, dataFileName, createType string, target interface{}, depOutVars map[string][]interface{}) (err error) {
	fileName := ""
	dataCount := 0
	tmpValue := Interface2Str(target)
	strList := strings.Split(tmpValue, ":")
	if len(strList) == 0 {
		err = fmt.Errorf(T("error.value_undefined"), createType)
		return
	} else {
		tmpFileName := strings.Split(strList[0], ".")
		if len(tmpFileName) < 2 {
			switch createType {
			case "csv":
				fileName = fmt.Sprintf("%s.csv", strList[0])
			case "txt":
				fileName = fmt.Sprintf("%s.txt", strList[0])
			case "xls":
				fileName = fmt.Sprintf("%s.xls", strList[0])
			}
		} else {
			fileName = strList[0]
		}

	}

	if len(strList) >= 2 {
		dataCount, _ = strconv.Atoi(strList[1])
	} else {
		dataCount = 100
	}

	if dataCount > 0 {
		filePath := fmt.Sprintf("%s/%s", UploadBasePath, fileName)
		var keyList []string
		splitTag := ","
		_, errFile := os.Stat(filePath)
		if errFile == nil {
			os.Remove(filePath)
		}

		content, errTmp := GetDataFileRawContent(dataFileName)
		if errTmp != nil {
			err = errTmp
			return
		}
		orderedKeys := GetBodyKeyOrder(content, strings.HasSuffix(dataFileName, ".json"))
		for i := 0; i < dataCount; i++ {
			bodys, errTmp := GetBodyFromRawContent(lang, dataFileName, content, depOutVars)
			if errTmp != nil {
				err = errTmp
				return
			}
			for index, item := range bodys {
				if i == 0 && index == 0 {
					for k, v := range item {
						keyList = append(keyList, k)
						vStr := Interface2Str(v)
						if strings.Contains(vStr, ",") {
							splitTag = "|"
						}
					}
					keyList = sortKeysByFileOrder(keyList, orderedKeys)
					_ = WriteDataInFile(filePath, createType, splitTag, keyList)

				}
				var valueList []string
				for _, k := range keyList {
					valueStr := Interface2Str(item[k])

					if len(valueStr) == 0 {
						valueList = append(valueList, " ")
					} else {
						valueList = append(valueList, valueStr)
					}
				}
				WriteDataInFile(filePath, createType, splitTag, valueList)
			}
		}
	}
	return
}

// GetBodyKeyOrder 从数据文件原始内容中，按书写顺序（从上到下）提取 body 字段名列表。
// 顺序优先级与 GetBody 合并逻辑一致：single.body 的字段在前，multi.body 中未重复的字段在后。
func GetBodyKeyOrder(content string, isJSON bool) []string {
	if isJSON {
		return bodyKeyOrderFromJSON(content)
	}
	return bodyKeyOrderFromYAML(content)
}

// bodyKeyOrderFromYAML 用 yaml.MapSlice 保序解析，取出 single.body / multi.body 的字段书写顺序。
func bodyKeyOrderFromYAML(content string) []string {
	var root yaml.MapSlice
	if err := yaml.Unmarshal([]byte(content), &root); err != nil {
		return nil
	}
	return mergeBodyKeyOrder(
		findYamlMap(root, "single", "body"),
		findYamlMap(root, "multi", "body"),
	)
}

// findYamlMap 沿 path 逐层在 yaml.MapSlice 中查找嵌套 map，找不到返回 nil。
func findYamlMap(root yaml.MapSlice, path ...string) yaml.MapSlice {
	cur := root
	for _, p := range path {
		var next yaml.MapSlice
		found := false
		for _, item := range cur {
			if Interface2Str(item.Key) == p {
				if v, ok := item.Value.(yaml.MapSlice); ok {
					next = v
					found = true
				}
				break
			}
		}
		if !found {
			return nil
		}
		cur = next
	}
	return cur
}

// mergeBodyKeyOrder 按给定 body 的先后顺序合并字段名，去重并保持首次出现位置。
func mergeBodyKeyOrder(bodies ...yaml.MapSlice) []string {
	var keys []string
	seen := map[string]bool{}
	for _, body := range bodies {
		for _, item := range body {
			k := Interface2Str(item.Key)
			if !seen[k] {
				seen[k] = true
				keys = append(keys, k)
			}
		}
	}
	return keys
}

// jsonNode 保序解析 JSON 的最小结构：只保留 object 的键序与子节点。
type jsonNode struct {
	isObj bool
	keys  []string
	m     map[string]jsonNode
}

// bodyKeyOrderFromJSON 用 encoding/json 的 Decoder 顺序读取，取出 single.body / multi.body 的字段书写顺序。
func bodyKeyOrderFromJSON(content string) []string {
	dec := json.NewDecoder(strings.NewReader(content))
	dec.UseNumber()
	root, err := parseJSONValue(dec)
	if err != nil || !root.isObj {
		return nil
	}
	var keys []string
	seen := map[string]bool{}
	for _, b := range []jsonNode{root.m["single"].m["body"], root.m["multi"].m["body"]} {
		if !b.isObj {
			continue
		}
		for _, k := range b.keys {
			if !seen[k] {
				seen[k] = true
				keys = append(keys, k)
			}
		}
	}
	return keys
}

// parseJSONValue 顺序解析一个 JSON 值；object 保序，array/scalar 仅消费不保留。
func parseJSONValue(dec *json.Decoder) (jsonNode, error) {
	tok, err := dec.Token()
	if err != nil {
		return jsonNode{}, err
	}
	delim, ok := tok.(json.Delim)
	if !ok {
		return jsonNode{}, nil // 标量
	}
	switch delim {
	case '{':
		n := jsonNode{isObj: true, m: map[string]jsonNode{}}
		for dec.More() {
			keyTok, _ := dec.Token()
			key, _ := keyTok.(string)
			child, err := parseJSONValue(dec)
			if err != nil {
				return jsonNode{}, err
			}
			n.keys = append(n.keys, key)
			n.m[key] = child
		}
		_, _ = dec.Token() // 消费 '}'
		return n, nil
	case '[':
		for dec.More() {
			if _, err := parseJSONValue(dec); err != nil {
				return jsonNode{}, err
			}
		}
		_, _ = dec.Token() // 消费 ']'
		return jsonNode{}, nil
	default:
		return jsonNode{}, nil
	}
}

// sortKeysByFileOrder 按文件书写顺序重排运行时字段：文件内字段按书写顺序在前，其余按字母序追加。
func sortKeysByFileOrder(runtimeKeys, fileOrderKeys []string) []string {
	if len(fileOrderKeys) == 0 {
		sort.Strings(runtimeKeys)
		return runtimeKeys
	}
	orderIdx := map[string]int{}
	for i, k := range fileOrderKeys {
		orderIdx[k] = i
	}
	sorted := make([]string, len(runtimeKeys))
	copy(sorted, runtimeKeys)
	sort.SliceStable(sorted, func(i, j int) bool {
		oi, iok := orderIdx[sorted[i]]
		oj, jok := orderIdx[sorted[j]]
		switch {
		case iok && jok:
			return oi < oj
		case iok && !jok:
			return true
		case !iok && jok:
			return false
		default:
			return sorted[i] < sorted[j]
		}
	})
	return sorted
}

func (df DataFile) RecordTargetFile(createType string, target interface{}) (err error) {
	fileName := ""
	tmpValue := Interface2Str(target)
	strList := strings.Split(tmpValue, ":")
	if len(strList) == 0 {
		err = fmt.Errorf(T("error.value_undefined"), createType)
		return
	} else {
		tmpFileName := strings.Split(strList[0], ".")
		if len(tmpFileName) < 2 {
			switch createType {
			case "csv":
				fileName = fmt.Sprintf("%s.csv", strList[0])
			case "txt":
				fileName = fmt.Sprintf("%s.txt", strList[0])
			case "sql":
				fileName = fmt.Sprintf("%s.sql", strList[0])
			case "xls":
				fileName = fmt.Sprintf("%s.xls", strList[0])
			}
		} else {
			fileName = strList[0]
		}

	}

	filePath := fmt.Sprintf("%s/%s", UploadBasePath, fileName)
	_, errFile := os.Stat(filePath)

	var keyList []string
	tmpBody := make(map[string]interface{})
	splitTag := ","
	bodys, _ := df.GetBody()
	tmpBody = bodys[0]
	for k, v := range tmpBody {
		keyList = append(keyList, k)
		vStr := Interface2Str(v)
		if strings.Contains(vStr, ",") {
			splitTag = "|"
		}
	}

	sort.Strings(keyList)

	if createType == "sql" {
		var parameterType, sqlStr string
		headStr := "CREATE TABLE IF NOT EXISTS default.all_type_data\n(\n"
		tailStr := "\n) PARTITIONED BY (ds STRING) STORED AS TEXTFILE;"
		midStr := ""
		for _, k := range keyList {
			varType := fmt.Sprintf("%T", tmpBody[k])
			switch varType {
			case "float64":
				parameterType = "DOUBLE"
			case "string":
				parameterType = "STRING"
			case "bool":
				parameterType = "BOOLEAN"
			case "int":
				parameterType = "INT"
			default:
				parameterType = "STRING"
			}

			if len(midStr) == 0 {
				midStr = fmt.Sprintf("  %s %s", k, parameterType)
			} else {
				midStr = fmt.Sprintf("%s,\n  %s %s", midStr, k, parameterType)
			}
		}
		sqlStr = fmt.Sprintf("%s%s%s", headStr, midStr, tailStr)

		_ = WriteDataInCommonFile(filePath, sqlStr)
		return
	}

	if os.IsNotExist(errFile) {
		WriteDataInFile(filePath, createType, splitTag, keyList)
	}

	for _, item := range bodys {
		var valueList []string
		for _, k := range keyList {
			valueStr := Interface2Str(item[k])
			valueList = append(valueList, valueStr)
		}
		WriteDataInFile(filePath, createType, splitTag, valueList)
	}

	return
}
