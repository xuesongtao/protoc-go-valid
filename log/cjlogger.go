package log

import (
	"fmt"
	"log"
	"os"
)

var (
	CjLog *defaultLogger
)

func init() {
	CjLog = NewCjLogger()
}

type defaultLogger struct {
	log *log.Logger
}

func NewCjLogger() *defaultLogger {
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

func (d *defaultLogger) Fatal(v ...interface{}) {
	d.Error(v...)
	os.Exit(1)
}

func (d *defaultLogger) Fatalf(format string, v ...interface{}) {
	d.Errorf(format, v...)
	os.Exit(1)
}

func (d *defaultLogger) Panic(v ...interface{}) {
	d.Error(v...)
	panic(fmt.Sprint(v...))
}

func (d *defaultLogger) Panicf(format string, v ...interface{}) {
	d.Errorf(format, v...)
	panic(fmt.Sprintf(format, v...))
}

// ============================= 常用方法封装 ===============================

func Info(v ...interface{}) {
	CjLog.Info(v...)
}

func Infof(format string, v ...interface{}) {
	CjLog.Infof(format, v...)
}

func Error(v ...interface{}) {
	CjLog.Error(v...)
}

func Errorf(format string, v ...interface{}) {
	CjLog.Errorf(format, v...)
}

func Fatal(v ...interface{}) {
	CjLog.Fatal(v...)
}

func Fatalf(format string, v ...interface{}) {
	CjLog.Fatalf(format, v...)
}

func Panic(v ...interface{}) {
	CjLog.Panic(v...)
}

func Panicf(format string, v ...interface{}) {
	CjLog.Panicf(format, v...)
}

