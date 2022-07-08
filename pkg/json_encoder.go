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
		if product[i].Kind() == AliasKind && x.Kind() == AliasKind {
			j.builder.WriteRune('"')
			j.builder.WriteString((product[i]).(Alias).Name)
			j.builder.WriteString(`":`)
			VisitValue(product[i].(Alias).Schema, x.(AliasValue).Value, j)
		} else {
			j.builder.WriteRune('"')
			j.builder.WriteString(fmt.Sprint(i))
			j.builder.WriteString(`":`)
			VisitValue(product[i], x, j)
		}
		if i < len(p)-1 {
			j.builder.WriteRune(',')
		}
	}
	j.builder.WriteRune('}')
}
func (j *JsonEncoder) VisitCoproduct(coproduct Coproduct, c CoproductValue) {
	j.builder.WriteString(`{"`)
	j.builder.WriteString(fmt.Sprint(c.Variant))
	j.builder.WriteString(`":`)
	VisitValue(coproduct[c.Variant], c.SchemaValue, j)
	j.builder.WriteRune('}')
}
func (j *JsonEncoder) VisitList(list List, l ListValue) {
	j.builder.WriteRune('[')
	for i, x := range l {
		VisitValue(list.Schema, x, j)
		if i < len(l)-1 {
			j.builder.WriteRune(',')
		}
	}
	j.builder.WriteRune(']')
}
func (j *JsonEncoder) VisitArray(array Array, a ArrayValue) {
	j.builder.WriteRune('[')
	for i, x := range a {
		VisitValue(array.Schema, x, j)
		if uint(i) < array.Count-1 {
			j.builder.WriteRune(',')
		}
	}
	j.builder.WriteRune(']')
}
func (j *JsonEncoder) VisitOptional(optional Optional, o OptionalValue) {
	if !o.IsValid {
		j.builder.WriteString("null")
	} else {
		VisitValue(optional.Schema, o.SchemaValue, j)
	}
}
func (j *JsonEncoder) VisitAlias(alias Alias, a AliasValue) {
	VisitValue(alias.Schema, a.Value, j)
}

func (j *JsonEncoder) String() string {
	return j.builder.String()
}
