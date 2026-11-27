package json_test

import (
	_ "embed"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/rprtr258/sixty"
	"github.com/rprtr258/sixty/encoding/json"
	"github.com/rprtr258/sixty/internal"
)

//go:embed bttv_emotes_db.json
var bttvEmotesDBJson string

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
		sixty.StringValue("hello world"),
	}
)

func TestBttvDBEncode(t *testing.T) {
	t.Parallel()

	got, err := json.Encode(internal.BttvEmotesDBSchema, internal.BttvEmotesDBData)
	assert.NoError(t, err)
	assert.Equal(t, bttvEmotesDBJson, got)
}

func TestBttvDBDecode(t *testing.T) {
	t.Parallel()

	got, err := json.Decode(internal.BttvEmotesDBSchema, bttvEmotesDBJson)
	assert.NoError(t, err)
	assert.Equal(t, internal.BttvEmotesDBData, got)
}

func TestAllTypesEncodeDecode(t *testing.T) {
	t.Parallel()

	encoded, err := json.Encode(allTypesSchema, allTypesValue)
	assert.NoError(t, err)
	got, err := json.Decode(allTypesSchema, encoded)
	assert.NoError(t, err)
	assert.Equal(t, allTypesValue, got)
}
