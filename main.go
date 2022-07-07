package main

import (
	"fmt"
	"log"
	"regexp"
	"strconv"

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

var TokenRegex = regexp.MustCompile(`(\[|\]|\(|\)|\{|\}|,|\d+|std:(string)|[a-zA-Z-_][a-zA-Z0-9-_]*|\.{2})`)

type SchemaReader struct {
	tokens [][]byte
	index  uint
}

// TODO: unexpected eof
func (r *SchemaReader) Peek() string {
	return string(r.tokens[r.index])
}

func (r *SchemaReader) Next() (res string) {
	res = string(r.tokens[r.index])
	r.index++
	return
}

func (r *SchemaReader) ParseCommaSeparatedTypes(endToken string) (res []sixty.Type) {
	for {
		res = append(res, r.ParseType())
		nextToken := r.Next()
		if nextToken == endToken {
			return
		} // TODO: else assert ","
	}
}

func (r *SchemaReader) ParseType() sixty.Type {
	token := r.Next()
	switch token {
	case "bit":
		return sixty.Bit{}
	case "int":
		return sixty.Int{}
	case "uint":
		return sixty.Uint{}
	case "rune":
		return sixty.Rune{}
	case "[":
		nextToken := r.Next()
		if nextToken == "]" {
			return sixty.List{Type: r.ParseType()}
		} else {
			arrayLength, err := strconv.Atoi(nextToken)
			if err != nil {
				log.Fatalln(err.Error())
			}
			r.Next() // TODO: assert "]"
			return sixty.Array{
				Type:  r.ParseType(),
				Count: uint(arrayLength),
			}
		}
	case "(":
		nextToken := r.Peek()
		if min, err := strconv.Atoi(nextToken); err == nil {
			r.Next()
			r.Next() // TODO: assert ..
			max, err := strconv.Atoi(r.Next())
			if err != nil {
				panic(err.Error())
			}
			return sixty.EnumInt{
				Min: min,
				Max: max,
			}
		}
		return sixty.Product(r.ParseCommaSeparatedTypes(")"))
	case "{":
		return sixty.Coproduct(r.ParseCommaSeparatedTypes("}"))
	// TODO: enum
	// TODO: alias
	case "std:string":
		return sixty.String
	default:
		panic(token)
	}
}

func Tokenize(s string) SchemaReader {
	return SchemaReader{
		tokens: TokenRegex.FindAll([]byte(s), -1),
		index:  0,
	}
}

func ParseSchema(s string) (sixty.Type, error) {
	reader := Tokenize(s)
	return reader.ParseType(), nil
}

func ExampleBttvEmotesDB() {
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

func main() {
	schema, err := ParseSchema(`[](std:string,[]std:string)`)
	if err != nil {
		log.Fatalln(err.Error())
	}

	tsv := sixty.NewPrettyPrinter()
	schema.Visit(&tsv)
	fmt.Println(tsv.String())
}
