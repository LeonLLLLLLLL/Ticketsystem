package tools

import (
	"fmt"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
)

// PrettyFormatter is a custom formatter for more readable logs
type PrettyFormatter struct {
	log.TextFormatter
}

func (f *PrettyFormatter) Format(entry *log.Entry) ([]byte, error) {
	// Use the TextFormatter for base formatting
	f.TextFormatter.DisableTimestamp = false
	f.TextFormatter.FullTimestamp = true
	f.TextFormatter.ForceColors = true

	// Store the original caller info
	caller := ""
	if entry.HasCaller() {
		// Short file path, just the filename and line number
		funcName := entry.Caller.Function
		fileName := entry.Caller.File
		lineNum := entry.Caller.Line

		// Extract just the filename from the full path
		parts := strings.Split(fileName, "/")
		shortFileName := parts[len(parts)-1]

		caller = fmt.Sprintf(" [%s:%d %s]", shortFileName, lineNum, funcName)
	}

	// Add caller info to the message
	if caller != "" {
		entry.Message = entry.Message + caller
	}

	return f.TextFormatter.Format(entry)
}

// Configure sets up the logger with pretty formatting
func Configure() {
	log.SetReportCaller(true)
	log.SetFormatter(&PrettyFormatter{})
	log.SetOutput(os.Stdout)
}
