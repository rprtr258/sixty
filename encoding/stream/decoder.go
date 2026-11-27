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

func isSignedKind(kind reflect.Kind) bool {
	return reflect.Int <= kind && kind <= reflect.Int64
}

func isUnsignedKind(kind reflect.Kind) bool {
	return reflect.Uint <= kind && kind <= reflect.Uint64
}

// TODO: error handling.
func (e *Decoder) appendErrorf(format string, args ...any) {
	panic(fmt.Sprintf("error at position %d: ", e.reader.Pos()) + fmt.Sprintf(format, args...))
}

func (e *Decoder) Read(t sixty.Schema, rv reflect.Value) {
	switch tt := t.(type) {
	case sixty.String:
		e.readString(rv)
	case sixty.Bit:
		e.readBit(rv)
	case sixty.Int:
		e.readInt(rv)
	case sixty.Uint:
		e.readUint(rv)
	case sixty.Rune:
		e.readRune(rv)
	case sixty.Enum:
		e.readEnum(tt, rv)
	case sixty.EnumInt:
		e.readEnumInt(tt, rv)
	case sixty.Product:
		e.readProduct(tt, rv)
	case sixty.Coproduct:
		e.readCoproduct(tt, rv)
	case sixty.List:
		e.readList(tt, rv)
	case sixty.Array:
		e.readArray(tt, rv)
	case sixty.Optional:
		e.readOptional(tt, rv)
	default:
		panic("something unimplemented")
	}
}

func (e *Decoder) readString(rv reflect.Value) {
	if rv.Kind() != reflect.String && rv.Kind() != reflect.Interface {
		panic("something unimplemented")
	}
	rv.Set(reflect.ValueOf(e.reader.ReadLengthString()))
}

func (e *Decoder) readBit(rv reflect.Value) {
	if rv.Kind() != reflect.Bool && rv.Kind() != reflect.Interface {
		panic("something unimplemented")
	}
	b := e.reader.ReadByte()
	var value bool
	switch b {
	case 0:
		value = false
	case 1:
		value = true
	default:
		e.appendErrorf("expected bit which is either 0 or 1, found: %d", b)
		return
	}
	if rv.Kind() == reflect.Bool {
		rv.SetBool(value)
	} else {
		rv.Set(reflect.ValueOf(value))
	}
}

func (e *Decoder) readInt(rv reflect.Value) {
	value := e.reader.ReadInt64()
	if isSignedKind(rv.Kind()) {
		rv.SetInt(value)
	} else {
		rv.Set(reflect.ValueOf(value))
	}
}

func (e *Decoder) readUint(rv reflect.Value) {
	value := e.reader.ReadUint64()
	if isUnsignedKind(rv.Kind()) {
		rv.SetUint(value)
	} else {
		rv.Set(reflect.ValueOf(value))
	}
}

func (e *Decoder) readRune(rv reflect.Value) {
	rv.Set(reflect.ValueOf(e.reader.ReadRune()))
}

func (e *Decoder) readEnum(tt sixty.Enum, rv reflect.Value) {
	n := e.reader.ReadUint64()
	if n >= uint64(len(tt.Values)) {
		e.appendErrorf("found %d'th enum value while there are only %d of them", n, len(tt.Values))
		return
	}
	if isUnsignedKind(rv.Kind()) {
		rv.SetUint(n)
	} else {
		rv.Set(reflect.ValueOf(uint(n)))
	}
}

func (e *Decoder) readEnumInt(tt sixty.EnumInt, rv reflect.Value) {
	n := e.reader.ReadUint64()
	if isSignedKind(rv.Kind()) {
		rv.SetInt(int64(n) + tt.Min)
	} else {
		rv.Set(reflect.ValueOf(int64(n) + tt.Min))
	}
}

func (e *Decoder) readProduct(tt sixty.Product, rv reflect.Value) {
	if rv.Kind() != reflect.Struct {
		panic("something unimplemented")
	}
	for i, x := range tt {
		e.Read(x.Schema, rv.Field(i))
	}
}

func (e *Decoder) readCoproduct(tt sixty.Coproduct, rv reflect.Value) {
	if rv.Kind() != reflect.Interface {
		panic("something unimplemented")
	}
	variant := e.reader.ReadUint64()
	e.Read(tt[variant].Schema, rv)
}

func (e *Decoder) readList(tt sixty.List, rv reflect.Value) {
	if rv.Kind() != reflect.Slice {
		panic("something unimplemented")
	}
	lst := reflect.MakeSlice(rv.Type(), 0, 0)
	hasNext := e.reader.ReadByte()
	for hasNext == 1 {
		x := reflect.New(rv.Type().Elem()).Elem()
		e.Read(tt.Schema, x)
		lst = reflect.Append(lst, x)
		hasNext = e.reader.ReadByte()
	}
	rv.Set(lst)
}

func (e *Decoder) readArray(tt sixty.Array, rv reflect.Value) {
	if rv.Kind() != reflect.Slice && rv.Kind() != reflect.Array {
		panic("something unimplemented")
	}
	var lst reflect.Value
	if rv.Kind() == reflect.Array {
		lst = reflect.New(rv.Type()).Elem()
	} else {
		lst = reflect.MakeSlice(rv.Type(), int(tt.Count), int(tt.Count))
	}
	for i := uint(0); i < tt.Count; i++ {
		e.Read(tt.Schema, lst.Index(int(i)))
	}
	rv.Set(lst)
}

func (e *Decoder) readOptional(tt sixty.Optional, rv reflect.Value) {
	isPresent := e.reader.ReadByte()
	switch isPresent {
	case 0:
		rv.Set(reflect.Zero(rv.Type()))
	case 1:
		if rv.Kind() == reflect.Pointer {
			rv.Set(reflect.New(rv.Type().Elem()))
			e.Read(tt.Schema, rv.Elem())
		} else {
			e.Read(tt.Schema, rv)
		}
	default:
		e.appendErrorf("option variant must be either 0 or 1, but it is %d", isPresent)
	}
}

func Decode(t sixty.Schema, data []byte, v any) error {
	buffer := bytes.NewBuffer(data)
	decoder := Decoder{reader: sixty.NewReader(buffer)}
	decoder.Read(t, reflect.ValueOf(v).Elem())
	return decoder.reader.Error
}
