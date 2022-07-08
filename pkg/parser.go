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

func (r *SchemaReader) ParseCommaSeparatedTypes(endToken string) (res []Type) {
	for {
		res = append(res, r.ParseType())
		nextToken := r.Next()
		if nextToken == endToken {
			return
		} // TODO: else assert ","
	}
}

func (r *SchemaReader) ParseType() Type {
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
		return Optional{Type: r.ParseType()}
	case "[":
		nextToken := r.Next()
		if nextToken == "]" {
			return List{Type: r.ParseType()}
		} else {
			arrayLength, err := strconv.Atoi(nextToken)
			if err != nil {
				log.Fatalln(err.Error())
			}
			r.Next() // TODO: assert "]"
			return Array{
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
			return EnumInt{
				Min: min,
				Max: max,
			}
		}
		return Product(r.ParseCommaSeparatedTypes(")"))
	case "{":
		return Coproduct(r.ParseCommaSeparatedTypes("}"))
	// TODO: enum
	// TODO: alias
	case "std:string":
		return String
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

func ParseSchema(s string) (Type, error) {
	reader := Tokenize(s)
	return reader.ParseType(), nil
}

func MustParseSchema(s string) Type {
	res, err := ParseSchema(s)
	if err != nil {
		panic(err)
	}
	return res
}
