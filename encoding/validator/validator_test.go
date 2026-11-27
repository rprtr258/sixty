package validator_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/rprtr258/sixty"
	"github.com/rprtr258/sixty/encoding/validator"
)

func TestValidSchema(t *testing.T) {
	t.Parallel()

	schema := sixty.Product{
		{"bit", sixty.Bit{}},
		{"int", sixty.Int{}},
		{"uint", sixty.Uint{}},
		{"rune", sixty.Rune{}},
		{"string", sixty.String{}},
		{"enum", sixty.Enum{Values: []string{"a", "b"}}},
		{"enum_int", sixty.EnumInt{Min: -5, Max: 5}},
		{"list", sixty.List{Schema: sixty.Bit{}}},
		{"array", sixty.Array{Schema: sixty.Bit{}, Count: 2}},
		{"optional", sixty.Optional{Schema: sixty.Bit{}}},
		{"coproduct", sixty.Coproduct{
			{"int", sixty.Int{}},
			{"string", sixty.String{}},
		}},
	}

	v := &validator.Validator{}
	v.Write(schema)
	require.Empty(t, v.Errs)
}

func TestInvalidSchemas(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		schema sixty.Schema
		reason string
	}{
		{sixty.EnumInt{Min: 5, Max: 5}, "incorrect enum int, max must be greater than min: min=5, max=5"},
		{sixty.Product{{"only", sixty.Bit{}}}, "wrong product field count, must be at least 2: 1"},
		// error inside a nested schema field is reported, not swallowed
		{
			sixty.Product{
				{"ok", sixty.Int{}},
				{"bad", sixty.Enum{Values: []string{"lonely"}}},
			},
			"wrong enum Values count, must be at least 2: [lonely]",
		},
		{sixty.Coproduct{{"only", sixty.Bit{}}}, "wrong coproduct field count, must be at least 2: 1"},
		{sixty.Array{Schema: sixty.Bit{}, Count: 1}, "wrong array length, must be at least 2: 1"},
	}
	for _, test := range testCases {
		v := &validator.Validator{}
		v.Write(test.schema)
		require.Len(t, v.Errs, 1, "schema %T", test.schema)
		require.Equal(t, test.reason, v.Errs[0].Error(), "schema %T", test.schema)
	}
}
