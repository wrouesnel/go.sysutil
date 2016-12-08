package logutil

import (
	"io"
	"bytes"
)

type LogWriter struct {
	logFunc func(args ...interface{})

	// Currently buffered line
	buffer []byte
}

// Creates a new io.Writer which writes to the log output. Takes a log function
// to use for writing output.
func NewLogWriter(logFunc func(args ...interface{})) io.Writer {
	this := new(LogWriter)
	this.logFunc = logFunc
	this.buffer = make([]byte, 0)

	return io.Writer(this)
}

func (this *LogWriter) Write(p []byte) (n int, err error) {
	// Handle trivial case of blank write
	if len(p) == 0 {
		return 0, nil
	}

	// Is there a new line in the incoming buffer?
	lines := bytes.Split(p, []byte("\n"))

	// Got 2 elements, means at least 1 newline split
	if len(lines) >= 2 {
		this.logFunc(string(append(this.buffer, lines[0]...)))
		// Initialize new buffer.
		this.buffer = make([]byte, 0)
		lines = lines[1:len(p)]

		// Log all remaining lines except the last one (because it's not
		// terminated yet)
		for len(lines) != 1 {
			this.logFunc(string(lines[0]))
			lines = lines[1:len(p)]
		}
	}

	// Append the last line to the current buffer
	this.buffer = append(this.buffer, lines[0]...)
	return len(p), nil
}