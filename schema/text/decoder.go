package text

import (
	"fmt"
	"regexp"
	"strconv"

	"go.uber.org/multierr"

	"github.com/rprtr258/sixty"
)

var (
	TokenRegex      = regexp.MustCompile(`(=|\[|\]|\(|\)|\{|\}|,|\?|\d+|:(string)|[a-zA-Z-_][a-zA-Z0-9-_]*|\.{2})`)
	IdentifierRegex = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_-]*$`)
)

// ParseError - schema text parsing failure, Pos counts scanned bytes.
type ParseError struct {
	Pos    uint
	Reason string
}

func (e ParseError) Error() string {
	return fmt.Sprintf("error at position %d: %s", e.Pos, e.Reason)
}

type SchemaReader struct {
	aliases map[string]sixty.Schema
	tokens  [][]byte
	index   uint
	// TODO: fix not counting space chars
	pos   uint
	Error error
}

func (r *SchemaReader) appendErrorf(format string, args ...any) {
	multierr.AppendInto(&r.Error, ParseError{
		Pos:    r.pos,
		Reason: fmt.Sprintf(format, args...),
	})
}

// shift gets next token and moves reader one step forward.
func (r *SchemaReader) shift() (res string) {
	if r.index == uint(len(r.tokens)) {
		r.appendErrorf("unexpected eof")
		return
	}
	res = string(r.tokens[r.index])
	r.pos += uint(len(res))
	r.index++
	return
}

// unshift moves reader one step back.
func (r *SchemaReader) unshift() {
	if r.index == 0 {
		r.appendErrorf("trying to shift when index is 0")
		return
	}
	if r.index < uint(len(r.tokens)) {
		r.pos -= uint(len(r.tokens[r.index]))
	}
	r.index--
}

func (r *SchemaReader) skip(expected string) {
	token := r.shift()
	if token != expected {
		r.appendErrorf("expected %q, found %q", expected, token)
		r.unshift()
	}
}

func (r *SchemaReader) ParseCommaSeparatedSchemas(endToken string) (res []sixty.NamedSchema) {
	for {
		if r.Error != nil {
			return
		}
		fieldName := r.shift()
		if !IdentifierRegex.Match([]byte(fieldName)) {
			r.appendErrorf("%s is incorrect identifier", fieldName)
			return
		}
		r.skip("=")
		schema := r.ParseSchema()
		res = append(res, sixty.NamedSchema{
			Name:   fieldName,
			Schema: schema,
		})
		nextToken := r.shift()
		switch nextToken {
		case endToken:
			return
		case ",":
			maybeEnd := r.shift()
			if maybeEnd == endToken {
				return
			}
			r.unshift()
		default:
			r.appendErrorf("expected ',' but found %q", nextToken)
			r.unshift()
		}
	}
}

// TODO: skip spaces.
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
				r.appendErrorf("error while parsing int: %v", err)
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
		enumValue := r.shift()
		if !IdentifierRegex.Match([]byte(enumValue)) {
			r.appendErrorf("%s is incorrect enum value", enumValue)
			return nil
		}
		vals = append(vals, enumValue)
		if r.shift() == "," {
			goto iter
		}
		return sixty.Enum{Values: vals}
	case string(sixty.StringKind):
		return sixty.String{}
	default:
		name := token
		nextToken := r.shift()
		if nextToken != "=" {
			r.unshift()
			return r.aliases[name]
		}
		r.aliases[name] = r.ParseSchema()
		if r.Error != nil {
			return nil
		}
		return r.ParseSchema()
	}
}

func Tokenize(s string) SchemaReader {
	return SchemaReader{
		tokens:  TokenRegex.FindAll([]byte(s), -1),
		aliases: make(map[string]sixty.Schema),
		index:   0,
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
