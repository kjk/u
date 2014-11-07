package u

import (
	"crypto/sha1"
	"fmt"
	"io"
	"os"
	"os/user"
	"runtime"
	"strings"
	"time"
)

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

func PanicIfErr(err error) {
	if err != nil {
		panic(err.Error())
	}
}

func IsLinux() bool {
	return runtime.GOOS == "linux"
}

func IsMac() bool {
	return runtime.GOOS == "darwin"
}

func UserHomeDir() string {
	// user.Current() returns nil if cross-compiled e.g. on mac for linux
	if usr, _ := user.Current(); usr != nil {
		return usr.HomeDir
	}
	return os.Getenv("HOME")
}

func ExpandTildeInPath(s string) string {
	if strings.HasPrefix(s, "~") {
		return UserHomeDir() + s[1:]
	}
	return s
}

func Sha1OfFile(path string) ([]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		//fmt.Printf("os.Open(%s) failed with %s\n", path, err.Error())
		return nil, err
	}
	defer f.Close()
	h := sha1.New()
	_, err = io.Copy(h, f)
	if err != nil {
		//fmt.Printf("io.Copy() failed with %s\n", err.Error())
		return nil, err
	}
	return h.Sum(nil), nil
}

func Sha1HexOfFile(path string) (string, error) {
	sha1, err := Sha1OfFile(path)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", sha1), nil
}

func Sha1HexOfBytes(data []byte) string {
	return fmt.Sprintf("%x", Sha1OfBytes(data))
}

func Sha1OfBytes(data []byte) []byte {
	h := sha1.New()
	h.Write(data)
	return h.Sum(nil)
}

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

func TimeSinceNowAsString(t time.Time) string {
	return DurationToString(time.Now().Sub(t))
}
