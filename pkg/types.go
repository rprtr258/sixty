package sixty

import (
	"fmt"
	"math"
)

var String = Alias{
	Type: List{Type: Rune{}},
	Name: "std:string",
}

type AccumVisitor[A any] interface {
	VisitBit(Bit) A
	VisitInt(Int) A
	VisitUint(Uint) A
	VisitRune(Rune) A
	VisitEnum(Enum) A
	VisitEnumInt(EnumInt) A
	VisitProduct([]A) A
	VisitCoproduct([]A) A
	VisitList(A) A
	VisitArray(A, uint) A
	VisitOptional(A) A
	VisitAlias(A, string) A
}

func DFSFactory[A any](visitor AccumVisitor[A]) func(Type) A {
	var res func(Type) A
	res = func(t Type) A {
		switch tt := t.(type) {
		case Bit:
			return visitor.VisitBit(tt)
		case Int:
			return visitor.VisitInt(tt)
		case Uint:
			return visitor.VisitUint(tt)
		case Rune:
			return visitor.VisitRune(tt)
		case Enum:
			return visitor.VisitEnum(tt)
		case EnumInt:
			return visitor.VisitEnumInt(tt)
		case Product:
			var product []A
			for _, x := range tt {
				product = append(product, res(x))
			}
			return visitor.VisitProduct(product)
		case Coproduct:
			var coproduct []A
			for _, x := range tt {
				coproduct = append(coproduct, res(x))
			}
			return visitor.VisitCoproduct(coproduct)
		case List:
			return visitor.VisitList(res(tt.Type))
		case Array:
			return visitor.VisitArray(res(tt.Type), tt.Count)
		case Optional:
			return visitor.VisitOptional(res(tt.Type))
		case Alias:
			return visitor.VisitAlias(res(tt.Type), tt.Name)
		default:
			panic(fmt.Sprintf("unknown type: %T", tt))
		}
	}
	return res
}

type TypeVisitor interface {
	VisitBit(Bit)
	VisitInt(Int)
	VisitUint(Uint)
	VisitRune(Rune)
	VisitEnum(Enum)
	VisitEnumInt(EnumInt)
	VisitProduct(Product)
	VisitCoproduct(Coproduct)
	VisitList(List)
	VisitArray(Array)
	VisitOptional(Optional)
	VisitAlias(Alias)
}

// TODO: make only one kind of visitor
type EncoderVisitor interface {
	VisitBit(Bit, BitValue)
	VisitInt(Int, IntValue)
	VisitUint(Uint, UintValue)
	VisitRune(Rune, RuneValue)
	VisitEnum(Enum, EnumValue)
	VisitEnumInt(EnumInt, EnumIntValue)
	VisitProduct(Product, ProductValue)
	VisitCoproduct(Coproduct, CoproductValue)
	VisitList(List, ListValue)
	VisitArray(Array, ArrayValue)
	VisitOptional(Optional, OptionalValue)
	VisitAlias(Alias, AliasValue)
}

type Kind uint

const (
	BitKind Kind = iota
	IntKind
	UintKind
	RuneKind
	EnumKind
	EnumIntKind
	ProductKind
	CoproductKind
	ListKind
	ArrayKind
	OptionalKind
	AliasKind
)

// TODO: rename to schema?
type Type interface {
	Visit(TypeVisitor)
	Kind() Kind
	// TODO: get value?
}

// TODO: is needed?
type TypeValue interface {
	isTypeValue()
	Kind() Kind
}

type Bit struct{}

func (bit Bit) Visit(visitor TypeVisitor) { visitor.VisitBit(bit) }
func (Bit) Kind() Kind                    { return BitKind }

type BitValue bool

func (BitValue) isTypeValue() {}
func (BitValue) Kind() Kind   { return BitKind }

type Int struct{}

func (intt Int) Visit(visitor TypeVisitor) { visitor.VisitInt(intt) }
func (Int) Kind() Kind                     { return IntKind }

type IntValue int

func (IntValue) isTypeValue() {}
func (IntValue) Kind() Kind   { return IntKind }

type Uint struct{}

func (uintt Uint) Visit(visitor TypeVisitor) { visitor.VisitUint(uintt) }
func (Uint) Kind() Kind                      { return UintKind }

type UintValue uint

func (UintValue) isTypeValue() {}
func (UintValue) Kind() Kind   { return UintKind }

type Rune struct{}

func (runee Rune) Visit(visitor TypeVisitor) { visitor.VisitRune(runee) }
func (Rune) Kind() Kind                      { return RuneKind }

type RuneValue rune

func (RuneValue) isTypeValue() {}
func (RuneValue) Kind() Kind   { return RuneKind }

type Enum struct {
	Values []string
}

func (enum Enum) Visit(visitor TypeVisitor) { visitor.VisitEnum(enum) }
func (Enum) Kind() Kind                     { return EnumKind }

type EnumValue uint

func (EnumValue) isTypeValue() {}
func (EnumValue) Kind() Kind   { return EnumKind }

type EnumInt struct {
	Min int
	Max int
}

func (enumInt EnumInt) Visit(visitor TypeVisitor) { visitor.VisitEnumInt(enumInt) }
func (EnumInt) Kind() Kind                        { return EnumIntKind }

// TODO: int64
type EnumIntValue int

func (EnumIntValue) isTypeValue() {}
func (EnumIntValue) Kind() Kind   { return EnumIntKind }

type Product []Type

func (product Product) Visit(visitor TypeVisitor) { visitor.VisitProduct(product) }
func (Product) Kind() Kind                        { return ProductKind }

type ProductValue []TypeValue

func (ProductValue) isTypeValue() {}
func (ProductValue) Kind() Kind   { return ProductKind }

type Coproduct []Type

func (coproduct Coproduct) Visit(visitor TypeVisitor) { visitor.VisitCoproduct(coproduct) }
func (Coproduct) Kind() Kind                          { return CoproductKind }

type CoproductValue struct {
	TypeValue
	Variant uint
}

func (CoproductValue) isTypeValue() {}
func (CoproductValue) Kind() Kind   { return CoproductKind }

type List struct {
	Type
}

func (list List) Visit(visitor TypeVisitor) { visitor.VisitList(list) }
func (List) Kind() Kind                     { return ListKind }

type ListValue []TypeValue

func (ListValue) isTypeValue() {}
func (ListValue) Kind() Kind   { return ListKind }

type Array struct {
	Type
	Count uint
}

func (array Array) Visit(visitor TypeVisitor) { visitor.VisitArray(array) }
func (Array) Kind() Kind                      { return ArrayKind }

type ArrayValue []TypeValue

func (ArrayValue) isTypeValue() {}
func (ArrayValue) Kind() Kind   { return ArrayKind }

type Optional struct {
	Type
}

func (optional Optional) Visit(visitor TypeVisitor) { visitor.VisitOptional(optional) }
func (Optional) Kind() Kind                         { return OptionalKind }

type OptionalValue struct {
	TypeValue
	IsValid bool
}

func (OptionalValue) isTypeValue() {}
func (OptionalValue) Kind() Kind   { return OptionalKind }

// TODO: alias in alias???
type Alias struct {
	Type
	Name string
}

func (alias Alias) Visit(visitor TypeVisitor) { visitor.VisitAlias(alias) }
func (Alias) Kind() Kind                      { return AliasKind }

type AliasValue struct {
	Name  string
	Value TypeValue
}

func (AliasValue) isTypeValue() {}
func (AliasValue) Kind() Kind   { return AliasKind }

// type TypeWithMetadata struct {
// 	Type
// 	KVMeta    map[string]string
// 	FlagsMeta map[string]struct{}
// }

// TODO: Metadata String

// TODO: stupid printer
// TODO: algebraic printer

// TODO: own Log2 on int and uint
func Log2(x int) uint {
	return uint(math.Ceil(math.Log2(float64(x))))
}
