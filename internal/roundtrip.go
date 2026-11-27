package internal

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/rprtr258/sixty"
)

// RoundTrip encodes value with schema and decodes it back with the given
// encoder/decoder pair; both operate on the same data type T.
func RoundTrip[T any](
	t *testing.T,
	schema sixty.Schema,
	value sixty.SchemaValue,
	encode func(sixty.Schema, sixty.SchemaValue) (T, error),
	decode func(sixty.Schema, T) (sixty.SchemaValue, error),
) sixty.SchemaValue {
	t.Helper()

	data, err := encode(schema, value)
	require.NoError(t, err)
	got, err := decode(schema, data)
	require.NoError(t, err)
	return got
}
