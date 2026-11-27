package internal

import (
	_ "embed"

	"github.com/rprtr258/sixty"
)

// idLength is the fixed length of bttv owner ids, mirroring real ids being
// fixed-width strings.
const idLength = 8

// TestCase pairs a schema with a value valid for that schema.
type TestCase struct {
	Schema sixty.Schema
	Value  sixty.SchemaValue
}

// PrimitiveTestCases covers every primitive and logical type; codecs extend
// the list with their composite/codec-specific cases.
func PrimitiveTestCases() []TestCase {
	return []TestCase{
		{sixty.Bit{}, sixty.BitValue(true)},
		{sixty.Bit{}, sixty.BitValue(false)},
		{sixty.Int{}, sixty.IntValue(-42)},
		{sixty.Int{}, sixty.IntValue(0)},
		{sixty.Uint{}, sixty.UintValue(42)},
		{sixty.Uint{}, sixty.UintValue(0)},
		{sixty.Rune{}, sixty.RuneValue('a')},
		{sixty.Rune{}, sixty.RuneValue('ж')},
		{sixty.Enum{Values: []string{"a", "b", "c"}}, sixty.EnumValue(0)},
		{sixty.Enum{Values: []string{"a", "b", "c"}}, sixty.EnumValue(2)},
		{sixty.EnumInt{Min: -5, Max: 5}, sixty.EnumIntValue(-5)},
		{sixty.EnumInt{Min: -5, Max: 5}, sixty.EnumIntValue(5)},
		{sixty.String{}, sixty.StringValue("")},
		{sixty.String{}, sixty.StringValue("ёжик")},
	}
}

// idArray splits an id string into the fixed-length rune array value.
func idArray(id string) sixty.ArrayValue {
	res := make(sixty.ArrayValue, 0, idLength)
	for _, r := range id {
		res = append(res, sixty.RuneValue(r))
	}
	return res
}

var (
	BttvEmotesDBSchema = sixty.List{Schema: sixty.Product{
		{"emote", sixty.String{}},
		{"ids", sixty.List{
			Schema: sixty.Array{Schema: sixty.Rune{}, Count: idLength},
		}},
	}}
	BttvEmotesDBData = sixty.ListValue{sixty.ProductValue{
		sixty.StringValue("Kappa"),
		sixty.ListValue{
			idArray("kappa_01"),
			idArray("kappa_02"),
		},
	}, sixty.ProductValue{
		sixty.StringValue("PogChamp"),
		sixty.ListValue{
			idArray("pogch_01"),
			idArray("pogch_02"),
		},
	}}
)
