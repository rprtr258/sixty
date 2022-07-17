package text

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/rprtr258/sixty"
)

type algebraicEncoder struct {
	builder strings.Builder
	aliases map[string]string
}

func NewAlgebraicEncoder(builder strings.Builder) algebraicEncoder {
	return algebraicEncoder{
		builder: builder,
		aliases: make(map[string]string),
	}
}

func (e *algebraicEncoder) writeFold(ts []sixty.NamedSchema, op rune) {
	i := 0
iter:
	e.writeSubExpr(ts[i].Schema)
	if i < len(ts)-1 {
		i++
		e.builder.WriteRune(' ')
		e.builder.WriteRune(op)
		e.builder.WriteRune(' ')
		goto iter
	}
}

func (e *algebraicEncoder) writeSubExpr(t sixty.Schema) {
	switch tt := t.(type) {
	case sixty.String:
		e.builder.WriteString(string(sixty.StringKind))
	case sixty.Bit:
		e.builder.WriteString("2")
	case sixty.Int:
		e.builder.WriteString("int")
	case sixty.Uint:
		e.builder.WriteString("uint")
	case sixty.Rune:
		e.builder.WriteString("rune")
	case sixty.Enum:
		e.builder.WriteString(strconv.Itoa(len(tt.Values)))
	case sixty.EnumInt:
		e.builder.WriteString(strconv.FormatInt(tt.Max-tt.Min, 10))
	case sixty.List:
		e.builder.WriteString("[]")
		e.writeSubExpr(tt.Schema)
	case sixty.Array:
		e.builder.WriteString(strconv.FormatUint(uint64(tt.Count), 10))
		e.builder.WriteString(" * ")
		e.writeSubExpr(tt.Schema)
	case sixty.Optional:
		e.builder.WriteString("1 + ")
		e.writeSubExpr(tt.Schema)
	default:
		fmt.Printf("%T\n", t)
		e.builder.WriteRune('(')
		e.Write(tt)
		e.builder.WriteRune(')')
	}
}

func (e *algebraicEncoder) Write(t sixty.Schema) {
	switch tt := t.(type) {
	case sixty.String, sixty.Bit, sixty.Int, sixty.Uint, sixty.Rune, sixty.Enum, sixty.EnumInt:
		e.writeSubExpr(tt)
	case sixty.Product:
		e.writeFold(tt, '*')
	case sixty.Coproduct:
		e.writeFold(tt, '+')
	case sixty.List, sixty.Array, sixty.Optional:
		e.writeSubExpr(tt)
	default:
		panic("something unimplemented")
	}
}

func (e *algebraicEncoder) String() string {
	var res strings.Builder
	for k, v := range e.aliases {
		res.WriteString(k)
		res.WriteString(" = ")
		res.WriteString(v)
	}
	res.WriteString(e.builder.String())
	return res.String()
}

func AlgebraicEncode(t sixty.Schema) string {
	tsv := NewAlgebraicEncoder(strings.Builder{})
	tsv.Write(t)
	return tsv.String()
}
