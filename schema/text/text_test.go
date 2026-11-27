package text_test

import (
	_ "embed"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/rprtr258/sixty"
	"github.com/rprtr258/sixty/internal"
	"github.com/rprtr258/sixty/schema/text"
)

func assertSchema(t *testing.T, schemaString string, expected sixty.Schema) {
	t.Helper()
	got, err := text.ParseSchema(schemaString)
	require.NoError(t, err)
	require.Equal(t, expected, got)
}

var bttvEmotesDBSchemaString = `[](emote = :string,ids = [][8]rune)`

func TestDecode_BttvDB(t *testing.T) {
	t.Parallel()

	assertSchema(t, bttvEmotesDBSchemaString, internal.BttvEmotesDBSchema)
}

func TestEncode_BttvDB(t *testing.T) {
	t.Parallel()

	require.Equal(t, bttvEmotesDBSchemaString, text.Encode(internal.BttvEmotesDBSchema))
}

var (
	VKInfoSchema = sixty.Coproduct{
		{"group", sixty.Product{
			{"id", sixty.Uint{}},
			{"alias", sixty.String{}},
			{"name", sixty.String{}},
		}},
		{"user", sixty.Product{
			{"id", sixty.Uint{}},
			{"first_name", sixty.String{}},
			{"last_name", sixty.String{}},
		}},
	}
	//go:embed vk_info_schema.txt
	VKInfoSchemaString string
)

func TestDecode_VKInfo(t *testing.T) {
	t.Parallel()

	assertSchema(t, VKInfoSchemaString, VKInfoSchema)
}

func TestEncode_VKInfo(t *testing.T) {
	t.Parallel()

	require.Equal(t, VKInfoSchemaString, text.Encode(VKInfoSchema))
}

func TestDecode_EnumInt(t *testing.T) {
	t.Parallel()

	assertSchema(t, `[][0..255]`, sixty.List{Schema: sixty.EnumInt{
		Min: 0,
		Max: 255,
	}})
}

func TestDecode_Everything(t *testing.T) {
	t.Parallel()

	assertSchema(
		t,
		`[][3][str=:string,prod=(ei=[-2..2],i=int,u=uint,b=bit,e={a,b,c})]`,
		sixty.List{Schema: sixty.Array{
			Count: 3,
			Schema: sixty.Coproduct{
				{"str", sixty.String{}},
				{"prod", sixty.Product{
					{"ei", sixty.EnumInt{Min: -2, Max: 2}},
					{"i", sixty.Int{}},
					{"u", sixty.Uint{}},
					{"b", sixty.Bit{}},
					{"e", sixty.Enum{Values: []string{"a", "b", "c"}}},
				}},
			},
		}},
	)
}

func TestDecode_NoEndComma(t *testing.T) {
	t.Parallel()

	assertSchema(t, `(a=int,b=int)`, sixty.Product{{"a", sixty.Int{}}, {"b", sixty.Int{}}})
}

func TestDecode_EndComma(t *testing.T) {
	t.Parallel()

	assertSchema(t, `(a=int,b=int,)`, sixty.Product{{"a", sixty.Int{}}, {"b", sixty.Int{}}})
}

func TestDecode_Alias(t *testing.T) {
	t.Parallel()

	assertSchema(t, `a = int
(first = a, second = a)`, sixty.Product{{"first", sixty.Int{}}, {"second", sixty.Int{}}})
}

func TestDecode_Spaces(t *testing.T) {
	t.Parallel()

	assertSchema(t, `(
  a = int,
  b = int,
)`, sixty.Product{{"a", sixty.Int{}}, {"b", sixty.Int{}}})
}

func TestDecode_AliasesOfAliases(t *testing.T) {
	t.Parallel()

	twoInts := sixty.Product{{"l", sixty.Int{}}, {"r", sixty.Int{}}}
	fourInts := sixty.Product{{"l", twoInts}, {"r", twoInts}}
	assertSchema(t, `
x0 = (l = int, r = int)
x1 = (l = x0, r = x0)
(first = x1, second = x1)`, sixty.Product{{"first", fourInts}, {"second", fourInts}})
}

func TestParseErrorPositionCountsSpaces(t *testing.T) {
	t.Parallel()

	// "b" ends at byte offset 9: "(a=int" plus two spaces
	_, err := text.ParseSchema("(a=int  b=int)")
	require.ErrorContains(t, err, "error at position 9:")
}

func TestParseErrorEofPosition(t *testing.T) {
	t.Parallel()

	_, err := text.ParseSchema("(a=int")
	require.ErrorContains(t, err, "error at position 6:")
}

func TestParseUnknownAlias(t *testing.T) {
	t.Parallel()

	_, err := text.ParseSchema("(a = undefined_alias)")
	require.ErrorContains(t, err, `unknown alias "undefined_alias"`)
}

func TestPrettyEncode_BttvDB(t *testing.T) {
	t.Parallel()

	require.Equal(t, `[](
  emote = :string,
  ids = [][8]rune,
)`, text.PrettyEncode(internal.BttvEmotesDBSchema))
}

func TestAlgebraicEncode_BttvDB(t *testing.T) {
	t.Parallel()

	require.Equal(t, `[](:string * []8 * rune)`, text.AlgebraicEncode(internal.BttvEmotesDBSchema))
}

func TestAlgebraicEncode_VKInfo(t *testing.T) {
	t.Parallel()

	require.Equal(t, `(uint * :string * :string) + (uint * :string * :string)`, text.AlgebraicEncode(VKInfoSchema))
}

func TestParseErrorMessage(t *testing.T) {
	t.Parallel()

	require.Equal(t, "error at position 3: oops", text.ParseError{Pos: 3, Reason: "oops"}.Error())
}

func TestParseEmpty(t *testing.T) {
	t.Parallel()

	_, err := text.ParseSchema("")
	require.ErrorContains(t, err, "unexpected eof")
	require.ErrorContains(t, err, "trying to shift when index is 0")
}

func TestParseAliasError(t *testing.T) {
	t.Parallel()

	got, err := text.ParseSchema("x = ")
	require.Error(t, err)
	require.Nil(t, got)
}

func TestParseUnknownName(t *testing.T) {
	t.Parallel()

	// "5" is not a primitive and has no "=", so it is treated as an alias lookup
	got, err := text.ParseSchema("5")
	require.Error(t, err)
	require.Nil(t, got)
}

func TestParseErrors(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		schemaString string
		reason       string
	}{
		{"(5=int)", "5 is incorrect identifier"},
		{"(a 5)", `expected "="`},
		{"(a=int b=int)", `expected ',' but found "b"`},
		{"{5}", "5 is incorrect enum value"},
		{"[0..a]", "error while parsing int"},
		{"[a..b]", `expected "="`},
		{"(a=int]", `expected ',' but found "]"`},
	}
	for _, test := range testCases {
		_, err := text.ParseSchema(test.schemaString)
		require.ErrorContains(t, err, test.reason, "schema %q", test.schemaString)
	}
}

func TestMustParsePanics(t *testing.T) {
	t.Parallel()

	require.Panics(t, func() { text.MustParse("5") })
}

func TestMustParse(t *testing.T) {
	t.Parallel()

	require.Equal(t, sixty.Int{}, text.MustParse("int"))
}

var allTypesSchemaText = sixty.Product{
	{"b", sixty.Bit{}},
	{"i", sixty.Int{}},
	{"u", sixty.Uint{}},
	{"r", sixty.Rune{}},
	{"e", sixty.Enum{Values: []string{"a", "b"}}},
	{"ei", sixty.EnumInt{Min: 0, Max: 5}},
	{"arr", sixty.Array{Schema: sixty.Bit{}, Count: 3}},
	{"opt", sixty.Optional{Schema: sixty.Int{}}},
	{"s", sixty.String{}},
	{"l", sixty.List{Schema: sixty.Int{}}},
	{"p", sixty.Product{{"x", sixty.Int{}}, {"y", sixty.Int{}}}},
	{"c", sixty.Coproduct{{"x", sixty.Int{}}, {"y", sixty.String{}}}},
}

func TestEncodeAllTypes(t *testing.T) {
	t.Parallel()

	// compact form: fields separated by commas, no indentation
	want := `(b = bit,i = int,u = uint,r = rune,e = {a, b},ei = [0..5],` +
		`arr = [3]bit,opt = ?int,s = :string,l = []int,` +
		`p = (x = int,y = int),c = [x = int,y = :string])`
	require.Equal(t, want, text.Encode(allTypesSchemaText))
}

func TestPrettyEncodeAllTypes(t *testing.T) {
	t.Parallel()

	require.Equal(t, `(
  b = bit,
  i = int,
  u = uint,
  r = rune,
  e = {a, b},
  ei = [0..5],
  arr = [3]bit,
  opt = ?int,
  s = :string,
  l = []int,
  p = (
    x = int,
    y = int,
),
  c = [
    x = int,
    y = :string,
],
)`, text.PrettyEncode(allTypesSchemaText))
}

func TestEncodeRoundTrips(t *testing.T) {
	t.Parallel()

	assertSchema(t, text.Encode(allTypesSchemaText), allTypesSchemaText)
}

func TestAlgebraicEncodeAllTypes(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		schema sixty.Schema
		want   string
	}{
		{sixty.String{}, ":string"},
		{sixty.Bit{}, "2"},
		{sixty.Int{}, "int"},
		{sixty.Uint{}, "uint"},
		{sixty.Rune{}, "rune"},
		{sixty.Enum{Values: []string{"a", "b", "c"}}, "3"},
		{sixty.EnumInt{Min: 0, Max: 10}, "10"},
		{sixty.Array{Schema: sixty.Bit{}, Count: 3}, "3 * 2"},
		{sixty.Optional{Schema: sixty.Bit{}}, "1 + 2"},
		{sixty.List{Schema: sixty.Int{}}, "[]int"},
		// nested products fold inside parens
		{sixty.Product{
			{"a", sixty.Product{{"x", sixty.Int{}}, {"y", sixty.Int{}}}},
			{"b", sixty.Int{}},
		}, "(int * int) * int"},
	}
	for _, test := range testCases {
		require.Equal(t, test.want, text.AlgebraicEncode(test.schema), "schema %T", test.schema)
	}
}
