package u

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"os"
	"os/user"
	"runtime"
	"strings"
	"time"
)

var (
	errInvalidBase64 = errors.New("Invalid base64 value")
)

// FmtArgs formats args as a string. First argument should be format string
// and the rest are arguments to the format
func FmtArgs(args ...interface{}) string {
	if len(args) == 0 {
		return ""
	}
	format := args[0].(string)
	if len(args) == 1 {
		return format
	}
	return fmt.Sprintf(format, args[1:]...)
}

func panicWithMsg(defaultMsg string, args ...interface{}) {
	s := FmtArgs(args...)
	if s == "" {
		s = defaultMsg
	}
	fmt.Printf("%s\n", s)
	panic(s)
}

// PanicIf panics if cond is true
func PanicIf(cond bool, args ...interface{}) {
	if !cond {
		return
	}
	panicWithMsg("PanicIf: condition failed", args...)
}

// PanicIfErr panics if err is not nil
func PanicIfErr(err error, args ...interface{}) {
	if err == nil {
		return
	}
	panicWithMsg(err.Error(), args...)
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

// UtcNow returns current time in UTC
func UtcNow() time.Time {
	return time.Now().UTC()
}

const base64Chars = "0123456789abcdefghijklmnopqrstuvwxyz"

// EncodeBase64 encodes n as base64
func EncodeBase64(n int) string {
	var buf [16]byte
	size := 0
	for {
		buf[size] = base64Chars[n%36]
		size++
		if n < 36 {
			break
		}
		n /= 36
	}
	end := size - 1
	for i := 0; i < end; i++ {
		b := buf[i]
		buf[i] = buf[end]
		buf[end] = b
		end--
	}
	return string(buf[:size])
}

// DecodeBase64 decodes base64 string
func DecodeBase64(s string) (int, error) {
	n := 0
	for _, c := range s {
		n *= 36
		i := strings.IndexRune(base64Chars, c)
		if i == -1 {
			return 0, errInvalidBase64
		}
		n += i
	}
	return n, nil
}
