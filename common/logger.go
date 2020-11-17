package common

import (
	"fmt"
	"os"
)

func Error(context string, err error) {
	if err != nil {
		Log("Error [%s]: %v", context, err)
	}
}

func Log(format string, details ...interface{}) {
	if !Silent {
		fmt.Fprintln(os.Stderr, fmt.Sprintf(format, details...))
		os.Stderr.Sync()
	}
}

func Debug(format string, details ...interface{}) error {
	if DebugFlag {
		fmt.Fprintln(os.Stderr, fmt.Sprintf(format, details...))
		os.Stderr.Sync()
	}
	return nil
}

func Trace(format string, details ...interface{}) error {
	if TraceFlag {
		fmt.Fprintln(os.Stderr, fmt.Sprintf(format, details...))
		os.Stderr.Sync()
	}
	return nil
}

func Out(format string, details ...interface{}) {
	fmt.Fprintf(os.Stdout, format, details...)
	os.Stdout.Sync()
}
