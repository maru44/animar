package domain

import (
	"animar/v1/configs"
	"animar/v1/internal/pkg/tools/tools"
	"encoding/json"
	"fmt"
	"net/http"
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
	// need fix without urgency
	LogWarn = "WAN"
	/******************
	   no probrem
	******************/
	LogInfo = "INF"
	/******************
	   no probrem
	******************/
	LogDebug = "DBG"
)

type Log struct {
	Kind  string    `json:"kind"`
	Time  time.Time `json:"time"`
	Level string    `json:"level"`
}

type LogE struct {
	Log
	Content string `json:"content"`
	Place   string `json:"place"`
}

type LogA struct {
	Log
	Address string `json:"address"`
	Method  string `json:"method"`
	Path    string `json:"path"`
}

func NewAccessLog() *LogA {
	// startのlogをここで作る
	alog := &LogA{
		Log: Log{
			Kind: "access",
		},
	}
	return alog
}

func NewErrorLog() *LogE {
	eLog := &LogE{
		Log: Log{
			Kind: "error",
		},
	}
	return eLog
}

// func (e *LogE) Alet(err error) {
// 	e.Logging(err, LogAlert)
// }
// func (e *LogE) Critical(err error) {
// 	e.Logging(err, LogCritical)
// }
// func (e *LogE) WARN(err error) {
// 	e.Logging(err, LogWarn)
// }

func ErrorAlert(err error) {
	e := &LogE{
		Log: Log{
			Kind: "error",
		},
	}
	e.Logging(err, LogAlert)
}

func ErrorCritical(err error) {
	e := &LogE{
		Log: Log{
			Kind: "error",
		},
	}
	e.Logging(err, LogCritical)
}

func ErrorWarn(err error) {
	e := &LogE{
		Log: Log{
			Kind: "error",
		},
	}
	e.Logging(err, LogWarn)
}

func (e *LogE) Logging(err error, level string) {
	if tools.IsProductionEnv() {
		e.write(err, level)
	} else {
		e.write(err, level)
	}
}

func (a *LogA) Logging(r *http.Request) {
	if tools.IsProductionEnv() {
		a.write(r)
	} else {
		a.write(r)
	}
}

func (a *LogA) write(r *http.Request) {
	today := time.Now().Format("20060102")

	a.Level = LogInfo
	a.Time = time.Now()
	a.Address = r.RemoteAddr
	a.Method = r.Method
	a.Path = r.URL.Path

	writeJsonFile(fmt.Sprintf("%s/log_%s.json", configs.ErrorLogDirectory, today), a)
}

func (e *LogE) write(err error, level string) {
	today := time.Now().Format("20060102")

	e.Level = level
	e.Content = err.Error()

	// auto
	_, file, line, _ := runtime.Caller(2)
	e.Place = fmt.Sprintf("%s:%d", file, line)
	e.Time = time.Now()

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
