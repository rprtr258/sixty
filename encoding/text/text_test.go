package text_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/rprtr258/sixty"
	"github.com/rprtr258/sixty/encoding/text"
	"github.com/rprtr258/sixty/internal"
)

const (
	// errExpectedIntQuoted is the parse error for a non-number token.
	errExpectedIntQuoted = `expected int, found "x"`
	// errExpectedIntBare is the parse error for an unquoted non-number token.
	errExpectedIntBare = `expected int, found x`
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

	got, err := text.Codec.Encode(internal.BttvEmotesDBSchema, internal.BttvEmotesDBData)
	require.NoError(t, err)
	require.Equal(t, `[{"Kappa",[['k','a','p','p','a','_','0','1'],['k','a','p','p','a','_','0','2']]},{"PogChamp",[['p','o','g','c','h','_','0','1'],['p','o','g','c','h','_','0','2']]}]`, got)
}

func TestBttvDBDecode(t *testing.T) {
	t.Parallel()

	bttvEmotesDBText := `[{"Kappa",[['k','a','p','p','a','_','0','1'],['k','a','p','p','a','_','0','2']]},
		{"PogChamp",[['p','o','g','c','h','_','0','1'],['p','o','g','c','h','_','0','2']]}]`
	got, err := text.Codec.Decode(internal.BttvEmotesDBSchema, bttvEmotesDBText)
	require.NoError(t, err)
	require.Equal(t, internal.BttvEmotesDBData, got)
}

func TestBttvDBEncodeDecode(t *testing.T) {
	t.Parallel()

	schema := internal.BttvEmotesDBSchema
	data := internal.BttvEmotesDBData
	enc, err := text.Codec.Encode(schema, data)
	require.NoError(t, err)
	dec, err := text.Codec.Decode(schema, enc)
	require.NoError(t, err)
	require.Equal(t, dec, data)
}

func TestAllTypesEncodeDecode(t *testing.T) {
	t.Parallel()

	schema := allTypesSchema
	data := allTypesValue
	enc, err := text.Codec.Encode(schema, data)
	require.NoError(t, err)
	dec, err := text.Codec.Decode(schema, enc)
	require.NoError(t, err)
	require.Equal(t, dec, data)
}

func TestDecodeGarbage(t *testing.T) {
	t.Parallel()

	_, err := text.Codec.Decode(internal.BttvEmotesDBSchema, "garbage")
	require.Error(t, err)
}

func TestEncodeRuneQuoteEscape(t *testing.T) {
	t.Parallel()

	got, err := text.Codec.Encode(sixty.Rune{}, sixty.RuneValue('\''))
	require.NoError(t, err)
	require.Equal(t, `'\''`, got)
}

func TestParseErrorMessage(t *testing.T) {
	t.Parallel()

	require.Equal(t, "error at position 3: oops", text.ParseError{Pos: 3, Reason: "oops"}.Error())
}

// unsupportedSchema drives the unsupported-schema error path.
type unsupportedSchema struct{}

func (unsupportedSchema) Kind() sixty.Kind { return sixty.BitKind }

func TestDecodeUnsupportedSchema(t *testing.T) {
	t.Parallel()

	_, err := text.Codec.Decode(unsupportedSchema{}, "true")
	require.EqualError(t, err, "error at position 0: unsupported schema type text_test.unsupportedSchema")
}

func TestDecodeEmptyInput(t *testing.T) {
	t.Parallel()

	// eof while expecting '[' triggers both the eof and the unshift-at-start errors
	_, err := text.Codec.Decode(sixty.List{Schema: sixty.Bit{}}, "")
	require.ErrorContains(t, err, "unexpected eof")
	require.ErrorContains(t, err, "trying to unshift at start")
}

func TestDecodeTrailingGarbage(t *testing.T) {
	t.Parallel()

	_, err := text.Codec.Decode(sixty.Int{}, "1 2")
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
		{sixty.Bit{}, "5", `expected bit, found "5"`},
		{sixty.Int{}, `"x"`, errExpectedIntQuoted},
		{sixty.Uint{}, "-5", "expected uint, found -5"},
		{sixty.Rune{}, "12", "expected rune, found 12"},
		{sixty.Rune{}, "'\xff'", "invalid rune"},
		{sixty.Enum{Values: []string{"a", "b"}}, "c", `unknown enum value "c"`},
		{sixty.EnumInt{}, `"x"`, `expected enum int, found "x"`},
	}
	for _, test := range testCases {
		_, err := text.Codec.Decode(test.schema, test.data)
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
		{"{1 2}", `expected ",", found "2"`},
		{"{1,2]", `expected "}", found "]"`},
		{"{1,2,", `expected "}", found ""`},
		{"{x,2}", errExpectedIntBare},
	}
	for _, test := range testCases {
		_, err := text.Codec.Decode(schema, test.data)
		require.ErrorContains(t, err, test.reason, "data %s", test.data)
	}

	// trailing comma tolerated: pretty encoder writes it
	got, err := text.Codec.Decode(schema, "{1,2,}")
	require.NoError(t, err)
	require.Equal(t, sixty.ProductValue{sixty.IntValue(1), sixty.IntValue(2)}, got)
}

func TestDecodeCoproductErrors(t *testing.T) {
	t.Parallel()

	schema := sixty.Coproduct{{"a", sixty.Int{}}, {"b", sixty.String{}}}

	testCases := []struct {
		data   string
		reason string
	}{
		{"1", `expected "{", found "1"`},
		{"{x:1}", `expected coproduct variant index, found "x"`},
		{"{2:1}", "found variant 2 while coproduct has only 2 variants"},
		{"{0:", `expected "}", found ""`},
		{"{0:x}", errExpectedIntBare},
		{"{0:1]", `expected "}", found "]"`},
	}
	for _, test := range testCases {
		_, err := text.Codec.Decode(schema, test.data)
		require.ErrorContains(t, err, test.reason, "data %s", test.data)
	}

	got, err := text.Codec.Decode(schema, "{0:5,}")
	require.NoError(t, err)
	require.Equal(t, sixty.CoproductValue{SchemaValue: sixty.IntValue(5), Variant: 0}, got)
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
		_, err := text.Codec.Decode(schema, test.data)
		require.ErrorContains(t, err, test.reason, "data %s", test.data)
	}

	got, err := text.Codec.Decode(schema, "[1,]")
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
		{"[1,2,", `expected "]", found ""`},
		{"[x,x]", errExpectedIntBare},
	}
	for _, test := range testCases {
		_, err := text.Codec.Decode(schema, test.data)
		require.ErrorContains(t, err, test.reason, "data %s", test.data)
	}
}

func TestDecodeRuneEscapedQuote(t *testing.T) {
	t.Parallel()

	got, err := text.Codec.Decode(sixty.Rune{}, `'\''`)
	require.NoError(t, err)
	require.Equal(t, sixty.RuneValue('\''), got)
}

func TestDecodeArrayMalformedOpen(t *testing.T) {
	t.Parallel()

	_, err := text.Codec.Decode(sixty.Array{Schema: sixty.Int{}, Count: 1}, "1")
	require.ErrorContains(t, err, `expected "["`)
}

func TestDecodeOptionalError(t *testing.T) {
	t.Parallel()

	_, err := text.Codec.Decode(sixty.Optional{Schema: sixty.Int{}}, "x")
	require.ErrorContains(t, err, "expected int, found x")
}
