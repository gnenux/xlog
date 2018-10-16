// Package log implements a simple logging package. It defines a type, Logger,
// with methods for formatting output. Logger accessible through helper functions
// Warn[f|ln], Info[f,ln], Error[f,ln], Fatal[f|ln], and Panic[f|ln], which are
// easier to use than creating a Logger manually.
// That logger writes to standard error and prints the date and time
// of each logged message.
// Every log message is output on a separate line: if the message being
// printed does not end in a newline, the logger will add one.
// The Fatal functions call os.Exit(1) after writing the log message.
// The Panic functions call panic after writing the log message.

package xlog

import (
	"io"
	"log"
	"os"
)

const (
	// DefaultLogFlag = log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile
	DefaultLogFlag = log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile
	warnSign       = "[warn] "
	infoSign       = "[info] "
	errorSign      = "[error] "
	fatalSign      = "[fatal] "
	panicSign      = "[panic] "
)

var (
	zeroInterface interface{}
)

type Logger struct {
	prefix string
	flag   int
	logger *log.Logger
}

func NewLogger(out io.Writer, prefix string, flag int) *Logger {
	logger := log.New(out, prefix, flag)
	l := new(Logger)
	l.prefix = prefix
	l.flag = flag
	l.logger = logger
	return l
}

func NewLoggerFromFileName(filename string, prefix string, flag int) *Logger {
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR|os.O_APPEND, os.ModePerm|os.ModeAppend)
	if err != nil {
		log.Fatal(err)
	}
	logger := log.New(f, prefix, flag)
	l := new(Logger)
	l.prefix = prefix
	l.flag = flag
	l.logger = logger
	return l
}

func (l Logger) Warn(v ...interface{}) {
	v = append(v, zeroInterface)
	copy(v[1:], v[0:])
	v[0] = warnSign
	l.logger.Print(v...)
}

func (l Logger) Warnf(format string, v ...interface{}) {
	v = append(v, zeroInterface)
	copy(v[1:], v[0:])
	v[0] = warnSign
	format = "%s" + format
	l.logger.Printf(format, v...)
}

func (l Logger) Warnln(v ...interface{}) {
	v = append(v, zeroInterface)
	copy(v[1:], v[0:])
	v[0] = warnSign
	l.logger.Println(v...)
}

func (l Logger) Info(v ...interface{}) {
	v = append(v, zeroInterface)
	copy(v[1:], v[0:])
	v[0] = infoSign
	l.logger.Print(v...)
}

func (l Logger) Infof(format string, v ...interface{}) {
	v = append(v, zeroInterface)
	copy(v[1:], v[0:])
	v[0] = infoSign
	format = "%s" + format
	l.logger.Printf(format, v...)
}

func (l Logger) Infoln(v ...interface{}) {
	v = append(v, zeroInterface)
	copy(v[1:], v[0:])
	v[0] = infoSign
	l.logger.Println(v...)
}

func (l Logger) Error(v ...interface{}) {
	v = append(v, zeroInterface)
	copy(v[1:], v[0:])
	v[0] = errorSign
	l.logger.Print(v...)
}

func (l Logger) Errorf(format string, v ...interface{}) {
	v = append(v, zeroInterface)
	copy(v[1:], v[0:])
	v[0] = errorSign
	format = "%s" + format
	l.logger.Printf(format, v...)
}

func (l Logger) Errorln(v ...interface{}) {
	v = append(v, zeroInterface)
	copy(v[1:], v[0:])
	v[0] = errorSign
	l.logger.Println(v...)
}

func (l Logger) Fatal(v ...interface{}) {
	v = append(v, zeroInterface)
	copy(v[1:], v[0:])
	v[0] = fatalSign
	l.logger.Fatal(v...)
}

func (l Logger) Fatalf(format string, v ...interface{}) {
	v = append(v, zeroInterface)
	copy(v[1:], v[0:])
	v[0] = fatalSign
	format = "%s" + format
	l.logger.Fatalf(format, v...)
}

func (l Logger) Fatalln(v ...interface{}) {
	v = append(v, zeroInterface)
	copy(v[1:], v[0:])
	v[0] = fatalSign
	l.logger.Fatalln(v...)
}

func (l Logger) Panic(v ...interface{}) {
	v = append(v, zeroInterface)
	copy(v[1:], v[0:])
	v[0] = panicSign
	l.logger.Panic(v...)
}

func (l Logger) Panicf(format string, v ...interface{}) {
	v = append(v, zeroInterface)
	copy(v[1:], v[0:])
	v[0] = panicSign
	format = "%s" + format
	l.logger.Panicf(format, v...)
}

func (l Logger) Panicln(v ...interface{}) {
	v = append(v, zeroInterface)
	copy(v[1:], v[0:])
	v[0] = panicSign
	l.logger.Panicln(v...)
}
