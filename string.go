package u

import (
	"bytes"
	"fmt"
	"sort"
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

// FmtSmart avoids formatting if only format is given
func FmtSmart(format string, args ...interface{}) string {
	if len(args) == 0 {
		return format
	}
	return fmt.Sprintf(format, args...)
}

// StringInSlice returns true if a string is present in slice
func StringInSlice(a []string, toCheck string) bool {
	for _, s := range a {
		if s == toCheck {
			return true
		}
	}
	return false
}

// StringsRemoveFirst removes first sstring from the slice
func StringsRemoveFirst(a []string) []string {
	n := len(a)
	if n > 0 {
		copy(a[:n-1], a[1:])
		a = a[:n-1]
	}
	return a
}

// StringRemoveFromSlice removes a given string from a slice of strings
// returns a (potentially) new slice and true if was removed
func StringRemoveFromSlice(a []string, toRemove string) ([]string, bool) {
	n := len(a)
	if n == 0 {
		return nil, false
	}
	res := make([]string, 0, n)
	for _, s := range a {
		if s != toRemove {
			res = append(res, s)
		}
	}
	didRemove := len(res) != len(a)
	if !didRemove {
		return a, false
	}
	return res, true
}

// RemoveDuplicateStrings removes duplicate strings from an array of strings.
// It's optimized for the case of no duplicates. It modifes a in place.
func RemoveDuplicateStrings(a []string) []string {
	if len(a) < 2 {
		return a
	}
	// sort and remove dupplicates (which are now grouped)
	sort.Strings(a)
	writeIdx := 1
	for i := 1; i < len(a); i++ {
		if a[i-1] == a[i] {
			continue
		}
		if writeIdx != i {
			a[writeIdx] = a[i]
		}
		writeIdx++
	}
	return a[:writeIdx]
}

// NormalizeNewLines changes CR and CRLF into LF
func NormalizeNewlines(d []byte) []byte {
	// replace CR LF \r\n (windows) with LF \n (unix)
	d = bytes.Replace(d, []byte{13, 10}, []byte{10}, -1)
	// replace CF \r (mac) with LF \n (unix)
	d = bytes.Replace(d, []byte{13}, []byte{10}, -1)
	return d
}

const (
	indentStr = "                                "
)

// IndentStr returns an indentation string which has (2*n) spaces
func IndentStr(n int) string {
	if n == 0 {
		return ""
	}
	n = n * 2
	if len(indentStr) >= n {
		return indentStr[:n]
	}
	s := indentStr
	for len(s) < n {
		s += "  "
	}
	return s
}
