package log

import (
	"io/ioutil"
	"log"
	"os"
	"sync"
)

// log levels
const (
	InfoLevel = iota
	ErrorLevel
	Disabled
)

// log variables
var (
	infoLog  = log.New(os.Stdout, "\033[34m[info ]\033[0m ", log.LstdFlags|log.Lshortfile)
	errorLog = log.New(os.Stdout, "\033[31m[error]\033[0m ", log.LstdFlags|log.Lshortfile)
	loggers  = []*log.Logger{infoLog, errorLog}
	mutex    sync.Mutex
)

// log methods
var (
	Infoln  = infoLog.Println
	Infof   = infoLog.Printf
	Errorln = errorLog.Println
	Errorf  = errorLog.Printf
)

// SetLevel is to set log level
func SetLevel(level int) {
	mutex.Lock()
	defer mutex.Unlock()

	for _, logger := range loggers {
		logger.SetOutput(os.Stdout)
	}

	if ErrorLevel < level {
		errorLog.SetOutput(ioutil.Discard)
	}

	if InfoLevel < level {
		infoLog.SetOutput(ioutil.Discard)
	}
}
