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
		// 一旦ローカルも同じ

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
	}
}
