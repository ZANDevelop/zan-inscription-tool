package app

import (
	"log"
	"os"
)

var logger = log.New(os.Stdout, "", log.LstdFlags)

func LogInfo(args ...interface{}) {
	logger.Println(args)
}

func LogInfof(template string, args ...interface{}) {
	logger.Printf(template, args...)
}

func LogErrorf(template string, args ...interface{}) {
	logger.Printf(template, args...)
}
