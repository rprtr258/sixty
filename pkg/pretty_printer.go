package sixty

import (
	"fmt"
	"strings"

	"github.com/rprtr258/go-flow/fun"
	"github.com/rprtr258/go-flow/stream"
)

type prettyPrinter struct {
	builder strings.Builder
	ident   uint
	aliases map[string]string
}

func NewPrettyPrinter() prettyPrinter {
	return prettyPrinter{
		ident:   0,
		aliases: make(map[string]string),
		builder: strings.Builder{},
	}
}

func (visitor *prettyPrinter) writeString(str string) {
	visitor.builder.WriteString(strings.Repeat(" ", int(visitor.ident*2)))
	visitor.builder.WriteString(str)
}

func (visitor *prettyPrinter) VisitBit(Bit) {
	visitor.writeString("bit")
}

func (visitor *prettyPrinter) VisitInt(Int) {
	visitor.writeString("int")
}

func (visitor *prettyPrinter) VisitUint(Uint) {
	visitor.writeString("uint")
}

func (visitor *prettyPrinter) VisitRune(Rune) {
	visitor.writeString("rune")
}

func (visitor *prettyPrinter) VisitProduct(product Product) {
	visitor.writeString("(\n")
	visitor.ident++

	for _, field := range product {
		field.Visit(visitor)
		visitor.builder.WriteString(",\n")
	}

	visitor.ident--
	visitor.writeString(")")
}

func (visitor *prettyPrinter) VisitCoproduct(coproduct Coproduct) {
	visitor.writeString("{\n")
	visitor.ident++

	for _, field := range coproduct {
		field.Visit(visitor)
		visitor.builder.WriteString(",\n")
	}

	visitor.ident--
	visitor.writeString("}")
}

func (visitor *prettyPrinter) VisitList(list List) {
	visitor.writeString("[]")
	tmp := visitor.ident
	visitor.ident = 0
	list.Schema.Visit(visitor)
	visitor.ident = tmp
}

func (visitor *prettyPrinter) VisitArray(array Array) {
	visitor.writeString(fmt.Sprintf("[%d]", array.Count))
	tmp := visitor.ident
	visitor.ident = 0
	array.Schema.Visit(visitor)
	visitor.ident = tmp
}

func (visitor *prettyPrinter) VisitOptional(optional Optional) {
	visitor.writeString("?")
	tmp := visitor.ident
	visitor.ident = 0
	optional.Schema.Visit(visitor)
	visitor.ident = tmp
}

func (visitor *prettyPrinter) VisitEnum(enum Enum) {
	visitor.writeString(fmt.Sprintf("{%s}", strings.Join(enum.Values, ", ")))
}

func (visitor *prettyPrinter) VisitEnumInt(enumInt EnumInt) {
	visitor.writeString(fmt.Sprintf("(%d..%d)", enumInt.Min, enumInt.Max))
}

func (visitor *prettyPrinter) VisitAlias(alias Alias) {
	visitor.writeString(alias.Name)
	if _, ok := visitor.aliases[alias.Name]; ok {
		return
	}
	pp := prettyPrinter{aliases: visitor.aliases}
	alias.Schema.Visit(&pp)
	visitor.aliases[alias.Name] = pp.String()
}

func (visitor *prettyPrinter) String() string {
	s1 := stream.FromMap(visitor.aliases)
	s2 := stream.Map(s1, func(kv fun.Pair[string, string]) string {
		// TODO: don't print std:types
		return fmt.Sprintf("%s = %s\n\n", kv.Left, kv.Right)
	})
	s3 := stream.CollectToSlice(s2)
	return strings.Join(s3, "") + visitor.builder.String()
}

func PrettyPrintSchema(t Schema) string {
	tsv := NewPrettyPrinter()
	t.Visit(&tsv)
	return tsv.String()
}
