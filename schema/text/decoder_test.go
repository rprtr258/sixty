package text_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/rprtr258/sixty"
	"github.com/rprtr258/sixty/internal"
	"github.com/rprtr258/sixty/schema/text"
)

func assertSchema(t *testing.T, schemaString string, expected sixty.Schema) {
	got, err := text.ParseSchema(schemaString)
	assert.NoError(t, err)
	assert.Equal(t, expected, got)
}

func TestDecode_BttvDB(t *testing.T) {
	assertSchema(t, internal.BttvEmotesDBSchemaString, internal.BttvEmotesDBSchema)
}

func TestDecode_VKInfo(t *testing.T) {
	assertSchema(t, internal.VKInfoSchemaString, internal.VKInfoSchema)
}

func TestDecode_EnumInt(t *testing.T) {
	assertSchema(t, `[][0..255]`, sixty.List{Schema: sixty.EnumInt{
		Min: 0,
		Max: 255,
	}})
}

func TestDecode_Everything(t *testing.T) {
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
	assertSchema(t, `(a=int,b=int)`, sixty.Product{{"a", sixty.Int{}}, {"b", sixty.Int{}}})
}

func TestDecode_EndComma(t *testing.T) {
	assertSchema(t, `(a=int,b=int,)`, sixty.Product{{"a", sixty.Int{}}, {"b", sixty.Int{}}})
}

func TestDecode_Alias(t *testing.T) {
	assertSchema(t, `a = int
(first = a, second = a)`, sixty.Product{{"first", sixty.Int{}}, {"second", sixty.Int{}}})
}

func TestDecode_Spaces(t *testing.T) {
	assertSchema(t, `(
  a = int,
  b = int,
)`, sixty.Product{{"a", sixty.Int{}}, {"b", sixty.Int{}}})
}
