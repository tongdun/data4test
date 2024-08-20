package biz

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

var (
	BASEPATH      string
	SERVER_PORT   int
	LOG_LEVEL     string
	CICD_HOST     string
	SWAGGER_PATH  string
	HOST_IP       string
	REDIRECT_PATH string
)

var (
	DataBasePath     string
	UploadBasePath   string
	DownloadBasePath string
	HistoryBasePath  string
	InfoLogPath      string
	ErrorLogPath     string
	DocFilePath      string
	OldFilePath      string
	CommonFilePath   string
	ApiFilePath      string
	CaseFilePath     string
)

type Config struct {
	FileBasePath string `json:"file_base_path"`
	ServerPort   int    `json:"server_port"`
	LogLevel     string `json:"log_level"`
	CicdHost     string `json:"cicd_host"`
	SwaggerPath  string `json:"swagger_path"`
	HostIp       string `json:"host_ip"`
	RedirectPath string `json:"redirect_path"`
}

func init() {
	content, err := ioutil.ReadFile("./config.json")
	if err != nil {
		info := fmt.Sprintf("Init Config Failed: %s", err)
		panic(info)
	}
	var config Config
	err = json.Unmarshal([]byte(content), &config)
	if err != nil {
		info := fmt.Sprintf("Init Config Failed: %s", err)
		panic(info)
	}
	if len(config.FileBasePath) == 0 {
		panic("Not Found file_base_path")
	}

	BASEPATH = config.FileBasePath
	SERVER_PORT = config.ServerPort
	LOG_LEVEL = config.LogLevel
	CICD_HOST = config.CicdHost
	SWAGGER_PATH = config.SwaggerPath
	HOST_IP = config.HostIp
	REDIRECT_PATH = config.RedirectPath

	DataBasePath = fmt.Sprintf("%s/data", BASEPATH)
	UploadBasePath = fmt.Sprintf("%s/upload", BASEPATH)
	DownloadBasePath = fmt.Sprintf("%s/download", BASEPATH)
	HistoryBasePath = fmt.Sprintf("%s/history", BASEPATH)
	InfoLogPath = fmt.Sprintf("%s/log/info.log", BASEPATH)
	ErrorLogPath = fmt.Sprintf("%s/log/error.log", BASEPATH)
	DocFilePath = fmt.Sprintf("%s/doc/file", BASEPATH)
	OldFilePath = fmt.Sprintf("%s/old", BASEPATH)
	CommonFilePath = fmt.Sprintf("%s/common", BASEPATH)
	ApiFilePath = fmt.Sprintf("%s/api", BASEPATH)
	CaseFilePath = fmt.Sprintf("%s/case", BASEPATH)

	_, err = os.Stat(BASEPATH)
	if err != nil {
		panic(err)
	}

	//_, err = os.Stat(BASEPATH)
	//if err != nil {
	//	if os.IsNotExist(err) {
	//		subErr := os.MkdirAll(BASEPATH, os.ModePerm)
	//		if subErr != nil {
	//			LogHandle.Printf("Error: %s", subErr)
	//		}
	//	}
	//}
	//
	//subPaths := [4]string{
	//	BASEPATH + "/" + "api",
	//	BASEPATH + "/" + "file",
	//	BASEPATH + "/" + "test",
	//	BASEPATH + "/" + "log",
	//}
	//
	//for _, item := range subPaths {
	//	_, subErr := os.Stat(item)
	//	if subErr != nil {
	//		if os.IsNotExist(err) {
	//			subSubErr := os.MkdirAll(item, os.ModePerm)
	//			if subSubErr != nil {
	//				LogHandle.Printf("Error: %s", subSubErr)
	//			}
	//		}
	//	}
	//}

}
