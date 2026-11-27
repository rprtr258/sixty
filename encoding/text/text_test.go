package text_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/rprtr258/sixty"
	"github.com/rprtr258/sixty/encoding/text"
	"github.com/rprtr258/sixty/internal"
)

var (
	allTypesSchema = sixty.Product{
		{"bit", sixty.Bit{}},
		{"int", sixty.Int{}},
		{"uint", sixty.Uint{}},
		{"rune", sixty.Rune{}},
		{"enum", sixty.Enum{Values: []string{"yes", "no"}}},
		{"enum_int", sixty.EnumInt{Min: -5, Max: 5}},
		{"optional_present", sixty.Optional{Schema: sixty.Int{}}},
		{"optional_absent", sixty.Optional{Schema: sixty.Int{}}},
		{"array", sixty.Array{Schema: sixty.Bit{}, Count: 2}},
		{"coproduct", sixty.Coproduct{
			{"int", sixty.Int{}},
			{"string", sixty.String{}},
		}},
		{"empty_list", sixty.List{Schema: sixty.Bit{}}},
		{"string", sixty.String{}},
	}
	allTypesValue = sixty.ProductValue{
		sixty.BitValue(true),
		sixty.IntValue(-42),
		sixty.UintValue(42),
		sixty.RuneValue('ж'),
		sixty.EnumValue(1),
		sixty.EnumIntValue(-3),
		sixty.OptionalValue{SchemaValue: sixty.IntValue(7), IsValid: true},
		sixty.OptionalValue{IsValid: false},
		sixty.ArrayValue{sixty.BitValue(false), sixty.BitValue(true)},
		sixty.CoproductValue{SchemaValue: sixty.StringValue("hi"), Variant: 1},
		sixty.ListValue{},
		sixty.StringValue(`he said "hi"`),
	}
)

func TestBttvDBEncode(t *testing.T) {
	t.Parallel()

	got, err := text.Encode(internal.BttvEmotesDBSchema, internal.BttvEmotesDBData)
	assert.NoError(t, err)
	assert.Equal(t, `[{"Kappa",["KappaOwner1","KappaOwner2"]},{"PogChamp",["PogChampOwner1","PogChampOwner2"]}]`, got)
}

func TestBttvDBDecode(t *testing.T) {
	t.Parallel()

	bttvEmotesDBText := `[{"Kappa",["KappaOwner1","KappaOwner2"]},
		{"PogChamp",["PogChampOwner1","PogChampOwner2"]}]`
	got, err := text.Decode(internal.BttvEmotesDBSchema, bttvEmotesDBText)
	assert.NoError(t, err)
	assert.Equal(t, internal.BttvEmotesDBData, got)
}

func TestBttvDBEncodeDecode(t *testing.T) {
	t.Parallel()

	encoded, err := text.Encode(internal.BttvEmotesDBSchema, internal.BttvEmotesDBData)
	assert.NoError(t, err)
	got, err := text.Decode(internal.BttvEmotesDBSchema, encoded)
	assert.NoError(t, err)
	assert.Equal(t, internal.BttvEmotesDBData, got)
}

func TestAllTypesEncodeDecode(t *testing.T) {
	t.Parallel()

	encoded, err := text.Encode(allTypesSchema, allTypesValue)
	assert.NoError(t, err)
	got, err := text.Decode(allTypesSchema, encoded)
	assert.NoError(t, err)
	assert.Equal(t, allTypesValue, got)
}

func TestDecodeGarbage(t *testing.T) {
	t.Parallel()

	_, err := text.Decode(internal.BttvEmotesDBSchema, "garbage")
	assert.Error(t, err)
}
