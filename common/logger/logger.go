package logger

import (
	"fmt"
	"log"
	"os"
	"time"
)

type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
	FATAL
)

var (
	logger *log.Logger
	level  LogLevel
)

func init() {
	// 初始化日志记录器
	file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("Failed to open log file:", err)
	}
	
	logger = log.New(file, "", log.LstdFlags)
	level = INFO
}

// SetLevel 设置日志级别
func SetLevel(l LogLevel) {
	level = l
}

// Debug 记录调试日志
func Debug(v ...interface{}) {
	if level <= DEBUG {
		logger.SetPrefix("[DEBUG] ")
		logger.Println(v...)
	}
}

// Info 记录信息日志
func Info(v ...interface{}) {
	if level <= INFO {
		logger.SetPrefix("[INFO] ")
		logger.Println(v...)
	}
}

// Warn 记录警告日志
func Warn(v ...interface{}) {
	if level <= WARN {
		logger.SetPrefix("[WARN] ")
		logger.Println(v...)
	}
}

// Error 记录错误日志
func Error(v ...interface{}) {
	if level <= ERROR {
		logger.SetPrefix("[ERROR] ")
		logger.Println(v...)
	}
}

// Fatal 记录致命错误日志
func Fatal(v ...interface{}) {
	if level <= FATAL {
		logger.SetPrefix("[FATAL] ")
		logger.Fatalln(v...)
	}
}

// Format 格式化日志
func Format(format string, v ...interface{}) string {
	return fmt.Sprintf(format, v...)
}

// WithTime 添加时间戳
func WithTime(msg string) string {
	return fmt.Sprintf("[%s] %s", time.Now().Format("2006-01-02 15:04:05"), msg)
}