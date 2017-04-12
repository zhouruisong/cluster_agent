package logger

import (
	"../common"
	log "github.com/Sirupsen/logrus"
	"path"
	"time"
)

func GetLogger(dir string, name string) *log.Logger {
	log_file := path.Join(dir, name+".log")
	output := &common.FileRotator{
		FileName:    log_file,
		MaxSize:     100 << 20,
		MaxDuration: 1 * time.Hour,
	}
	formatter := &common.ClassicFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
		FieldsDelimiter: ", ",
	}
	logger := &log.Logger{
		Out:       output,
		Formatter: formatter,
		Hooks:     nil,
		// TODO: log level configurable
		Level: log.DebugLevel,
	}
	return logger
}
