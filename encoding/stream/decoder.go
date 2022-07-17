package stream

import (
	"bytes"
	"fmt"
	"reflect"

	"github.com/rprtr258/sixty"
)

type Decoder struct {
	reader sixty.Reader
}

// TODO: error
func (e *Decoder) appendError(format string, args ...any) {
	panic(fmt.Sprintf("error at position %d: ", e.reader.Pos()) + fmt.Sprintf(format, args...))
}

// TODO: out arg
func (e *Decoder) Read(t sixty.Schema, rv reflect.Value) {
	if _, ok := t.(sixty.String); ok && rv.Kind() == reflect.String {
		rv.Set(reflect.ValueOf(e.reader.ReadLengthString()))
	} else if _, _, ok := sixty.AssertSchemas[sixty.Bit, bool](t, rv); ok {
		b := e.reader.ReadByte()
		switch b {
		case 0:
			rv.SetBool(false)
		case 1:
			rv.SetBool(true)
		default:
			e.appendError("expected bit which is either 0 or 1, found: %d", b)
		}
	} else if _, _, ok := sixty.AssertSchemas[sixty.Int, int64](t, rv); ok {
		rv.SetInt(e.reader.ReadInt64())
	} else if _, _, ok := sixty.AssertSchemas[sixty.Int, any](t, rv); ok {
		rv.Set(reflect.ValueOf(e.reader.ReadInt64()))
	} else if _, _, ok := sixty.AssertSchemas[sixty.Uint, uint64](t, rv); ok {
		rv.SetUint(e.reader.ReadUint64())
	} else if _, _, ok := sixty.AssertSchemas[sixty.Rune, rune](t, rv); ok {
		rv.Set(reflect.ValueOf(e.reader.ReadRune()))
	} else if tt, _, ok := sixty.AssertSchemas[sixty.Enum, uint](t, rv); ok {
		n := e.reader.ReadUint64()
		if n < uint64(len(tt.Values)) {
			rv.SetUint(n)
			return
		}
		e.appendError("found %d'th enum value while there are only %d of them", n, len(tt.Values))
	} else if tt, _, ok := sixty.AssertSchemas[sixty.EnumInt, int64](t, rv); ok {
		n := e.reader.ReadUint64()
		rv.SetInt(int64(n) + tt.Min)
	} else if tt, ok := t.(sixty.Product); ok && rv.Kind() == reflect.Struct {
		for i, x := range tt {
			e.Read(x.Schema, rv.Field(i))
		}
	} else if tt, ok := t.(sixty.Coproduct); ok && rv.Kind() == reflect.Interface {
		variant := e.reader.ReadUint64()
		e.Read(tt[variant].Schema, rv)
	} else if tt, ok := t.(sixty.List); ok && rv.Kind() == reflect.Slice {
		lst := reflect.MakeSlice(rv.Type(), 0, 0)
		hasNext := e.reader.ReadByte()
		for hasNext == 1 {
			x := reflect.New(rv.Type().Elem()).Elem()
			e.Read(tt.Schema, x)
			lst = reflect.Append(lst, x)
			hasNext = e.reader.ReadByte()
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
	} else {
		panic("something unimplemented")
	}
}

func Decode(t sixty.Schema, data []byte, v any) error {
	buffer := bytes.NewBuffer(data)
	decoder := Decoder{reader: sixty.NewReader(buffer)}
	decoder.Read(t, reflect.ValueOf(v).Elem())
	return decoder.reader.Error
}
