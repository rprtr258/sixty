package sixty

func AssertSchemas[S Schema, V SchemaValue](t Schema, v SchemaValue) (tt S, vv V, ok bool) {
	tt, okSchema := t.(S)
	vv, okValue := v.(V)
	ok = okSchema && okValue
	return
}
