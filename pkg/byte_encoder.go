package sixty

import (
	"bytes"
)

// TODO: bit encoder/decoder
// TODO: byte decoder
type ByteEncoder struct {
	builder bytes.Buffer
}

func (j *ByteEncoder) writeUint(x uint) {
	// TODO: implement
	switch {
	case x < 128:
		j.builder.WriteByte(byte(x))
	}
}
func (j *ByteEncoder) writeInt(x int) {
	if x > 0 {
		j.writeUint(uint(x) * 2)
	} else {
		// TODO: fix and test overflow/steal from protobuf
		j.writeUint(uint(-2*x) - 1)
	}
}

func (j *ByteEncoder) VisitString(s string) {
	j.writeUint(uint(len(s)))
	j.builder.WriteString(s)
}
func (j *ByteEncoder) VisitBit(_ Bit, bit BitValue) {
	if bit {
		j.builder.WriteByte(1)
	} else {
		j.builder.WriteByte(0)
	}
}
func (j *ByteEncoder) VisitInt(_ Int, n IntValue) {
	j.writeInt(int(n))
}
func (j *ByteEncoder) VisitUint(_ Uint, n UintValue) {
	j.writeUint(uint(n))
}
func (j *ByteEncoder) VisitRune(_ Rune, r RuneValue) {
	j.builder.WriteRune(rune(r))
}
func (j *ByteEncoder) VisitEnum(enum Enum, n EnumValue) {
	j.writeUint(Log2(int(n)))
}
func (j *ByteEncoder) VisitEnumInt(enumInt EnumInt, n EnumIntValue) {
	j.writeUint(Log2(int(n) - enumInt.Min))
}
func (j *ByteEncoder) VisitProduct(product Product, p ProductValue) {
	for i, x := range p {
		VisitValue(product[i], x, j)
	}
}
func (j *ByteEncoder) VisitCoproduct(coproduct Coproduct, c CoproductValue) {
	j.writeUint(Log2(int(c.Variant)))
	VisitValue(coproduct[c.Variant], c.TypeValue, j)
}
func (j *ByteEncoder) VisitList(list List, l ListValue) {
	j.writeUint(uint(len(l)))
	for _, x := range l {
		VisitValue(list.Type, x, j)
	}
}
func (j *ByteEncoder) VisitArray(array Array, a ArrayValue) {
	for _, x := range a {
		VisitValue(array.Type, x, j)
	}
}
func (j *ByteEncoder) VisitOptional(optional Optional, o OptionalValue) {
	if !o.IsValid {
		j.writeUint(0)
	} else {
		j.writeUint(1)
		VisitValue(optional.Type, o.TypeValue, j)
	}
}
func (j *ByteEncoder) VisitAlias(alias Alias, a AliasValue) {
	VisitValue(alias.Type, a.Value, j)
}

func (j *ByteEncoder) Bytes() []byte {
	return j.builder.Bytes()
}
