// Copyright 2021 Axlrose. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file

// Package logkit implements logging with severity levels and message categories.
package logkit

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"
)

// RFC5424 log message levels.
const (
	SeverityEmergency Severity = iota
	SeverityAlert
	SeverityCritical
	SeverityError
	SeverityWarning
	SeverityNotice
	SeverityInfo
	SeverityDebug
)

// Severity describes the level of a log message.
type Severity int

// SeverityLabel map severity to name
var SeverityName = map[Severity]string{
	SeverityEmergency: "Emergency",
	SeverityAlert:     "Alert",
	SeverityCritical:  "Critical",
	SeverityError:     "Error",
	SeverityWarning:   "Warning",
	SeverityNotice:    "Notice",
	SeverityInfo:      "Info",
	SeverityDebug:     "Debug",
}

// String() return the label of Severity
func (s Severity) String() string {
	if label, ok := SeverityName[s]; ok {
		return label
	}
	return "Unknown"
}

// Entry represents a log entry
type Entry struct {
	Severity         Severity
	Category         string
	Message          string
	Time             time.Time
	CallStack        string
	FormattedMessage string
}

// String return the string representation of a entry
func (e *Entry) String() string {
	return e.FormattedMessage
}

type Target interface {
	Open(errWriter io.Writer) error
	Process(*Entry)
	Close()
}

type core struct {
	lock    sync.Mutex
	open    bool        // whether the log is open
	entries chan *Entry // log entries

	errorWriter     io.Writer // the writer used to write errors during it's startup
	bufferSize      int       // the size of channel store log entries
	callStackDepth  int       // the number of call stack information to be logged for each message
	callStackFilter string    // the filter of call stack information
	maxLevel        Severity  // the maximum level of message to be logged
	targets         []Target  // targets for sending log messages to
}

// Formatter format the a log message into a string
type Formatter func(*LogKit, *Entry) string

// LogKit records log messages and dispatches the to targets for further processing
type LogKit struct {
	*core
	Category  string
	Formatter Formatter
}

// Option is the type all options need to adhere to
type Option func(c *core)

// WithErrorWriter set the error writer default os.Stderr
func WithErrorWriter(w io.Writer) Option {
	return func(c *core) {
		c.errorWriter = w
	}
}

// WithBuffer set the size of channel which store log entries
func WithBuffer(s int) Option {
	return func(c *core) {
		c.bufferSize = s
	}
}

// WithCallStackDepth set the depth of call stack information
func WithCallStackDepth(d int) Option {
	return func(c *core) {
		c.callStackDepth = d
	}
}

// WithStackFilter set the filter of call stack information
func WithStackFilter(f string) Option {
	return func(c *core) {
		c.callStackFilter = f
	}
}

// WithMaxLevel set the maximum level of message to be logged
func WithMaxLevel(l Severity) Option {
	return func(c *core) {
		c.maxLevel = l
	}
}

// WithTargets set targets for sending log messages to
func WithTargets(t ...Target) Option {
	return func(c *core) {
		c.targets = t
	}
}

// NewLogKit constructs a instance of LogKit, with default option
func NewLogKit() *LogKit {
	core := &core{
		errorWriter: os.Stderr,
		bufferSize:  1024,
		maxLevel:    SeverityDebug,
		targets:     make([]Target, 0),
	}
	return &LogKit{
		core:      core,
		Category:  "app",
		Formatter: DefaultFormatter,
	}
}

// NewLogKitByOptions constructs a new instance of LogKit, with any option you specify
func NewLogKitByOptions(opt ...Option) *LogKit {
	core := &core{
		errorWriter:    os.Stderr,
		bufferSize:     1024,
		callStackDepth: 3,
		// callStackFilter: "",
		maxLevel: SeverityDebug,
		targets:  make([]Target, 0),
	}
	for _, o := range opt {
		o(core)
	}
	return &LogKit{
		core:      core,
		Category:  "app",
		Formatter: DefaultFormatter,
	}
}

// CopyLogKit get new LogKit instance with option
// Option category default the LogKit's old category, if not set
// Option formatter default the DefaultFormatter, if not set
func (l *LogKit) Copy(category string, formatter Formatter) *LogKit {
	lk := &LogKit{
		core:      l.core,
		Category:  category,
		Formatter: formatter,
	}
	return lk
}

// Open init core instance. and listen entries channel,
// if has value, send to target for processing
// Open must be called before any message can be logged.
func (c *core) Open() error {
	c.lock.Lock()
	defer c.lock.Unlock()

	if c.open {
		return nil
	}

	if c.errorWriter == nil {
		return errors.New("LogKit.ErrorWriter must be set. ")
	}
	if c.bufferSize < 0 {
		return errors.New("LogKit.BufferSize must be no less than 0. ")
	}
	if c.callStackDepth < 0 {
		return errors.New("Logger.CallStackDepth must be no less than 0. ")
	}
	c.entries = make(chan *Entry, c.bufferSize)

	// set targets
	var targets []Target
	for _, target := range c.targets {
		if err := target.Open(c.errorWriter); err != nil {
			_, _ = fmt.Fprintf(c.errorWriter, "Failed to open target:%v ", err)
		} else {
			targets = append(targets, target)
		}
	}
	c.targets = targets

	go c.process()

	c.open = true

	return nil
}

// process get message from channel entries
// and send the message to every target for processing
func (c *core) process() {
	for {
		entry, ok := <-c.entries
		for _, target := range c.targets {
			target.Process(entry)
		}
		if !ok {
			return
		}
	}
}

// Close close the LogKit
// Set LogKit status to false. close entries channel
// And call the function Close of every target defined
func (c *core) Close() {
	c.lock.Lock()
	defer c.lock.Unlock()

	if !c.open {
		return
	}
	c.open = false

	close(c.entries)

	for _, target := range c.targets {
		target.Close()
	}
}

// DefaultFormatter is the default formatter used to format log message.
func DefaultFormatter(l *LogKit, e *Entry) string {
	return fmt.Sprintf("%v [%v][%v] %v%v", e.Time.Format(time.RFC3339), e.Severity, e.Category, e.Message, e.CallStack)
}

// Debug logs a message for debugging purpose.
// Please refer to Error() for how to use this method.
func (l *LogKit) Debug(format string, a ...interface{}) {
	l.Log(SeverityDebug, format, a...)
}

// Info logs a message for informational purpose.
// Please refer to Error() for how to use this method.
func (l *LogKit) Info(format string, a ...interface{}) {
	l.Log(SeverityInfo, format, a...)
}

// Notice logs a message meaning normal but significant condition.
// Please refer to Error() for how to use this method.
func (l *LogKit) Notice(format string, a ...interface{}) {
	l.Log(SeverityNotice, format, a...)
}

// Warning logs a message indicating a warning condition.
// Please refer to Error() for how to use this method.
func (l *LogKit) Warning(format string, a ...interface{}) {
	l.Log(SeverityWarning, format, a...)
}

// Error logs a message indicating an error condition.
// This method takes one or multiple parameters. If a single parameter
// is provided, it will be treated as the log message. If multiple parameters
// are provided, they will be passed to fmt.Sprintf() to generate the log message.
func (l *LogKit) Error(format string, a ...interface{}) {
	l.Log(SeverityError, format, a...)
}

// Critical logs a message indicating critical conditions.
// Please refer to Error() for how to use this method.
func (l *LogKit) Critical(format string, a ...interface{}) {
	l.Log(SeverityCritical, format, a...)
}

// Alert logs a message indicating action must be taken immediately.
// Please refer to Error() for how to use this method.
func (l *LogKit) Alert(format string, a ...interface{}) {
	l.Log(SeverityAlert, format, a...)
}

// Emergency logs a message indicating the system is unusable.
// Please refer to Error() for how to use this method.
func (l *LogKit) Emergency(format string, a ...interface{}) {
	l.Log(SeverityEmergency, format, a...)
}

// Log log send message to entries channel
func (l *LogKit) Log(level Severity, format string, a ...interface{}) {
	if level > l.maxLevel || !l.open {
		return
	}
	message := format
	if len(a) > 0 {
		message = fmt.Sprintf(format, a...)
	}
	entry := &Entry{
		Severity: level,
		Category: l.Category,
		Message:  message,
		Time:     time.Now(),
	}
	// call stack information
	if l.callStackDepth > 0 {
		entry.CallStack = GetCallStack(3, l.callStackDepth, l.callStackFilter)
	}
	// format information
	entry.FormattedMessage = l.Formatter(l, entry)
	// send entry to entries channel
	l.entries <- entry
}

// GetCallStack returns the current call stack information as a string
// The skip parameter specifies how many top frames should be skipped.
// The maxStackDepth specifics at most how many frames should be return
func GetCallStack(skip int, maxStackDepth int, filter string) string {
	buf := new(bytes.Buffer)

	for i, depth := skip, 0; depth < maxStackDepth; i++ {
		// get file and line number information about function
		_, file, line, ok := runtime.Caller(i)
		// break if recover information failed
		if !ok {
			break
		}

		// if not set filter or file path include filter
		// write file information to buf
		if filter == "" || strings.Contains(file, filter) {
			_, _ = fmt.Fprintf(buf, "\n%s:%d", file, line)
			depth++
		}
	}
	return buf.String()
}
