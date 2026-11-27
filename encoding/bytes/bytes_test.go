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

	for _, test := range internal.TestCases {
		require.Equal(t, test.Value, internal.RoundTrip(t, test.Schema, test.Value, bytes.Codec), "schema %T", test.Schema)
	}
}

func TestDecodeErrors(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		schema sixty.Schema
		data   []byte
		reason string
	}{
		{sixty.Enum{Values: []string{"a", "b"}}, []byte{0x02}, "found 2'th enum value while there are only 2 of them"},
		{sixty.EnumInt{Min: 0, Max: 1}, []byte{0x05}, "out of range"},
		{sixty.Optional{Schema: sixty.Int{}}, []byte{0x02}, "option variant must be either 0 or 1, but it is 2"},
		// variant 3 while coproduct has only 2 variants
		{sixty.Coproduct{{"a", sixty.Int{}}, {"b", sixty.String{}}}, []byte{0x03}, "found variant 3 while coproduct has only 2 of them"},
		// 1-byte varint prefix + invalid bit: error reported after 2 scanned bytes
		{sixty.List{Schema: sixty.Bit{}}, []byte{0x01, 0x05}, "error at position 2:"},
	}
	for _, test := range testCases {
		_, err := bytes.Codec.Decode(test.schema, test.data)
		require.ErrorContains(t, err, test.reason, "schema %T data %v", test.schema, test.data)
	}
}
