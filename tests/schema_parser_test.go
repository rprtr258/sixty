package sixty_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	sixty "github.com/rprtr258/sixty/pkg"
)

var bttvEmotesDB = sixty.List{Type: sixty.Product{
	sixty.String,                   // TODO: name as "emote"
	sixty.List{Type: sixty.String}, // TODO: name as "ids"
}}

// TODO: use to encode string and other primitive types
func stringToString(s string) sixty.TypeValue {
	var l sixty.ListValue
	for _, c := range s {
		l = append(l, sixty.RuneValue(c))
	}
	return sixty.AliasValue{
		Name:  "string",
		Value: l,
	}
}

func TestBttvDBJsonEncoder(t *testing.T) {
	je := &sixty.JsonEncoder{}
	data := sixty.ListValue{sixty.ProductValue{
		stringToString("Kappa"),
		sixty.ListValue{
			stringToString("KappaOwner1"),
			stringToString("KappaOwner2"),
		},
	}, sixty.ProductValue{
		stringToString("PogChamp"),
		sixty.ListValue{
			stringToString("PogChampOwner1"),
			stringToString("PogChampOwner2"),
		},
	}}
	sixty.VisitValue(bttvEmotesDB, data, je)
	assert.Equal(t, `[{"0":"Kappa","1":["KappaOwner1","KappaOwner2"]},{"0":"PogChamp","1":["PogChampOwner1","PogChampOwner2"]}]`, je.String())
}

func TestBttvDBSize(t *testing.T) {
	assert.Equal(t, sixty.Inf, sixty.SizeVisitor(bttvEmotesDB))
}

func TestBttvDB(t *testing.T) {
	schema, err := sixty.ParseSchema(`[](std:string,[]std:string)`)
	assert.NoError(t, err)
	assert.Equal(t, bttvEmotesDB, schema)
}

func TestBttvDBPrettyPrinter(t *testing.T) {
	assert.Equal(t, `std:string = []rune

[](
  std:string,
  []std:string,
)`, sixty.PrettyPrintSchema(bttvEmotesDB))
}

func TestEnumInt(t *testing.T) {
	schema, err := sixty.ParseSchema(`[](0..255)`)
	assert.NoError(t, err)
	assert.Equal(t, sixty.List{Type: sixty.EnumInt{
		Min: 0,
		Max: 255,
	}}, schema)

}
