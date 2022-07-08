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

type JsonEncoder struct {
	builder strings.Builder
}

func (j *JsonEncoder) VisitString(s string) {
	j.builder.WriteString(fmt.Sprintf(`"%s"`, s))
}
func (j *JsonEncoder) VisitBit(_ Bit, bit BitValue) {
	if bit {
		j.builder.WriteString("true")
	} else {
		j.builder.WriteString("false")
	}
}
func (j *JsonEncoder) VisitInt(_ Int, n IntValue) {
	j.builder.WriteString(fmt.Sprint(n))
}
func (j *JsonEncoder) VisitUint(_ Uint, n UintValue) {
	j.builder.WriteString(fmt.Sprint(n))
}
func (j *JsonEncoder) VisitRune(_ Rune, r RuneValue) {
	j.builder.WriteString(fmt.Sprintf(`"%c"`, r))
}
func (j *JsonEncoder) VisitEnum(enum Enum, n EnumValue) {
	j.builder.WriteString(enum.Values[n])
}
func (j *JsonEncoder) VisitEnumInt(_ EnumInt, n EnumIntValue) {
	j.builder.WriteString(fmt.Sprint(n))
}
func (j *JsonEncoder) VisitProduct(product Product, p ProductValue) {
	j.builder.WriteString("{")
	for i, x := range p {
		// TODO: move out std:types as a separate case
		if checkTypeValue[Alias, AliasValue](product[i], x) && product[i].(Alias).Name[:4] != "std:" {
			j.builder.WriteString(fmt.Sprintf(`"%s":`, (product[i]).(Alias).Name))
			VisitValue(product[i].(Alias).Type, x.(AliasValue).Value, j)
		} else {
			j.builder.WriteString(fmt.Sprintf(`"%d":`, i))
			VisitValue(product[i], x, j)
		}
		if i < len(p)-1 {
			j.builder.WriteString(",")
		}
	}
	j.builder.WriteString("}")
}
func (j *JsonEncoder) VisitCoproduct(Coproduct, CoproductValue) {
	// TODO: implement
}
func (j *JsonEncoder) VisitList(list List, l ListValue) {
	j.builder.WriteString("[")
	for i, x := range l {
		VisitValue(list.Type, x, j)
		if i < len(l)-1 {
			j.builder.WriteString(",")
		}
	}
	j.builder.WriteString("]")
}
func (j *JsonEncoder) VisitArray(array Array, a ArrayValue) {
	j.builder.WriteString("[")
	for i, x := range a {
		VisitValue(array.Type, x, j)
		if uint(i) < array.Count-1 {
			j.builder.WriteString(",")
		}
	}
	j.builder.WriteString("]")
}
func (j *JsonEncoder) VisitOptional(optional Optional, o OptionalValue) {
	if !o.IsValid {
		j.builder.WriteString("null")
	} else {
		VisitValue(optional.Type, o.TypeValue, j)
	}
}
func (j *JsonEncoder) VisitAlias(alias Alias, a AliasValue) {
	VisitValue(alias.Type, a.Value, j)
}

func (j *JsonEncoder) String() string {
	return j.builder.String()
}
