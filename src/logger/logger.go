package logger

import (
	"fmt"
	"os"
	"time"
)

type Logger struct {
	log_file *os.File
}

var instance *Logger

// alias
var p = fmt.Println

func Resp(str string, logger *Logger) {
	p(str)
	logger.Res(str)
}

func InitLogger() {
	instance = new(Logger)
	const format = "20060102150405"
	t := time.Now()
	file_name := t.Format(format) + ".log"
	// instance.log_file ,_ =os.Create("/home/gohan/workspace/shogi01/" + file_name)
	instance.log_file, _ = os.Create(file_name)
}

func GetLogger() *Logger {
	return instance
}

func (l Logger) Req(msg string) {
	l.Trace(">   " + msg)
}

func (l Logger) Res(msg string) {
	l.Trace("  < " + msg)
}

func (l Logger) Trace(msg string) {
	l.log_file.WriteString(msg + "\n")
	l.log_file.Sync()
}

func (l Logger) Close() {
	l.log_file.Close()
}
