package internal

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/rprtr258/sixty"
)

// RoundTrip encodes value with schema and decodes it back with codec.
func RoundTrip[T any](
	t *testing.T,
	schema sixty.Schema,
	value sixty.SchemaValue,
	codec sixty.Codec[T],
) sixty.SchemaValue {
	t.Helper()

	data, err := codec.Encode(schema, value)
	require.NoError(t, err)
	got, err := codec.Decode(schema, data)
	require.NoError(t, err)
	return got
}
