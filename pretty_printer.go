package main

import (
	"fmt"
	"strings"

	"github.com/rprtr258/go-flow/fun"
	"github.com/rprtr258/go-flow/stream"
)

type PrettyPrinter struct {
	strings.Builder
	ident   uint
	aliases map[string]string
}

func (visitor *PrettyPrinter) writeString(str string) {
	visitor.Builder.WriteString(strings.Repeat(" ", int(visitor.ident*2)))
	visitor.Builder.WriteString(str)
}

func (visitor *PrettyPrinter) VisitBit(Bit) {
	visitor.writeString("bit")
}

func (visitor *PrettyPrinter) VisitInt(Int) {
	visitor.writeString("int")
}

func (visitor *PrettyPrinter) VisitUint(Uint) {
	visitor.writeString("uint")
}

func (visitor *PrettyPrinter) VisitRune(Rune) {
	visitor.writeString("rune")
}

func (visitor *PrettyPrinter) VisitProduct(product Product) {
	visitor.writeString("(\n")
	visitor.ident++

	n := len(product)
	for i, field := range product {
		field.Visit(visitor)
		if i < n-1 {
			visitor.Builder.WriteString(",\n")
		}
	}

	visitor.ident--
	visitor.writeString("\n)")
}

func (visitor *PrettyPrinter) VisitCoproduct(coproduct Coproduct) {
	visitor.writeString("[\n")
	visitor.ident++

	n := len(coproduct)
	for i, field := range coproduct {
		field.Visit(visitor)
		if i < n-1 {
			visitor.Builder.WriteString(",")
		}
	}

	visitor.ident--
	visitor.writeString("\n]")
}

func (visitor *PrettyPrinter) VisitList(list List) {
	visitor.writeString("[]")
	tmp := visitor.ident
	visitor.ident = 0
	list.Type.Visit(visitor)
	visitor.ident = tmp
}

func (visitor *PrettyPrinter) VisitArray(array Array) {
	visitor.writeString(fmt.Sprintf("[%d]", array.Count))
	tmp := visitor.ident
	visitor.ident = 0
	array.Type.Visit(visitor)
	visitor.ident = tmp
}

func (visitor *PrettyPrinter) VisitOptional(optional Optional) {
	visitor.writeString("?")
	tmp := visitor.ident
	visitor.ident = 0
	optional.Type.Visit(visitor)
	visitor.ident = tmp
}

func (visitor *PrettyPrinter) VisitEnum(enum Enum) {
	visitor.writeString(fmt.Sprintf("{%s}", strings.Join(enum.Values, ", ")))
}

func (visitor *PrettyPrinter) VisitEnumInt(enumInt EnumInt) {
	visitor.writeString(fmt.Sprintf("[%d..%d]", enumInt.Min, enumInt.Max))
}

func (visitor *PrettyPrinter) VisitAlias(alias Alias) {
	visitor.writeString(alias.Name)
	if _, ok := visitor.aliases[alias.Name]; ok {
		return
	}
	pp := PrettyPrinter{aliases: visitor.aliases}
	alias.Type.Visit(&pp)
	visitor.aliases[alias.Name] = pp.String()
}

func (visitor *PrettyPrinter) String() string {
	s1 := stream.FromMap(visitor.aliases)
	s2 := stream.Map(s1, func(kv fun.Pair[string, string]) string {
		return fmt.Sprintf("%s = %s\n\n", kv.Left, kv.Right)
	})
	s3 := stream.CollectToSlice(s2)
	return strings.Join(s3, "") + visitor.Builder.String()
}
