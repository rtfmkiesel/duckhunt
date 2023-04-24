package logger

import (
	"duckhunt/pkg/file"
	"fmt"
	"os"
	"time"
)

var (
	// Log file name
	logFile = "duckhunt.log"
)

const (
	// Time format for the console and the log file
	TimeFormat = "2006-01-02 15:04:05"
)

// LogWrite() will write text to the log file
func LogWrite(text string) {
	// Assemble line
	fmt.Printf("%s %s\n", getTime(), text)
	// Append to file
	err := file.Append(logFile, fmt.Sprintf(
		"%s %s\n", getTime(), text),
	)
	if err != nil {
		CatchErr(err)
	}
}

// Write() will write to the console in the specified format
func Write(text string) {
	fmt.Printf("%s %s\n", getTime(), text)
}

// LogInit() will initalise the log file
func LogInit(filename string) error {
	// Set config file
	logFile = filename

	// Check if file exists
	if !file.Exists(logFile) {
		// File does not exists so create
		err := file.Create(logFile)
		if err != nil {
			return err
		}

		// Test write
		LogWrite("Log file created")
	}

	return nil
}

// CatchErr() will catch errors and display them
// into the terminal and append them to the log file
func CatchErr(err error) {
	var msg = fmt.Sprintf("%s ERROR: %s\n", getTime(), err)
	fmt.Print(msg)
	_ = file.Append(logFile, msg)
}

// CatchCritErr() handles critical errors
// by displaying them into the terminal and
// appending them to the log
// Quits afterswards with exit code 1
func CatchCritErr(err error) {
	var msg = fmt.Sprintf("%s CRITICAL: %s\n", getTime(), err)
	fmt.Print(msg)
	_ = file.Append(logFile, msg)
	os.Exit(1)
}

// getTime() returns the current time as a string
// based on the hardcoded format
func getTime() string {
	return time.Now().Format(TimeFormat)
}
