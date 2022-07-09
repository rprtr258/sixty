package text

import (
	"errors"
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

// shift gets next token and moves reader one step forward
func (r *SchemaReader) shift() (res string) {
	if r.index == uint(len(r.tokens)) {
		r.appendError("unexpected eof")
		return
	}
	res = string(r.tokens[r.index])
	r.pos += uint(len(res))
	r.index++
	return
}

// unshift moves reader one step back
func (r *SchemaReader) unshift() {
	if r.index == 0 {
		r.appendError("trying to shift when index is 0")
		return
	}
	if r.index < uint(len(r.tokens)) {
		r.pos -= uint(len(r.tokens[r.index]))
	}
	r.index--
}

func (r *SchemaReader) appendError(format string, args ...any) {
	multierr.AppendInto(&r.Error, errors.New(
		fmt.Sprintf("error at position %d: ", r.pos)+fmt.Sprintf(format, args...)),
	)
}

func (r *SchemaReader) skip(expected string) {
	token := r.shift()
	if token != expected {
		r.appendError("expected %s, found %s", token, expected)
		r.unshift()
	}
}

// TODO: trailing comma
func (r *SchemaReader) ParseCommaSeparatedSchemas(endToken string) (res []sixty.Schema) {
	for {
		if r.Error != nil {
			return
		}
		res = append(res, r.ParseSchema())
		nextToken := r.shift()
		if nextToken == endToken {
			return
		} else if nextToken != "," {
			r.appendError("expected ',' but found %q", nextToken)
			r.unshift()
		}
	}
}

// TODo: skip spaces
func (r *SchemaReader) ParseSchema() sixty.Schema {
	token := r.shift()
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
		nextToken := r.shift()
		if nextToken == "]" {
			return sixty.List{Schema: r.ParseSchema()}
		} else if num, err := strconv.ParseInt(nextToken, 10, 64); err == nil {
			if r.shift() != ".." {
				r.unshift()
				r.skip("]")
				return sixty.Array{
					Schema: r.ParseSchema(),
					Count:  uint(num),
				}
			}
			max, err := strconv.ParseInt(r.shift(), 10, 64)
			if err != nil {
				r.appendError("error while parsing int: %v", err)
				return nil
			}
			r.skip("]")
			return sixty.EnumInt{
				Min: num,
				Max: max,
			}
		} else {
			r.unshift()
			return sixty.Coproduct(r.ParseCommaSeparatedSchemas("]"))
		}
	case "(":
		return sixty.Product(r.ParseCommaSeparatedSchemas(")"))
	case "{":
		var vals []string
	iter:
		if r.Error != nil {
			return nil
		}
		// TODO: assert identifier
		vals = append(vals, r.shift())
		if r.shift() == "," {
			goto iter
		}
		return sixty.Enum{Values: vals}
	// TODO: alias
	case string(sixty.StringKind):
		return sixty.String{}
	default:
		r.appendError("unexpected token %q", token)
		return nil
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
	return reader.ParseSchema(), reader.Error
}

func MustParse(s string) sixty.Schema {
	res, err := ParseSchema(s)
	if err != nil {
		panic(err)
	}
	return res
}
