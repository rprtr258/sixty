package stream

import (
	"bytes"
	"fmt"
	"reflect"

	"go.uber.org/multierr"

	"github.com/rprtr258/sixty"
)

type Decoder struct {
	reader sixty.Reader
	// Error - first decoding failure
	Error error
}

func isSignedKind(kind reflect.Kind) bool {
	return reflect.Int <= kind && kind <= reflect.Int64
}

func isUnsignedKind(kind reflect.Kind) bool {
	return reflect.Uint <= kind && kind <= reflect.Uint64
}

// DecodeError - data decoding failure, Pos counts scanned bytes.
type DecodeError struct {
	Pos    uint
	Reason string
}

func (e DecodeError) Error() string {
	return fmt.Sprintf("error at position %d: %s", e.Pos, e.Reason)
}

func (e *Decoder) appendErrorf(format string, args ...any) {
	multierr.AppendInto(&e.Error, DecodeError{
		Pos:    e.reader.Pos(),
		Reason: fmt.Sprintf(format, args...),
	})
}

// failed reports whether decoding must stop because of a decode or read error.
func (e *Decoder) failed() bool {
	return e.Error != nil || e.reader.Error != nil
}

func (e *Decoder) Read(t sixty.Schema, rv reflect.Value) {
	if e.failed() {
		return
	}
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
		e.appendErrorf("unsupported schema type %T", t)
	}
}

func (e *Decoder) readString(rv reflect.Value) {
	if rv.Kind() != reflect.String && rv.Kind() != reflect.Interface {
		e.appendErrorf("expected string, got %s", rv.Kind())
		return
	}
	rv.Set(reflect.ValueOf(e.reader.ReadLengthString()))
}

func (e *Decoder) readBit(rv reflect.Value) {
	if rv.Kind() != reflect.Bool && rv.Kind() != reflect.Interface {
		e.appendErrorf("expected bool, got %s", rv.Kind())
		return
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
		return
	}
	if rv.Kind() != reflect.Interface {
		e.appendErrorf("expected signed int, got %s", rv.Kind())
		return
	}
	rv.Set(reflect.ValueOf(value))
}

func (e *Decoder) readUint(rv reflect.Value) {
	value := e.reader.ReadUint64()
	if isUnsignedKind(rv.Kind()) {
		rv.SetUint(value)
		return
	}
	if rv.Kind() != reflect.Interface {
		e.appendErrorf("expected unsigned int, got %s", rv.Kind())
		return
	}
	rv.Set(reflect.ValueOf(value))
}

func (e *Decoder) readRune(rv reflect.Value) {
	if rv.Kind() != reflect.Int32 && rv.Kind() != reflect.Interface {
		e.appendErrorf("expected rune, got %s", rv.Kind())
		return
	}
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
		return
	}
	if rv.Kind() != reflect.Interface {
		e.appendErrorf("expected enum, got %s", rv.Kind())
		return
	}
	rv.Set(reflect.ValueOf(uint(n)))
}

func (e *Decoder) readEnumInt(tt sixty.EnumInt, rv reflect.Value) {
	n := e.reader.ReadUint64()
	if int64(n) > tt.Max-tt.Min {
		e.appendErrorf("enum int value %d out of range [%d..%d]", int64(n)+tt.Min, tt.Min, tt.Max)
		return
	}
	if isSignedKind(rv.Kind()) {
		rv.SetInt(int64(n) + tt.Min)
		return
	}
	if rv.Kind() != reflect.Interface {
		e.appendErrorf("expected enum int, got %s", rv.Kind())
		return
	}
	rv.Set(reflect.ValueOf(int64(n) + tt.Min))
}

func (e *Decoder) readProduct(tt sixty.Product, rv reflect.Value) {
	if rv.Kind() != reflect.Struct {
		e.appendErrorf("expected struct, got %s", rv.Kind())
		return
	}
	for i, x := range tt {
		e.Read(x.Schema, rv.Field(i))
		if e.failed() {
			return
		}
	}
}

func (e *Decoder) readCoproduct(tt sixty.Coproduct, rv reflect.Value) {
	if rv.Kind() != reflect.Interface {
		e.appendErrorf("expected interface, got %s", rv.Kind())
		return
	}
	variant := e.reader.ReadUint64()
	if variant >= uint64(len(tt)) {
		e.appendErrorf("found variant %d while coproduct has only %d of them", variant, len(tt))
		return
	}
	e.Read(tt[variant].Schema, rv)
}

func (e *Decoder) readList(tt sixty.List, rv reflect.Value) {
	if rv.Kind() != reflect.Slice {
		e.appendErrorf("expected slice, got %s", rv.Kind())
		return
	}
	lst := reflect.MakeSlice(rv.Type(), 0, 0)
	for {
		switch hasNext := e.reader.ReadByte(); hasNext {
		case 0:
			rv.Set(lst)
			return
		case 1:
			x := reflect.New(rv.Type().Elem()).Elem()
			e.Read(tt.Schema, x)
			if e.failed() {
				return
			}
			lst = reflect.Append(lst, x)
		default:
			e.appendErrorf("expected list marker to be either 0 or 1, but it is %d", hasNext)
			return
		}
	}
}

func (e *Decoder) readArray(tt sixty.Array, rv reflect.Value) {
	if rv.Kind() != reflect.Slice && rv.Kind() != reflect.Array {
		e.appendErrorf("expected slice or array, got %s", rv.Kind())
		return
	}
	var lst reflect.Value
	if rv.Kind() == reflect.Array {
		lst = reflect.New(rv.Type()).Elem()
	} else {
		lst = reflect.MakeSlice(rv.Type(), int(tt.Count), int(tt.Count))
	}
	for i := uint(0); i < tt.Count; i++ {
		e.Read(tt.Schema, lst.Index(int(i)))
		if e.failed() {
			return
		}
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
	if decoder.Error != nil {
		return decoder.Error
	}
	return decoder.reader.Error
}
