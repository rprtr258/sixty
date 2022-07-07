package main

import (
	"fmt"

	sixty "github.com/rprtr258/sixty/pkg"
)

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

func main() {
	BttvEmotesDB := sixty.List{Type: sixty.Product{
		sixty.String,
		sixty.List{Type: sixty.String},
	}}

	tsv := sixty.NewPrettyPrinter()
	BttvEmotesDB.Visit(&tsv)
	fmt.Println(tsv.String())

	fmt.Println(sixty.SizeVisitor(BttvEmotesDB))

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
	sixty.VisitValue(BttvEmotesDB, data, je)
	fmt.Println(je.String())
}
