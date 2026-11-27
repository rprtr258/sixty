package common

import (
	"fmt"
	"reflect"

	"github.com/rprtr258/sixty"
)

func decode(t sixty.Schema, value sixty.SchemaValue, rv reflect.Value) error {
	switch tt := t.(type) {
	case sixty.String:
		return decodeString(value, rv)
	case sixty.Bit:
		return decodeBit(value, rv)
	case sixty.Int:
		return decodeInt(value, rv)
	case sixty.Uint:
		return decodeUint(value, rv)
	case sixty.Rune:
		return decodeRune(value, rv)
	case sixty.Enum:
		return decodeEnum(value, rv)
	case sixty.EnumInt:
		return decodeEnumInt(value, rv)
	case sixty.Product:
		return decodeProduct(tt, value, rv)
	case sixty.Coproduct:
		return decodeCoproduct(tt, value, rv)
	case sixty.List:
		return decodeList(tt.Schema, value, rv)
	case sixty.Array:
		return decodeArray(tt, value, rv)
	case sixty.Optional:
		return decodeOptional(tt.Schema, value, rv)
	default:
		return MismatchError{Reason: fmt.Sprintf("unsupported schema %T", t)}
	}
}

func decodeString(value sixty.SchemaValue, rv reflect.Value) error {
	vv, ok := value.(sixty.StringValue)
	if !ok {
		return valueMismatch("string", value)
	}
	if rv.Kind() == reflect.String {
		rv.SetString(string(vv))
		return nil
	}
	if rv.Kind() == reflect.Interface {
		rv.Set(reflect.ValueOf(string(vv)))
		return nil
	}
	return kindMismatch("string or any", rv.Kind())
}

func decodeBit(value sixty.SchemaValue, rv reflect.Value) error {
	vv, ok := value.(sixty.BitValue)
	if !ok {
		return valueMismatch("bit", value)
	}
	if rv.Kind() == reflect.Bool {
		rv.SetBool(bool(vv))
		return nil
	}
	if rv.Kind() == reflect.Interface {
		rv.Set(reflect.ValueOf(bool(vv)))
		return nil
	}
	return kindMismatch("bool or any", rv.Kind())
}

func decodeInt(value sixty.SchemaValue, rv reflect.Value) error {
	vv, ok := value.(sixty.IntValue)
	if !ok {
		return valueMismatch("int", value)
	}
	if isSignedKind(rv.Kind()) {
		rv.SetInt(int64(vv))
		return nil
	}
	if rv.Kind() == reflect.Interface {
		rv.Set(reflect.ValueOf(int64(vv)))
		return nil
	}
	return kindMismatch("signed integer or any", rv.Kind())
}

func decodeUint(value sixty.SchemaValue, rv reflect.Value) error {
	vv, ok := value.(sixty.UintValue)
	if !ok {
		return valueMismatch("uint", value)
	}
	if isUnsignedKind(rv.Kind()) {
		rv.SetUint(uint64(vv))
		return nil
	}
	if rv.Kind() == reflect.Interface {
		rv.Set(reflect.ValueOf(uint64(vv)))
		return nil
	}
	return kindMismatch("unsigned integer or any", rv.Kind())
}

func decodeRune(value sixty.SchemaValue, rv reflect.Value) error {
	vv, ok := value.(sixty.RuneValue)
	if !ok {
		return valueMismatch("rune", value)
	}
	if rv.Kind() == reflect.Int32 {
		rv.SetInt(int64(vv))
		return nil
	}
	if rv.Kind() == reflect.Interface {
		rv.Set(reflect.ValueOf(rune(vv)))
		return nil
	}
	return kindMismatch("rune or any", rv.Kind())
}

func decodeEnum(value sixty.SchemaValue, rv reflect.Value) error {
	vv, ok := value.(sixty.EnumValue)
	if !ok {
		return valueMismatch("enum", value)
	}
	if isUnsignedKind(rv.Kind()) {
		rv.SetUint(uint64(vv))
		return nil
	}
	if rv.Kind() == reflect.Interface {
		rv.Set(reflect.ValueOf(uint(vv)))
		return nil
	}
	return kindMismatch("unsigned integer or any", rv.Kind())
}

func decodeEnumInt(value sixty.SchemaValue, rv reflect.Value) error {
	vv, ok := value.(sixty.EnumIntValue)
	if !ok {
		return valueMismatch("enum int", value)
	}
	if isSignedKind(rv.Kind()) {
		rv.SetInt(int64(vv))
		return nil
	}
	if rv.Kind() == reflect.Interface {
		rv.Set(reflect.ValueOf(int64(vv)))
		return nil
	}
	return kindMismatch("signed integer or any", rv.Kind())
}

func decodeProduct(tt sixty.Product, value sixty.SchemaValue, rv reflect.Value) error {
	vv, ok := value.(sixty.ProductValue)
	if !ok {
		return valueMismatch("product", value)
	}
	if rv.Kind() != reflect.Struct {
		return kindMismatch("struct", rv.Kind())
	}
	if rv.NumField() != len(tt) {
		return MismatchError{Reason: fmt.Sprintf("expected struct with %d fields, got %d", len(tt), rv.NumField())}
	}
	if len(vv) != len(tt) {
		return MismatchError{Reason: fmt.Sprintf("expected product value with %d fields, got %d", len(tt), len(vv))}
	}
	for i, field := range tt {
		if err := decode(field.Schema, vv[i], rv.Field(i)); err != nil {
			return PathError{Path: fmt.Sprintf("field %q", field.Name), Err: err}
		}
	}
	return nil
}

func decodeCoproduct(tt sixty.Coproduct, value sixty.SchemaValue, rv reflect.Value) error {
	vv, ok := value.(sixty.CoproductValue)
	if !ok {
		return valueMismatch("coproduct", value)
	}
	if rv.Kind() != reflect.Interface {
		return kindMismatch("interface", rv.Kind())
	}
	if vv.Variant >= uint(len(tt)) {
		return MismatchError{Reason: fmt.Sprintf("coproduct variant %d out of range %d", vv.Variant, len(tt))}
	}
	return decode(tt[vv.Variant].Schema, vv.SchemaValue, rv)
}

func decodeElements(schema sixty.Schema, vv []sixty.SchemaValue, rv reflect.Value) error {
	for i, elem := range vv {
		if err := decode(schema, elem, rv.Index(i)); err != nil {
			return PathError{Path: fmt.Sprintf("element %d", i), Err: err}
		}
	}
	return nil
}

func decodeList(schema sixty.Schema, value sixty.SchemaValue, rv reflect.Value) error {
	vv, ok := value.(sixty.ListValue)
	if !ok {
		return valueMismatch("list", value)
	}
	if rv.Kind() != reflect.Slice {
		return kindMismatch("slice", rv.Kind())
	}
	lst := reflect.MakeSlice(rv.Type(), len(vv), len(vv))
	if err := decodeElements(schema, vv, lst); err != nil {
		return err
	}
	rv.Set(lst)
	return nil
}

func decodeArray(tt sixty.Array, value sixty.SchemaValue, rv reflect.Value) error {
	vv, ok := value.(sixty.ArrayValue)
	if !ok {
		return valueMismatch("array", value)
	}
	if len(vv) != int(tt.Count) {
		return MismatchError{Reason: fmt.Sprintf("expected array value of %d elements, got %d", tt.Count, len(vv))}
	}
	if rv.Kind() == reflect.Slice {
		lst := reflect.MakeSlice(rv.Type(), len(vv), len(vv))
		if err := decodeElements(tt.Schema, vv, lst); err != nil {
			return err
		}
		rv.Set(lst)
		return nil
	}
	if rv.Kind() == reflect.Array {
		return decodeElements(tt.Schema, vv, rv)
	}
	return kindMismatch("slice or array", rv.Kind())
}

func decodeOptional(schema sixty.Schema, value sixty.SchemaValue, rv reflect.Value) error {
	vv, ok := value.(sixty.OptionalValue)
	if !ok {
		return valueMismatch("optional", value)
	}
	if !vv.IsValid {
		rv.Set(reflect.Zero(rv.Type()))
		return nil
	}
	if rv.Kind() == reflect.Pointer {
		if rv.IsNil() {
			rv.Set(reflect.New(rv.Type().Elem()))
		}
		rv = rv.Elem()
	}
	return decode(schema, vv.SchemaValue, rv)
}

// Decode fills the value pointed to by v with the native Go representation of value.
func Decode(t sixty.Schema, value sixty.SchemaValue, v any) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Pointer {
		return MismatchError{Reason: fmt.Sprintf("decode target must be a pointer, got %s", rv.Kind())}
	}
	return decode(t, value, rv.Elem())
}
