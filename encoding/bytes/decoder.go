package bytes

import (
	"bytes"
	"fmt"
	"reflect"

	"github.com/rprtr258/sixty"
)

type Decoder struct {
	buffer sixty.Reader
}

func (e *Decoder) readUint() uint64 {
	// TODO: implement
	b := e.buffer.ReadByte()
	switch {
	case b < 128:
		return uint64(b)
	}
	panic("not implemented decoding uints larger than 127")
}

func (e *Decoder) readInt() int64 {
	return sixty.FromZigZag(e.readUint())
}

// TODO: error
func (e *Decoder) appendError(format string, args ...any) {
	panic(fmt.Sprintf("error at position %d: ", e.buffer.Pos()) + fmt.Sprintf(format, args...))
}

// TODO: out arg
func (e *Decoder) Read(t sixty.Schema, rv reflect.Value) {
	if _, ok := t.(sixty.String); ok && rv.Kind() == reflect.String {
		stringLength := e.readUint()
		sb := bytes.Buffer{}
		sb.Grow(int(stringLength))
		for i := uint64(0); i < stringLength; i++ {
			sb.WriteRune(e.buffer.ReadRune())
		}
		rv.Set(reflect.ValueOf(sb.String()))
	} else if _, _, ok := sixty.AssertSchemas[sixty.Bit, bool](t, rv); ok {
		b := e.buffer.ReadByte()
		switch b {
		case 0:
			rv.SetBool(false)
		case 1:
			rv.SetBool(true)
		default:
			e.appendError("expected bit which is either 0 or 1, found: %d", b)
		}
	} else if _, _, ok := sixty.AssertSchemas[sixty.Int, int64](t, rv); ok {
		rv.SetInt(e.readInt())
	} else if _, _, ok := sixty.AssertSchemas[sixty.Uint, uint64](t, rv); ok {
		rv.SetUint(e.readUint())
	} else if _, _, ok := sixty.AssertSchemas[sixty.Rune, rune](t, rv); ok {
		rv.Set(reflect.ValueOf(e.buffer.ReadRune()))
	} else if tt, _, ok := sixty.AssertSchemas[sixty.Enum, uint](t, rv); ok {
		n := e.readUint()
		if n < uint64(len(tt.Values)) {
			rv.SetUint(n)
			return
		}
		e.appendError("found %d'th enum value while there are only %d of them", n, len(tt.Values))
	} else if tt, _, ok := sixty.AssertSchemas[sixty.EnumInt, int64](t, rv); ok {
		n := e.readUint()
		rv.SetInt(int64(n) + tt.Min)
	} else if tt, ok := t.(sixty.Product); ok && rv.Kind() == reflect.Struct {
		for i, x := range tt {
			e.Read(x, rv.Field(i))
		}
		// case sixty.Coproduct:
		// 	variant := e.readUint()
		// 	return sixty.CoproductValue{
		// 		SchemaValue: e.Read(tt[variant]),
		// 		Variant:     uint(variant),
		// 	}
	} else if tt, ok := t.(sixty.List); ok && rv.Kind() == reflect.Slice {
		listLen := int(e.readUint())
		lst := reflect.MakeSlice(rv.Type(), listLen, listLen)
		for i := 0; i < listLen; i++ {
			e.Read(tt.Schema, lst.Index(i))
		}
		rv.Set(lst)
		// case sixty.Array:
		// 	res := make([]sixty.SchemaValue, 0, tt.Count)
		// 	for i := uint(0); i < tt.Count; i++ {
		// 		res = append(res, e.Read(tt.Schema))
		// 	}
		// 	return sixty.ArrayValue(res)
		// case sixty.Optional:
		// 	isPresent := e.buffer.ReadByte()
		// 	switch isPresent {
		// 	case 0:
		// 		return sixty.OptionalValue{IsValid: false}
		// 	case 1:
		// 		return sixty.OptionalValue{
		// 			SchemaValue: e.Read(tt.Schema),
		// 			IsValid:     false,
		// 		}
		// 	default:
		// 		e.appendError("option variant must be either 0 or 1, but it is %d", isPresent)
		// 	}
		// case sixty.Alias:
		// 	return sixty.AliasValue{
		// 		Name:  tt.Name,
		// 		Value: e.Read(tt.Schema),
		// 	}
		// default:
		// 	panic("something unimplemented")
		// }
	}
}

func Decode(t sixty.Schema, data []byte, v any) error {
	buffer := bytes.NewBuffer(data)
	decoder := Decoder{buffer: sixty.NewReader(buffer)}
	decoder.Read(t, reflect.ValueOf(v).Elem())
	return decoder.buffer.Error
}
