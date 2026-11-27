package bytes_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/rprtr258/sixty"
	"github.com/rprtr258/sixty/encoding/bytes"
	"github.com/rprtr258/sixty/internal"
)

func TestRoundTripAllTypes(t *testing.T) {
	t.Parallel()

	testCases := append(internal.PrimitiveTestCases(),
		internal.TestCase{
			Schema: sixty.Product{{"x", sixty.Int{}}, {"b", sixty.Bit{}}},
			Value:  sixty.ProductValue{sixty.IntValue(5), sixty.BitValue(true)},
		},
		internal.TestCase{
			Schema: sixty.Coproduct{{"int", sixty.Int{}}, {"string", sixty.String{}}},
			Value:  sixty.CoproductValue{SchemaValue: sixty.IntValue(42), Variant: 0},
		},
		internal.TestCase{
			Schema: sixty.Coproduct{{"int", sixty.Int{}}, {"string", sixty.String{}}},
			Value:  sixty.CoproductValue{SchemaValue: sixty.StringValue("hi"), Variant: 1},
		},
		internal.TestCase{
			Schema: sixty.List{Schema: sixty.Bit{}},
			Value:  sixty.ListValue{sixty.BitValue(true)},
		},
		internal.TestCase{
			Schema: sixty.List{Schema: sixty.Bit{}},
			Value:  sixty.ListValue{},
		},
		internal.TestCase{
			Schema: sixty.Array{Schema: sixty.Int{}, Count: 2},
			Value:  sixty.ArrayValue{sixty.IntValue(1), sixty.IntValue(2)},
		},
		internal.TestCase{
			Schema: sixty.Optional{Schema: sixty.Int{}},
			Value:  sixty.OptionalValue{IsValid: false},
		},
		internal.TestCase{
			Schema: sixty.Optional{Schema: sixty.Int{}},
			Value:  sixty.OptionalValue{SchemaValue: sixty.IntValue(7), IsValid: true},
		},
	)
	for _, test := range testCases {
		require.Equal(t, test.Value, internal.RoundTrip(t, test.Schema, test.Value, bytes.Codec), "schema %T", test.Schema)
	}
}

func TestDecodeEnumOutOfRange(t *testing.T) {
	t.Parallel()

	_, err := bytes.Codec.Decode(sixty.Enum{Values: []string{"a", "b"}}, []byte{0x02})
	require.ErrorContains(t, err, "found 2'th enum value while there are only 2 of them")
}

func TestDecodeEnumIntOutOfRange(t *testing.T) {
	t.Parallel()

	_, err := bytes.Codec.Decode(sixty.EnumInt{Min: 0, Max: 1}, []byte{0x05})
	require.ErrorContains(t, err, "out of range")
}

func TestDecodeOptionalInvalidMarker(t *testing.T) {
	t.Parallel()

	_, err := bytes.Codec.Decode(sixty.Optional{Schema: sixty.Int{}}, []byte{0x02})
	require.ErrorContains(t, err, "option variant must be either 0 or 1, but it is 2")
}

func TestDecodeCoproductVariantOutOfRange(t *testing.T) {
	t.Parallel()

	// variant 3 while coproduct has only 2 variants
	_, err := bytes.Codec.Decode(
		sixty.Coproduct{{"a", sixty.Int{}}, {"b", sixty.String{}}},
		[]byte{0x03},
	)
	require.ErrorContains(t, err, "found variant 3 while coproduct has only 2 of them")
}

func TestDecodeErrorPosition(t *testing.T) {
	t.Parallel()

	// 1-byte varint prefix + invalid bit: error reported after 2 scanned bytes
	_, err := bytes.Codec.Decode(sixty.List{Schema: sixty.Bit{}}, []byte{0x01, 0x05})
	require.ErrorContains(t, err, "error at position 2:")
}
