package bytes_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/rprtr258/sixty"
	"github.com/rprtr258/sixty/encoding/bytes"
	"github.com/rprtr258/sixty/internal"
)

var bttvEmotesDBBytes = []byte{
	0x02,
	0x05, 0x4b, 0x61, 0x70, 0x70, 0x61, // "Kappa"
	0x02,
	0x0b, 0x4b, 0x61, 0x70, 0x70, 0x61, 0x4f, 0x77, 0x6e, 0x65, 0x72, 0x31, // "KappaOwner1"
	0x0b, 0x4b, 0x61, 0x70, 0x70, 0x61, 0x4f, 0x77, 0x6e, 0x65, 0x72, 0x32, // "KappaOwner2"
	0x08, 0x50, 0x6f, 0x67, 0x43, 0x68, 0x61, 0x6d, 0x70, // "PogChamp"
	0x02,
	0x0e, 0x50, 0x6f, 0x67, 0x43, 0x68, 0x61, 0x6d, 0x70, 0x4f, 0x77, 0x6e, 0x65, 0x72, 0x31, // "PogChampOwner1"
	0x0e, 0x50, 0x6f, 0x67, 0x43, 0x68, 0x61, 0x6d, 0x70, 0x4f, 0x77, 0x6e, 0x65, 0x72, 0x32, // "PogChampOwner2"
}

func TestBttvDBDecode(t *testing.T) {
	t.Parallel()

	got, err := bytes.Decode(internal.BttvEmotesDBSchema, bttvEmotesDBBytes)
	require.NoError(t, err)
	require.Equal(t, internal.BttvEmotesDBData, got)
}

func TestBttvDBEncode(t *testing.T) {
	t.Parallel()

	got, err := bytes.Encode(internal.BttvEmotesDBSchema, internal.BttvEmotesDBData)
	require.NoError(t, err)
	require.Equal(t, bttvEmotesDBBytes, got)
}

// roundTrip encodes value with schema and decodes it back.
func roundTrip(t *testing.T, schema sixty.Schema, value sixty.SchemaValue) sixty.SchemaValue {
	t.Helper()

	data, err := bytes.Encode(schema, value)
	require.NoError(t, err)
	got, err := bytes.Decode(schema, data)
	require.NoError(t, err)
	return got
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
		{sixty.Uint{}, sixty.UintValue(300)}, // multi-byte varint
		{sixty.Rune{}, sixty.RuneValue('ж')}, // multi-byte utf-8
		{sixty.Enum{Values: []string{"a", "b", "c"}}, sixty.EnumValue(2)},
		{sixty.EnumInt{Min: -5, Max: 5}, sixty.EnumIntValue(-3)},
		{
			sixty.Product{{"x", sixty.Int{}}, {"b", sixty.Bit{}}},
			sixty.ProductValue{sixty.IntValue(5), sixty.BitValue(true)},
		},
		{
			sixty.Coproduct{{"int", sixty.Int{}}, {"string", sixty.String{}}},
			sixty.CoproductValue{SchemaValue: sixty.IntValue(42), Variant: 0},
		},
		{
			sixty.Coproduct{{"int", sixty.Int{}}, {"string", sixty.String{}}},
			sixty.CoproductValue{SchemaValue: sixty.StringValue("hi"), Variant: 1},
		},
		{sixty.List{Schema: sixty.Bit{}}, sixty.ListValue{sixty.BitValue(true)}},
		{sixty.List{Schema: sixty.Bit{}}, sixty.ListValue{}},
		{
			sixty.Array{Schema: sixty.Int{}, Count: 2},
			sixty.ArrayValue{sixty.IntValue(1), sixty.IntValue(2)},
		},
		{sixty.Optional{Schema: sixty.Int{}}, sixty.OptionalValue{IsValid: false}},
		{
			sixty.Optional{Schema: sixty.Int{}},
			sixty.OptionalValue{SchemaValue: sixty.IntValue(7), IsValid: true},
		},
		{sixty.String{}, sixty.StringValue("κ산☺")}, // rune-count prefix, multi-byte payload
	}
	for _, test := range testCases {
		require.Equal(t, test.value, roundTrip(t, test.schema, test.value), "schema %T", test.schema)
	}
}

func TestDecodeInvalidBitPanics(t *testing.T) {
	t.Parallel()

	require.Panics(t, func() {
		_, _ = bytes.Decode(sixty.Bit{}, []byte{0x02})
	})
}

func TestDecodeEnumOutOfRangePanics(t *testing.T) {
	t.Parallel()

	require.Panics(t, func() {
		_, _ = bytes.Decode(sixty.Enum{Values: []string{"a", "b"}}, []byte{0x02})
	})
}

func TestDecodeOptionalInvalidMarkerPanics(t *testing.T) {
	t.Parallel()

	require.Panics(t, func() {
		_, _ = bytes.Decode(sixty.Optional{Schema: sixty.Int{}}, []byte{0x02})
	})
}
