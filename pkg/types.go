package sixty

import (
	"fmt"
	"math"
)

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
	Logical
)

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

func DFSFactory[A any](visitor AccumVisitor[A]) func(Schema) A {
	var res func(Schema) A
	res = func(t Schema) A {
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
			return visitor.VisitList(res(tt.Schema))
		case Array:
			return visitor.VisitArray(res(tt.Schema), tt.Count)
		case Optional:
			return visitor.VisitOptional(res(tt.Schema))
		case Alias:
			return visitor.VisitAlias(res(tt.Schema), tt.Name)
		default:
			panic(fmt.Sprintf("unknown type: %T", tt))
		}
	}
	return res
}

type SchemaVisitor interface {
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
	VisitString(String)
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

type Schema interface {
	Visit(SchemaVisitor)
	Kind() Kind
	// TODO: get value?
}

// TODO: is needed?
type SchemaValue interface {
	isSchemaValue()
	Kind() Kind
}

type Bit struct{}

func (bit Bit) Visit(visitor SchemaVisitor) { visitor.VisitBit(bit) }
func (Bit) Kind() Kind                      { return BitKind }

type BitValue bool

func (BitValue) isSchemaValue() {}
func (BitValue) Kind() Kind     { return BitKind }

type Int struct{}

func (intt Int) Visit(visitor SchemaVisitor) { visitor.VisitInt(intt) }
func (Int) Kind() Kind                       { return IntKind }

type IntValue int

func (IntValue) isSchemaValue() {}
func (IntValue) Kind() Kind     { return IntKind }

type Uint struct{}

func (uintt Uint) Visit(visitor SchemaVisitor) { visitor.VisitUint(uintt) }
func (Uint) Kind() Kind                        { return UintKind }

type UintValue uint

func (UintValue) isSchemaValue() {}
func (UintValue) Kind() Kind     { return UintKind }

// TODO: optimise out passing by
type Rune struct{}

func (runee Rune) Visit(visitor SchemaVisitor) { visitor.VisitRune(runee) }
func (Rune) Kind() Kind                        { return RuneKind }

type RuneValue rune

func (RuneValue) isSchemaValue() {}
func (RuneValue) Kind() Kind     { return RuneKind }

type Enum struct {
	Values []string
}

func (enum Enum) Visit(visitor SchemaVisitor) { visitor.VisitEnum(enum) }
func (Enum) Kind() Kind                       { return EnumKind }

type EnumValue uint

func (EnumValue) isSchemaValue() {}
func (EnumValue) Kind() Kind     { return EnumKind }

type EnumInt struct {
	Min int
	Max int
}

func (enumInt EnumInt) Visit(visitor SchemaVisitor) { visitor.VisitEnumInt(enumInt) }
func (EnumInt) Kind() Kind                          { return EnumIntKind }

// TODO: int64
type EnumIntValue int

func (EnumIntValue) isSchemaValue() {}
func (EnumIntValue) Kind() Kind     { return EnumIntKind }

type Product []Schema

func (product Product) Visit(visitor SchemaVisitor) { visitor.VisitProduct(product) }
func (Product) Kind() Kind                          { return ProductKind }

type ProductValue []SchemaValue

func (ProductValue) isSchemaValue() {}
func (ProductValue) Kind() Kind     { return ProductKind }

type Coproduct []Schema

func (coproduct Coproduct) Visit(visitor SchemaVisitor) { visitor.VisitCoproduct(coproduct) }
func (Coproduct) Kind() Kind                            { return CoproductKind }

type CoproductValue struct {
	SchemaValue
	Variant uint
}

func (CoproductValue) isSchemaValue() {}
func (CoproductValue) Kind() Kind     { return CoproductKind }

type List struct {
	Schema
}

func (list List) Visit(visitor SchemaVisitor) { visitor.VisitList(list) }
func (List) Kind() Kind                       { return ListKind }

type ListValue []SchemaValue

func (ListValue) isSchemaValue() {}
func (ListValue) Kind() Kind     { return ListKind }

type Array struct {
	Schema
	Count uint
}

func (array Array) Visit(visitor SchemaVisitor) { visitor.VisitArray(array) }
func (Array) Kind() Kind                        { return ArrayKind }

type ArrayValue []SchemaValue

func (ArrayValue) isSchemaValue() {}
func (ArrayValue) Kind() Kind     { return ArrayKind }

type Optional struct {
	Schema
}

func (optional Optional) Visit(visitor SchemaVisitor) { visitor.VisitOptional(optional) }
func (Optional) Kind() Kind                           { return OptionalKind }

type OptionalValue struct {
	SchemaValue
	IsValid bool
}

func (OptionalValue) isSchemaValue() {}
func (OptionalValue) Kind() Kind     { return OptionalKind }

// TODO: alias in alias???
type Alias struct {
	Schema
	Name string
}

func (alias Alias) Visit(visitor SchemaVisitor) { visitor.VisitAlias(alias) }
func (Alias) Kind() Kind                        { return AliasKind }

type AliasValue struct {
	Name  string
	Value SchemaValue
}

func (AliasValue) isSchemaValue() {}
func (AliasValue) Kind() Kind     { return AliasKind }

// type SchemaWithMetadata struct {
// 	Schema
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

type LogicalScheme interface {
	Schema() Schema
	LogicalKind() LogicalKind
}

type LogicalKind string

const (
	StringKind LogicalKind = "std:string"
	// TODO: DateKind     LogicalKind = "std:date"
	// TODO: DatetimeKind LogicalKind = "std:datetime"
	// TODO: datetime with millis, micros, nanos
	// TODO: decimal
)

func (kind LogicalKind) Schema() Schema {
	switch kind {
	case StringKind:
		return List{Schema: Rune{}}
	// case DateKind:
	// 	return Int{}
	// case DatetimeKind:
	// 	return Int{}
	default:
		panic(fmt.Sprintf("unknown kind: %v", kind))
	}
}

type String struct{}

func (s String) Visit(visitor SchemaVisitor) { visitor.VisitString(s) }
func (String) Kind() Kind                    { return ListKind }

func (String) Schema() Schema           { return List{Schema: Rune{}} }
func (String) LogicalKind() LogicalKind { return StringKind }

type StringValue string

func (StringValue) isSchemaValue() {}
func (StringValue) Kind() Kind     { return ListKind }
