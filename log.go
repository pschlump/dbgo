package dbgo

import (
	"fmt"
	"io"
	"os"

	"github.com/pschlump/filelib"
)

type LogSeverity int

type LogData interface {
	// Setup
	OutputDestination(FilePtr io.Writer, isTerminal bool) // Output destinations are collected
	StatusCode(code int)                                  // 404 etc.
	ErrorMessageFormat()                                  // Set an error message format, text/plain, html, application/json etc.

	// Message
	OutputDestinationWeb( /* www, req */ )      // per call output destination.
	Severity(x LogSeverity)                     //
	Depth(d int)                                // Default to -1, for normal LF but can be set.
	MessageQuery(s string, data ...interface{}) // SQL query if applicable.
	MessageData(s string, data ...interface{})  // Bind variables
	Message(s string, data ...interface{})      // fmt.Printf format, followed by items, %! - set to secret for next item.
	UserMessage(s string, data ...interface{})  // fmt.Printf format, followed by items
	// MessageQueryOutput(s string, data ...interface{}) // Bind variables
	//ColorSet(...)

	// Finalize
	Send() // Completed Log Message
}

type LogDataIntermediate struct {
	LineNo int
}

type LogToFile struct {
	LineNo int
	Files  []*os.File
}

func NewLogDataFile(fn string) (rv LogData) {
	fp, err := filelib.Fopen(fn, "a")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to open %s for append: %s\n", fn, err)
		os.Exit(1)
	}
	return &LogToFile{
		LineNo: 0,
		Files:  []*os.File{fp},
	}
}

func (lf *LogToFile) OutputDestination(FilePtr io.Writer, isTerminal bool) { // Output destinations are collected
	// lf.Files = append(lf.Files, FilePtr)
}
func (lf *LogToFile) StatusCode(code int) { // 404 etc.
}
func (lf *LogToFile) ErrorMessageFormat() { // Set an error message format, text/plain, html, application/json etc.
}

// Message
func (lf *LogToFile) OutputDestinationWeb( /* www, req */ ) { // per call output destination.
}
func (lf *LogToFile) Severity(x LogSeverity) { //
}
func (lf *LogToFile) Depth(d int) { // Default to -1, for normal LF but can be set.
}
func (lf *LogToFile) MessageQuery(s string, data ...interface{}) { // SQL query if applicable.
}
func (lf *LogToFile) MessageData(s string, data ...interface{}) { // Bind variables
}
func (lf *LogToFile) Message(s string, data ...interface{}) { // fmt.Printf format, followed by items, %! - set to secret for next item.
}
func (lf *LogToFile) UserMessage(s string, data ...interface{}) { // fmt.Printf format, followed by items
}

// Finalize
func (lf *LogToFile) Send() { // Completed Log Message
}
