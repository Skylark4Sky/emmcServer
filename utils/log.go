package utils

import (
	"fmt"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	logrus "github.com/sirupsen/logrus"
	"os"
	"path"
	"time"
)

var logger = logrus.New()

func init() {

	if logConfig := GetConfig().GetSystem().GetLog(); logConfig != nil {
		logFilePath := logConfig.Filepath
		logFileName := logConfig.Filename
		fileName := path.Join(logFilePath, logFileName)
		file, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
		if err != nil {
			fmt.Println("err", err)
			return
		}

		logrus.SetOutput(file)
		logrus.SetLevel(logrus.DebugLevel)

		// 设置 rotatelogs
		logWriter, err := rotatelogs.New(
			// 分割后的文件名称
			fileName + ".%Y%m%d.log",
			// 生成软链，指向最新日志文件
			rotatelogs.WithLinkName(fileName),
			// 设置最大保存时间(7天)
			rotatelogs.WithMaxAge(7*24*time.Hour),
			// 设置日志切割时间间隔(1天)
			rotatelogs.WithRotationTime(10*time.Minute),
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

		lfHook := lfshook.NewHook(writeMap, &logrus.JSONFormatter{
			TimestampFormat:GetConfig().GetSystem().Timeformat,
		})

		// 新增钩子
		logger.AddHook(lfHook)
	}
}

func PrintInfo(args ...interface{}) {
	logger.Info(args...)
	//logger.WithFields(logrus.Fields{
	//	"status_code"  : statusCode,
	//	"latency_time" : latencyTime,
	//	"client_ip"    : clientIP,
	//	"req_method"   : reqMethod,
	//	"req_uri"      : reqUri,
	//}).Info()
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
