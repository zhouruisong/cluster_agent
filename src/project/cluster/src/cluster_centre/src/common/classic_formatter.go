package common

import (
	"bytes"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"sort"
)

type ClassicFormatter struct {
	TimestampFormat string
	FieldsDelimiter string
}

func (f *ClassicFormatter) Format(entry *log.Entry) ([]byte, error) {
	b := &bytes.Buffer{}
	// write [%time] [%level] %message
	if f.TimestampFormat == "" {
		f.TimestampFormat = log.DefaultTimestampFormat
	}
	if f.FieldsDelimiter == "" {
		f.FieldsDelimiter = " "
	}
	fmt.Fprintf(b, "[%s] [%s] %s", entry.Time.Format(f.TimestampFormat),
		entry.Level.String(), entry.Message)
	// sort fields
	keys := make([]string, 0, len(entry.Data))
	for key := range entry.Data {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	// append fields
	for idx := range keys {
		fmt.Fprint(b, f.FieldsDelimiter)
		appendKeyValue(b, keys[idx], entry.Data[keys[idx]])
	}
	b.WriteByte('\n')
	return b.Bytes(), nil
}

func needsQuoting(text string) bool {
	for _, ch := range text {
		if !((ch >= 'a' && ch <= 'z') ||
			(ch >= 'A' && ch <= 'Z') ||
			(ch >= '0' && ch <= '9') ||
			ch == '-' || ch == '.') {
			return false
		}
	}
	return true
}

func appendKeyValue(b *bytes.Buffer, key, value interface{}) {
	switch value.(type) {
	case string:
		if needsQuoting(value.(string)) {
			fmt.Fprintf(b, "%v=%s", key, value)
		} else {
			fmt.Fprintf(b, "%v=%q", key, value)
		}
	case error:
		if needsQuoting(value.(error).Error()) {
			fmt.Fprintf(b, "%v=%s", key, value)
		} else {
			fmt.Fprintf(b, "%v=%q", key, value)
		}
	default:
		fmt.Fprintf(b, "%v=%v", key, value)
	}
}
