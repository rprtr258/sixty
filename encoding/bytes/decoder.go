package bytes

import (
	"bytes"
	"fmt"

	"github.com/rprtr258/sixty"
)

// TODO: bit encoder/decoder
// TODO: byte decoder
type Decoder struct {
	// TODO: make private
	// TODO: change to struct based on io.Writer
	Buffer sixty.Reader
}

func (e *Decoder) readUint() uint64 {
	// TODO: implement
	// TODO: error handling
	b := e.Buffer.ReadByte()
	switch {
	case b < 128:
		return uint64(b)
	}
	panic("not implemented")
}

// TODO: fix overflow
func (e *Decoder) readInt() int64 {
	zigZag := e.readUint()
	if zigZag%2 == 0 {
		return int64(zigZag / 2)
	} else {
		return int64(zigZag)/(-2) - 1
	}
}

// TODO: change to read into Pointer?
func (e *Decoder) Read(t sixty.Schema) sixty.SchemaValue {
	switch tt := t.(type) {
	case sixty.String:
		stringLength := e.readUint()
		s := make([]rune, 0, stringLength)
		for i := uint64(0); i < stringLength; i++ {
			// TODO: rewrite
			r := e.Buffer.ReadRune()
			s = append(s, r)
		}
		return sixty.StringValue(s)
	case sixty.Bit:
		b := e.Buffer.ReadByte()
		switch b {
		case 0:
			return sixty.BitValue(false)
		case 1:
			return sixty.BitValue(true)
		default:
			// TODO: position
			panic(fmt.Sprintf("expected bit which is either 0 or 1, found: %d", b))
		}
	case sixty.Int:
		return sixty.IntValue(e.readInt())
	case sixty.Uint:
		return sixty.UintValue(e.readUint())
	case sixty.Rune:
		r := e.Buffer.ReadRune()
		return sixty.RuneValue(r)
	case sixty.Enum:
		n := e.readUint()
		if n < uint64(len(tt.Values)) {
			return sixty.EnumValue(n)
		}
		panic(fmt.Sprintf("found %d'th enum value while there are only %d of them", n, len(tt.Values)))
	case sixty.EnumInt:
		n := e.readUint()
		return sixty.EnumIntValue(int(n) + tt.Min)
	case sixty.Product:
		res := make([]sixty.SchemaValue, 0, len(tt))
		for _, x := range tt {
			res = append(res, e.Read(x))
		}
		return sixty.ProductValue(res)
	case sixty.Coproduct:
		variant := e.readUint()
		return sixty.CoproductValue{
			SchemaValue: e.Read(tt[variant]),
			Variant:     uint(variant),
		}
	case sixty.List:
		listLen := e.readUint()
		res := make([]sixty.SchemaValue, 0, listLen)
		for i := uint64(0); i < listLen; i++ {
			res = append(res, e.Read(tt.Schema))
		}
		return sixty.ListValue(res)
	case sixty.Array:
		res := make([]sixty.SchemaValue, 0, tt.Count)
		for i := uint(0); i < tt.Count; i++ {
			res = append(res, e.Read(tt.Schema))
		}
		return sixty.ArrayValue(res)
	case sixty.Optional:
		isPresent := e.Buffer.ReadByte()
		switch isPresent {
		case 0:
			return sixty.OptionalValue{IsValid: false}
		case 1:
			return sixty.OptionalValue{
				SchemaValue: e.Read(tt.Schema),
				IsValid:     false,
			}
		default:
			panic(fmt.Sprintf("option variant must be either 0 or 1, but it is %d", isPresent))
		}
	case sixty.Alias:
		return sixty.AliasValue{
			Name:  tt.Name,
			Value: e.Read(tt.Schema),
		}
	default:
		panic("something unimplemented")
	}
}

func Decode(t sixty.Schema, data []byte) sixty.SchemaValue {
	buffer := bytes.NewBuffer(data)
	return (&Decoder{Buffer: sixty.NewReader(buffer)}).Read(t)
}
