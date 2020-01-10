package mlogger

import (
	"log"
	"os"
	"sync"
)

type matchaLogger struct {
	*log.Logger
	filename string
}

var mlogger *matchaLogger
var once sync.Once

//GetInstance returns an instance of the logger
func GetInstance() *matchaLogger {
	once.Do(func() {
		mlogger = createLogger("matchaLogger.log")
	})
	return mlogger
}

//createLogger creates an instance of the logger
func createLogger(fname string) *matchaLogger {
	file, _ := os.OpenFile(fname, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)

	return &matchaLogger{
		Logger:   log.New(file, "Matcha ", log.Lshortfile),
		filename: fname,
	}
}
