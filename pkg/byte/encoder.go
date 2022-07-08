package bytesenc

import (
	sixty "github.com/rprtr258/sixty/pkg"
)

// TODO: bit encoder/decoder
// TODO: byte decoder
// TODO: compress list/array of bits
// TODO: streaming encoder/decoder, with size of element for lists (if it is not known)
type ByteEncoder struct {
	// TODO: make unexported
	Buffer interface {
		WriteByte(byte) error
		WriteString(string) (int, error)
		WriteRune(rune) (int, error)
	}
}

func (v *ByteEncoder) writeUint(x uint) {
	// TODO: implement
	switch {
	case x < 128:
		v.Buffer.WriteByte(byte(x))
	}
}
func (v *ByteEncoder) writeInt(x int) {
	if x > 0 {
		v.writeUint(uint(x) * 2)
	} else {
		// TODO: fix and test overflow/steal from protobuf
		v.writeUint(uint(-2*x) - 1)
	}
}

func (v *ByteEncoder) VisitString(s string) {
	v.writeUint(uint(len(s)))
	v.Buffer.WriteString(s)
}
func (v *ByteEncoder) VisitBit(_ sixty.Bit, bit sixty.BitValue) {
	if bit {
		v.Buffer.WriteByte(1)
	} else {
		v.Buffer.WriteByte(0)
	}
}
func (v *ByteEncoder) VisitInt(_ sixty.Int, n sixty.IntValue) {
	v.writeInt(int(n))
}
func (v *ByteEncoder) VisitUint(_ sixty.Uint, n sixty.UintValue) {
	v.writeUint(uint(n))
}
func (v *ByteEncoder) VisitRune(_ sixty.Rune, r sixty.RuneValue) {
	v.Buffer.WriteRune(rune(r))
}
func (v *ByteEncoder) VisitEnum(enum sixty.Enum, n sixty.EnumValue) {
	v.writeUint(uint(n))
}
func (v *ByteEncoder) VisitEnumInt(enumInt sixty.EnumInt, n sixty.EnumIntValue) {
	v.writeUint(uint(int(n) - enumInt.Min))
}
func (v *ByteEncoder) VisitProduct(product sixty.Product, p sixty.ProductValue) {
	for i, x := range p {
		sixty.VisitValue(product[i], x, v)
	}
}
func (v *ByteEncoder) VisitCoproduct(coproduct sixty.Coproduct, c sixty.CoproductValue) {
	v.writeUint(c.Variant)
	sixty.VisitValue(coproduct[c.Variant], c.SchemaValue, v)
}
func (v *ByteEncoder) VisitList(list sixty.List, l sixty.ListValue) {
	v.writeUint(uint(len(l)))
	for _, x := range l {
		sixty.VisitValue(list.Schema, x, v)
	}
}
func (v *ByteEncoder) VisitArray(array sixty.Array, a sixty.ArrayValue) {
	for _, x := range a {
		sixty.VisitValue(array.Schema, x, v)
	}
}
func (v *ByteEncoder) VisitOptional(optional sixty.Optional, o sixty.OptionalValue) {
	if !o.IsValid {
		v.writeUint(0)
	} else {
		v.writeUint(1)
		sixty.VisitValue(optional.Schema, o.SchemaValue, v)
	}
}
func (v *ByteEncoder) VisitAlias(alias sixty.Alias, a sixty.AliasValue) {
	sixty.VisitValue(alias.Schema, a.Value, v)
}
