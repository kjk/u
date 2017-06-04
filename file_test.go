package u

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFile(t *testing.T) {
	assert.True(t, FileExists("file.go"))
	assert.False(t, FileExists("file_that_doesnt_exist.go"))
	assert.True(t, DirExists("."))
	assert.False(t, DirExists("dir_that_doesnt_exist"))

	lines, err := ReadLinesFromFile("for_tests.txt")
	assert.Nil(t, err)
	assert.Equal(t, 3, len(lines))
	assert.Equal(t, "line 1", lines[0])
}
