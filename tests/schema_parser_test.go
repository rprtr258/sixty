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

func stoS(s string) sixty.AliasValue {
	var l sixty.ListValue
	for _, c := range s {
		l = append(l, sixty.RuneValue(c))
	}
	return sixty.AliasValue{
		Name:  "string",
		Value: l,
	}
}

func TestBttvDBSize(t *testing.T) {
	assert.Equal(t, sixty.Inf, sixty.SizeVisitor(bttvEmotesDB))
}

func TestBttvDBJsonEncoder(t *testing.T) {
	je := &sixty.JsonEncoder{}
	data := sixty.ListValue{sixty.ProductValue{
		stoS("Kappa"),
		sixty.ListValue{
			stoS("KappaOwner1"),
			stoS("KappaOwner2"),
		},
	}, sixty.ProductValue{
		stoS("PogChamp"),
		sixty.ListValue{
			stoS("PogChampOwner1"),
			stoS("PogChampOwner2"),
		},
	}}
	sixty.VisitValue(bttvEmotesDB, data, je)
	assert.Equal(t, `[{"0":"Kappa","1":["KappaOwner1","KappaOwner2"]},{"0":"PogChamp","1":["PogChampOwner1","PogChampOwner2"]}]`, je.String())
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
