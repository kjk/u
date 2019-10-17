package u

import (
	"bytes"
	"sort"
)

// StringsRemoveFirst removes first sstring from the slice
func StringsRemoveFirst(a []string) []string {
	n := len(a)
	if n > 0 {
		copy(a[:n-1], a[1:])
		a = a[:n-1]
	}
	return a
}

// RemoveDuplicateStrings removes duplicate strings from an array of strings.
// It's optimized for the case of no duplicates. It modifes a in place.
func RemoveDuplicateStrings(a []string) []string {
	if len(a) < 2 {
		return a
	}
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
