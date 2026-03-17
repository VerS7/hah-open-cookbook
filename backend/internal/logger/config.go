package logger

import (
	"bytes"
	"fmt"
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

type LogFormatter struct {
	logrus.TextFormatter
}

var levels = map[string]string{
	"trace":   "TRACE",
	"debug":   "DEBUG",
	"info":    "INFO ",
	"warning": "WARN ",
	"error":   "ERROR",
	"fatal":   "FATAL",
	"panic":   "PANIC",
	"unknown": "UNKWN",
}

var Default = &logrus.Logger{
	Out:          os.Stderr,
	Formatter:    &LogFormatter{},
	Level:        logrus.DebugLevel,
	ReportCaller: true,
}

func (f *LogFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	level := levels[entry.Level.String()]
	frame := entry.Caller

	msg := bytes.Buffer{}

	msg.WriteString(fmt.Sprintf("[%s] [%s] [%s:%d] - ",
		timestamp,
		level,
		frame.Function,
		frame.Line,
	))

	if len(entry.Message) > 0 {
		msg.WriteString(entry.Message + " ")
	}

	if len(entry.Data) > 0 {
		for k, v := range entry.Data {
			msg.WriteString(fmt.Sprintf("%s=%v ", k, v))
		}
	}

	msg.WriteString("\n")

	return msg.Bytes(), nil
}
