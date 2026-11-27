package sixty_test

import (
	"fmt"
	"testing"

	"github.com/rprtr258/sixty"
	"github.com/stretchr/testify/require"
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
	// TODO: fix overflow
	// {orig: -0x7FFFFFFFFFFFFFF1, zigZag: 0xFFFFFFFFFFFFFFF1},
	// {orig: -0x7FFFFFFFFFFFFFFF, zigZag: 0xFFFFFFFFFFFFFFFF},
}

func TestFromZigZag(t *testing.T) {
	t.Parallel()

	for _, test := range testCases {
		t.Run(fmt.Sprintf("from %d(%[1]x) %x", test.orig, test.zigZag), func(t *testing.T) {
			t.Parallel()
			require.Equal(t, test.orig, sixty.FromZigZag(test.zigZag))
		})
	}
}

func TestToZigZag(t *testing.T) {
	t.Parallel()

	for _, test := range testCases {
		t.Run(fmt.Sprintf("to %d(%[1]x) %x", test.orig, test.zigZag), func(t *testing.T) {
			t.Parallel()
			require.Equal(t, test.zigZag, sixty.ToZigZag(test.orig))
		})
	}
}

func TestAssertSchemas(t *testing.T) {
	t.Parallel()

	tt, vv, ok := sixty.AssertSchemas[sixty.Int, sixty.IntValue](sixty.Int{}, sixty.IntValue(3))
	require.True(t, ok)
	require.Equal(t, sixty.Int{}, tt)
	require.Equal(t, sixty.IntValue(3), vv)

	_, _, ok = sixty.AssertSchemas[sixty.Int, sixty.IntValue](sixty.Bit{}, sixty.IntValue(3))
	require.False(t, ok)

	_, _, ok = sixty.AssertSchemas[sixty.Int, sixty.IntValue](sixty.Int{}, sixty.UintValue(3))
	require.False(t, ok)

	_, _, ok = sixty.AssertSchemas[sixty.Int, sixty.IntValue](sixty.Bit{}, sixty.UintValue(3))
	require.False(t, ok)
}
