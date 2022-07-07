package main

import (
	"fmt"

	sixty "github.com/rprtr258/sixty/pkg"
)

func main() {
	String := sixty.Alias{
		Type: sixty.List{Type: sixty.Rune{}},
		Name: "string",
	}
	BttvEmotesDB := sixty.List{Type: sixty.Product([]sixty.Type{
		String,
		sixty.List{Type: String},
	})}

	tsv := sixty.NewPrettyPrinter()
	BttvEmotesDB.Visit(&tsv)
	fmt.Println(tsv.String())
	fmt.Println(sixty.SizeVisitor(BttvEmotesDB))
}
