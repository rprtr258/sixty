package text_test

import (
	_ "embed"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/rprtr258/sixty"
	"github.com/rprtr258/sixty/internal"
	"github.com/rprtr258/sixty/schema/text"
)

func assertSchema(t *testing.T, schemaString string, expected sixty.Schema) {
	t.Helper()
	got, err := text.ParseSchema(schemaString)
	assert.NoError(t, err)
	assert.Equal(t, expected, got)
}

var bttvEmotesDBSchemaString = `[](emote = :string,ids = []:string)`

func TestDecode_BttvDB(t *testing.T) {
	t.Parallel()

	assertSchema(t, bttvEmotesDBSchemaString, internal.BttvEmotesDBSchema)
}

func TestEncode_BttvDB(t *testing.T) {
	t.Parallel()

	assert.Equal(t, bttvEmotesDBSchemaString, text.Encode(internal.BttvEmotesDBSchema))
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

	assert.Equal(t, VKInfoSchemaString, text.Encode(VKInfoSchema))
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

// TODO: store aliases in aliases cache
// func TestDecode_Alias2(t *testing.T) {
// 	assertSchema(t, `
// x0 = (l = int, r = int)
// x1 = (l = x0, r = x0)
// x2 = (l = x1, r = x1)
// x3 = (l = x2, r = x2)
// x4 = (l = x3, r = x3)
// x5 = (l = x4, r = x4)
// x6 = (l = x5, r = x5)
// x7 = (l = x6, r = x6)
// x8 = (l = x7, r = x7)
// (first = x8, second = x8)`, sixty.Product{{"first", sixty.Int{}}, {"second", sixty.Int{}}})
// }

func TestPrettyEncode_BttvDB(t *testing.T) {
	t.Parallel()

	assert.Equal(t, `[](
  emote = :string,
  ids = []:string,
)`, text.PrettyEncode(internal.BttvEmotesDBSchema))
}

func TestAlgebraicEncode_BttvDB(t *testing.T) {
	t.Parallel()

	assert.Equal(t, `[](:string * []:string)`, text.AlgebraicEncode(internal.BttvEmotesDBSchema))
}

func TestAlgebraicEncode_VKInfo(t *testing.T) {
	t.Parallel()

	assert.Equal(t, `(uint * :string * :string) + (uint * :string * :string)`, text.AlgebraicEncode(VKInfoSchema))
}
