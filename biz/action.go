package biz

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
	//"github.com/derekgr/hivething"
	"github.com/tealeg/xlsx"
	"os"
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

func (df DataFile) SetSleepAction() (err error) {
	if len(df.Action) > 0 {
		for _, item := range df.Action {
			if item.Type == "sleep" {
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
		}
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
						err = fmt.Errorf("change_output值定义有误，请确认~")
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

		//fileName := ""
		//tmpValue := Interface2Str(csvValue)
		//if len(tmpValue) == 0 {
		//	err = fmt.Errorf("record_csv的值未定义，请先定义")
		//	return
		//} else {
		//	if !strings.Contains(tmpValue, ".csv") {
		//		fileName = fmt.Sprintf("%s.csv", tmpValue)
		//	} else {
		//		fileName = tmpValue
		//	}
		//}
		//
		//filePath := fmt.Sprintf("%s/%s", UploadBasePath, fileName)
		//var keyList []string
		//tmpBody := make(map[string]interface{})
		//splitTag := ","
		//
		//tmpBody = bodys[0]
		//for k, v := range tmpBody {
		//	keyList = append(keyList, k)
		//	vStr := Interface2Str(v)
		//	if strings.Contains(vStr, ",") {
		//		splitTag = "|"
		//	}
		//}
		//sort.Strings(keyList)
		//
		//_, err := os.Stat(filePath)
		//if os.IsNotExist(err) {
		//	keyStr := ""
		//	for index, k := range keyList {
		//		if index == 0 {
		//			keyStr = k
		//		} else {
		//			keyStr = fmt.Sprintf("%s%s%s", keyStr, splitTag, k)
		//		}
		//
		//	}
		//	tStr := fmt.Sprintf("%s\n", keyStr)
		//	_ = WriteDataInCommonFile(filePath, tStr)
		//
		//}
		//
		//for _, item := range bodys {
		//	valueStr := ""
		//	for _, key := range keyList {
		//		if len(valueStr) == 0 {
		//			valueStr = fmt.Sprintf("%v", item[key])
		//		} else {
		//			valueStr = fmt.Sprintf("%s%s%s", valueStr, splitTag, item[key])
		//		}
		//	}
		//	tStr := fmt.Sprintf("%s\n", valueStr)
		//	_ = WriteDataInCommonFile(filePath, tStr)
		//}
	}

	if isRecordXLS {
		df.RecordTargetFile("xls", xlsValue)
		//fileName := ""
		//tmpValue := Interface2Str(xlsValue)
		//if len(tmpValue) == 0 {
		//	err = fmt.Errorf("record_excel的值未定义，请先定义")
		//	return
		//}
		//if !(strings.Contains(tmpValue, ".xlsx") || strings.Contains(tmpValue, ".xls")) {
		//	fileName = fmt.Sprintf("%s.xlsx", tmpValue)
		//} else {
		//	fileName = tmpValue
		//}
		//filePath := fmt.Sprintf("%s/%s", UploadBasePath, fileName)
		//var keyList []string
		//tmpBody := make(map[string]interface{})
		//
		//tmpBody = bodys[0]
		//for k, _ := range tmpBody {
		//	keyList = append(keyList, k)
		//}
		//sort.Strings(keyList)
		//
		//_, err := os.Stat(filePath)
		//if os.IsNotExist(err) {
		//	for _, k := range keyList {
		//		keyList = append(keyList, k)
		//	}
		//	_ = WriteDataInXls(filePath, keyList)
		//}
		//
		//for _, item := range bodys {
		//	var valueList []string
		//	for _, key := range keyList {
		//		vStr := Interface2Str(item[key])
		//		valueList = append(valueList, vStr)
		//	}
		//	_ = WriteDataInXls(filePath, valueList)
		//}

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
			err = fmt.Errorf("modify_file的值未定义，请先定义")
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
						err = fmt.Errorf("未找到变量[%s]定义，请先定义或关联", dataName)
						Logger.Error("%s", err)
						return
					}
				}
			} else {
				Logger.Warning("[%s] 没有需替换的数据，请核对~", targetName)
			}
		} else {
			err = fmt.Errorf("modify_file的值未定义完整，请先定义， e.g: name.xml:name_{uniValueVarName}.xml")
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
					err = fmt.Errorf("未找到变量[%s]定义，请先定义或关联", dataName)
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
			Logger.Warning("[%s] 没有需要替换的数据，如有需要，请先进行占位符定义", templateFilePath)
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
		err = fmt.Errorf("%s的值未定义，请先定义", createType)
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
		tmpBody := make(map[string]interface{})
		splitTag := ","

		_, errFile := os.Stat(filePath)

		content, errTmp := GetDataFileRawContent(dataFileName)
		if errTmp != nil {
			err = errTmp
			return
		}

		for i := 0; i < dataCount; i++ {
			bodys, errTmp := GetBodyFromRawContent(lang, dataFileName, content, depOutVars)
			if errTmp != nil {
				err = errTmp
				return
			}

			if i == 0 {
				tmpBody = bodys[0]
				for k, v := range tmpBody {
					keyList = append(keyList, k)
					vStr := Interface2Str(v)
					if strings.Contains(vStr, ",") {
						splitTag = "|"
					}
				}
				sort.Strings(keyList)
				if os.IsNotExist(errFile) {
					_ = WriteDataInFile(filePath, createType, splitTag, keyList)
				}
			}

			for _, item := range bodys {
				var valueList []string
				for _, k := range keyList {
					valueStr := Interface2Str(item[k])
					valueList = append(valueList, valueStr)
				}
				WriteDataInFile(filePath, createType, splitTag, valueList)
			}
		}
	}
	return
}

func (df DataFile) RecordTargetFile(createType string, target interface{}) (err error) {
	fileName := ""
	tmpValue := Interface2Str(target)
	strList := strings.Split(tmpValue, ":")
	if len(strList) == 0 {
		err = fmt.Errorf("%s的值未定义，请先定义", createType)
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
