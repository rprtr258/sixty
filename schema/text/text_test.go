package text_test

import (
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

var VKInfoSchema = sixty.Coproduct{
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

func TestPrettyEncode_BttvDB(t *testing.T) {
	t.Parallel()

	require.Equal(t, `[](
  emote = :string,
  ids = [][8]rune,
)`, text.PrettyEncode(internal.BttvEmotesDBSchema))
}

func TestAlgebraicEncode_VKInfo(t *testing.T) {
	t.Parallel()

	require.Equal(t, `(uint * :string * :string) + (uint * :string * :string)`, text.AlgebraicEncode(VKInfoSchema))
}

func TestParseErrors(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		schemaString string
		reasons      []string
	}{
		{"", []string{"unexpected eof", "trying to shift when index is 0"}},
		// trailing "=" declares an alias without a value
		{"x = ", []string{`unknown alias ""`}},
		{"(5=int)", []string{"5 is incorrect identifier"}},
		{"(a 5)", []string{`expected "="`}},
		{"(a=int b=int)", []string{`expected ',' but found "b"`}},
		{"{5}", []string{"5 is incorrect enum value"}},
		{"[0..a]", []string{"error while parsing int"}},
		{"[a..b]", []string{`expected "="`}},
		{"(a=int]", []string{`expected ',' but found "]"`}},
	}
	for _, test := range testCases {
		_, err := text.ParseSchema(test.schemaString)
		require.Error(t, err, "schema %q", test.schemaString)
		for _, reason := range test.reasons {
			require.ErrorContains(t, err, reason, "schema %q", test.schemaString)
		}
	}
}

func TestMustParse(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		schemaString string
		want         sixty.Schema
		panics       bool
	}{
		{"int", sixty.Int{}, false},
		// "5" is not a valid schema, MustParse must panic
		{"5", nil, true},
	}
	for _, test := range testCases {
		if test.panics {
			require.Panics(t, func() { text.MustParse(test.schemaString) }, "schema %q", test.schemaString)
			continue
		}
		require.Equal(t, test.want, text.MustParse(test.schemaString), "schema %q", test.schemaString)
	}
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
