package u

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStringArrayRemoveFirst(t *testing.T) {
	tests := [][]string{
		nil, nil,
		[]string{"a"}, []string{},
		[]string{"a", "b"}, []string{"b"},
	}
	n := len(tests) / 2
	for i := 0; i < n; i++ {
		got := StringArrayRemoveFirst(tests[i*2])
		exp := tests[i*2+1]
		assert.Equal(t, exp, got)
	}
}
