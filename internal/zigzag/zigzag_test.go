package zigzag_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/rprtr258/sixty/internal/zigzag"
)

var testCases = []struct {
	orig   int64
	zigZag uint64
}{
	{orig: 0, zigZag: 0},
	{orig: -1, zigZag: 1},
	{orig: 1, zigZag: 2},
	{orig: -2, zigZag: 3},
	{orig: 2, zigZag: 4},
	{orig: 0x7FFFFFFFFFFFFFF8, zigZag: 0xFFFFFFFFFFFFFFF0},
	{orig: 0x7FFFFFFFFFFFFFFF, zigZag: 0xFFFFFFFFFFFFFFFE},
	{orig: -0x7FFFFFFFFFFFFFF1, zigZag: 0xFFFFFFFFFFFFFFE1},
	{orig: -0x7FFFFFFFFFFFFFFF, zigZag: 0xFFFFFFFFFFFFFFFD},
	// math.MinInt64 must survive the round trip despite overflowing int64 in
	// intermediate arithmetic
	{orig: -0x8000000000000000, zigZag: 0xFFFFFFFFFFFFFFFF},
}

func TestEncode(t *testing.T) {
	t.Parallel()

	for _, test := range testCases {
		t.Run(fmt.Sprintf("to %d(%[1]x) %x", test.orig, test.zigZag), func(t *testing.T) {
			t.Parallel()
			require.Equal(t, test.zigZag, zigzag.Encode(test.orig))
		})
	}
}

func TestDecode(t *testing.T) {
	t.Parallel()

	for _, test := range testCases {
		t.Run(fmt.Sprintf("from %d(%[1]x) %x", test.orig, test.zigZag), func(t *testing.T) {
			t.Parallel()
			require.Equal(t, test.orig, zigzag.Decode(test.zigZag))
		})
	}
}

func TestRoundTrip(t *testing.T) {
	t.Parallel()

	for x := -1000; x < 1000; x++ {
		require.Equal(t, int64(x), zigzag.Decode(zigzag.Encode(int64(x))))
	}
}
