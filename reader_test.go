package sixty_test

import (
	"bytes"
	"errors"
	"testing"

	"github.com/rprtr258/sixty"
	"github.com/stretchr/testify/require"
)

func Test_ReadByte(t *testing.T) {
	t.Parallel()

	reader := sixty.NewReader(bytes.NewBufferString(""))
	require.Equal(t, byte(0), reader.ReadByte())
	require.NoError(t, reader.Error)
}

func TestNewReaderBufferSize_Empty(t *testing.T) {
	t.Parallel()

	reader := sixty.NewReader(bytes.NewBufferString(""))
	require.Equal(t, byte(0), reader.ReadByte())
	require.Equal(t, rune(0), reader.ReadRune())
	require.NoError(t, reader.Error)
}

func Test_ReadString(t *testing.T) {
	t.Parallel()

	reader := sixty.NewReader(bytes.NewBuffer([]byte{0x03, 0x41, 0x42, 0x43}))
	require.Equal(t, "ABC", reader.ReadLengthString())
	require.NoError(t, reader.Error)
}

func TestReadRuneMultiByte(t *testing.T) {
	t.Parallel()

	// 2-, 3- and 4-byte UTF-8 runes
	for _, r := range []rune{'ж', '☺', '🦄'} {
		reader := sixty.NewReader(bytes.NewBufferString(string(r)))
		require.Equal(t, r, reader.ReadRune(), "rune %q", r)
		require.NoError(t, reader.Error)
	}
}

func TestReadRuneTooLong(t *testing.T) {
	t.Parallel()

	// 0xF8 claims 5 continuation bytes, but UTF-8 runes are at most 4 bytes
	reader := sixty.NewReader(bytes.NewBuffer([]byte{0xF8}))
	require.Equal(t, rune(0), reader.ReadRune())
	var tooLong sixty.TooLongRuneError
	require.ErrorAs(t, reader.Error, &tooLong)
	require.Equal(t, uint(1), tooLong.Pos)
	require.Equal(t, "TOO LONG RUNE at 1", tooLong.Error())
}

func TestReadRuneInvalidContinuation(t *testing.T) {
	t.Parallel()

	// lead byte 0xC2 promises a continuation byte but 0x00 has no leading ones:
	// decoding must yield the raw bits without error
	reader := sixty.NewReader(bytes.NewBuffer([]byte{0xC2, 0x00}))
	require.Equal(t, rune(0xC2&0x3F)<<6, reader.ReadRune())
	require.NoError(t, reader.Error)
}

func TestReadUint64MultiByte(t *testing.T) {
	t.Parallel()

	reader := sixty.NewReader(bytes.NewBuffer([]byte{0xAC, 0x02}))
	require.Equal(t, uint64(300), reader.ReadUint64())
	require.NoError(t, reader.Error)
}

func TestPos(t *testing.T) {
	t.Parallel()

	reader := sixty.NewReader(bytes.NewBufferString("abc"))
	require.Equal(t, uint(0), reader.Pos())
	reader.ReadByte()
	reader.ReadByte()
	require.Equal(t, uint(2), reader.Pos())
}

// errDiskOnFire is a non-EOF read error, distinct from io.EOF.
var errDiskOnFire = errors.New("disk on fire")

// failingReader returns an error other than io.EOF.
type failingReader struct{ err error }

func (r failingReader) Read([]byte) (int, error) { return 0, r.err }

func TestUnderlyingErrorPropagates(t *testing.T) {
	t.Parallel()

	reader := sixty.NewReader(failingReader{errDiskOnFire})
	require.Equal(t, byte(0), reader.ReadByte())
	require.ErrorIs(t, reader.Error, errDiskOnFire)
}
