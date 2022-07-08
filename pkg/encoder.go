package sixty

import (
	"fmt"
	"strings"
)

func checkTypeValue[T Type, V TypeValue](t Type, val TypeValue) bool {
	_, okType := t.(T)
	_, okValue := val.(V)
	return okType && okValue
}

func VisitValue(t Type, val TypeValue, visitor interface {
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
	VisitString(string)
}) {
	switch {
	case t == String:
		// TODO: check val is string
		s := val.(AliasValue)
		var b strings.Builder
		for _, c := range s.Value.(ListValue) {
			cc := c.(RuneValue)
			b.WriteRune(rune(cc))
		}
		visitor.VisitString(b.String())
	case checkTypeValue[Bit, BitValue](t, val):
		visitor.VisitBit(t.(Bit), val.(BitValue))
	case checkTypeValue[Int, IntValue](t, val):
		visitor.VisitInt(t.(Int), val.(IntValue))
	case checkTypeValue[Uint, UintValue](t, val):
		visitor.VisitUint(t.(Uint), val.(UintValue))
	case checkTypeValue[Rune, RuneValue](t, val):
		visitor.VisitRune(t.(Rune), val.(RuneValue))
	case checkTypeValue[Enum, EnumValue](t, val):
		visitor.VisitEnum(t.(Enum), val.(EnumValue))
	case checkTypeValue[EnumInt, EnumIntValue](t, val):
		visitor.VisitEnumInt(t.(EnumInt), val.(EnumIntValue))
	case checkTypeValue[Product, ProductValue](t, val):
		visitor.VisitProduct(t.(Product), val.(ProductValue))
	case checkTypeValue[Coproduct, CoproductValue](t, val):
		visitor.VisitCoproduct(t.(Coproduct), val.(CoproductValue))
	case checkTypeValue[List, ListValue](t, val):
		visitor.VisitList(t.(List), val.(ListValue))
	case checkTypeValue[Array, ArrayValue](t, val):
		visitor.VisitArray(t.(Array), val.(ArrayValue))
	case checkTypeValue[Optional, OptionalValue](t, val):
		visitor.VisitOptional(t.(Optional), val.(OptionalValue))
	case checkTypeValue[Alias, AliasValue](t, val):
		visitor.VisitAlias(t.(Alias), val.(AliasValue))
	default:
		// TODO: error: wrong schema
		panic(fmt.Sprintf("%T(%[1]v) is not %T(%[2]v) value", val, t))
	}
}
