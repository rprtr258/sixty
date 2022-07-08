package bytesenc

import (
	"bytes"

	sixty "github.com/rprtr258/sixty/pkg"
)

// TODO: bit encoder/decoder
// TODO: byte decoder
// TODO: compress list/array of bits
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
func (j *ByteEncoder) VisitBit(_ sixty.Bit, bit sixty.BitValue) {
	if bit {
		j.builder.WriteByte(1)
	} else {
		j.builder.WriteByte(0)
	}
}
func (j *ByteEncoder) VisitInt(_ sixty.Int, n sixty.IntValue) {
	j.writeInt(int(n))
}
func (j *ByteEncoder) VisitUint(_ sixty.Uint, n sixty.UintValue) {
	j.writeUint(uint(n))
}
func (j *ByteEncoder) VisitRune(_ sixty.Rune, r sixty.RuneValue) {
	j.builder.WriteRune(rune(r))
}
func (j *ByteEncoder) VisitEnum(enum sixty.Enum, n sixty.EnumValue) {
	j.writeUint(uint(n))
}
func (j *ByteEncoder) VisitEnumInt(enumInt sixty.EnumInt, n sixty.EnumIntValue) {
	j.writeUint(uint(int(n) - enumInt.Min))
}
func (j *ByteEncoder) VisitProduct(product sixty.Product, p sixty.ProductValue) {
	for i, x := range p {
		sixty.VisitValue(product[i], x, j)
	}
}
func (j *ByteEncoder) VisitCoproduct(coproduct sixty.Coproduct, c sixty.CoproductValue) {
	j.writeUint(c.Variant)
	sixty.VisitValue(coproduct[c.Variant], c.TypeValue, j)
}
func (j *ByteEncoder) VisitList(list sixty.List, l sixty.ListValue) {
	j.writeUint(uint(len(l)))
	for _, x := range l {
		sixty.VisitValue(list.Type, x, j)
	}
}
func (j *ByteEncoder) VisitArray(array sixty.Array, a sixty.ArrayValue) {
	for _, x := range a {
		sixty.VisitValue(array.Type, x, j)
	}
}
func (j *ByteEncoder) VisitOptional(optional sixty.Optional, o sixty.OptionalValue) {
	if !o.IsValid {
		j.writeUint(0)
	} else {
		j.writeUint(1)
		sixty.VisitValue(optional.Type, o.TypeValue, j)
	}
}
func (j *ByteEncoder) VisitAlias(alias sixty.Alias, a sixty.AliasValue) {
	sixty.VisitValue(alias.Type, a.Value, j)
}

func (j *ByteEncoder) Bytes() []byte {
	return j.builder.Bytes()
}
