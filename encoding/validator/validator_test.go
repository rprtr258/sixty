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

func TestInvalidEnumInt(t *testing.T) {
	t.Parallel()

	v := &validator.Validator{}
	v.Write(sixty.EnumInt{Min: 5, Max: 5})
	require.Len(t, v.Errs, 1)
	require.Equal(t, "incorrect enum int, max must be greater than min: min=5, max=5", v.Errs[0].Error())
}

func TestInvalidProduct(t *testing.T) {
	t.Parallel()

	v := &validator.Validator{}
	v.Write(sixty.Product{{"only", sixty.Bit{}}})
	require.Len(t, v.Errs, 1)
	require.Equal(t, "wrong product field count, must be at least 2: 1", v.Errs[0].Error())
}

func TestProductNestedError(t *testing.T) {
	t.Parallel()

	v := &validator.Validator{}
	v.Write(sixty.Product{
		{"ok", sixty.Int{}},
		{"bad", sixty.Enum{Values: []string{"lonely"}}},
	})
	require.Len(t, v.Errs, 1)
	require.Equal(t, "wrong enum Values count, must be at least 2: [lonely]", v.Errs[0].Error())
}

func TestInvalidCoproduct(t *testing.T) {
	t.Parallel()

	v := &validator.Validator{}
	v.Write(sixty.Coproduct{{"only", sixty.Bit{}}})
	require.Len(t, v.Errs, 1)
	require.Equal(t, "wrong coproduct field count, must be at least 2: 1", v.Errs[0].Error())
}

func TestInvalidArray(t *testing.T) {
	t.Parallel()

	v := &validator.Validator{}
	v.Write(sixty.Array{Schema: sixty.Bit{}, Count: 1})
	require.Len(t, v.Errs, 1)
	require.Equal(t, "wrong array length, must be at least 2: 1", v.Errs[0].Error())
}
