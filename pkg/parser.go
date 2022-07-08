package sixty

import (
	"log"
	"regexp"
	"strconv"
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

func (r *SchemaReader) ParseCommaSeparatedSchemas(endToken string) (res []Schema) {
	for {
		res = append(res, r.ParseSchema())
		nextToken := r.Next()
		if nextToken == endToken {
			return
		} // TODO: else assert ","
	}
}

func (r *SchemaReader) ParseSchema() Schema {
	token := r.Next()
	switch token {
	case "bit":
		return Bit{}
	case "int":
		return Int{}
	case "uint":
		return Uint{}
	case "rune":
		return Rune{}
	case "?":
		return Optional{Schema: r.ParseSchema()}
	case "[":
		nextToken := r.Next()
		if nextToken == "]" {
			return List{Schema: r.ParseSchema()}
		} else {
			arrayLength, err := strconv.Atoi(nextToken)
			if err != nil {
				log.Fatalln(err.Error())
			}
			r.Next() // TODO: assert "]"
			return Array{
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
			return EnumInt{
				Min: min,
				Max: max,
			}
		}
		return Product(r.ParseCommaSeparatedSchemas(")"))
	case "{":
		return Coproduct(r.ParseCommaSeparatedSchemas("}"))
	// TODO: enum
	// TODO: alias
	case string(StringKind):
		return String{}
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

func ParseSchema(s string) (Schema, error) {
	reader := Tokenize(s)
	return reader.ParseSchema(), nil
}

func MustParseSchema(s string) Schema {
	res, err := ParseSchema(s)
	if err != nil {
		panic(err)
	}
	return res
}
