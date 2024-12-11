package library

import (
	"context"
	"reflect"

	log "github.com/sirupsen/logrus"
)

const LogLevelDefault = log.ErrorLevel

type LogInterface interface {
	GetLogger() *log.Entry
}

type LogStruct struct {
	LogEntry *log.Entry `json:"-"` // log entry
	LogInterface
}

// get default log entry, never nil
func (source *LogStruct) GetLogger() *log.Entry {
	return GetLogger(source)
}

func GetLogger(source *LogStruct) *log.Entry {
	if source != nil && source.LogEntry != nil {
		return source.LogEntry
	}

	logentry := log.WithContext(context.Background())
	logentry.Level = LogLevelDefault
	logentry.Infof("generating new log entry for %s, with level: %s", reflect.TypeOf(source), logentry.Level)

	if source != nil {
		source.LogEntry = logentry
	}

	return logentry
}
