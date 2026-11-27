package sixty_test

import (
	"bytes"
	"testing"

	"github.com/rprtr258/sixty"
	"github.com/stretchr/testify/assert"
)

func Test_ReadByte(t *testing.T) {
	t.Parallel()

	reader := sixty.NewReader(bytes.NewBufferString(""))
	assert.Equal(t, byte(0), reader.ReadByte())
	assert.NoError(t, reader.Error)
}

func TestNewReaderBufferSize_Empty(t *testing.T) {
	t.Parallel()

	reader := sixty.NewReader(bytes.NewBufferString(""))
	assert.Equal(t, byte(0), reader.ReadByte())
	assert.Equal(t, rune(0), reader.ReadRune())
	assert.NoError(t, reader.Error)
}

func Test_ReadString(t *testing.T) {
	t.Parallel()

	reader := sixty.NewReader(bytes.NewBuffer([]byte{0x03, 0x41, 0x42, 0x43}))
	assert.Equal(t, "ABC", reader.ReadLengthString())
	assert.NoError(t, reader.Error)
}
