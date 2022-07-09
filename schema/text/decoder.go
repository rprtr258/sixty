package text

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/rprtr258/sixty"
	"go.uber.org/multierr"
)

var TokenRegex = regexp.MustCompile(`(\[|\]|\(|\)|\{|\}|,|\?|\d+|std:(string)|[a-zA-Z-_][a-zA-Z0-9-_]*|\.{2})`)

type SchemaReader struct {
	tokens [][]byte
	index  uint
	pos    uint
	Error  error
}

// TODO: unexpected eof
// TODO: get back instead?
func (r *SchemaReader) Peek() string {
	return string(r.tokens[r.index])
}

func (r *SchemaReader) Next() (res string) {
	res = string(r.tokens[r.index])
	r.index++
	r.pos += uint(len(res))
	return
}

func (r *SchemaReader) appendError(format string, args ...any) {
	args2 := make([]any, len(args)+1)
	args2[0] = r.pos
	copy(args2[1:], args)
	multierr.AppendInto(&r.Error, fmt.Errorf("at position %d: "+format, args2...))
}

func (r *SchemaReader) skip(expected string) {
	token := r.Peek()
	if token != expected {
		r.appendError("expected %s, found %s", token, expected)
	} else {
		r.Next()
	}
}

// TODO: trailing comma
func (r *SchemaReader) ParseCommaSeparatedSchemas(endToken string) (res []sixty.Schema) {
	for {
		res = append(res, r.ParseSchema())
		nextToken := r.Next()
		if nextToken == endToken {
			return
		} else if nextToken != "," {
			r.appendError("expected ',' but found %q", nextToken)
		}
	}
}

// TODo: skip spaces
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
		}
		arrayLength, err := strconv.Atoi(nextToken)
		if err != nil {
			r.appendError("error while parsing int: %w", err)
		}
		r.skip("]")
		return sixty.Array{
			Schema: r.ParseSchema(),
			Count:  uint(arrayLength),
		}
	case "(":
		nextToken := r.Peek()
		min, err := strconv.ParseInt(nextToken, 10, 64)
		if err != nil {
			return sixty.Product(r.ParseCommaSeparatedSchemas(")"))
		}
		r.Next()
		r.skip("..")
		max, err := strconv.ParseInt(r.Next(), 10, 64)
		if err != nil {
			r.appendError("error while parsing int: %w", err)
		}
		return sixty.EnumInt{
			Min: min,
			Max: max,
		}
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

func MustParse(s string) sixty.Schema {
	res, err := ParseSchema(s)
	if err != nil {
		panic(err)
	}
	return res
}
