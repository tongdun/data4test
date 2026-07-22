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
	REDIRECT_PATH map[string]string
)

var (
	DataBasePath      string
	UploadBasePath    string
	DownloadBasePath  string
	HistoryBasePath   string
	InfoLogPath       string
	ErrorLogPath      string
	DocFilePath       string
	OldFilePath       string
	CommonFilePath    string
	ApiFilePath       string
	CaseFilePath      string
	AiDataBasePath    string
	KnowledgeBasePath string
)

type Config struct {
	FileBasePath string `json:"file_base_path"`
	ServerPort   int    `json:"server_port"`
	LogLevel     string `json:"log_level"`
	CicdHost     string `json:"cicd_host"`
	SwaggerPath  string `json:"swagger_path"`
	HostIp       string `json:"host_ip"`
	RedirectPath string `json:"redirect_path"`
	Language     string `json:"language"`
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
	REDIRECT_PATH = make(map[string]string)
	err = json.Unmarshal([]byte(config.RedirectPath), &REDIRECT_PATH)
	if err != nil {
		Logger.Warning("重定向路径定义异常: %s", err)
	}

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
	AiDataBasePath = fmt.Sprintf("%s/ai_data", BASEPATH)
	KnowledgeBasePath = fmt.Sprintf("%s/knowledge", BASEPATH)
	_, err = os.Stat(BASEPATH)
	if err != nil {
		panic(err)
	}

	// 初始化 i18n
	locale := config.Language
	if locale == "cn" {
		locale = "zh-CN"
	} else if locale == "en" {
		locale = "en-US"
	}
	if locale == "" {
		locale = "zh-CN"
	}
	fmt.Printf("当前 go-admin 语言设置: %s, 初始化项目翻译器为: %s\n", config.Language, locale)
	InitI18n(locale)
}
