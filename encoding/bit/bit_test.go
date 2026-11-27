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

	for _, test := range internal.TestCases {
		require.Equal(t, test.Value, internal.RoundTrip(t, test.Schema, test.Value, bit.Codec), "schema %T", test.Schema)
	}
}

func TestDecodeErrors(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		schema sixty.Schema
		data   []byte
		reason string
	}{
		// 2-bit tag "11" = 3 while enum has only 3 values (indices 0..2)
		{sixty.Enum{Values: []string{"a", "b", "c"}}, []byte{0xC0}, "found 3'th enum value while there are only 3 of them"},
		// gamma(3) = "00100" decodes to optional variant 3
		{sixty.Optional{Schema: sixty.Int{}}, []byte{0x20}, "option variant must be either 0 or 1, but it is 3"},
		// 2-bit tag "11" = variant 3 while coproduct has only 2 variants
		{
			sixty.Coproduct{{"a", sixty.Int{}}, {"b", sixty.String{}}},
			[]byte{0xC0},
			"found variant 3 while coproduct has only 2 of them",
		},
		// [-5..5] takes 4 bits: "1111" = 15 > 10 = Max-Min
		{sixty.EnumInt{Min: -5, Max: 5}, []byte{0xF0}, "out of range"},
		// 64 zero bits followed by a terminator bit: a valid gamma code has at
		// most 63 leading zeros
		{sixty.Uint{}, append(make([]byte, 8), 0x80), "uint is too long"},
		{sixty.String{}, []byte{}, "EOF"},
		// 0x80 is a continuation byte, not a valid utf-8 lead byte
		{sixty.Rune{}, []byte{0x80}, "invalid utf-8 lead byte"},
	}
	for _, test := range testCases {
		_, err := bit.Codec.Decode(test.schema, test.data)
		require.ErrorContains(t, err, test.reason, "schema %T data %v", test.schema, test.data)
	}
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
