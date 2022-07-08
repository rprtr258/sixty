package sixty

import (
	"fmt"
	"strings"
)

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
