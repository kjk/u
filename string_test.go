package u

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStringsRemoveFirst(t *testing.T) {
	tests := [][]string{
		nil, nil,
		[]string{"a"}, []string{},
		[]string{"a", "b"}, []string{"b"},
	}
	n := len(tests) / 2
	for i := 0; i < n; i++ {
		got := StringsRemoveFirst(tests[i*2])
		exp := tests[i*2+1]
		assert.Equal(t, exp, got)
	}
}

func TestRemoveDuplicateStrings(t *testing.T) {
	// note: the fact that arrays are sorted after RemoveDuplicateString
	// is accidental. We could make the tests more robust by writing
	// doStringArraysHaveTheSameContent(a1, a2 []string)
	tests := [][]string{
		nil, nil,
		[]string{"a"}, []string{"a"},
		[]string{"b", "a"}, []string{"a", "b"},
		[]string{"a", "a"}, []string{"a"},
		[]string{"a", "b", "a"}, []string{"a", "b"},
		[]string{"ab", "ba", "ab", "ab", "cd"}, []string{"ab", "ba", "cd"},
	}
	n := len(tests) / 2
	for i := 0; i < n; i++ {
		got := RemoveDuplicateStrings(tests[i*2])
		exp := tests[i*2+1]
		assert.Equal(t, exp, got)
	}
}
