package biz

import (
	"bufio"
	//"fmt"
	//"github.com/derekgr/hivething"
	"github.com/tealeg/xlsx"
	"os"
)

func WriteDataInCSV(filePath, content string) (err error) {
	fileHandle, _ := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	defer fileHandle.Close()

	write := bufio.NewWriter(fileHandle)
	_, _ = write.WriteString(content)
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
