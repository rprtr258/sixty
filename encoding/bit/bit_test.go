package bit_test

import (
	"testing"

	"github.com/rprtr258/sixty"
	"github.com/rprtr258/sixty/encoding/bit"
	"github.com/rprtr258/sixty/internal"
	"github.com/stretchr/testify/require"
)

var bttvEmotesDBBits = []byte{
	0xad, 0x52, 0xd8, 0x5c, 0x1c, 0x18, 0x6b, 0xad, 0x2d, 0x85, 0xc1, 0xc1, 0x85, 0x3d, 0xdd, 0xb9,
	0x95, 0xc8, 0xc7, 0xad, 0x2d, 0x85, 0xc1, 0xc1, 0x85, 0x3d, 0xdd, 0xb9, 0x95, 0xc8, 0xcb, 0xa1,
	0x41, 0xbd, 0x9d, 0x0d, 0xa1, 0x85, 0xb5, 0xc2, 0xbb, 0x94, 0x1b, 0xd9, 0xd0, 0xda, 0x18, 0x5b,
	0x5c, 0x13, 0xdd, 0xdb, 0x99, 0x5c, 0x8c, 0x7b, 0x94, 0x1b, 0xd9, 0xd0, 0xda, 0x18, 0x5b, 0x5c,
	0x13, 0xdd, 0xdb, 0x99, 0x5c, 0x8c, 0x80,
}

func TestBttvDBDecode(t *testing.T) {
	t.Parallel()

	got, err := bit.Decode(internal.BttvEmotesDBSchema, bttvEmotesDBBits)
	require.NoError(t, err)
	require.Equal(t, internal.BttvEmotesDBData, got)
}

func TestBttvDBEncode(t *testing.T) {
	t.Parallel()

	got, err := bit.Encode(internal.BttvEmotesDBSchema, internal.BttvEmotesDBData)
	require.NoError(t, err)
	require.Equal(t, bttvEmotesDBBits, got)
}

func TestRoundTripAllTypes(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		schema sixty.Schema
		value  sixty.SchemaValue
	}{
		{sixty.Bit{}, sixty.BitValue(true)},
		{sixty.Bit{}, sixty.BitValue(false)},
		{sixty.Int{}, sixty.IntValue(-42)},
		{sixty.Int{}, sixty.IntValue(0)},
		{sixty.Uint{}, sixty.UintValue(42)},
		{sixty.Rune{}, sixty.RuneValue('a')},
		{sixty.Enum{Values: []string{"a", "b", "c"}}, sixty.EnumValue(0)},
		{sixty.Enum{Values: []string{"a", "b", "c"}}, sixty.EnumValue(2)},
		{sixty.EnumInt{Min: -5, Max: 5}, sixty.EnumIntValue(0)},
		{sixty.EnumInt{Min: -5, Max: 5}, sixty.EnumIntValue(5)},
		{sixty.Product{{"x", sixty.Int{}}, {"b", sixty.Bit{}}}, sixty.ProductValue{sixty.IntValue(5), sixty.BitValue(true)}},
		{sixty.List{Schema: sixty.Bit{}}, sixty.ListValue{sixty.BitValue(true), sixty.BitValue(false)}},
		{sixty.List{Schema: sixty.Bit{}}, sixty.ListValue{}},
		{sixty.Array{Schema: sixty.Int{}, Count: 2}, sixty.ArrayValue{sixty.IntValue(1), sixty.IntValue(2)}},
		{sixty.Optional{Schema: sixty.Int{}}, sixty.OptionalValue{IsValid: false}},
		{sixty.String{}, sixty.StringValue("Kappa")},
	}
	for _, test := range testCases {
		require.Equal(t, test.value, internal.RoundTrip(t, test.schema, test.value, bit.Encode, bit.Decode), "schema %T", test.schema)
	}
}

func TestDecodeOptionalPresent(t *testing.T) {
	t.Parallel()

	// marker byte 0x01 (present) + zigzag(42)=84 in unary format: 111111 0 1010100
	got, err := bit.Decode(sixty.Optional{Schema: sixty.Int{}}, []byte{0x01, 0xFD, 0x50})
	require.NoError(t, err)
	require.Equal(t, sixty.OptionalValue{SchemaValue: sixty.IntValue(42), IsValid: true}, got)
}

func TestDecodeCoproduct(t *testing.T) {
	t.Parallel()

	// variant tag "00" (variant 0 of 2) + zigzag(42)=84 unary: 111111 0 1010100
	got, err := bit.Decode(
		sixty.Coproduct{{"a", sixty.Int{}}, {"b", sixty.String{}}},
		[]byte{0x3F, 0x54},
	)
	require.NoError(t, err)
	require.Equal(t,
		sixty.CoproductValue{SchemaValue: sixty.IntValue(42), Variant: 0},
		got,
	)
}

func TestCoproductRoundTrip(t *testing.T) {
	t.Parallel()

	schema := sixty.Coproduct{{"a", sixty.Int{}}, {"b", sixty.String{}}}
	for _, test := range []sixty.CoproductValue{
		{SchemaValue: sixty.IntValue(42), Variant: 0},
		{SchemaValue: sixty.StringValue("hi"), Variant: 1},
	} {
		require.Equal(t, test, internal.RoundTrip(t, schema, test, bit.Encode, bit.Decode))
	}
}

func TestDecodeEnumOutOfRangePanics(t *testing.T) {
	t.Parallel()

	// 2-bit tag "11" = 3 while enum has only 3 values (indices 0..2)
	require.Panics(t, func() {
		_, _ = bit.Decode(sixty.Enum{Values: []string{"a", "b", "c"}}, []byte{0xC0})
	})
}

func TestDecodeOptionalInvalidMarkerPanics(t *testing.T) {
	t.Parallel()

	require.Panics(t, func() {
		_, _ = bit.Decode(sixty.Optional{Schema: sixty.Int{}}, []byte{0x02})
	})
}

func TestDecodeTruncated(t *testing.T) {
	t.Parallel()

	_, err := bit.Decode(sixty.String{}, []byte{})
	require.Error(t, err)
}

// encoding a present optional works; decoding its own output panics: the
// encoder writes the marker as a writeUint bit-run (1-2 bits) while the
// decoder reads it as a full byte.
func TestEncodeOptionalPresent(t *testing.T) {
	t.Parallel()

	data, err := bit.Encode(sixty.Optional{Schema: sixty.Int{}},
		sixty.OptionalValue{SchemaValue: sixty.IntValue(42), IsValid: true})
	require.NoError(t, err)
	// absent-marker bit 0 + present-marker bit 1 + zigzag(42)=84 unary
	require.Equal(t, []byte{0x7F, 0x54}, data)
}
