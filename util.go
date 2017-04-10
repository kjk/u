package u

import (
	"crypto/sha1"
	"fmt"
	"os"
	"os/user"
	"runtime"
	"strings"
	"time"
)

// PanicIf panics if cond is true
func PanicIf(cond bool, args ...interface{}) {
	if !cond {
		return
	}
	msg := "invalid state"
	if len(args) > 0 {
		s, ok := args[0].(string)
		if ok {
			msg = s
			if len(s) > 1 {
				msg = fmt.Sprintf(msg, args[1:]...)
			}
		}
	}
	panic(msg)
}

// PanicIfErr panics if err is not nil
func PanicIfErr(err error) {
	if err != nil {
		panic(err.Error())
	}
}

// IsLinux returns true if running on linux
func IsLinux() bool {
	return runtime.GOOS == "linux"
}

// IsMac returns true if running on mac
func IsMac() bool {
	return runtime.GOOS == "darwin"
}

// UserHomeDir returns $HOME diretory of the user
func UserHomeDir() string {
	// user.Current() returns nil if cross-compiled e.g. on mac for linux
	if usr, _ := user.Current(); usr != nil {
		return usr.HomeDir
	}
	return os.Getenv("HOME")
}

// ExpandTildeInPath converts ~ to $HOME
func ExpandTildeInPath(s string) string {
	if strings.HasPrefix(s, "~") {
		return UserHomeDir() + s[1:]
	}
	return s
}

// Sha1HexOfBytes returns 40-byte hex sha1 of bytes
func Sha1HexOfBytes(data []byte) string {
	return fmt.Sprintf("%x", Sha1OfBytes(data))
}

// Sha1OfBytes returns 20-byte sha1 of bytes
func Sha1OfBytes(data []byte) []byte {
	h := sha1.New()
	h.Write(data)
	return h.Sum(nil)
}

// DurationToString converts duration to a string
func DurationToString(d time.Duration) string {
	minutes := int(d.Minutes()) % 60
	hours := int(d.Hours())
	days := hours / 24
	hours = hours % 24
	if days > 0 {
		return fmt.Sprintf("%dd %dhr", days, hours)
	}
	if hours > 0 {
		return fmt.Sprintf("%dhr %dm", hours, minutes)
	}
	return fmt.Sprintf("%dm", minutes)
}

// TimeSinceNowAsString returns string version of time since a ginve timestamp
func TimeSinceNowAsString(t time.Time) string {
	return DurationToString(time.Now().Sub(t))
}
