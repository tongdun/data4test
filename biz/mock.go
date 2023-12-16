package biz

import (
	"fmt"
	"io/ioutil"
	"os"
)

func GetFileContent(lang, name string) (content string) {
	fileName := fmt.Sprintf("%s/%s", DownloadBasePath, name)
	contentTmp, errTmp := ioutil.ReadFile(fileName)
	if os.IsNotExist(errTmp) {
		fileName := fmt.Sprintf("%s/%s", UploadBasePath, name)
		contentTmp, subErr := ioutil.ReadFile(fileName)
		if os.IsNotExist(subErr) {
			if errTmp != nil {
				Logger.Error("%s", errTmp)
			}
			return
		}
		content = string(contentTmp)
	} else {
		content = string(contentTmp)
	}
	content = GetSpecialStr(lang, content)
	return
}
