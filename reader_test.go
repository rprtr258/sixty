package sixty_test

import (
	"bytes"
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
