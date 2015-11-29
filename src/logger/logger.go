package logger

import (
	"os"
)

type Logger struct {
	log_file *os.File
}

func GetLogger() *Logger {
	l := new(Logger)
	// 実行時のカレントディレクトリに作成される。
	l.log_file, _ = os.Create("log")
	return l
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
