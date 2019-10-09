package u

import (
	"fmt"
	"io"
)

var (
	LogFile io.Writer
)

// a centralized place allows us to tweak logging, if need be
func Logf(format string, args ...interface{}) {
	if len(args) == 0 {
		fmt.Print(format)
		if LogFile != nil {
			_, _ = fmt.Fprint(LogFile, format)
		}
		return
	}
	fmt.Printf(format, args...)
	if LogFile != nil {
		_, _ = fmt.Fprintf(LogFile, format, args...)
	}
}
