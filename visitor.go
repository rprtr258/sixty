package sixty

// TODO: move to internal / package

func AssertSchemas[S Schema, V SchemaValue](t Schema, v SchemaValue) (tt S, vv V, ok bool) {
	tt, okSchema := t.(S)
	vv, okValue := v.(V)
	ok = okSchema && okValue
	return
}

func FromZigZag(x uint64) int64 {
	if x%2 == 0 {
		return int64(x / 2)
	} else {
		return int64(x)/(-2) - 1
	}
}

func ToZigZag(x int64) uint64 {
	return uint64((x >> 63) ^ (x << 1))
}
