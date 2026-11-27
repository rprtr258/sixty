package bit_test

import (
	"testing"

	"github.com/rprtr258/sixty"
	"github.com/rprtr258/sixty/encoding/bit"
	"github.com/rprtr258/sixty/internal"
	"github.com/stretchr/testify/require"
)

func TestRoundTripAllTypes(t *testing.T) {
	t.Parallel()

	testCases := append(internal.PrimitiveTestCases(),
		internal.TestCase{
			Schema: sixty.Product{{"x", sixty.Int{}}, {"b", sixty.Bit{}}},
			Value:  sixty.ProductValue{sixty.IntValue(5), sixty.BitValue(true)},
		},
		internal.TestCase{
			Schema: sixty.Coproduct{{"a", sixty.Int{}}, {"b", sixty.String{}}},
			Value:  sixty.CoproductValue{SchemaValue: sixty.StringValue("hi"), Variant: 1},
		},
		internal.TestCase{
			Schema: sixty.List{Schema: sixty.Bit{}},
			Value:  sixty.ListValue{sixty.BitValue(true), sixty.BitValue(false)},
		},
		internal.TestCase{
			Schema: sixty.Array{Schema: sixty.Int{}, Count: 2},
			Value:  sixty.ArrayValue{sixty.IntValue(1), sixty.IntValue(2)},
		},
		internal.TestCase{
			Schema: sixty.Optional{Schema: sixty.Int{}},
			Value:  sixty.OptionalValue{SchemaValue: sixty.IntValue(7), IsValid: true},
		},
		internal.TestCase{
			Schema: sixty.Optional{Schema: sixty.Int{}},
			Value:  sixty.OptionalValue{IsValid: false},
		},
	)
	for _, test := range testCases {
		require.Equal(t, test.Value, internal.RoundTrip(t, test.Schema, test.Value, bit.Codec), "schema %T", test.Schema)
	}
}

func TestDecodeEnumOutOfRange(t *testing.T) {
	t.Parallel()

	// 2-bit tag "11" = 3 while enum has only 3 values (indices 0..2)
	_, err := bit.Codec.Decode(sixty.Enum{Values: []string{"a", "b", "c"}}, []byte{0xC0})
	require.ErrorContains(t, err, "found 3'th enum value while there are only 3 of them")
}

func TestDecodeOptionalInvalidMarker(t *testing.T) {
	t.Parallel()

	// gamma(3) = "00100" decodes to optional variant 3
	_, err := bit.Codec.Decode(sixty.Optional{Schema: sixty.Int{}}, []byte{0x20})
	require.ErrorContains(t, err, "option variant must be either 0 or 1, but it is 3")
}

func TestDecodeCoproductVariantOutOfRange(t *testing.T) {
	t.Parallel()

	// 2-bit tag "11" = variant 3 while coproduct has only 2 variants
	_, err := bit.Codec.Decode(
		sixty.Coproduct{{"a", sixty.Int{}}, {"b", sixty.String{}}},
		[]byte{0xC0},
	)
	require.ErrorContains(t, err, "found variant 3 while coproduct has only 2 of them")
}

func TestDecodeEnumIntOutOfRange(t *testing.T) {
	t.Parallel()

	// [0..1] takes 1 bit, value "1" is valid but 2-bit tags cannot appear...
	// [-5..5] takes 4 bits: "1111" = 15 > 10 = Max-Min
	_, err := bit.Codec.Decode(sixty.EnumInt{Min: -5, Max: 5}, []byte{0xF0})
	require.ErrorContains(t, err, "out of range")
}

func TestDecodeUintTooLong(t *testing.T) {
	t.Parallel()

	// 64 zero bits followed by a terminator bit: a valid gamma code has at
	// most 63 leading zeros
	_, err := bit.Codec.Decode(sixty.Uint{}, append(make([]byte, 8), 0x80))
	require.ErrorContains(t, err, "uint is too long")
}

func TestDecodeTruncated(t *testing.T) {
	t.Parallel()

	_, err := bit.Codec.Decode(sixty.String{}, []byte{})
	require.Error(t, err)
}

func TestDecodeInvalidRune(t *testing.T) {
	t.Parallel()

	// 0x80 is a continuation byte, not a valid utf-8 lead byte
	_, err := bit.Codec.Decode(sixty.Rune{}, []byte{0x80})
	require.ErrorContains(t, err, "invalid utf-8 lead byte")
}

func TestUnicodeRoundTrip(t *testing.T) {
	t.Parallel()

	// string length prefix counts runes: "ёжик" is 5 runes but 9 bytes
	require.Equal(t,
		sixty.StringValue("ёжик"),
		internal.RoundTrip(t, sixty.String{}, sixty.StringValue("ёжик"), bit.Codec),
	)
	require.Equal(t,
		sixty.StringValue("κ산☺"),
		internal.RoundTrip(t, sixty.String{}, sixty.StringValue("κ산☺"), bit.Codec),
	)
	require.Equal(t,
		sixty.RuneValue('ж'),
		internal.RoundTrip(t, sixty.Rune{}, sixty.RuneValue('ж'), bit.Codec),
	)
}

func TestEncodeMaxUintPanics(t *testing.T) {
	t.Parallel()

	// gamma coding needs value+1 to fit in uint64
	require.Panics(t, func() {
		_, _ = bit.Codec.Encode(sixty.Int{}, sixty.IntValue(-0x8000000000000000))
	})
}
