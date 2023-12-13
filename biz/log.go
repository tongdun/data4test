package biz

import (
	"github.com/gamelife1314/logging"
	//"io"
	//"log"
	"os"
)

var (
	Logger *logging.Logger
)

func init() {

	infoLogFile, err := os.OpenFile(InfoLogPath, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	errorLogFile, err := os.OpenFile(ErrorLogPath, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)

	var logLevel logging.MessageLevel

	if LOG_LEVEL == "debug" {
		logLevel = logging.DEBUG
		Logger = &logging.Logger{
			Level: logLevel,
			StreamHandler: &logging.StreamMessageHandler{
				Level: logLevel,
				Formatter: &logging.MessageFormatter{
					Format:     `{{.Color}}{{.Time}} {{.LevelString | printf "%4s"}} {{.ShortFileName}}:{{.Line}} {{.ColorClear}} {{.Message}}`,
					TimeFormat: "2006-01-02T15:04:05.001+0800",
				},
				Destination: os.Stdout,
			},
			FileHandler: &logging.FileMessageHandler{
				Level: logLevel,
				Formatter: &logging.MessageFormatter{
					Format:     `{{.Color}}{{.Time}} {{.LevelString | printf "%4s"}} {{.ShortFileName}}:{{.Line}} {{.ColorClear}} {{.Message}}`,
					TimeFormat: "2006-01-02T15:04:05.001+0800",
				},
				Destination: errorLogFile,
			},
		}
	} else {
		logLevel = logging.INFO
		Logger = &logging.Logger{
			Level: logLevel,
			StreamHandler: &logging.StreamMessageHandler{
				Level: logLevel,
				Formatter: &logging.MessageFormatter{
					Format:     `{{.Color}}{{.Time}} {{.LevelString | printf "%4s"}} {{.ShortFileName}}:{{.Line}} {{.ColorClear}} {{.Message}}`,
					TimeFormat: "2006-01-02T15:04:05.001+0800",
				},
				Destination: os.Stdout,
			},
			FileHandler: &logging.FileMessageHandler{
				Level: logLevel,
				Formatter: &logging.MessageFormatter{
					Format:     `{{.Color}}{{.Time}} {{.LevelString | printf "%4s"}} {{.ShortFileName}}:{{.Line}} {{.ColorClear}} {{.Message}}`,
					TimeFormat: "2006-01-02T15:04:05.001+0800",
				},
				Destination: infoLogFile,
			},
		}
	}

	return
}
