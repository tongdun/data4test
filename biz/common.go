package biz

import (
	"bytes"
	"data4perf/models"
	"encoding/json"
	"fmt"
	"github.com/Knetic/govaluate"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/gohobby/deepcopy"
	"github.com/mritd/chinaid"
	uuid "github.com/satori/go.uuid"
	chIdNo "github.com/sleagon/chinaid"
	"io/ioutil"
	"log"
	"math/rand"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"syreclabs.com/go/faker"
	"time"
)

func CopyMapInterface(src map[string]interface{}) (dst map[string]interface{}) {
	dst = make(map[string]interface{}, len(src))
	for k, v := range src {
		dst[k] = v
	}
	return dst
}

func CopyMapList(src map[string][]interface{}) (dst map[string][]interface{}) {
	dst = make(map[string][]interface{}, len(src))
	for k, v := range src {
		dst[k] = v
	}
	return dst
}

func CopyList(src []interface{}) (dst []interface{}) {
	dst = make([]interface{}, 0, len(src))
	for _, v := range src {
		varType := fmt.Sprintf("%T", v)
		if varType == "map[string]interface{}" {
			vMap := v.(map[string]interface{})
			tMap := CopyMap(vMap)
			dst = append(dst, tMap)
		} else if varType == "map[interface {}]interface {}" {
			tMap := make(map[string]interface{})
			vMap := v.(map[interface{}]interface{})
			for subK, subV := range vMap {
				key := subK.(string)
				tMap[key] = subV
			}
			dst = append(dst, tMap)
		} else {
			dst = deepcopy.Slice(src).Clone()
		}

	}
	return
}

func Interface2Str(value interface{}) (strValue string, err error) {
	varType := fmt.Sprintf("%T", value)
	switch varType {
	case "float64":
		tmpVar := value.(float64)
		strValue = strconv.FormatFloat(tmpVar, 'f', 0, 64)
	case "string":
		strValue = value.(string)
	case "bool":
		strValue = strconv.FormatBool(value.(bool))
	case "int":
		strValue = strconv.Itoa(value.(int))
	//case "interface {}":
	//	newValue, errTmp := json.Marshal(value)
	//	if errTmp != nil {
	//		Logger.Error("%s", errTmp)
	//		return
	//	}
	//	strValue = string(newValue)
	//	Logger.Debug("strValue: %v", strValue)
	//case "[]interface {}":
	//	items := value.([]interface{})
	//	if len(items) > 0 {
	//		subValue := items[0]
	//		subVarType := fmt.Sprintf("%T", subValue)
	//		switch subVarType {
	//		case "float64":
	//			tmpVar := subValue.(float64)
	//			strValue = strconv.FormatFloat(tmpVar, 'f', 0, 64)
	//		case "string":
	//			strValue = subValue.(string)
	//		case "bool":
	//			strValue = strconv.FormatBool(subValue.(bool))
	//		case "int":
	//			strValue = strconv.Itoa(subValue.(int))
	//		default:
	//			err := fmt.Errorf("不支持类型: %v 的转换, 值: %v,  如有需要，请联系管理员~", subVarType, subValue)
	//			Logger.Warning("%s", err)
	//		}
	//	}
	//	//for _, item := range value.([]interface{}) {
	//	//	subVarType := fmt.Sprintf("%T", item)
	//	//	switch subVarType {
	//	//	case "float64":
	//	//		tmpVar := value.(float64)
	//	//		strValue = strconv.FormatFloat(tmpVar, 'f', 0, 64)
	//	//	case "string":
	//	//		strValue = value.(string)
	//	//	case "bool":
	//	//		strValue = strconv.FormatBool(value.(bool))
	//	//	case "int":
	//	//		strValue = strconv.Itoa(value.(int))
	//	//	}
	//	//	break
	//	//	//newValue, errTmp := json.Marshal(item)
	//	//	//if errTmp != nil {
	//	//	//	Logger.Error("%s", errTmp)
	//	//	//	return
	//	//	//}
	//	//	//strValue = item.(string)
	//	//	//strValue = string(newValue)
	//	//	//strValue = strValue + string(newValue)
	//	//}
	//	//Logger.Debug("strValue: %v", strValue)
	default:
		err = fmt.Errorf("不支持类型: %v 的转换, 值: %v,  原样请求，如有需要，请联系管理员~", varType, value)
		Logger.Warning("%s", err)
	}
	return
}

func GetSliceMinLen(info map[string][]interface{}) (minLen int) {
	tag := 0
	for _, v := range info {
		if tag == 0 {
			minLen = len(v)
		} else {
			if minLen > len(v) {
				minLen = len(v)
			}
		}
		tag = tag + 1
	}
	return
}

func StrComparion(sType, curStr, targetStr string) (b bool, err error) {
	switch sType {
	case "=", "equal":
		_, err1 := strconv.ParseFloat(curStr, 64)
		if err1 == nil {
			Expression := "value1==value2"
			Expression = strings.Replace(Expression, "value1", curStr, -1)
			Expression = strings.Replace(Expression, "value2", targetStr, -1)
			expr, err2 := govaluate.NewEvaluableExpression(Expression)
			if err2 != nil {
				err = err2
				Logger.Error("四则运算表达式书写有误，%s", err2)
				return
			}
			result, err3 := expr.Evaluate(nil)
			if err3 != nil {
				err = err3
				Logger.Error("四则运算表达式书写有误，%s", err3)
				return
			}
			if result == true {
				b = true
			} else {
				b = false
			}
		} else {
			//b = false
			if curStr == targetStr {
				b = true
			} else {
				b = false
			}
		}

	case "!=", "not_equal":
		if curStr != targetStr {
			b = true
		} else {
			b = false
		}
	case "in", "contain":
		if strings.Contains(curStr, targetStr) {
			b = true
		} else {
			b = false
		}
	case "!in", "not_in", "not_contain":
		if !strings.Contains(curStr, targetStr) {
			b = true
		} else {
			b = false
		}
	case "re", "regex", "regexp":
		re := regexp.MustCompile(targetStr)
		result := re.FindStringSubmatch(curStr)
		if len(result) > 0 {
			//Logger.Debug("result: %q", result)
			b = true
		} else {
			b = false
		}
	case "null", "empty":
		if len(curStr) == 0 {
			b = true
		} else {
			b = false
		}
	case "!null", "!empty", "not_null", "not_empty":
		if len(curStr) > 0 {
			b = true
		} else {
			b = false
		}
	case ">", "larger_than", "greater_than":
		targetInt, err1 := strconv.Atoi(targetStr)
		curInt, err2 := strconv.Atoi(curStr)
		if err1 != nil || err2 != nil {
			b = false
		} else {
			if curInt > targetInt {
				b = true
			} else {
				b = false
			}
		}
	case ">=", "larger_equal", "greater_equal":
		targetInt, err1 := strconv.Atoi(targetStr)
		curInt, err2 := strconv.Atoi(curStr)
		if err1 != nil || err2 != nil {
			b = false
		} else {
			if curInt >= targetInt {
				b = true
			} else {
				b = false
			}
		}
	case "<", "less_than":
		targetInt, err1 := strconv.Atoi(targetStr)
		curInt, err2 := strconv.Atoi(curStr)
		if err1 != nil || err2 != nil {
			b = false
		} else {
			if curInt < targetInt {
				b = true
			} else {
				b = false
			}
		}
	case "<=", "less_equal":
		targetInt, err1 := strconv.Atoi(targetStr)
		curInt, err2 := strconv.Atoi(curStr)
		if err1 != nil || err2 != nil {
			b = false
		} else {
			if curInt <= targetInt {
				b = true
			} else {
				b = false
			}
		}
	default:
		err = fmt.Errorf("不支持%s类型的比较，如有需要请反馈致相关人员", sType)
		Logger.Error("%s", err)
	}
	//Logger.Info("断言: [%v] %v [%v] 结果:pass", curStr, sType, targetStr)  // 由%s占位符更换成%v, 是否能解决空指针问题
	if !b {
		err = fmt.Errorf("断言: [%v] %v [%v] 结果:fail", curStr, sType, targetStr)
		Logger.Error("%s", err)
	}

	return
}

func GetInStrDef(rawStr string) (allDef map[string]string) {
	strByte := []byte(rawStr)
	allDef = make(map[string]string)
	// 匹配字符串
	strReg := regexp.MustCompile(`\{(\w+)\}`)
	strMatch := strReg.FindAllSubmatch(strByte, -1)
	//支持${xx}写法不替换
	strReg2 := regexp.MustCompile(`\$\{(\w+)\}`)
	strMatch2 := strReg2.FindAllSubmatch(strByte, -1)
	for _, item := range strMatch {
		key := string(item[1])
		rawStrDef := string(item[0])
		isNeed := true
		for _, item2 := range strMatch2 {
			key2 := string(item2[1])
			if key == key2 {
				isNeed = false
			}
		}
		if isNeed == true {
			allDef[key] = rawStrDef
		}
	}

	return
}

func GetSpecialStr(lang, rawStr string) (newStr string) {
	if lang == "en" {
		newStr = GetEnData(rawStr)
		newStr = GetChData(newStr)
	} else {
		newStr = GetChData(rawStr)
		newStr = GetEnData(newStr)
	}

	newStr = GetCommonData(newStr)
	newStr = GetTimeData(newStr)
	newStr = GetLengthData(newStr)
	newStr = GetRangeData(newStr)
	newStr = GetTimeFormatData(newStr)

	allDef := GetInStrDef(newStr)
	for k, v := range allDef {
		value := GetValueFromSysParameter(lang, k)
		if len(value) > 0 {
			newStr = strings.ReplaceAll(newStr, v, value)
		}
	}

	return
}

func GetStrType(lang, s string) (t int, v string, allDef map[string]string) {
	// 判断参数是否为其他用例提供的参数
	if s == "{self}" {
		t = 1
		v = s
		return
	}
	// 判断程序自动生成的一些特征数据
	tmpStr := GetSpecialStr(lang, s)
	if tmpStr != s {
		t = 2
		v = tmpStr
		allDef = GetInStrDef(v)
		if len(allDef) > 0 {
			t = 3
		}
		return
	}
	// 判断字符串中依赖output变量的信息
	allDef = GetInStrDef(s)
	if len(allDef) > 0 {
		t = 3
		return
	}
	return
}

func GetTimeFormatData(rawStr string) (newStr string) {
	strByte := []byte(rawStr)
	newStr = rawStr
	timeFormatReg := regexp.MustCompile(`\{TimeFormat\((.+)\)\}`)
	timeFormatMatch := timeFormatReg.FindAllSubmatch(strByte, -1)
	timeNow := time.Now()
	if len(timeFormatMatch) > 0 {
		for i := range timeFormatMatch {
			formatDef := string(timeFormatMatch[i][1])
			rawStrDef := string(timeFormatMatch[i][0])
			ret := fmt.Sprintf(timeNow.Format(formatDef))
			newStr = strings.Replace(newStr, rawStrDef, ret, -1)
		}
	}

	return
}

func GetHistoryDataDirName(fileName string) (dirName string) {
	b, num := IsStrEndWithTimeFormat(fileName)
	suffix := GetStrSuffix(fileName)

	if b {
		dirName = fileName[:len(fileName)-num-len(suffix)]
	} else {
		dirName = fileName[:len(fileName)-len(suffix)]
	}
	return
}

func GetStrSuffix(s string) (suffix string) {
	if strings.HasSuffix(s, ".json") {
		suffix = ".json"
	} else if strings.HasSuffix(s, ".yaml") {
		suffix = ".yaml"
	} else if strings.HasSuffix(s, ".yml") {
		suffix = ".yml"
	}

	return
}

func IsStrEndWithTimeFormat(s string) (b bool, num int) {
	timeReg, err := regexp.Compile(`_[0-9]{8}_[0-9]{6}\.[0-9]+\.`)
	if err != nil {
		Logger.Error("%v", err)
	}
	timeMatch := timeReg.FindString(s)
	if len(timeMatch) > 0 {
		b = true
		num = len(timeMatch) - 1
	}

	return
}

func GetLengthData(rawStr string) (newStr string) {
	strByte := []byte(rawStr)
	newStr = rawStr
	// 匹配字符串
	comReg := regexp.MustCompile(`\{([A-Z][A-Za-z]+)\((-*\d+)\)\}`)
	comMatch := comReg.FindAllSubmatch(strByte, -1)
	if len(comMatch) > 0 {
		for i := range comMatch {
			var ret string
			dataType := string(comMatch[i][1])
			num, errTmp := strconv.Atoi(string(comMatch[i][2]))
			if errTmp != nil {
				Logger.Error("%v", num)
			}
			rawStrDef := string(comMatch[i][0])
			switch dataType {
			case "Rune":
				ret = GetRandomRune(num)
			case "Str":
				ret = GetRandomStr(num, "")
			case "UpperStr":
				ret = GetRandomStr(num, "upper")
			case "LowerStr":
				ret = GetRandomStr(num, "lower")
			case "IntStr":
				ret = GetRandomStr(num, "int")
			case "UpperLowerStr":
				ret = GetRandomStr(num, "upperLower")
			case "Timestamp":
				curTimestamp := time.Now().Unix()
				targetTimestamp := (int(curTimestamp) + 86400*num) * 1000
				ret = strconv.Itoa(targetTimestamp)
			case "Date":
				curTimestamp := time.Now().Unix() + int64(86400*num)
				ret = fmt.Sprintf(time.Unix(curTimestamp, 0).Format("2006-01-02"))
			case "Time":
				curTimestamp := time.Now().Unix() + int64(86400*num)
				ret = fmt.Sprintf(time.Unix(curTimestamp, 0).Format("2006-01-02 15:04:05"))
			case "Minute":
				curTimestamp := time.Now().Unix() + int64(60*num)
				ret = fmt.Sprintf(time.Unix(curTimestamp, 0).Format("2006-01-02 15:04:05"))
			case "Hour":
				curTimestamp := time.Now().Unix() + int64(3600*num)
				ret = fmt.Sprintf(time.Unix(curTimestamp, 0).Format("2006-01-02 15:04:05"))
			case "Second":
				curTimestamp := time.Now().Unix() + int64(1*num)
				ret = fmt.Sprintf(time.Unix(curTimestamp, 0).Format("2006-01-02 15:04:05"))
			case "Month":
				curTimestamp := time.Now().Unix() + int64(86400*30*num)
				ret = fmt.Sprintf(time.Unix(curTimestamp, 0).Format("2006-01"))
			case "Year":
				curTimestamp := time.Now().Unix() + int64(86400*30*12*num)
				ret = fmt.Sprintf(time.Unix(curTimestamp, 0).Format("2006"))
			case "TimeGMT":
				curTimestamp := time.Now().Unix() + int64(86400*num)
				ret = fmt.Sprintf(time.Unix(curTimestamp, 0).Format("Mon Jan 02 2006 15:04:05 GMT+0800"))
			case "DayBegin":
				curTimestamp := time.Now().Unix() + int64(86400*num)
				retBegin := fmt.Sprintf(time.Unix(curTimestamp, 0).Format("2006-01-02"))
				ret = fmt.Sprintf("%s 00:00:00", retBegin)
			case "DayStampBegin":
				curTimestamp := time.Now().Unix() + int64(86400*num)
				retReplace := fmt.Sprintf(time.Unix(curTimestamp, 0).Format("2006-01-02"))
				retBegin := fmt.Sprintf("%s 00:00:00", retReplace)
				loc, _ := time.LoadLocation("Asia/Shanghai")
				tt, _ := time.ParseInLocation("2006-01-02 15:04:05", retBegin, loc) //2006-01-02 15:04:05是转换的格式如php的"Y-m-d H:i:s"
				ret = fmt.Sprintf("%v", tt.UnixNano()/1e6)
			case "DayEnd":
				curTimestamp := time.Now().Unix() + int64(86400*num)
				retBegin := fmt.Sprintf(time.Unix(curTimestamp, 0).Format("2006-01-02"))
				ret = fmt.Sprintf("%s 23:59:59", retBegin)
			case "DayStampEnd":
				curTimestamp := time.Now().Unix() + int64(86400*num)
				retBegin := fmt.Sprintf(time.Unix(curTimestamp, 0).Format("2006-01-02"))
				retRepalce := fmt.Sprintf("%s 23:59:59", retBegin)
				loc, _ := time.LoadLocation("Asia/Shanghai")
				tt, _ := time.ParseInLocation("2006-01-02 15:04:05", retRepalce, loc) //2006-01-02 15:04:05是转换的格式如php的"Y-m-d H:i:s"
				ret = fmt.Sprintf("%v", tt.UnixNano()/1e6)
			case "MonthBegin":
				month := time.Now().Format("01")
				intMonth, _ := strconv.Atoi(month)
				targetMonth := intMonth + num
				year := time.Now().Format("2006")
				ret = fmt.Sprintf("%s-%2d-01 00:00:00", year, targetMonth)
			case "MonthStampBegin":
				month := time.Now().Format("01")
				intMonth, _ := strconv.Atoi(month)
				targetMonth := intMonth + num
				year := time.Now().Format("2006")
				retReplace := fmt.Sprintf("%s-%2d-01 00:00:00", year, targetMonth)
				loc, _ := time.LoadLocation("Asia/Shanghai")
				tt, _ := time.ParseInLocation("2006-01-02 15:04:05", retReplace, loc) //2006-01-02 15:04:05是转换的格式如php的"Y-m-d H:i:s"
				ret = fmt.Sprintf("%v", tt.UnixNano()/1e6)
			case "MonthEnd":
				month := time.Now().Format("01")
				intMonth, _ := strconv.Atoi(month)
				targetMonth := intMonth + num
				year := time.Now().Format("2006")
				ret = fmt.Sprintf("%s-%2d-31 23:59:59", year, targetMonth)
			case "MonthStampEnd":
				month := time.Now().Format("01")
				intMonth, _ := strconv.Atoi(month)
				targetMonth := intMonth + num
				year := time.Now().Format("2006")
				retReplace := fmt.Sprintf("%s-%2d-31 23:59:59", year, targetMonth)
				loc, _ := time.LoadLocation("Asia/Shanghai")
				tt, _ := time.ParseInLocation("2006-01-02 15:04:05", retReplace, loc) //2006-01-02 15:04:05是转换的格式如php的"Y-m-d H:i:s"
				ret = fmt.Sprintf("%v", tt.UnixNano()/1e6)
			case "YearBegin":
				year := time.Now().Format("2006")
				intYear, _ := strconv.Atoi(year)
				targetYear := intYear + num
				ret = fmt.Sprintf("%d-01-01 00:00:00", targetYear)
			case "YearStampBegin":
				year := time.Now().Format("2006")
				intYear, _ := strconv.Atoi(year)
				targetYear := intYear + num
				retReplace := fmt.Sprintf("%d-01-01 00:00:00", targetYear)
				loc, _ := time.LoadLocation("Asia/Shanghai")
				tt, _ := time.ParseInLocation("2006-01-02 15:04:05", retReplace, loc) //2006-01-02 15:04:05是转换的格式如php的"Y-m-d H:i:s"
				ret = fmt.Sprintf("%v", tt.UnixNano()/1e6)
			case "YearEnd":
				year := time.Now().Format("2006")
				intYear, _ := strconv.Atoi(year)
				targetYear := intYear + num
				ret = fmt.Sprintf("%d-12-31 23:59:59", targetYear)
			case "YearStampEnd":
				year := time.Now().Format("2006")
				intYear, _ := strconv.Atoi(year)
				targetYear := intYear + num
				retReplace := fmt.Sprintf("%d-12-31 23:59:59", targetYear)
				loc, _ := time.LoadLocation("Asia/Shanghai")
				tt, _ := time.ParseInLocation("2006-01-02 15:04:05", retReplace, loc) //2006-01-02 15:04:05是转换的格式如php的"Y-m-d H:i:s"
				ret = fmt.Sprintf("%v", tt.UnixNano()/1e6)
			}
			if len(ret) > 0 {
				newStr = strings.Replace(newStr, rawStrDef, ret, -1)
			}

		}
	}

	return
}

func GetEnData(rawStr string) (newStr string) {
	strByte := []byte(rawStr)
	newStr = rawStr
	// 匹配字符串
	comReg := regexp.MustCompile(`\{([A-Z][A-Za-z]+)\}`)
	comMatch := comReg.FindAllSubmatch(strByte, -1)
	if len(comMatch) > 0 {
		for i := range comMatch {
			var ret string
			dataType := string(comMatch[i][1])
			rawStrDef := string(comMatch[i][0])
			switch dataType {
			case "Name":
				ret = gofakeit.Name()
			case "Email":
				ret = gofakeit.Email()
			case "Phone", "Mobile", "Contact":
				ret = gofakeit.Phone()
			case "Gender", "Sex":
				ret = gofakeit.Gender()
			case "CardNo", "BankNo":
				ret = faker.Finance().CreditCard(faker.CC_VISA)
			case "SSN", "IDNo":
				ret = gofakeit.Person().SSN
			case "Address":
				ret = faker.Address().String()
			case "Company":
				ret = gofakeit.Company()
			case "Country":
				ret = gofakeit.Country()
			case "Province", "State":
				ret = gofakeit.State()
			case "City":
				ret = gofakeit.City()
			}

			if len(ret) > 0 {
				newStr = strings.Replace(newStr, rawStrDef, ret, -1)
			}

		}
	}

	return
}

func GetCommonData(rawStr string) (newStr string) {
	strByte := []byte(rawStr)
	newStr = rawStr
	// 匹配字符串
	comReg := regexp.MustCompile(`\{([A-Z][A-Za-z]+)\}`)
	comMatch := comReg.FindAllSubmatch(strByte, -1)
	if len(comMatch) > 0 {
		for i := range comMatch {
			var value string
			dataType := string(comMatch[i][1])
			rawStrDef := string(comMatch[i][0])
			switch dataType {
			case "QQ":
				value = GetRandomStr(10, "int")
			case "Age":
				value = strconv.Itoa(GetRandomInt(2, 130))
			case "IntStr10":
				value = GetRandomStr(10, "int")
			case "Int":
				value = strconv.Itoa(GetRandomInt(0, 100))
			case "Income":
				value = strconv.Itoa(GetRandomInt(1850, 50000))
			case "Uuid":
				u2 := uuid.NewV4()
				value = u2.String()
			case "HexColor":
				value = gofakeit.HexColor()
			case "DomainName":
				value = gofakeit.DomainName()
			case "IPv4":
				value = gofakeit.IPv4Address()
			case "IPv6":
				value = gofakeit.IPv6Address()
			case "MacAddress", "MacAddr":
				value = gofakeit.MacAddress()
			case "Emoji":
				value = gofakeit.Emoji()
			}

			if len(value) > 0 {
				newStr = strings.Replace(newStr, rawStrDef, value, -1)
			}

		}
	}

	return
}

func GetTimeData(rawStr string) (newStr string) {
	strByte := []byte(rawStr)
	newStr = rawStr
	// 匹配字符串
	comReg := regexp.MustCompile(`\{([A-Z][A-Za-z]+)\}`)
	comMatch := comReg.FindAllSubmatch(strByte, -1)
	if len(comMatch) > 0 {
		for i := range comMatch {
			var value string
			dataType := string(comMatch[i][1])
			rawStrDef := string(comMatch[i][0])
			switch dataType {
			case "Time":
				value = fmt.Sprintf(time.Unix(time.Now().Unix(), 0).Format("2006-01-02 15:04:05"))
			case "Date":
				value = fmt.Sprintf(time.Unix(time.Now().Unix(), 0).Format("2006-01-02")) // 带时分秒的时间由Time实现
			case "Month":
				value = fmt.Sprintf(time.Unix(time.Now().Unix(), 0).Format("2006-01"))
			case "Year":
				value = fmt.Sprintf(time.Unix(time.Now().Unix(), 0).Format("2006"))
			case "Timestamp":
				curTimestamp := (int(time.Now().Unix())) * 1000
				value = strconv.Itoa(curTimestamp)
			case "DayBegin":
				day := time.Now().Format("2006-01-02")
				value = fmt.Sprintf("%s 00:00:00", day)
			case "DayEnd":
				day := time.Now().Format("2006-01-02")
				value = fmt.Sprintf("%s 23:59:59", day)
			case "MonthBegin":
				day := time.Now().Format("2006-01")
				value = fmt.Sprintf("%s-01 00:00:00", day)
			case "MonthEnd":
				day := time.Now().Format("2006-01")
				value = fmt.Sprintf("%s-31 23:59:59", day)
			case "YearBegin":
				day := time.Now().Format("2006")
				value = fmt.Sprintf("%s-01-01 00:00:00", day)
			case "YearEnd":
				day := time.Now().Format("2006")
				value = fmt.Sprintf("%s-12-31 23:59:59", day)
			case "DayStampBegin":
				day := time.Now().Format("2006-01-02")
				valueTmp := fmt.Sprintf("%s 00:00:00", day)
				loc, _ := time.LoadLocation("Asia/Shanghai")
				tt, _ := time.ParseInLocation("2006-01-02 15:04:05", valueTmp, loc) //2006-01-02 15:04:05是转换的格式如php的"Y-m-d H:i:s"
				value = fmt.Sprintf("%v", tt.UnixNano()/1e6)
			case "DayStampEnd":
				day := time.Now().Format("2006-01-02")
				valueTmp := fmt.Sprintf("%s 23:59:59", day)
				loc, _ := time.LoadLocation("Asia/Shanghai")
				tt, _ := time.ParseInLocation("2006-01-02 15:04:05", valueTmp, loc) //2006-01-02 15:04:05是转换的格式如php的"Y-m-d H:i:s"
				value = fmt.Sprintf("%v", tt.UnixNano()/1e6)
			case "MonthStampBegin":
				day := time.Now().Format("2006-01")
				valueTmp := fmt.Sprintf("%s-01 00:00:00", day)
				loc, _ := time.LoadLocation("Asia/Shanghai")
				tt, _ := time.ParseInLocation("2006-01-02 15:04:05", valueTmp, loc) //2006-01-02 15:04:05是转换的格式如php的"Y-m-d H:i:s"
				value = fmt.Sprintf("%v", tt.UnixNano()/1e6)
			case "MonthStampEnd":
				day := time.Now().Format("2006-01")
				valueTmp := fmt.Sprintf("%s-31 23:59:59", day)
				loc, _ := time.LoadLocation("Asia/Shanghai")
				tt, _ := time.ParseInLocation("2006-01-02 15:04:05", valueTmp, loc) //2006-01-02 15:04:05是转换的格式如php的"Y-m-d H:i:s"
				value = fmt.Sprintf("%v", tt.UnixNano()/1e6)
			case "YearStampBegin":
				day := time.Now().Format("2006")
				valueTmp := fmt.Sprintf("%s-01-01 00:00:00", day)
				loc, _ := time.LoadLocation("Asia/Shanghai")
				tt, _ := time.ParseInLocation("2006-01-02 15:04:05", valueTmp, loc) //2006-01-02 15:04:05是转换的格式如php的"Y-m-d H:i:s"
				value = fmt.Sprintf("%v", tt.UnixNano()/1e6)
			case "YearStampEnd":
				day := time.Now().Format("2006")
				valueTmp := fmt.Sprintf("%s-12-31 23:59:59", day)
				loc, _ := time.LoadLocation("Asia/Shanghai")
				tt, _ := time.ParseInLocation("2006-01-02 15:04:05", valueTmp, loc) //2006-01-02 15:04:05是转换的格式如php的"Y-m-d H:i:s"
				value = fmt.Sprintf("%v", tt.UnixNano()/1e6)
			case "TimeGMT":
				timeNow := time.Now()
				value = fmt.Sprintf(timeNow.Format("Mon Jan 02 2006 15:04:05 GMT+0800"))
			}

			if len(value) > 0 {
				newStr = strings.Replace(newStr, rawStrDef, value, -1)
			}

		}
	}

	return
}

func GetInfoFromIDNo(idno string) (idInfo chIdNo.IDCardDetail) {
	id := chIdNo.IDCard(idno)
	idInfo, err := id.Decode()
	if err != nil {
		log.Println("非法身份证号")
		return
	}

	return
}

func GetChData(rawStr string) (newStr string) {
	strByte := []byte(rawStr)
	newStr = rawStr
	// 匹配字符串
	comReg := regexp.MustCompile(`\{([A-Z][A-Za-z]+)\}`)
	comMatch := comReg.FindAllSubmatch(strByte, -1)
	if len(comMatch) > 0 {
		for i := range comMatch {
			var value string
			dataType := string(comMatch[i][1])
			rawStrDef := string(comMatch[i][0])
			switch dataType {
			case "Name":
				value = chinaid.Name()
			case "Email":
				value = chinaid.Email()
			case "Mobile", "Phone", "Contact":
				value = chinaid.Mobile()
			case "Gender", "Sex":
				count := RandInt(0, 4)
				value = [...]string{"男", "女", "未知"}[count]
			case "CardNo", "BankNo":
				value = chinaid.BankNo()
			case "Address":
				value = chinaid.Address()
			case "SSN", "IDNo":
				value = chinaid.IDNo()
			case "Company":
				count := RandInt(0, 10)
				prefix := GetRandomRune(4)
				suffix := [...]string{"集团", "控股有限公司", "科技有限公司", "出口公司", "进口公司", "证券公司", "保险公司", "投资公司", "证券公司", "银行"}[count]
				value = fmt.Sprintf("%s%s", prefix, suffix)
			case "Country":
				count := RandInt(0, 10)
				value = [...]string{"中国", "美国", "英国", "俄罗斯", "日本", "韩国", "马来西亚", "新加坡", "菲律宾", "澳大利亚"}[count]
			case "Province", "City", "State":
				value = chinaid.ProvinceAndCity()

			}

			if len(value) > 0 {
				newStr = strings.Replace(newStr, rawStrDef, value, -1)
			}

		}
	}

	return
}

func GetRangeData(rawStr string) (newStr string) {
	strByte := []byte(rawStr)
	newStr = rawStr
	// 匹配字符串
	comReg := regexp.MustCompile(`\{([A-Z][A-Za-z]+)\((-*\d+)\,(-*\d+)\)\}`)
	comMatch := comReg.FindAllSubmatch(strByte, -1)
	if len(comMatch) > 0 {
		for i := range comMatch {
			var ret string
			dataType := string(comMatch[i][1])
			start, _ := strconv.Atoi(string(comMatch[i][2]))
			end, _ := strconv.Atoi(string(comMatch[i][3]))
			rawStrDef := string(comMatch[i][0])
			switch dataType {
			case "Int":
				ret = strconv.Itoa(GetRandomInt(start, end))
			case "Time":
				ret = GetRandomTime(start, end)
			case "Minute":
				ret = GetRandomMin(start, end)
			case "Hour":
				ret = GetRandomHour(start, end)
			case "Second":
				ret = GetRandomSecond(start, end)
			case "Date":
				ret = GetRandomDate(start, end)
			case "Timestamp":
				ret = GetRandomTimestamp(start, end)
			}
			newStr = strings.Replace(newStr, rawStrDef, ret, -1)
		}
	}
	return
}

func GetSlicesIndex(src string) (isHit bool, index int, keyName string) {
	strByte := []byte(src)
	indexReg := regexp.MustCompile(`(\w+)\[(\W*\d+)\]`)
	indexMatch := indexReg.FindAllSubmatch(strByte, -1)
	if len(indexMatch) > 0 {
		index, _ = strconv.Atoi(string(indexMatch[0][2]))
		keyName = string(indexMatch[0][1])
		isHit = true
	}
	return
}

func Is2Split(src string) (b bool) {
	starIndex := strings.Index(src, "*")
	barIndex := strings.Index(src, "-")
	if (starIndex < barIndex) && (starIndex != -1) {
		return
	}
	if barIndex >= 0 {
		b = true
	}
	return
}

func GetRandomInt(min, max int) int {
	// Intn returns, as an int, a non-negative pseudo-random number in [0,n)
	return min + rand.Intn(max-min)
}

func GetRandomDate(min, max int) (value string) {
	interval := min + rand.Intn(max-min)
	curTimestamp := time.Now().Unix() + int64(86400*interval)
	value = fmt.Sprintf(time.Unix(curTimestamp, 0).Format("2006-01-02"))
	return
}

func GetRandomTime(min, max int) (value string) {
	interval := min + rand.Intn(max-min)
	curTimestamp := time.Now().Unix() + int64(86400*interval)
	value = fmt.Sprintf(time.Unix(curTimestamp, 0).Format("2006-01-02 15:04:05"))
	return
}

func GetRandomMin(min, max int) (value string) {
	interval := min + rand.Intn(max-min)
	curTimestamp := time.Now().Unix() + int64(60*interval)
	value = fmt.Sprintf(time.Unix(curTimestamp, 0).Format("2006-01-02 15:04:05"))
	return
}

func GetRandomHour(min, max int) (value string) {
	interval := min + rand.Intn(max-min)
	curTimestamp := time.Now().Unix() + int64(3600*interval)
	value = fmt.Sprintf(time.Unix(curTimestamp, 0).Format("2006-01-02 15:04:05"))
	return
}

func GetRandomSecond(min, max int) (value string) {
	interval := min + rand.Intn(max-min)
	curTimestamp := time.Now().Unix() + int64(1*interval)
	value = fmt.Sprintf(time.Unix(curTimestamp, 0).Format("2006-01-02 15:04:05"))
	return
}

func GetRandomTimestamp(min, max int) (value string) {
	interval := min + rand.Intn(max-min)
	curTimestamp := time.Now().Unix() + int64(86400*interval)
	targetTimeStr := fmt.Sprintf(time.Unix(curTimestamp, 0).Format("2006-01-02 15:04:05"))
	loc, _ := time.LoadLocation("Asia/Shanghai")
	tt, _ := time.ParseInLocation("2006-01-02 15:04:05", targetTimeStr, loc) //2006-01-02 15:04:05是转换的格式如php的"Y-m-d H:i:s"
	value = fmt.Sprintf("%v", tt.UnixNano()/1e6)
	return
}

func GetRandomRune(runeLen int) string {
	a := make([]rune, runeLen)
	for i := range a {
		a[i] = rune(RandInt(19968, 40869))
	}
	return string(a)
}

func RandInt(min, max int64) int64 {
	rand.Seed(time.Now().UnixNano())
	return min + rand.Int63n(max-min)
}

func GetRandomStr(length int, strType string) (ranStr string) {
	var baseStr string
	if strType == "upper" {
		baseStr = "ABCDEFGHIGKLMNOPQRSTUVWXYZ"
	} else if strType == "lower" {
		baseStr = "abcdefghigklmnopqrstuvwxyz"
	} else if strType == "upperLower" {
		baseStr = "ABCDEFGHIGKLMNOPQRSTUVWXYZabcdefghigklmnopqrstuvwxyz"
	} else if strType == "int" {
		baseStr = "0123456789"
	} else {
		baseStr = "ABCDEFGHIGKLMNOPQRSTUVWXYZabcdefghigklmnopqrstuvwxyz0123456789"
	}

	bytes := []byte(baseStr)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < length; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	ranStr = string(result)
	return
}

func GetSliceIndex(values []string, value string) (index int) {
	index = -1
	for p, v := range values {
		tmps := strings.Split(v, "/")
		lastV := tmps[len(tmps)-1]
		if lastV == value {
			index = p
			return
		}
	}
	return
}

func ExecCommand(strCommand string) (result string, err error) {
	cmd := exec.Command("/bin/bash", "-c", strCommand)
	out_bytes, err := cmd.CombinedOutput()
	if err != nil {
		Logger.Info("cmd: %s", strCommand)
		Logger.Warning("%s", err)
		return
	}

	//stdout, _ := cmd.StdoutPipe()
	//err = cmd.Start()
	//if err != nil {
	//	Logger.Info("cmd: %s", strCommand)
	//	Logger.Error("%s", err)
	//	return "", err
	//}
	//
	//out_bytes, err := ioutil.ReadAll(stdout)
	//stdout.Close()
	//
	//if err != nil {
	//	Logger.Error("%s", err)
	//	return
	//}
	//if err := cmd.Wait(); err != nil {
	//	Logger.Info("cmd: %s", strCommand)
	//	Logger.Error("%s", err)
	//	return "", err
	//}

	result = strings.Replace(string(out_bytes), "\n", "", -1)

	return
}

func WriteJson(data []byte, path string) (err error) {
	var str bytes.Buffer
	_ = json.Indent(&str, data, "", "    ")

	err = ioutil.WriteFile(path, []byte(str.String()), 0644)
	if err != nil {
		Logger.Error("%s", err)
	}
	return
}

// GetBetweenDates 根据开始日期和结束日期计算出时间段内所有日期
// 参数为日期格式，如：2020-01-01
func GetBetweenDates(sdate, edate string) []string {
	d := []string{}
	timeFormatTpl := "2006-01-02 15:04:05"
	if len(timeFormatTpl) != len(sdate) {
		timeFormatTpl = timeFormatTpl[0:len(sdate)]
	}
	date, err := time.Parse(timeFormatTpl, sdate)
	if err != nil {
		// 时间解析，异常
		return d
	}
	date2, err := time.Parse(timeFormatTpl, edate)
	if err != nil {
		// 时间解析，异常
		return d
	}
	if date2.Before(date) {
		// 如果结束时间小于开始时间，异常
		return d
	}
	// 输出日期格式固定
	timeFormatTpl = "2006-01-02"
	date2Str := date2.Format(timeFormatTpl)
	d = append(d, date.Format(timeFormatTpl))
	for {
		date = date.AddDate(0, 0, 1)
		dateStr := date.Format(timeFormatTpl)
		d = append(d, dateStr)
		if dateStr == date2Str {
			break
		}
	}
	return d
}

func Reverse(arr *[]string) {
	var temp string
	length := len(*arr)
	for i := 0; i < length/2; i++ {
		temp = (*arr)[i]
		(*arr)[i] = (*arr)[length-1-i]
		(*arr)[length-1-i] = temp
	}
}

func RawStr2MadeStr(lang, keyName, value string, order int, depOutVars map[string][]interface{}) (afterStr interface{}, err error) {
	t, subV, allDef := GetStrType(lang, value)
	if len(keyName) == 0 {
		keyName = subV
	}

	//if t == 1 || t == 3 {
	//	Logger.Debug("t: %v, keyName: %s, subV: %v, allDef: %v", t, keyName, subV, allDef)
	//}

	if t == 1 {
		if value, ok := depOutVars[keyName]; ok {
			afterStr = value[0]
		} else {
			err = fmt.Errorf("未找到变量[%s]定义，请先定义或关联", keyName)
			Logger.Error("%s", err)
			Logger.Debug("t: %v, keyName: %s, subV: %v, allDef: %v", t, keyName, subV, allDef)
			return
		}
	} else if t == 2 {
		afterStr = subV
	} else if t == 3 {
		var tmpKey, tmpStr string
		count := 0
		for defKey, defValue := range allDef {
			if value, ok := depOutVars[defKey]; ok {
				if len(value) > order {
					if order < 0 {
						tmpKey, _ = Interface2Str(value[len(value)+order])
					} else {
						tmpKey, _ = Interface2Str(value[order])
					}
				} else {
					tmpKey, _ = Interface2Str(value[0])
				}

			} else {
				err = fmt.Errorf("未找到变量[%s]定义，请先定义或关联", defKey)
				Logger.Error("%s", err)
				Logger.Debug("t: %v, keyName: %s, subV: %v, allDef: %v", t, keyName, subV, allDef)
				return
			}

			if count == 0 {
				if len(subV) == 0 {
					tmpStr = strings.Replace(value, defValue, tmpKey, -1)
				} else {
					tmpStr = strings.Replace(subV, defValue, tmpKey, -1)
				}
			} else {
				tmpStr = strings.Replace(tmpStr, defValue, tmpKey, -1)
			}
			count++
		}
		afterStr = tmpStr
	} else {
		return value, err
	}

	return
}

func GetAbDef() (inDef map[string]string) {
	inDef = make(map[string]string)
	inDef["intAb"] = "65536,-1"
	inDef["intNor"] = "10,1"
	inDef["strAb"] = GetRandomStr(256, "") + "," + GetRandomStr(65, "")
	inDef["strNor"] = GetRandomStr(255, "") + "," + GetRandomStr(8, "")
	inDef["arrAb"] = "65536,-1"
	inDef["arrNor"] = "1,2"
	inDef["objAb"] = "{}"
	inDef["bool"] = "true,false"
	return
}

func CopyMap(src map[string]interface{}) (dst map[string]interface{}) {
	dst = make(map[string]interface{}, len(src))
	for k, v := range src {
		dst[k] = v
	}
	return dst
}

func GetTreeDataTag(rawStr string) (treeDataKey string, deep int) {
	strByte := []byte(rawStr)
	// 匹配字符串
	comReg := regexp.MustCompile(`\{TreeData_([A-Z][A-Za-z]+)\[(\d+)\]\}`)
	comMatch := comReg.FindAllSubmatch(strByte, -1)
	if len(comMatch) > 0 {
		treeDataKey = string(comMatch[0][1])
		deepStr := string(comMatch[0][2])
		var err error
		deep, err = strconv.Atoi(deepStr)
		if err != nil {
			Logger.Error("%s", err)
		}
	}
	return
}

func GetTreeDataValue(keyName string, deep int, first, second string) (after1, after2, after3 string) {
	var chinaData ChinaData
	var sysPara SysParameter
	models.Orm.Table("sys_parameter").Where("name = ?", keyName).Find(&sysPara)
	if len(sysPara.ValueList) == 0 {
		Logger.Error("未找到参数: %s的定义，请核对", keyName)
		return
	}
	err := json.Unmarshal([]byte((sysPara.ValueList)), &chinaData)
	if err != nil {
		Logger.Debug("%s:%v", keyName, sysPara.ValueList)
		Logger.Error("层级参数定义有误，请核对: %v", err)
		return
	}
	if deep == 1 {
		after1 = chinaData.AddName
	} else if deep == 2 {
		after1 = first
		randNum := int(RandInt(0, int64(len(chinaData.Children))))
		for k, data := range chinaData.Children {
			if k == randNum {
				after2 = data.AddName
				break
			}
		}
	} else if deep == 3 {
		after1 = first
		after2 = second
		for _, data := range chinaData.Children {
			if data.AddName == after2 {
				thirdLen := len(data.Children)
				if thirdLen == 0 {
					after3 = after2
					break
				} else {
					randNum := int(RandInt(0, int64(thirdLen)))
					for subK, thirdData := range data.Children {
						if subK == randNum {
							after3 = thirdData.AddName
							break
						}
					}
				}

			}
		}
	}

	return
}

func RawStrComparion(sType, resp string, target interface{}) (b bool, err error) {
	targetStr, _ := Interface2Str(target)
	switch sType {
	case "in", "contain":
		b = strings.Contains(resp, targetStr)
	case "!in", "not_in", "not_contain":
		b = !strings.Contains(resp, targetStr)
	case "re", "regex", "regexp":
		re := regexp.MustCompile(targetStr)
		result := re.FindStringSubmatch(resp)
		if len(result) > 0 {
			b = true
		} else {
			b = false
		}
	case "=", "equal":
		if resp == target {
			b = true
		} else {
			b = false
		}
	case "!=", "not_equal":
		if resp != targetStr {
			b = true
		} else {
			b = false
		}
	case "null", "empty":
		if len(resp) == 0 {
			b = true
		} else {
			b = false
		}
	case "!null", "!empty", "not_null", "not_empty":
		if len(resp) > 0 {
			b = true
		} else {
			b = false
		}
	default:
		err = fmt.Errorf("不支持%s类型的比较，如有需要请反馈致相关人员", sType)
		Logger.Error("%s", err)
	}

	if !b {
		err = fmt.Errorf("断言: Response %s [%s] 结果:fail", sType, targetStr)
		Logger.Error("%s", err)
	}

	return
}

func GetRequestLangage(info map[string]interface{}) (lang string) {
	contentTypeRaw, _ := Interface2Str(info["Cookie"])
	if len(contentTypeRaw) == 0 {
		contentTypeRaw, _ = Interface2Str(info["cookie"])
	}

	if strings.Contains(contentTypeRaw, "lang=en") {
		lang = "en"
	}
	return
}

func GetOneValueFromStringList(in string) (out string) {
	var values []string
	tmps := strings.Split(in, ",")
	for _, item := range tmps {
		info := strings.TrimSpace(item)
		if len(info) > 0 {
			values = append(values, info)
		}
	}
	index := GetRandomInt(0, len(values))
	out = values[index]
	return
}
