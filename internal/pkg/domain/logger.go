package domain

import (
	"animar/v1/configs"
	"animar/v1/internal/pkg/tools/tools"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"runtime"
	"time"
)

const (
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

type SError struct {
	Level   string    `json:"level"`
	Content string    `json:"content"`
	Place   string    `json:"place"`
	Time    time.Time `json:"time"`
}

func ErrorLog(err error, mode string) {
	if tools.IsProductionEnv() {
		today := time.Now().Format("20060102")
		logFile, _ := os.OpenFile(fmt.Sprintf("%s/log_%s.log", configs.ErrorLogDirectory, today), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
		defer logFile.Close()

		// get called place
		_, file, line, _ := runtime.Caller(1)

		log.SetOutput(logFile)
		log.SetFlags(log.LstdFlags)

		// set prefix
		var level string
		switch mode {
		case "":
			level = LogInfo
		default:
			level = mode
		}
		log.SetPrefix(level)

		log.Println(fmt.Sprintf("%s:%d: %v", file, line, err))

		/***************************
		         json log
		***************************/
		e := &SError{
			Level:   level,
			Content: err.Error(),
			Place:   fmt.Sprintf("%s:%d", file, line),
			Time:    time.Now(),
		}
		writeJsonFile(fmt.Sprintf("%s/log_%s.json", configs.ErrorLogDirectory, today), e)
	} else {
		log.Print(err)
	}
}

func writeJsonFile(fileName string, object interface{}) {
	file, _ := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0600)
	defer file.Close()
	fi, _ := file.Stat()
	leng := fi.Size()

	json_, _ := json.Marshal(object)

	if leng == 0 {
		file.Write([]byte(fmt.Sprintf(`[%s]`, json_)))
	} else {
		file.WriteAt([]byte(fmt.Sprintf(`,%s]`, json_)), leng-1)
	}
}
