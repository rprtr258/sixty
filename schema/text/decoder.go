package text

import (
	"log"
	"regexp"
	"strconv"

	"github.com/rprtr258/sixty"
)

var TokenRegex = regexp.MustCompile(`(\[|\]|\(|\)|\{|\}|,|\?|\d+|std:(string)|[a-zA-Z-_][a-zA-Z0-9-_]*|\.{2})`)

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

func (r *SchemaReader) ParseCommaSeparatedSchemas(endToken string) (res []sixty.Schema) {
	for {
		res = append(res, r.ParseSchema())
		nextToken := r.Next()
		if nextToken == endToken {
			return
		} // TODO: else assert ","
	}
}

func (r *SchemaReader) ParseSchema() sixty.Schema {
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
	case "?":
		return sixty.Optional{Schema: r.ParseSchema()}
	case "[":
		nextToken := r.Next()
		if nextToken == "]" {
			return sixty.List{Schema: r.ParseSchema()}
		} else {
			arrayLength, err := strconv.Atoi(nextToken)
			if err != nil {
				log.Fatalln(err.Error())
			}
			r.Next() // TODO: assert "]"
			return sixty.Array{
				Schema: r.ParseSchema(),
				Count:  uint(arrayLength),
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
		return sixty.Product(r.ParseCommaSeparatedSchemas(")"))
	case "{":
		return sixty.Coproduct(r.ParseCommaSeparatedSchemas("}"))
	// TODO: enum
	// TODO: alias
	case string(sixty.StringKind):
		return sixty.String{}
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

func ParseSchema(s string) (sixty.Schema, error) {
	reader := Tokenize(s)
	return reader.ParseSchema(), nil
}

func MustParseSchema(s string) sixty.Schema {
	res, err := ParseSchema(s)
	if err != nil {
		panic(err)
	}
	return res
}
