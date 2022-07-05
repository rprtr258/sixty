package main

import (
	"fmt"
	"strings"

	"github.com/rprtr258/go-flow/fun"
	"github.com/rprtr258/go-flow/stream"
)

// TODO: own Log2 on uint

type TypeVisitor interface {
	VisitBit(Bit)
	VisitInt(Int)
	VisitUint(Uint)
	VisitRune(Rune)
	VisitProduct(Product)
	VisitCoproduct(Coproduct)
	VisitList(List)
	VisitArray(Array)
	VisitOptional(Optional)
	VisitEnum(Enum)
	VisitEnumInt(EnumInt)
	VisitAlias(Alias)
}

type Type interface {
	Visit(TypeVisitor)
}

type Bit struct{}

func (bit Bit) Visit(visitor TypeVisitor) { visitor.VisitBit(bit) }

type Int struct{}

func (intt Int) Visit(visitor TypeVisitor) { visitor.VisitInt(intt) }

type Uint struct{}

func (uintt Uint) Visit(visitor TypeVisitor) { visitor.VisitUint(uintt) }

type Rune struct{}

func (runee Rune) Visit(visitor TypeVisitor) { visitor.VisitRune(runee) }

type Product struct {
	Fields []Type
}

func (product Product) Visit(visitor TypeVisitor) { visitor.VisitProduct(product) }

type Coproduct struct {
	Variants []Type
}

func (coproduct Coproduct) Visit(visitor TypeVisitor) { visitor.VisitCoproduct(coproduct) }

type List struct {
	Type
}

func (list List) Visit(visitor TypeVisitor) { visitor.VisitList(list) }

type Array struct {
	Type
	Count uint
}

func (array Array) Visit(visitor TypeVisitor) { visitor.VisitArray(array) }

type Optional struct {
	Type
}

func (optional Optional) Visit(visitor TypeVisitor) { visitor.VisitOptional(optional) }

type Enum struct {
	Values []string
}

func (enum Enum) Visit(visitor TypeVisitor) { visitor.VisitEnum(enum) }

type EnumInt struct {
	Min int
	Max int
}

func (enumInt EnumInt) Visit(visitor TypeVisitor) { visitor.VisitEnumInt(enumInt) }

type Alias struct {
	Type
	Name string
}

func (alias Alias) Visit(visitor TypeVisitor) { visitor.VisitAlias(alias) }

// TODO: move to visitor
// func (kind Kind) Size() fun.Option[uint] {
// 	switch kind {
// 	case Bit:
// 		return fun.Some[uint](1)
// 	case Int:
// 		return fun.Option[uint]{}
// 	case Uint:
// 		return fun.Option[uint]{}
// 	case Rune:
// 		return fun.Option[uint]{}
// 	default:
// 		panic(fmt.Sprintf("Unknown kind: %d", kind))
// 	}
// }

// func (product Product) Size() fun.Option[uint] {
// 	if len(product.Fields) == 0 {
// 		panic("Product with zero fields")
// 	}
// 	return stream.Reduce(
// 		fun.Some[uint](1),
// 		func(acc fun.Option[uint], t Type) fun.Option[uint] {
// 			tSize := t.Size()
// 			if acc.IsNone() || tSize.IsNone() {
// 				return fun.None[uint]()
// 			}
// 			return fun.Some(acc.Unwrap() * tSize.Unwrap())
// 		},
// 		stream.FromSlice(product.Fields),
// 	)
// }

// func (coproduct Coproduct) Size() fun.Option[uint] {
// 	if len(coproduct.Variants) == 0 {
// 		panic("Coproduct with zero fields")
// 	}
// 	maxSize := stream.Reduce(
// 		fun.Some[uint](1),
// 		func(acc fun.Option[uint], t Type) fun.Option[uint] {
// 			tSize := t.Size()
// 			if acc.IsNone() || tSize.IsNone() {
// 				return fun.None[uint]()
// 			}
// 			if acc.Unwrap() > tSize.Unwrap() {
// 				return fun.Some(acc.Unwrap())
// 			} else {
// 				return fun.Some(tSize.Unwrap())
// 			}
// 		},
// 		stream.FromSlice(coproduct.Variants),
// 	)
// 	return fun.Map(maxSize, func(x uint) uint {
// 		return x + uint(math.Ceil(math.Log2(float64(len(coproduct.Variants)))))
// 	})
// }

// func (List) Size() fun.Option[uint] { return fun.None[uint]() }

// func (array Array) Size() fun.Option[uint] {
// 	return fun.Map(
// 		array.Type.Size(),
// 		func(x uint) uint { return x * array.Count },
// 	)
// }

// func (optional Optional) Size() fun.Option[uint] {
// 	return fun.Map(
// 		optional.Type.Size(),
// 		func(x uint) uint { return x + 1 },
// 	)
// }

// func (enum Enum) Size() fun.Option[uint] {
// 	return fun.Some(uint(math.Ceil(math.Log2(float64(len(enum.Values))))))
// }

// func (enumInt EnumInt) Size() fun.Option[uint] {
// 	if enumInt.Max < enumInt.Min {
// 		panic(fmt.Sprintf("Incorrect EnumInt: min=%d max=%d", enumInt.Min, enumInt.Max))
// 	}
// 	return fun.Some(uint(math.Ceil(math.Log2(float64(enumInt.Max - enumInt.Min)))))
// }

// type TypeWithMetadata struct {
// 	Type
// 	KVMeta    map[string]string
// 	FlagsMeta map[string]struct{}
// }

// TODO: Metadata String

// TODO: stupid printer
// TODO: algebraic printer
type PrettyPrinter struct {
	strings.Builder
	ident   uint
	aliases map[string]string
}

func (visitor *PrettyPrinter) WriteString(str string) {
	visitor.Builder.WriteString(strings.Repeat(" ", int(visitor.ident*2)))
	visitor.Builder.WriteString(str)
}

func (visitor *PrettyPrinter) VisitBit(Bit) {
	visitor.WriteString("bit")
}

func (visitor *PrettyPrinter) VisitInt(Int) {
	visitor.WriteString("int")
}

func (visitor *PrettyPrinter) VisitUint(Uint) {
	visitor.WriteString("uint")
}

func (visitor *PrettyPrinter) VisitRune(Rune) {
	visitor.WriteString("rune")
}

func (visitor *PrettyPrinter) VisitProduct(product Product) {
	visitor.WriteString("(\n")
	visitor.ident++

	n := len(product.Fields)
	for i, field := range product.Fields {
		field.Visit(visitor)
		if i < n-1 {
			visitor.Builder.WriteString(",\n")
		}
	}

	visitor.ident--
	visitor.WriteString("\n)")
}

func (visitor *PrettyPrinter) VisitCoproduct(coproduct Coproduct) {
	visitor.WriteString("[\n")
	visitor.ident++

	n := len(coproduct.Variants)
	for i, field := range coproduct.Variants {
		field.Visit(visitor)
		if i < n-1 {
			visitor.Builder.WriteString(",")
		}
	}

	visitor.ident--
	visitor.WriteString("\n]")
}

func (visitor *PrettyPrinter) VisitList(list List) {
	visitor.WriteString("[]")
	tmp := visitor.ident
	visitor.ident = 0
	list.Type.Visit(visitor)
	visitor.ident = tmp
}

func (visitor *PrettyPrinter) VisitArray(array Array) {
	visitor.WriteString(fmt.Sprintf("[%d]", array.Count))
	tmp := visitor.ident
	visitor.ident = 0
	array.Type.Visit(visitor)
	visitor.ident = tmp
}

func (visitor *PrettyPrinter) VisitOptional(optional Optional) {
	visitor.WriteString("?")
	tmp := visitor.ident
	visitor.ident = 0
	optional.Type.Visit(visitor)
	visitor.ident = tmp
}

func (visitor *PrettyPrinter) VisitEnum(enum Enum) {
	visitor.WriteString(fmt.Sprintf("{%s}", strings.Join(enum.Values, ", ")))
}

func (visitor *PrettyPrinter) VisitEnumInt(enumInt EnumInt) {
	visitor.WriteString(fmt.Sprintf("[%d..%d]", enumInt.Min, enumInt.Max))
}

func (visitor *PrettyPrinter) VisitAlias(alias Alias) {
	visitor.WriteString(alias.Name)
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

func main() {
	String := Alias{List{Rune{}}, "string"}
	BttvEmotesDB := List{Product{[]Type{
		String,
		List{String},
	}}}

	tsv := PrettyPrinter{
		aliases: make(map[string]string),
	}
	BttvEmotesDB.Visit(&tsv)
	fmt.Println(tsv.String())
}
