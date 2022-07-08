package bytesenc

import (
	"fmt"

	sixty "github.com/rprtr258/sixty/pkg"
)

// TODO: bit encoder/decoder
// TODO: byte decoder
type ByteDecoder struct {
	// TODO: make private
	Buffer interface {
		ReadByte() (byte, error)
		ReadRune() (rune, int, error)
	}
}

func (j *ByteDecoder) readUint() uint {
	// TODO: implement
	// TODO: error handling
	b, _ := j.Buffer.ReadByte()
	switch {
	case b < 128:
		return uint(b)
	}
	panic("not implemented")
}

// TODO: fix overflow
func (j *ByteDecoder) readInt() int {
	zigZag := j.readUint()
	if zigZag%2 == 0 {
		return int(zigZag / 2)
	} else {
		return int(zigZag)/(-2) - 1
	}
}

func (j *ByteDecoder) VisitString() string {
	stringLength := j.readUint()
	s := make([]rune, 0, stringLength)
	for i := uint(0); i < stringLength; i++ {
		// TODO: rewrite
		r, _, _ := j.Buffer.ReadRune()
		s = append(s, r)
	}
	return string(s)
}
func (j *ByteDecoder) VisitBit(_ sixty.Bit) sixty.BitValue {
	b, _ := j.Buffer.ReadByte()
	switch b {
	case 0:
		return false
	case 1:
		return true
	default:
		panic(fmt.Sprintf("expected bit, found: %d", b))
	}
}
func (j *ByteDecoder) VisitInt(_ sixty.Int) sixty.IntValue {
	return sixty.IntValue(j.readInt())
}
func (j *ByteDecoder) VisitUint(_ sixty.Uint) sixty.UintValue {
	return sixty.UintValue(j.readUint())
}
func (j *ByteDecoder) VisitRune(_ sixty.Rune) sixty.RuneValue {
	r, _, _ := j.Buffer.ReadRune()
	return sixty.RuneValue(r)
}
func (j *ByteDecoder) VisitEnum(enum sixty.Enum) sixty.EnumValue {
	n := j.readUint()
	if n < uint(len(enum.Values)) {
		return sixty.EnumValue(n)
	}
	panic(fmt.Sprintf("found %d'th enum value while there are only %d of them", n, len(enum.Values)))
}
func (j *ByteDecoder) VisitEnumInt(enumInt sixty.EnumInt) sixty.EnumIntValue {
	n := j.readUint()
	return sixty.EnumIntValue(int(n) + enumInt.Min)
}
func (j *ByteDecoder) VisitProduct(product sixty.Product) sixty.ProductValue {
	res := make([]sixty.TypeValue, 0, len(product))
	for _, x := range product {
		res = append(res, sixty.VisitValue2(x, j))
	}
	return res
}
func (j *ByteDecoder) VisitCoproduct(coproduct sixty.Coproduct) sixty.CoproductValue {
	variant := j.readUint()
	return sixty.CoproductValue{
		TypeValue: sixty.VisitValue2(coproduct[variant], j),
		Variant:   variant,
	}
}
func (j *ByteDecoder) VisitList(list sixty.List) sixty.ListValue {
	listLen := j.readUint()
	res := make([]sixty.TypeValue, 0, listLen)
	for i := uint(0); i < listLen; i++ {
		res = append(res, sixty.VisitValue2(list.Type, j))
	}
	return res
}
func (j *ByteDecoder) VisitArray(array sixty.Array) sixty.ArrayValue {
	res := make([]sixty.TypeValue, 0, array.Count)
	for i := uint(0); i < array.Count; i++ {
		res = append(res, sixty.VisitValue2(array.Type, j))
	}
	return res
}
func (j *ByteDecoder) VisitOptional(optional sixty.Optional) sixty.OptionalValue {
	isPresent, _ := j.Buffer.ReadByte()
	switch isPresent {
	case 0:
		return sixty.OptionalValue{IsValid: false}
	case 1:
		return sixty.OptionalValue{
			TypeValue: sixty.VisitValue2(optional.Type, j),
			IsValid:   false,
		}
	default:
		panic(fmt.Sprintf("option variant must be either 0 or 1, but it is %d", isPresent))
	}
}
func (j *ByteDecoder) VisitAlias(alias sixty.Alias) sixty.AliasValue {
	return sixty.AliasValue{
		Name:  alias.Name,
		Value: sixty.VisitValue2(alias.Type, j),
	}
}
