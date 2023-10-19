package startup

import (
	"log"
	"os"
)

type Logger struct {
	InfoLog  *log.Logger
	ErrorLog *log.Logger
}

func NewLogger(loggerName string) *log.Logger {
	//infoLog := log.New(os.Stdout, "INFO\t"+loggerName+"\t", log.Ldate|log.Ltime)
	//errorLog := log.New(os.Stderr, "ERROR\t"+loggerName+"\t", log.Ldate|log.Ltime|log.Lshortfile)
	logger := log.New(os.Stderr, loggerName+"\t", log.Ldate|log.Ltime|log.Lshortfile)
	return logger
}
