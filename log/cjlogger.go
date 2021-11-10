package log

import (
	"log"
	"os"
)

var (
	cjLog *defaultLogger
)

func init() {
	cjLog = newCjLogger()
}

type defaultLogger struct {
	log *log.Logger
}

func newCjLogger() *defaultLogger {
	return &defaultLogger{
		log: log.New(os.Stderr, "", log.LstdFlags),
	}
}

func (d *defaultLogger) Info(v ...interface{}) {
	d.log.Println(append([]interface{}{"[INFO]"}, v...)...)
}

func (d *defaultLogger) Infof(format string, v ...interface{}) {
	d.log.Printf("[INFO] "+format, v...)
}

func (d *defaultLogger) Error(v ...interface{}) {
	d.log.Println(append([]interface{}{"[ERRO]"}, v...)...)
}

func (d *defaultLogger) Errorf(format string, v ...interface{}) {
	d.log.Printf("[ERRO] "+format, v...)
}

// ============================= 常用方法封装 ===============================

func Info(v ...interface{}) {
	cjLog.Info(v...)
}

func Infof(format string, v ...interface{}) {
	cjLog.Infof(format, v...)
}

func Error(v ...interface{}) {
	cjLog.Error(v...)
}

func Errorf(format string, v ...interface{}) {
	cjLog.Errorf(format, v...)
}
