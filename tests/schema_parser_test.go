package sixty_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	sixty "github.com/rprtr258/sixty/pkg"
)

var bttvEmotesDB = sixty.List{Type: sixty.Product{
	sixty.String,
	sixty.List{Type: sixty.String},
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

func exampleBttvEmotesDB() {
	tsv := sixty.NewPrettyPrinter()
	bttvEmotesDB.Visit(&tsv)
	fmt.Println(tsv.String())

	fmt.Println(sixty.SizeVisitor(bttvEmotesDB))

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
	fmt.Println(je.String())
}

func TestBttvDBSchemaParse(t *testing.T) {
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

// schema, err := sixty.ParseSchema(`[](0..255)`)
// if err != nil {
// 	log.Fatalln(err.Error())
// }
