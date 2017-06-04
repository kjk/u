package u

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func testEncodeBase64(t *testing.T, n int) {
	s := EncodeBase64(n)
	n2, err := DecodeBase64(s)
	assert.Nil(t, err)
	assert.Equal(t, n, n2)
}

func TestEncodeBase64(t *testing.T) {
	testEncodeBase64(t, 1404040)
	testEncodeBase64(t, 0)
	testEncodeBase64(t, 1)
	testEncodeBase64(t, 35)
	testEncodeBase64(t, 36)
	testEncodeBase64(t, 37)
	testEncodeBase64(t, 123413343)
	_, err := DecodeBase64("azasdf!")
	assert.Error(t, err)
}
