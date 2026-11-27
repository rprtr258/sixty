package bit_test

import (
	"testing"

	"github.com/rprtr258/sixty"
	"github.com/rprtr258/sixty/encoding/bit"
	"github.com/rprtr258/sixty/internal"
	"github.com/stretchr/testify/require"
)

var bttvEmotesDBBits = []byte{
	0x66, 0x4b, 0x61, 0x70, 0x70, 0x61, 0x6d, 0x6c, 0x2e, 0x0e, 0x0c, 0x2b, 0xe6, 0x06, 0x2d, 0x6c,
	0x2e, 0x0e, 0x0c, 0x2b, 0xe6, 0x06, 0x42, 0x54, 0x1b, 0xd9, 0xd0, 0xda, 0x18, 0x5b, 0x5c, 0x1b,
	0x83, 0x7b, 0x3b, 0x1b, 0x42, 0xf9, 0x81, 0x8b, 0x83, 0x7b, 0x3b, 0x1b, 0x42, 0xf9, 0x81, 0x90,
}

func TestBttvDBDecode(t *testing.T) {
	t.Parallel()

	got, err := bit.Codec.Decode(internal.BttvEmotesDBSchema, bttvEmotesDBBits)
	require.NoError(t, err)
	require.Equal(t, internal.BttvEmotesDBData, got)
}

func TestBttvDBEncode(t *testing.T) {
	t.Parallel()

	got, err := bit.Codec.Encode(internal.BttvEmotesDBSchema, internal.BttvEmotesDBData)
	require.NoError(t, err)
	require.Equal(t, bttvEmotesDBBits, got)
}

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

func TestWriteUintGamma(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		value    uint64
		encoding []byte
	}{
		// Elias gamma of value+1: leading zeros then the binary of value+1
		{0, []byte{0x80}},        // 1
		{1, []byte{0x40}},        // 010
		{2, []byte{0x60}},        // 011
		{42, []byte{0x05, 0x60}}, // 00000 101011
	}
	for _, test := range testCases {
		got, err := bit.Codec.Encode(sixty.Uint{}, sixty.UintValue(test.value))
		require.NoError(t, err)
		require.Equal(t, test.encoding, got, "value %d", test.value)
	}
}

func TestDecodeOptionalPresent(t *testing.T) {
	t.Parallel()

	// present marker gamma(2)="010" + zigzag(42)=84 as gamma(85)="0000001010101"
	got, err := bit.Codec.Decode(sixty.Optional{Schema: sixty.Int{}}, []byte{0x40, 0x55})
	require.NoError(t, err)
	require.Equal(t, sixty.OptionalValue{SchemaValue: sixty.IntValue(42), IsValid: true}, got)
}

func TestDecodeCoproduct(t *testing.T) {
	t.Parallel()

	// variant tag "00" (variant 0 of 2) + zigzag(42)=84 as gamma(85)
	got, err := bit.Codec.Decode(
		sixty.Coproduct{{"a", sixty.Int{}}, {"b", sixty.String{}}},
		[]byte{0x00, 0xaa},
	)
	require.NoError(t, err)
	require.Equal(t,
		sixty.CoproductValue{SchemaValue: sixty.IntValue(42), Variant: 0},
		got,
	)
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

// regression: encoding 0 used to write 1 bit while decoding consumed 2 bits,
// silently desyncing composite values.
func TestZeroUintInsideProduct(t *testing.T) {
	t.Parallel()

	schema := sixty.Product{{"a", sixty.Uint{}}, {"b", sixty.Uint{}}}
	value := sixty.ProductValue{sixty.UintValue(0), sixty.UintValue(1)}
	require.Equal(t, value, internal.RoundTrip(t, schema, value, bit.Codec))
}

func TestRoundTripNearMaxUint(t *testing.T) {
	t.Parallel()

	require.Equal(t,
		sixty.UintValue(0xFFFFFFFFFFFFFFFE),
		internal.RoundTrip(t, sixty.Uint{}, sixty.UintValue(0xFFFFFFFFFFFFFFFE), bit.Codec),
	)
}

func TestEncodeMaxUintPanics(t *testing.T) {
	t.Parallel()

	// gamma coding needs value+1 to fit in uint64
	require.Panics(t, func() {
		_, _ = bit.Codec.Encode(sixty.Int{}, sixty.IntValue(-0x8000000000000000))
	})
}
