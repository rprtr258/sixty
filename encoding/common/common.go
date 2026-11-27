package common

import "reflect"

func isSignedKind(kind reflect.Kind) bool {
	return reflect.Int <= kind && kind <= reflect.Int64
}

func isUnsignedKind(kind reflect.Kind) bool {
	return reflect.Uint <= kind && kind <= reflect.Uint64
}

func isNilKind(kind reflect.Kind) bool {
	return reflect.Chan <= kind && kind <= reflect.Slice
}
