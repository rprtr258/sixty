package size_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/rprtr258/sixty"
	"github.com/rprtr258/sixty/internal"
	"github.com/rprtr258/sixty/schema/size"
)

func TestBttvDB(t *testing.T) {
	t.Parallel()

	got := size.SchemaSize(internal.BttvEmotesDBSchema)
	require.Equal(t, size.Inf, got)
}

func TestSchemaSizePrimitives(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		schema sixty.Schema
		size   size.Size
	}{
		// strings, ints and runes are variable length
		{sixty.String{}, size.Inf},
		{sixty.Int{}, size.Inf},
		{sixty.Uint{}, size.Inf},
		{sixty.Rune{}, size.Inf},
		{sixty.List{Schema: sixty.Bit{}}, size.Inf},
		{sixty.Bit{}, 1},
		// size uses bits.Len64: log2-ish, off-by-one vs true log2
		{sixty.Enum{Values: []string{"a", "b", "c", "d"}}, 3}, // bits.Len64(4)
		{sixty.Enum{Values: []string{"a", "b", "c"}}, 2},
		{sixty.Enum{Values: []string{"a"}}, 1},
		// enum int is inclusive [min..max], span = max-min = 7, bits.Len64(7) = 3
		{sixty.EnumInt{Min: 0, Max: 7}, 3},
		{sixty.EnumInt{Min: -4, Max: 3}, 3},
		{sixty.Optional{Schema: sixty.Bit{}}, 2}, // presence bit + value bit
		{sixty.Array{Schema: sixty.Bit{}, Count: 5}, 5},
		{sixty.Array{Schema: sixty.String{}, Count: 3}, size.Inf},
		{sixty.Product{
			{"a", sixty.Bit{}},
			{"b", sixty.Bit{}},
			{"c", sixty.Enum{Values: []string{"x", "y"}}},
		}, 4},
		// coproduct: max branch size + bits.Len64(len(tt)) tag bits
		{sixty.Coproduct{
			{"a", sixty.Bit{}},
			{"b", sixty.Bit{}},
		}, 3},
		{sixty.Coproduct{
			{"a", sixty.Bit{}},
			{"b", sixty.Enum{Values: []string{"x", "y", "z", "w"}}},
		}, 5}, // max(1, 3) + bits.Len64(2) = 5
		// Inf branch does not dominate the max since Inf is encoded as -1
		{sixty.Coproduct{
			{"a", sixty.Bit{}},
			{"b", sixty.String{}},
		}, 3},
		{internal.BttvEmotesDBSchema, size.Inf},
	}
	for _, test := range testCases {
		require.Equal(t, test.size, size.SchemaSize(test.schema), "schema %T", test.schema)
	}
}

func TestSizeString(t *testing.T) {
	t.Parallel()

	require.Equal(t, "Inf", size.Inf.String())
	require.Equal(t, "7", size.Size(7).String())
}

func TestSizeAdd(t *testing.T) {
	t.Parallel()

	x := size.Size(2)
	require.Equal(t, size.Size(5), *x.Add(3))
	require.Equal(t, size.Size(5), *x.Add(0))

	y := size.Size(3)
	// adding Inf to a finite size is a no-op, not saturation
	require.Equal(t, size.Size(3), *y.Add(size.Inf))

	z := size.Inf
	require.Equal(t, size.Inf, *z.Add(4))
}

func TestSizeMultiply(t *testing.T) {
	t.Parallel()

	require.Equal(t, size.Size(6), size.Size(2).Multiply(3))
	require.Equal(t, size.Size(0), size.Size(2).Multiply(0))
	require.Equal(t, size.Inf, size.Inf.Multiply(3))
}
