package utils

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
)

var log = logrus.New()

// 自定义 Formatter，只打印消息
type PlainFormatter struct{}

func (f *PlainFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	return []byte(fmt.Sprintf("[%s] %s\n", entry.Level, entry.Message)), nil // 只返回消息 + 换行
}

func init() {
	// 创建日志目录
	if err := os.MkdirAll("log", 0755); err != nil {
		panic(err)
	}

	// 创建日志文件
	logFile, err := os.OpenFile("log/sync.log", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		panic(err)
	}

	// 配置日志格式
	// log.SetFormatter(&logrus.TextFormatter{
	// 	DisableTimestamp: true,
	// })
	log.SetFormatter(&PlainFormatter{})

	// 设置日志级别
	log.SetLevel(logrus.InfoLevel)

	// 设置输出
	log.SetOutput(os.Stdout)
	log.AddHook(&MultiWriterHook{
		Writer: logFile,
	})
}

// MultiWriterHook 是一个自定义的 hook，用于同时写入多个输出
type MultiWriterHook struct {
	Writer *os.File
}

func (hook *MultiWriterHook) Fire(entry *logrus.Entry) error {
	line, err := entry.String()
	if err != nil {
		return err
	}
	_, err = hook.Writer.Write([]byte(line))
	return err
}

func (hook *MultiWriterHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

// GetLogger 返回配置好的日志实例
func GetLogger() *logrus.Logger {
	return log
}
