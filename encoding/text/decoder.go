package text

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"

	"go.uber.org/multierr"

	"github.com/rprtr258/sixty"
)

// TokenRegex matches one token of the text data encoding: a string with doubled
// quote escaping, a rune with backslash escaping, a decimal number, an identifier
// or a punctuation character.
var TokenRegex = regexp.MustCompile(`"(?:[^"]|"")*"|'(?:\\.|[^'])'|-?\d+|[a-zA-Z_][a-zA-Z0-9_-]*|[{}\[\],:]`)

// ParseError - data text parsing failure, Pos is byte offset of the offending token.
type ParseError struct {
	Pos    uint
	Reason string
}

func (e ParseError) Error() string {
	return fmt.Sprintf("error at position %d: %s", e.Pos, e.Reason)
}

type decoder struct {
	source string
	tokens []string
	// starts - byte offset of every token inside source
	starts []int
	// index - next token to read
	index int
	// pos - byte offset of the current token
	pos   uint
	Error error
}

func (d *decoder) appendErrorf(format string, args ...any) {
	multierr.AppendInto(&d.Error, ParseError{
		Pos:    d.pos,
		Reason: fmt.Sprintf(format, args...),
	})
}

// shift returns the next token and moves the reader forward.
func (d *decoder) shift() string {
	if d.index == len(d.tokens) {
		d.pos = uint(len(d.source))
		d.appendErrorf("unexpected eof")
		return ""
	}
	d.pos = uint(d.starts[d.index])
	token := d.tokens[d.index]
	d.index++
	return token
}

// unshift moves the reader one token back.
func (d *decoder) unshift() {
	if d.index == 0 {
		d.appendErrorf("trying to unshift at start")
		return
	}
	d.index--
	d.pos = uint(d.starts[d.index])
}

func (d *decoder) expect(expected string) {
	token := d.shift()
	if token != expected {
		d.appendErrorf("expected %q, found %q", expected, token)
		d.unshift()
	}
}

// expectEnd consumes the closing token of a collection, tolerating one trailing
// comma as written by the pretty encoder.
func (d *decoder) expectEnd(end string) {
	token := d.shift()
	switch token {
	case end:
		return
	case ",":
		d.expect(end)
	default:
		d.appendErrorf("expected %q, found %q", end, token)
		d.unshift()
	}
}

func (d *decoder) Read(t sixty.Schema) sixty.SchemaValue {
	if d.Error != nil {
		return nil
	}
	switch tt := t.(type) {
	case sixty.String:
		return d.readString()
	case sixty.Bit:
		return d.readBit()
	case sixty.Int:
		return d.readInt()
	case sixty.Uint:
		return d.readUint()
	case sixty.Rune:
		return d.readRune()
	case sixty.Enum:
		return d.readEnum(tt)
	case sixty.EnumInt:
		return d.readEnumInt()
	case sixty.Product:
		return d.readProduct(tt)
	case sixty.Coproduct:
		return d.readCoproduct(tt)
	case sixty.List:
		return d.readList(tt.Schema)
	case sixty.Array:
		return d.readArray(tt)
	case sixty.Optional:
		return d.readOptional(tt)
	default:
		d.appendErrorf("unsupported schema type %T", t)
		return nil
	}
}

func (d *decoder) readString() sixty.StringValue {
	token := d.shift()
	if len(token) < 2 || token[0] != '"' || token[len(token)-1] != '"' {
		d.appendErrorf("expected string, found %s", token)
		return ""
	}
	return sixty.StringValue(strings.ReplaceAll(token[1:len(token)-1], `""`, `"`))
}

func (d *decoder) readBit() sixty.BitValue {
	switch token := d.shift(); token {
	case "true":
		return true
	case "false":
		return false
	default:
		d.appendErrorf("expected bit, found %q", token)
		return false
	}
}

func (d *decoder) readInt() sixty.IntValue {
	token := d.shift()
	value, err := strconv.ParseInt(token, 10, 64)
	if err != nil {
		d.appendErrorf("expected int, found %s", token)
		return 0
	}
	return sixty.IntValue(value)
}

func (d *decoder) readUint() sixty.UintValue {
	token := d.shift()
	value, err := strconv.ParseUint(token, 10, 64)
	if err != nil {
		d.appendErrorf("expected uint, found %s", token)
		return 0
	}
	return sixty.UintValue(value)
}

func (d *decoder) readRune() sixty.RuneValue {
	token := d.shift()
	if len(token) < 3 || token[0] != '\'' || token[len(token)-1] != '\'' {
		d.appendErrorf("expected rune, found %s", token)
		return 0
	}
	inner := token[1 : len(token)-1]
	// encoder escapes only the quote character, so a lone backslash is a backslash rune
	if len(inner) > 1 && inner[0] == '\\' {
		inner = inner[1:]
	}
	r, size := utf8.DecodeRuneInString(inner)
	if size == 0 || (r == utf8.RuneError && size == 1) || size != len(inner) {
		d.appendErrorf("invalid rune %s", token)
		return 0
	}
	return sixty.RuneValue(r)
}

func (d *decoder) readEnum(tt sixty.Enum) sixty.EnumValue {
	token := d.shift()
	for i, name := range tt.Values {
		if name == token {
			return sixty.EnumValue(i)
		}
	}
	d.appendErrorf("unknown enum value %q", token)
	return 0
}

func (d *decoder) readEnumInt() sixty.EnumIntValue {
	token := d.shift()
	value, err := strconv.ParseInt(token, 10, 64)
	if err != nil {
		d.appendErrorf("expected enum int, found %s", token)
		return 0
	}
	return sixty.EnumIntValue(value)
}

func (d *decoder) readProduct(tt sixty.Product) sixty.ProductValue {
	d.expect("{")
	if d.Error != nil {
		return nil
	}
	res := make(sixty.ProductValue, 0, len(tt))
	for i, field := range tt {
		res = append(res, d.Read(field.Schema))
		if d.Error != nil {
			return nil
		}
		if i < len(tt)-1 {
			d.expect(",")
		}
	}
	d.expectEnd("}")
	if d.Error != nil {
		return nil
	}
	return res
}

func (d *decoder) readCoproduct(tt sixty.Coproduct) sixty.CoproductValue {
	d.expect("{")
	if d.Error != nil {
		return sixty.CoproductValue{}
	}
	variantToken := d.shift()
	variant, err := strconv.ParseUint(variantToken, 10, 64)
	if err != nil {
		d.appendErrorf("expected coproduct variant index, found %q", variantToken)
		return sixty.CoproductValue{}
	}
	if variant >= uint64(len(tt)) {
		d.appendErrorf("found variant %d while coproduct has only %d variants", variant, len(tt))
		return sixty.CoproductValue{}
	}
	d.expect(":")
	value := d.Read(tt[variant].Schema)
	d.expectEnd("}")
	if d.Error != nil {
		return sixty.CoproductValue{}
	}
	return sixty.CoproductValue{
		SchemaValue: value,
		Variant:     uint(variant),
	}
}

func (d *decoder) readList(schema sixty.Schema) sixty.ListValue {
	d.expect("[")
	if d.Error != nil {
		return nil
	}
	res := make(sixty.ListValue, 0)
	token := d.shift()
	if token == "]" {
		return res
	}
	d.unshift()
	for {
		res = append(res, d.Read(schema))
		if d.Error != nil {
			return nil
		}
		switch token := d.shift(); token {
		case "]":
			return res
		case ",":
			if d.shift() == "]" { // tolerate trailing comma
				return res
			}
			d.unshift()
		default:
			d.appendErrorf("expected ',' or ']', found %q", token)
			return nil
		}
	}
}

func (d *decoder) readArray(tt sixty.Array) sixty.ArrayValue {
	d.expect("[")
	if d.Error != nil {
		return nil
	}
	res := make(sixty.ArrayValue, 0, tt.Count)
	for i := uint(0); i < tt.Count; i++ {
		res = append(res, d.Read(tt.Schema))
		if d.Error != nil {
			return nil
		}
		if i < tt.Count-1 {
			d.expect(",")
		}
	}
	d.expectEnd("]")
	if d.Error != nil {
		return nil
	}
	return res
}

func (d *decoder) readOptional(tt sixty.Optional) sixty.OptionalValue {
	token := d.shift()
	if token == "null" {
		return sixty.OptionalValue{IsValid: false}
	}
	d.unshift()
	value := d.Read(tt.Schema)
	if d.Error != nil {
		return sixty.OptionalValue{}
	}
	return sixty.OptionalValue{
		SchemaValue: value,
		IsValid:     true,
	}
}

func newDecoder(data string) *decoder {
	locations := TokenRegex.FindAllStringIndex(data, -1)
	tokens := make([]string, len(locations))
	starts := make([]int, len(locations))
	for i, location := range locations {
		tokens[i] = data[location[0]:location[1]]
		starts[i] = location[0]
	}
	return &decoder{source: data, tokens: tokens, starts: starts}
}

// Decode parses data written by text.Encode.
func decode(t sixty.Schema, data string) (sixty.SchemaValue, error) {
	decoder := newDecoder(data)
	res := decoder.Read(t)
	if decoder.Error == nil && decoder.index != len(decoder.tokens) {
		decoder.appendErrorf("expected eof, found %s", decoder.tokens[decoder.index])
	}
	return res, decoder.Error
}
