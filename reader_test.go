package sixty_test

import (
	"bytes"
	"testing"

	"github.com/rprtr258/sixty"
	"github.com/stretchr/testify/assert"
)

func TestNewReaderBufferSize_Empty(t *testing.T) {
	reader := sixty.NewReader(bytes.NewBufferString(""))
	assert.Equal(t, byte(0), reader.ReadByte())
	assert.Equal(t, rune(0), reader.ReadRune())
	assert.NoError(t, reader.Error)
}
