package bytes

import (
	"bytes"
	"fmt"

	"github.com/rprtr258/sixty"
)

type Decoder struct {
	reader sixty.Reader
}

// TODO: error
func (e *Decoder) appendError(format string, args ...any) {
	panic(fmt.Sprintf("error at position %d: ", e.reader.Pos()) + fmt.Sprintf(format, args...))
}

func (e *Decoder) Read(t sixty.Schema) sixty.SchemaValue {
	switch tt := t.(type) {
	case sixty.String:
		stringLength := e.reader.ReadUint64()
		sb := bytes.Buffer{}
		sb.Grow(int(stringLength))
		for i := uint64(0); i < stringLength; i++ {
			sb.WriteRune(e.reader.ReadRune())
		}
		return sixty.StringValue(sb.String())
	case sixty.Bit:
		b := e.reader.ReadByte()
		switch b {
		case 0:
			return sixty.BitValue(false)
		case 1:
			return sixty.BitValue(true)
		default:
			e.appendError("expected bit which is either 0 or 1, found: %d", b)
		}
	case sixty.Int:
		return sixty.IntValue(e.reader.ReadUint64())
	case sixty.Uint:
		return sixty.UintValue(e.reader.ReadUint64())
	case sixty.Rune:
		r := e.reader.ReadRune()
		return sixty.RuneValue(r)
	case sixty.Enum:
		n := e.reader.ReadUint64()
		if n < uint64(len(tt.Values)) {
			return sixty.EnumValue(n)
		}
		e.appendError("found %d'th enum value while there are only %d of them", n, len(tt.Values))
	case sixty.EnumInt:
		n := e.reader.ReadUint64()
		return sixty.EnumIntValue(int64(n) + tt.Min)
	case sixty.Product:
		res := make([]sixty.SchemaValue, 0, len(tt))
		for _, x := range tt {
			res = append(res, e.Read(x.Schema))
		}
		return sixty.ProductValue(res)
	case sixty.Coproduct:
		variant := e.reader.ReadUint64()
		return sixty.CoproductValue{
			SchemaValue: e.Read(tt[variant].Schema),
			Variant:     uint(variant),
		}
	case sixty.List:
		listLen := e.reader.ReadUint64()
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
		isPresent := e.reader.ReadByte()
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
	default:
		panic("something unimplemented")
	}
	return nil
}

func Decode(t sixty.Schema, data []byte) (sixty.SchemaValue, error) {
	buffer := bytes.NewBuffer(data)
	decoder := Decoder{reader: sixty.NewReader(buffer)}
	res := decoder.Read(t)
	return res, decoder.reader.Error
}
