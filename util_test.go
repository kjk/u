package u

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func testShortenID(t *testing.T, n int) {
	s := EncodeBase64(n)
	n2, err := DecodeBase64(s)
	assert.Nil(t, err)
	assert.Equal(t, n, n2)
}

func TestShortenId(t *testing.T) {
	testShortenID(t, 1404040)
	testShortenID(t, 0)
	testShortenID(t, 1)
	testShortenID(t, 35)
	testShortenID(t, 36)
	testShortenID(t, 37)
	testShortenID(t, 123413343)
	_, err := DecodeBase64("azasdf!")
	assert.Error(t, err)
}
