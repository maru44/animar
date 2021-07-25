package domain

import (
	"animar/v1/configs"
	"animar/v1/internal/pkg/tools/tools"
	"fmt"
	"log"
	"os"
	"runtime"
	"time"
)

var (
	/******************
	    probrem
	******************/
	// urgency (at once)
	LogAlert = "ALT"
	// urgency (at many times)
	LogCritical = "CRT"
	// need fix without urgency (at once)
	LogError = "ERZR"
	// need fix without urgency (at many times)
	LogWarn = "WAN"
	/******************
	   no probrem
	******************/
	LogNotice = "NOT"
	LogInfo   = "INF"
	/******************
	   no probrem
	******************/
	LogDebug = "DBG"
)

// @TODO:ADD mode to args
func LogWriter(str string) {
	if tools.IsProductionEnv() {
		today := time.Now().Format("20060102")
		// @TODO:EDIT err --> _
		logFile, _ := os.OpenFile(fmt.Sprintf("%s/log_%s.log", configs.LogDirectory, today), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
		// jsonFile, err := os.OpenFile(fmt.Sprint())
		defer logFile.Close()

		// get called place
		_, file, line, _ := runtime.Caller(1)

		log.SetOutput(logFile)
		log.SetFlags(log.LstdFlags)
		log.Println(fmt.Sprintf("%s:%d: %s", file, line, str))
	} else {
		log.Print(str)
	}
}

func ErrorLog(err error, mode string) {
	if tools.IsProductionEnv() {
		today := time.Now().Format("20060102")
		// @TODO:EDIT err --> _
		logFile, _ := os.OpenFile(fmt.Sprintf("%s/log_%s.log", configs.LogDirectory, today), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
		// jsonFile, err := os.OpenFile(fmt.Sprint())
		defer logFile.Close()

		// get called place
		_, file, line, _ := runtime.Caller(1)

		log.SetOutput(logFile)
		log.SetFlags(log.LstdFlags)

		// set prefix
		switch mode {
		case "":
			log.SetPrefix(LogInfo)
		default:
			log.SetPrefix(mode)
		}

		log.Println(fmt.Sprintf("%s:%d: %v", file, line, err))
	} else {
		log.Print(err)
	}
}
