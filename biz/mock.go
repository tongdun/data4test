package biz

import (
	"fmt"
	"io/ioutil"
)

func GetFileContent(name string) (content string) {
	fileName := fmt.Sprintf("%s/%s", UploadBasePath, name)
	contentTmp, errTmp := ioutil.ReadFile(fileName)
	content = string(contentTmp)
	if errTmp != nil {
		Logger.Error("%s", errTmp)
	}
	return
}
