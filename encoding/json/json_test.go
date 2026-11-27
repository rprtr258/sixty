package json_test

import (
	_ "embed"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/rprtr258/sixty"
	"github.com/rprtr258/sixty/encoding/json"
	"github.com/rprtr258/sixty/internal"
)

const (
	// errExpectedIntQuoted is the parse error for a non-number token.
	errExpectedIntQuoted = `expected int, found "x"`
	// errExpectedIntBare is the parse error for an unquoted non-number token.
	errExpectedIntBare = `expected int, found x`
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

	got, err := json.Codec.Encode(internal.BttvEmotesDBSchema, internal.BttvEmotesDBData)
	require.NoError(t, err)
	require.Equal(t, bttvEmotesDBJson, got)
}

func TestBttvDBDecode(t *testing.T) {
	t.Parallel()

	got, err := json.Codec.Decode(internal.BttvEmotesDBSchema, bttvEmotesDBJson)
	require.NoError(t, err)
	require.Equal(t, internal.BttvEmotesDBData, got)
}

func TestAllTypesEncodeDecode(t *testing.T) {
	t.Parallel()

	schema := allTypesSchema
	data := allTypesValue
	enc, err := json.Codec.Encode(schema, data)
	require.NoError(t, err)
	dec, err := json.Codec.Decode(schema, enc)
	require.NoError(t, err)
	require.Equal(t, dec, data)
}

// unsupportedSchema drives the unsupported-schema error path.
type unsupportedSchema struct{}

func (unsupportedSchema) Kind() sixty.Kind { return sixty.BitKind }

func TestDecodeUnsupportedSchema(t *testing.T) {
	t.Parallel()

	_, err := json.Codec.Decode(unsupportedSchema{}, "true")
	require.EqualError(t, err, "error at position 0: unsupported schema type json_test.unsupportedSchema")
}

func TestDecodeEmptyInput(t *testing.T) {
	t.Parallel()

	// eof while expecting '[' triggers both the eof and the unshift-at-start errors
	_, err := json.Codec.Decode(sixty.List{Schema: sixty.Bit{}}, "")
	require.ErrorContains(t, err, "unexpected eof")
	require.ErrorContains(t, err, "trying to unshift at start")
}

func TestDecodeTrailingGarbage(t *testing.T) {
	t.Parallel()

	_, err := json.Codec.Decode(sixty.Int{}, "1 2")
	require.ErrorContains(t, err, "expected eof, found 2")
}

func TestDecodePrimitiveErrors(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		schema sixty.Schema
		data   string
		reason string
	}{
		{sixty.String{}, "42", "expected string, found 42"},
		{sixty.Bit{}, `"yes"`, `expected bit, found "\"yes\""`},
		{sixty.Int{}, `"x"`, errExpectedIntQuoted},
		{sixty.Uint{}, "-5", "expected uint, found -5"},
		{sixty.Rune{}, "12", "expected rune, found 12"},
		{sixty.Rune{}, `"ab"`, `expected single character rune, found "ab"`},
		{sixty.Enum{Values: []string{"a", "b"}}, `"c"`, `unknown enum value "\"c\""`},
		{sixty.EnumInt{}, `"x"`, `expected enum int, found "x"`},
	}
	for _, test := range testCases {
		_, err := json.Codec.Decode(test.schema, test.data)
		require.ErrorContains(t, err, test.reason, "schema %T data %s", test.schema, test.data)
	}
}

func TestDecodeProductErrors(t *testing.T) {
	t.Parallel()

	schema := sixty.Product{{"a", sixty.Int{}}, {"b", sixty.Int{}}}

	testCases := []struct {
		data   string
		reason string
	}{
		{"1", `expected "{", found "1"`},
		{"{}", "expected 2 fields, found none"},
		{"{a:1,b:2}", "expected field name, found a"},
		{`{"a":1,"c":2}`, `expected field "b", found "c"`},
		{`{"a":1 "b":2}`, `expected ",", found "\"b\""`},
		{`{"a":1,"b":2`, `expected '}', found ""`},
		{`{"a":"x","b":2}`, errExpectedIntQuoted},
	}
	for _, test := range testCases {
		_, err := json.Codec.Decode(schema, test.data)
		require.ErrorContains(t, err, test.reason, "data %s", test.data)
	}
}

func TestDecodeCoproductErrors(t *testing.T) {
	t.Parallel()

	schema := sixty.Coproduct{{"a", sixty.Int{}}, {"b", sixty.String{}}}

	testCases := []struct {
		data   string
		reason string
	}{
		{"1", `expected "{", found "1"`},
		{"{}", `expected coproduct variant key, found }`},
		{"{x:1}", `expected coproduct variant key, found x`},
		{`{"x":1}`, `expected coproduct variant index, found "x"`},
		{`{"2":1}`, "found variant 2 while coproduct has only 2 variants"},
		{`{"0":5`, `expected '}', found ""`},
		{`{"0":"x"}`, errExpectedIntQuoted},
	}
	for _, test := range testCases {
		_, err := json.Codec.Decode(schema, test.data)
		require.ErrorContains(t, err, test.reason, "data %s", test.data)
	}
}

func TestDecodeListErrors(t *testing.T) {
	t.Parallel()

	schema := sixty.List{Schema: sixty.Int{}}

	testCases := []struct {
		data   string
		reason string
	}{
		{"[1 2]", `expected ',' or ']', found "2"`},
		{"[x]", errExpectedIntBare},
	}
	for _, test := range testCases {
		_, err := json.Codec.Decode(schema, test.data)
		require.ErrorContains(t, err, test.reason, "data %s", test.data)
	}

	// trailing comma is tolerated
	got, err := json.Codec.Decode(schema, "[1,]")
	require.NoError(t, err)
	require.Equal(t, sixty.ListValue{sixty.IntValue(1)}, got)
}

func TestDecodeArrayErrors(t *testing.T) {
	t.Parallel()

	schema := sixty.Array{Schema: sixty.Int{}, Count: 2}

	testCases := []struct {
		data   string
		reason string
	}{
		{"[1]", `expected ",", found "]"`},
		{"[1,2", `expected ']', found ""`},
		{"[x,x]", errExpectedIntBare},
	}
	for _, test := range testCases {
		_, err := json.Codec.Decode(schema, test.data)
		require.ErrorContains(t, err, test.reason, "data %s", test.data)
	}
}

func TestDecodeOptionalError(t *testing.T) {
	t.Parallel()

	_, err := json.Codec.Decode(sixty.Optional{Schema: sixty.Int{}}, "x")
	require.ErrorContains(t, err, errExpectedIntBare)
}

func TestDecodeEmptyProductSchema(t *testing.T) {
	t.Parallel()

	got, err := json.Codec.Decode(sixty.Product{}, "{}")
	require.NoError(t, err)
	require.Equal(t, sixty.ProductValue{}, got)
}

func TestDecodeArrayMalformedOpen(t *testing.T) {
	t.Parallel()

	_, err := json.Codec.Decode(sixty.Array{Schema: sixty.Int{}, Count: 1}, "1")
	require.ErrorContains(t, err, `expected "["`)
}
