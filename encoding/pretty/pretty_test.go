package pretty_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/rprtr258/sixty"
	"github.com/rprtr258/sixty/encoding/pretty"
	"github.com/rprtr258/sixty/internal"
)

func TestBttvDBEncode(t *testing.T) {
	t.Parallel()

	got, err := pretty.Codec.Encode(internal.BttvEmotesDBSchema, internal.BttvEmotesDBData)
	require.NoError(t, err)
	require.Equal(t, `[
  {
    "Kappa",
    [
      [
      'k',
      'a',
      'p',
      'p',
      'a',
      '_',
      '0',
      '1',
      ],
      [
      'k',
      'a',
      'p',
      'p',
      'a',
      '_',
      '0',
      '2',
      ],
    ],
  },
  {
    "PogChamp",
    [
      [
      'p',
      'o',
      'g',
      'c',
      'h',
      '_',
      '0',
      '1',
      ],
      [
      'p',
      'o',
      'g',
      'c',
      'h',
      '_',
      '0',
      '2',
      ],
    ],
  },
]`, got)
}

func TestBttvDBEncodeDecode(t *testing.T) {
	t.Parallel()

	schema := internal.BttvEmotesDBSchema
	data := internal.BttvEmotesDBData
	enc, err := pretty.Codec.Encode(schema, data)
	require.NoError(t, err)
	dec, err := pretty.Codec.Decode(schema, enc)
	require.NoError(t, err)
	require.Equal(t, dec, data)
}

func TestRuneQuoteEscape(t *testing.T) {
	t.Parallel()

	schema := sixty.Rune{}
	data := sixty.RuneValue('\'')
	enc, err := pretty.Codec.Encode(schema, data)
	require.NoError(t, err)
	require.Equal(t, `'\''`, enc)
	dec, err := pretty.Codec.Decode(schema, enc)
	require.NoError(t, err)
	require.Equal(t, dec, data)
}

func TestAllTypesEncode(t *testing.T) {
	t.Parallel()

	schema := sixty.Product{
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
	data := sixty.ProductValue{
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
	enc, err := pretty.Codec.Encode(schema, data)
	require.NoError(t, err)
	dec, err := pretty.Codec.Decode(schema, enc)
	require.NoError(t, err)
	require.Equal(t, dec, data)
}
