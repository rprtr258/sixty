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

func TestReadRune(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		input    string
		wantRune rune
		errPos   uint // non-zero: ReadRune must fail with TooLongRuneError at this position
	}{
		// 2-, 3- and 4-byte UTF-8 runes
		{"ж", 'ж', 0},
		{"☺", '☺', 0},
		{"🦄", '🦄', 0},
		// 0xF8 claims 5 continuation bytes, but UTF-8 runes are at most 4 bytes
		{"\xF8", 0, 1},
	}
	for _, test := range testCases {
		reader := sixty.NewReader(bytes.NewBufferString(test.input))
		require.Equal(t, test.wantRune, reader.ReadRune(), "input %q", test.input)
		if test.errPos == 0 {
			require.NoError(t, reader.Error, "input %q", test.input)
			continue
		}
		var tooLong sixty.TooLongRuneError
		require.ErrorAs(t, reader.Error, &tooLong, "input %q", test.input)
		require.Equal(t, test.errPos, tooLong.Pos, "input %q", test.input)
		require.Equal(t, "TOO LONG RUNE at 1", tooLong.Error(), "input %q", test.input)
	}
}
