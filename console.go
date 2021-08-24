// Copyright 2021 Axlrose. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file

package logkit

import (
	"errors"
	"fmt"
	"io"
	"os"
	"runtime"
)

type colorConsole func(string) string

func newColorConsole(format string) colorConsole {
	return func(text string) string {
		return "\033[" + format + "m" + text + "\033[0m"
	}
}

var colorMap = map[Severity]colorConsole{
	SeverityEmergency: newColorConsole("1;95"), // bold light magenta
	SeverityAlert:     newColorConsole("1;91"), // bold light red
	SeverityCritical:  newColorConsole("35"),   // magenta
	SeverityError:     newColorConsole("31"),   // red
	SeverityWarning:   newColorConsole("33"),   // yellow
	SeverityNotice:    newColorConsole("36"),   // cyan
	SeverityInfo:      newColorConsole("32"),   // green
	SeverityDebug:     newColorConsole("39"),   // default
}

type ConsoleTarget struct {
	ColorMode bool
	Writer    io.Writer
	close     chan bool
}

func NewConsoleTarget() *ConsoleTarget {
	return &ConsoleTarget{
		ColorMode: true,
		Writer:    os.Stdout,
		close:     make(chan bool),
	}
}

func (t *ConsoleTarget) Open(w io.Writer) error {
	// filter
	if t.Writer == nil {
		return errors.New("ConsoleTarget.Writer cannot be nil. ")
	}
	if runtime.GOOS == "windows" {
		t.ColorMode = false
	}
	return nil
}

func (t *ConsoleTarget) Process(entry *Entry) {
	if entry == nil {
		t.close <- true
		return
	}

	msg := entry.String()

	if t.ColorMode {
		if color, ok := colorMap[entry.Severity]; ok {
			msg = color(msg)
		}
	}

	_, _ = fmt.Fprintln(t.Writer, msg)
}

// Close close the current console target
func (t *ConsoleTarget) Close() {
	<-t.close
}
