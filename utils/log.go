package utils

import (
	"fmt"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	logrus "github.com/sirupsen/logrus"
	"os"
	"path"
	"path/filepath"
	"time"
)

var logger = logrus.New()

type CustomFormatter struct{}

func (s *CustomFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	//strings.ToUpper(entry.Level.String())
	msg := fmt.Sprintf("%s\n", entry.Message)
	return []byte(msg), nil
}

func init() {

	if logConfig := GetConfig().GetSystem().GetLog(); logConfig != nil {
		var fileName string
		logFilePath := logConfig.Filepath
		logFileName := logConfig.Filename

		//相对路径
		if path.IsAbs(logFilePath) == false {
			filePath, _ := filepath.Abs("./")
			logFilePath = path.Join(filePath, logFilePath)
			fileName = path.Join(logFilePath, logFileName)
		} else {
			fileName = path.Join(logFilePath, logFileName)
		}

		fmt.Println("logFilePath:", fileName)

		os.Remove(fileName)

		file, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.ModeAppend)
		if err != nil {
			fmt.Println("err", err)
			return
		}

		logrus.SetOutput(file)
		logrus.SetLevel(logrus.DebugLevel)

		// 设置 rotatelogs
		logWriter, err := rotatelogs.New(
			// 分割后的文件名称
			fileName+".%Y%m%d-%H%M%S.log",
			// 生成软链，指向最新日志文件
			rotatelogs.WithLinkName(fileName),
			// 设置最大保存时间(7天)
			rotatelogs.WithMaxAge(7*24*time.Hour),
			// 设置日志切割时间间隔(1天)
			rotatelogs.WithRotationTime(24*time.Hour),
		)

		// 设置钩子
		writeMap := lfshook.WriterMap{
			logrus.InfoLevel:  logWriter,
			logrus.FatalLevel: logWriter,
			logrus.DebugLevel: logWriter,
			logrus.WarnLevel:  logWriter,
			logrus.ErrorLevel: logWriter,
			logrus.PanicLevel: logWriter,
		}

		lfHook := lfshook.NewHook(writeMap, new(CustomFormatter))

		//lfHook := lfshook.NewHook(writeMap, &logrus.TextFormatter{
		//	TimestampFormat: GetConfig().GetSystem().Timeformat,
		//	DisableTimestamp:true,
		//	DisableLevelTruncation:true,
		//})

		// 新增钩子
		logger.AddHook(lfHook)
	}
}

func PrintInfo(args ...interface{}) {
	if GetConfig().GetSystem().GetLog().Enabel == true {
		logger.Info(args...)
	}
}

// 日志记录到 MongoDB
func LoggerToMongo() {

}

// 日志记录到 ES
func LoggerToES() {
}

// 日志记录到 MQ
func LoggerToMQ() {

}
