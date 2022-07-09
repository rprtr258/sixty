package bit

import (
	"bytes"
	"fmt"

	"github.com/icza/bitio"
	"github.com/rprtr258/sixty"
)

type Decoder struct {
	buffer *bitio.CountReader
}

func (e *Decoder) readUint() uint64 {
	n := 1
	for e.buffer.TryReadBool() {
		n++
	}
	res := e.buffer.TryReadBits(uint8(n))
	return res
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

// TODO: error
func (e *Decoder) appendError(format string, args ...any) {
	panic(fmt.Sprintf("error at position %d: ", e.buffer.BitsCount) + fmt.Sprintf(format, args...))
}

// TODO: change to read into pointer?
func (e *Decoder) Read(t sixty.Schema) sixty.SchemaValue {
	switch tt := t.(type) {
	case sixty.String:
		stringLength := e.readUint()
		s := make([]byte, stringLength)
		// TODO: fix reading unicode
		e.buffer.TryRead(s)
		return sixty.StringValue(string(s))
	case sixty.Bit:
		return sixty.BitValue(e.buffer.TryReadBool())
	case sixty.Int:
		return sixty.IntValue(e.readInt())
	case sixty.Uint:
		return sixty.UintValue(e.readUint())
	case sixty.Rune:
		// TODO: read rune
		r := e.buffer.TryReadByte()
		return sixty.RuneValue(r)
	case sixty.Enum:
		n := e.readUint()
		if n < uint64(len(tt.Values)) {
			return sixty.EnumValue(n)
		}
		e.appendError("found %d'th enum value while there are only %d of them", n, len(tt.Values))
	case sixty.EnumInt:
		n := e.readUint()
		return sixty.EnumIntValue(int64(n) + tt.Min)
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
		isPresent := e.buffer.TryReadByte()
		switch isPresent {
		case 0:
			return sixty.OptionalValue{IsValid: false}
		case 1:
			return sixty.OptionalValue{
				SchemaValue: e.Read(tt.Schema),
				IsValid:     false,
			}
		default:
			e.appendError("option variant must be either 0 or 1, but it is %d", isPresent)
		}
	case sixty.Alias:
		return sixty.AliasValue{
			Name:  tt.Name,
			Value: e.Read(tt.Schema),
		}
	default:
		panic("something unimplemented")
	}
	return nil
}

func Decode(t sixty.Schema, data []byte) (sixty.SchemaValue, error) {
	buffer := bytes.NewBuffer(data)
	reader := bitio.NewCountReader(buffer)
	decoder := Decoder{buffer: reader}
	res := decoder.Read(t)
	return res, reader.TryError
}
