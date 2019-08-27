package tasks

import (
	"fmt"
	"runtime"
	"time"
)

// Error defines an error in a process.
type Error struct {
	ProcessName string        // Name of the process where the error was executed
	Time        time.Time     // Time when the error occurred
	Err         error         // Error itself
	Frame       runtime.Frame // Frame where the error occurred
}

func (err *Error) Error() string {
	return fmt.Sprintf("[%v] %s (file %s, function %s, line %d)", err.Time, err.Err, err.Frame.File, err.Frame.Function, err.Frame.Line)
}
