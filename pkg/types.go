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

// TODO: rename to schema?
type Type interface {
	Visit(TypeVisitor)
	// TODO: get value?
}

// TODO: is needed?
type TypeValue interface {
	isTypeValue()
	// TODO: get kind?
}

type Bit struct{}

func (bit Bit) Visit(visitor TypeVisitor) { visitor.VisitBit(bit) }

type BitValue bool

func (BitValue) isTypeValue() {}

type Int struct{}

func (intt Int) Visit(visitor TypeVisitor) { visitor.VisitInt(intt) }

type IntValue int

func (IntValue) isTypeValue() {}

type Uint struct{}

func (uintt Uint) Visit(visitor TypeVisitor) { visitor.VisitUint(uintt) }

type UintValue uint

func (UintValue) isTypeValue() {}

type Rune struct{}

func (runee Rune) Visit(visitor TypeVisitor) { visitor.VisitRune(runee) }

type RuneValue rune

func (RuneValue) isTypeValue() {}

type Enum struct {
	Values []string
}

func (enum Enum) Visit(visitor TypeVisitor) { visitor.VisitEnum(enum) }

type EnumValue uint

func (EnumValue) isTypeValue() {}

type EnumInt struct {
	Min int
	Max int
}

func (enumInt EnumInt) Visit(visitor TypeVisitor) { visitor.VisitEnumInt(enumInt) }

type EnumIntValue int

func (EnumIntValue) isTypeValue() {}

type Product []Type

func (product Product) Visit(visitor TypeVisitor) { visitor.VisitProduct(product) }

type ProductValue []TypeValue

func (ProductValue) isTypeValue() {}

type Coproduct []Type

func (coproduct Coproduct) Visit(visitor TypeVisitor) { visitor.VisitCoproduct(coproduct) }

type CoproductValue struct {
	TypeValue
	Variant uint
}

func (CoproductValue) isTypeValue() {}

type List struct {
	Type
}

func (list List) Visit(visitor TypeVisitor) { visitor.VisitList(list) }

type ListValue []TypeValue

func (ListValue) isTypeValue() {}

type Array struct {
	Type
	Count uint
}

func (array Array) Visit(visitor TypeVisitor) { visitor.VisitArray(array) }

type ArrayValue []TypeValue

func (ArrayValue) isTypeValue() {}

type Optional struct {
	Type
}

func (optional Optional) Visit(visitor TypeVisitor) { visitor.VisitOptional(optional) }

type OptionalValue struct {
	TypeValue
	IsValid bool
}

func (OptionalValue) isTypeValue() {}

// TODO: alias in alias???
type Alias struct {
	Type
	Name string
}

func (alias Alias) Visit(visitor TypeVisitor) { visitor.VisitAlias(alias) }

type AliasValue struct {
	Name  string
	Value TypeValue
}

func (AliasValue) isTypeValue() {}

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
