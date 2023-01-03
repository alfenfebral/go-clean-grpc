package logger

import "github.com/sirupsen/logrus"

func Error(err error) {
	logrus.Error(err)
}

func Println(args ...interface{}) {
	logrus.Println(args...)
}

func Printf(format string, args ...interface{}) {
	logrus.Printf(format, args...)
}

func Info(args ...interface{}) {
	logrus.Info(args...)
}
