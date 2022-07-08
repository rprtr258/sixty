package bytesenc

import (
	"bytes"

	sixty "github.com/rprtr258/sixty/pkg"
)

// TODO: bit encoder/decoder
// TODO: byte decoder
type ByteDecoder struct {
	buf bytes.Buffer
}

func (j *ByteDecoder) writeUint(x uint) {
	// TODO: implement
	switch {
	case x < 128:
		j.buf.WriteByte(byte(x))
	}
}
func (j *ByteDecoder) writeInt(x int) {
	if x > 0 {
		j.writeUint(uint(x) * 2)
	} else {
		// TODO: fix and test overflow/steal from protobuf
		j.writeUint(uint(-2*x) - 1)
	}
}

func (j *ByteDecoder) VisitString(s string) {
	j.writeUint(uint(len(s)))
	j.buf.WriteString(s)
}
func (j *ByteDecoder) VisitBit(_ sixty.Bit, bit sixty.BitValue) {
	if bit {
		j.buf.WriteByte(1)
	} else {
		j.buf.WriteByte(0)
	}
}
func (j *ByteDecoder) VisitInt(_ sixty.Int, n sixty.IntValue) {
	j.writeInt(int(n))
}
func (j *ByteDecoder) VisitUint(_ sixty.Uint, n sixty.UintValue) {
	j.writeUint(uint(n))
}
func (j *ByteDecoder) VisitRune(_ sixty.Rune, r sixty.RuneValue) {
	j.buf.WriteRune(rune(r))
}
func (j *ByteDecoder) VisitEnum(enum sixty.Enum, n sixty.EnumValue) {
	j.writeUint(sixty.Log2(int(n)))
}
func (j *ByteDecoder) VisitEnumInt(enumInt sixty.EnumInt, n sixty.EnumIntValue) {
	j.writeUint(sixty.Log2(int(n) - enumInt.Min))
}
func (j *ByteDecoder) VisitProduct(product sixty.Product, p sixty.ProductValue) {
	for i, x := range p {
		sixty.VisitValue(product[i], x, j)
	}
}
func (j *ByteDecoder) VisitCoproduct(coproduct sixty.Coproduct, c sixty.CoproductValue) {
	j.writeUint(sixty.Log2(int(c.Variant)))
	sixty.VisitValue(coproduct[c.Variant], c.TypeValue, j)
}
func (j *ByteDecoder) VisitList(list sixty.List, l sixty.ListValue) {
	j.writeUint(uint(len(l)))
	for _, x := range l {
		sixty.VisitValue(list.Type, x, j)
	}
}
func (j *ByteDecoder) VisitArray(array sixty.Array, a sixty.ArrayValue) {
	for _, x := range a {
		sixty.VisitValue(array.Type, x, j)
	}
}
func (j *ByteDecoder) VisitOptional(optional sixty.Optional, o sixty.OptionalValue) {
	if !o.IsValid {
		j.writeUint(0)
	} else {
		j.writeUint(1)
		sixty.VisitValue(optional.Type, o.TypeValue, j)
	}
}
func (j *ByteDecoder) VisitAlias(alias sixty.Alias, a sixty.AliasValue) {
	sixty.VisitValue(alias.Type, a.Value, j)
}

func (j *ByteDecoder) Bytes() []byte {
	return j.buf.Bytes()
}
