package domain

import (
	"animar/v1/configs"
	"animar/v1/internal/pkg/tools/tools"
	"encoding/json"
	"fmt"
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

type Log struct {
	Kind string    `json:"kind"`
	Time time.Time `json:"time"`
}

type LogE struct {
	Log
	Level   string `json:"level"`
	Content string `json:"content"`
	Place   string `json:"place"`
}

type LogA struct {
	Log
	Address string `json:"address"`
	Method  string `json:"method"`
	Path    string `json:"path"`
}

func NewAccessLog(addr, method, path string) *LogA {
	alog := &LogA{
		Log: Log{
			Kind: "access",
			Time: time.Now(),
		},
		Address: addr,
		Method:  method,
		Path:    path,
	}
	return alog
}

func NewErrorLog(content, level string) *LogE {
	_, file, line, _ := runtime.Caller(1)
	var lev string
	if level == "" {
		lev = LogWarn
	} else {
		lev = level
	}

	eLog := &LogE{
		Log: Log{
			Kind: "error",
			Time: time.Now(),
		},
		Level:   lev,
		Content: content,
		Place:   fmt.Sprintf("%s:%d", file, line),
	}
	return eLog
}

func (e *LogE) Logging() {
	if tools.IsProductionEnv() {
		e.write()
	} else {
		e.write()
	}
}

func (a *LogA) Logging() {
	if tools.IsProductionEnv() {
		a.write()
	} else {
		a.write()
	}
}

func (a *LogA) write() {
	today := time.Now().Format("20060102")
	// logFile, _ := os.OpenFile(fmt.Sprintf("%s/%s_%s.log", configs.ErrorLogDirectory, l.Kind, today), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	logFile, _ := os.OpenFile(fmt.Sprintf("%s/log_%s.log", configs.ErrorLogDirectory, today), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	defer logFile.Close()

	writeJsonFile(fmt.Sprintf("%s/log_%s.json", configs.ErrorLogDirectory, today), a)
}

func (e *LogE) write() {
	today := time.Now().Format("20060102")
	// logFile, _ := os.OpenFile(fmt.Sprintf("%s/%s_%s.log", configs.ErrorLogDirectory, l.Kind, today), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	logFile, _ := os.OpenFile(fmt.Sprintf("%s/log_%s.log", configs.ErrorLogDirectory, today), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	defer logFile.Close()

	writeJsonFile(fmt.Sprintf("%s/log_%s.json", configs.ErrorLogDirectory, today), e)
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
