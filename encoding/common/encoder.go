package common

import (
	"fmt"
	"reflect"

	"github.com/rprtr258/sixty"
)

// Encode converts a native Go value to its schema-tagged representation.
func Encode(t sixty.Schema, v any) (sixty.SchemaValue, error) {
	return encode(t, reflect.ValueOf(v))
}

func encode(t sixty.Schema, rv reflect.Value) (sixty.SchemaValue, error) {
	switch tt := t.(type) {
	case sixty.String:
		return encodeString(rv)
	case sixty.Bit:
		return encodeBit(rv)
	case sixty.Int:
		return encodeInt(rv)
	case sixty.Uint:
		return encodeUint(rv)
	case sixty.Rune:
		return encodeRune(rv)
	case sixty.Enum:
		return encodeEnum(tt, rv)
	case sixty.EnumInt:
		return encodeEnumInt(tt, rv)
	case sixty.Product:
		return encodeProduct(tt, rv)
	case sixty.Coproduct:
		return encodeCoproduct(tt, rv)
	case sixty.List:
		return encodeList(tt.Schema, rv)
	case sixty.Array:
		return encodeArray(tt, rv)
	case sixty.Optional:
		return encodeOptional(tt.Schema, rv)
	default:
		return nil, MismatchError{Reason: fmt.Sprintf("unsupported schema %T", t)}
	}
}

func encodeString(rv reflect.Value) (sixty.SchemaValue, error) {
	if rv.Kind() != reflect.String {
		return nil, kindMismatch("string", rv.Kind())
	}
	return sixty.StringValue(rv.String()), nil
}

func encodeBit(rv reflect.Value) (sixty.SchemaValue, error) {
	if rv.Kind() != reflect.Bool {
		return nil, kindMismatch("bool", rv.Kind())
	}
	return sixty.BitValue(rv.Bool()), nil
}

func encodeInt(rv reflect.Value) (sixty.SchemaValue, error) {
	if !isSignedKind(rv.Kind()) {
		return nil, kindMismatch("signed integer", rv.Kind())
	}
	return sixty.IntValue(rv.Int()), nil
}

func encodeUint(rv reflect.Value) (sixty.SchemaValue, error) {
	if !isUnsignedKind(rv.Kind()) {
		return nil, kindMismatch("unsigned integer", rv.Kind())
	}
	return sixty.UintValue(rv.Uint()), nil
}

func encodeRune(rv reflect.Value) (sixty.SchemaValue, error) {
	if rv.Kind() != reflect.Int32 {
		return nil, kindMismatch("rune", rv.Kind())
	}
	return sixty.RuneValue(rune(rv.Int())), nil
}

func encodeEnum(tt sixty.Enum, rv reflect.Value) (sixty.SchemaValue, error) {
	if !isUnsignedKind(rv.Kind()) {
		return nil, kindMismatch("unsigned integer", rv.Kind())
	}
	n := rv.Uint()
	if n >= uint64(len(tt.Values)) {
		return nil, MismatchError{Reason: fmt.Sprintf("enum index %d out of range %d", n, len(tt.Values))}
	}
	return sixty.EnumValue(n), nil
}

func encodeEnumInt(tt sixty.EnumInt, rv reflect.Value) (sixty.SchemaValue, error) {
	if !isSignedKind(rv.Kind()) {
		return nil, kindMismatch("signed integer", rv.Kind())
	}
	n := rv.Int()
	if n < tt.Min || n > tt.Max {
		return nil, MismatchError{Reason: fmt.Sprintf("value %d out of enum int range [%d..%d]", n, tt.Min, tt.Max)}
	}
	return sixty.EnumIntValue(n), nil
}

func encodeProduct(tt sixty.Product, rv reflect.Value) (sixty.SchemaValue, error) {
	if rv.Kind() != reflect.Struct {
		return nil, kindMismatch("struct", rv.Kind())
	}
	if rv.NumField() != len(tt) {
		return nil, MismatchError{Reason: fmt.Sprintf("expected struct with %d fields, got %d", len(tt), rv.NumField())}
	}
	res := make(sixty.ProductValue, 0, len(tt))
	for i, field := range tt {
		fieldValue, err := encode(field.Schema, rv.Field(i))
		if err != nil {
			return nil, PathError{Path: fmt.Sprintf("field %q", field.Name), Err: err}
		}
		res = append(res, fieldValue)
	}
	return res, nil
}

// First matching variant wins.
func encodeCoproduct(tt sixty.Coproduct, rv reflect.Value) (sixty.SchemaValue, error) {
	// top-level Encode receives bare concrete values (reflect.ValueOf unwraps
	// the any), values nested in struct fields arrive interface-typed
	inner := rv
	if inner.Kind() == reflect.Interface {
		inner = inner.Elem()
	}
	for i, variant := range tt {
		variantValue, err := encode(variant.Schema, inner)
		if err == nil {
			return sixty.CoproductValue{SchemaValue: variantValue, Variant: uint(i)}, nil
		}
	}
	return nil, MismatchError{Reason: fmt.Sprintf("no coproduct variant matches value of kind %s", inner.Kind())}
}

func encodeElements(schema sixty.Schema, rv reflect.Value) (sixty.ListValue, error) {
	res := make(sixty.ListValue, 0, rv.Len())
	for i := range rv.Len() {
		elem, err := encode(schema, rv.Index(i))
		if err != nil {
			return nil, PathError{Path: fmt.Sprintf("element %d", i), Err: err}
		}
		res = append(res, elem)
	}
	return res, nil
}

func encodeList(schema sixty.Schema, rv reflect.Value) (sixty.SchemaValue, error) {
	if rv.Kind() != reflect.Slice && rv.Kind() != reflect.Array {
		return nil, kindMismatch("slice", rv.Kind())
	}
	elems, err := encodeElements(schema, rv)
	if err != nil {
		return nil, err
	}
	return elems, nil
}

func encodeArray(tt sixty.Array, rv reflect.Value) (sixty.SchemaValue, error) {
	if rv.Kind() != reflect.Slice && rv.Kind() != reflect.Array {
		return nil, kindMismatch("array", rv.Kind())
	}
	if rv.Len() != int(tt.Count) {
		return nil, MismatchError{Reason: fmt.Sprintf("expected array of %d elements, got %d", tt.Count, rv.Len())}
	}
	elems, err := encodeElements(tt.Schema, rv)
	if err != nil {
		return nil, err
	}
	return sixty.ArrayValue(elems), nil
}

func encodeOptional(schema sixty.Schema, rv reflect.Value) (sixty.SchemaValue, error) {
	// nil pointer/interface/slice/map/chan and invalid value mean "absent",
	// present pointers are dereferenced
	if !rv.IsValid() || isNilKind(rv.Kind()) && rv.IsNil() {
		return sixty.OptionalValue{IsValid: false}, nil
	}
	if rv.Kind() == reflect.Pointer {
		rv = rv.Elem()
	}
	inner, err := encode(schema, rv)
	if err != nil {
		return nil, err
	}
	return sixty.OptionalValue{SchemaValue: inner, IsValid: true}, nil
}
