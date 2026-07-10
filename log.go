package dbgo

// Copyright (C) Philip Schlump, 2014-2024.
//
// WARNING: This file is EXPERIMENTAL and INCOMPLETE.
//
// LogData and LogToFile sketch a structured, multi-destination logger with
// severity, status codes, SQL query/data capture and a per-call web
// destination. The interface is still being designed and the LogToFile
// implementation is a set of no-op placeholders: today every method does
// nothing and Send() writes nothing. Do not depend on this API yet; it will
// change.

import (
	"fmt"
	"io"
	"os"

	"github.com/pschlump/filelib"
)

// LogSeverity labels the importance of a log message. No severity levels are
// defined yet; this is part of the experimental logging API.
type LogSeverity int

// LogData is the experimental structured-logging interface. A caller configures
// a message with the Setup and Message methods, then calls Send to emit it.
// All methods are part of the experimental API and subject to change.
type LogData interface {
	// Setup
	OutputDestination(FilePtr io.Writer, isTerminal bool) // Output destinations are collected
	StatusCode(code int)                                  // 404 etc.
	ErrorMessageFormat()                                  // Set an error message format: text/plain, html, application/json, etc.

	// Message
	OutputDestinationWeb( /* www, req */ )      // per-call output destination
	Severity(x LogSeverity)                     // set the message severity
	Depth(d int)                                // Default to -1, for normal LF but can be set
	MessageQuery(s string, data ...interface{}) // SQL query if applicable
	MessageData(s string, data ...interface{})  // bind variables
	Message(s string, data ...interface{})      // fmt.Printf-style format and items; %! sets the next item secret
	UserMessage(s string, data ...interface{})  // fmt.Printf-style format and items

	// Finalize
	Send() // Completed Log Message
}

// LogDataIntermediate is a partial base for LogData implementations. It is part
// of the experimental logging API.
type LogDataIntermediate struct {
	LineNo int
}

// LogToFile is an experimental LogData implementation backed by a file. The
// intended behavior is to append formatted log records to Files, but at present
// every method is a no-op placeholder and Send() writes nothing.
type LogToFile struct {
	LineNo int
	Files  []*os.File
}

// NewLogDataFile opens fn for append and returns a LogToFile backed by it. The
// file is created via filelib.Fopen. If the file cannot be opened the returned
// error is non-nil and the LogData value is nil; the caller decides how to
// handle the failure (this function does not exit the process).
func NewLogDataFile(fn string) (rv LogData, err error) {
	fp, err := filelib.Fopen(fn, "a")
	if err != nil {
		return nil, fmt.Errorf("dbgo.NewLogDataFile: opening %q for append: %w", fn, err)
	}
	return &LogToFile{
		LineNo: 0,
		Files:  []*os.File{fp},
	}, nil
}

// The following methods are all no-op placeholders for the experimental
// LogToFile implementation (see the file-level warning). They exist to satisfy
// the LogData interface and will gain real implementations as the API firms up.

func (lf *LogToFile) OutputDestination(FilePtr io.Writer, isTerminal bool) {}
func (lf *LogToFile) StatusCode(code int)                                  {}
func (lf *LogToFile) ErrorMessageFormat()                                  {}
func (lf *LogToFile) OutputDestinationWeb()                                {}
func (lf *LogToFile) Severity(x LogSeverity)                               {}
func (lf *LogToFile) Depth(d int)                                          {}
func (lf *LogToFile) MessageQuery(s string, data ...interface{})           {}
func (lf *LogToFile) MessageData(s string, data ...interface{})            {}
func (lf *LogToFile) Message(s string, data ...interface{})                {}
func (lf *LogToFile) UserMessage(s string, data ...interface{})            {}
func (lf *LogToFile) Send()                                                {}
