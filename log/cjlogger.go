package log

import (
	"fmt"
	"log"
	"os"
)

const (
	LevelDebug = 1 << iota
	LevelInfo
	LevelWarn
	LevelError
	LevelPanic
	LevelFatal
)

var (
	// level 的前缀字符串
	defaultLevelPrefixes = map[int]string{
		LevelDebug: "DEBU",
		LevelInfo:  "INFO",
		LevelWarn:  "WARN",
		LevelError: "ERRO",
		LevelPanic: "PANI",
		LevelFatal: "FATA",
	}

	// level 颜色, 颜色参数格式: 格式：\033[显示方式;前景色;背景色m
	levelColor = map[int]string{
		-1: "\033[0m", // 重置
		// LevelDebug: "DEBU",
		LevelInfo:  "\033[32m",   // 绿色
		LevelWarn:  "\033[33m",   // 黄色
		LevelError: "\033[31m",   // 红色
		LevelPanic: "\033[35m",   // 紫红色
		LevelFatal: "\033[1;35m", // 加粗紫红色
	}
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
	d.log.Println(append([]interface{}{d.getLevelPrefix(LevelInfo)}, v...)...)
}

func (d *defaultLogger) Infof(format string, v ...interface{}) {
	d.log.Printf(d.getLevelPrefix(LevelInfo)+" "+format, v...)
}

func (d *defaultLogger) Error(v ...interface{}) {
	d.log.Println(append([]interface{}{d.getLevelPrefix(LevelError)}, v...)...)
}

func (d *defaultLogger) Errorf(format string, v ...interface{}) {
	d.log.Printf(d.getLevelPrefix(LevelError)+" "+format, v...)
}

func (d *defaultLogger) Warning(v ...interface{}) {
	d.log.Println(append([]interface{}{d.getLevelPrefix(LevelWarn)}, v...)...)
}

func (d *defaultLogger) Warningf(format string, v ...interface{}) {
	d.log.Printf(d.getLevelPrefix(LevelWarn)+" "+format, v...)
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

// getLevelPrefix
func (x *defaultLogger) getLevelPrefix(level int) string {
	str := levelColor[level] + defaultLevelPrefixes[level] + levelColor[-1] // 同时需要重置下
	return "[" + str + "]"
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

func Warning(v ...interface{}) {
	CjLog.Warning(v...)
}

func Warningf(format string, v ...interface{}) {
	CjLog.Warningf(format, v...)
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
