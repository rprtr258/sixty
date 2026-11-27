package sixty_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/rprtr258/sixty"
)

func TestSchemaKinds(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		schema sixty.Schema
		kind   sixty.Kind
	}{
		{sixty.Bit{}, sixty.BitKind},
		{sixty.Int{}, sixty.IntKind},
		{sixty.Uint{}, sixty.UintKind},
		{sixty.Rune{}, sixty.RuneKind},
		{sixty.Enum{}, sixty.EnumKind},
		{sixty.EnumInt{}, sixty.EnumIntKind},
		{sixty.Product{}, sixty.ProductKind},
		{sixty.Coproduct{}, sixty.CoproductKind},
		{sixty.List{}, sixty.ListKind},
		{sixty.Array{}, sixty.ArrayKind},
		{sixty.Optional{}, sixty.OptionalKind},
		// :string is a logical type over []rune, so its schema kind is list
		{sixty.String{}, sixty.ListKind},
	}
	for _, test := range testCases {
		require.Equal(t, test.kind, test.schema.Kind())
	}
}

func TestSchemaValuesKinds(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		value sixty.SchemaValue
		kind  sixty.Kind
	}{
		{sixty.BitValue(true), sixty.BitKind},
		{sixty.IntValue(-1), sixty.IntKind},
		{sixty.UintValue(1), sixty.UintKind},
		{sixty.RuneValue('a'), sixty.RuneKind},
		{sixty.EnumValue(0), sixty.EnumKind},
		{sixty.EnumIntValue(0), sixty.EnumIntKind},
		{sixty.ProductValue{}, sixty.ProductKind},
		{sixty.CoproductValue{}, sixty.CoproductKind},
		{sixty.ListValue{}, sixty.ListKind},
		{sixty.ArrayValue{}, sixty.ArrayKind},
		{sixty.OptionalValue{}, sixty.OptionalKind},
		{sixty.StringValue(""), sixty.ListKind},
	}
	for _, test := range testCases {
		require.Equal(t, test.kind, test.value.Kind())
		// every value must satisfy the SchemaValue interface
		var _ sixty.SchemaValue = test.value
	}
}

func TestStringLogicalKind(t *testing.T) {
	t.Parallel()

	require.Equal(t, sixty.StringKind, sixty.String{}.PredefKind())
	var _ sixty.LogicalSchema = sixty.String{}
}
