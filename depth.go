package sixty

func maxSchemaDepth(ts []Schema) Size {
	max := Size(0)
	for _, x := range ts {
		depth := SchemaDepth(x)
		if max < depth {
			max = depth
		}
	}
	return max
}

func SchemaDepth(t Schema) Size {
	switch tt := t.(type) {
	case String, Bit, Int, Uint, Rune, Enum, EnumInt:
		return 0
	case Product:
		return maxSchemaDepth(tt)
	case Coproduct:
		return maxSchemaDepth(tt)
	case List:
		return 1 + SchemaDepth(tt.Schema)
	case Array:
		return 1 + SchemaDepth(tt.Schema)
	case Optional:
		return 1 + SchemaDepth(tt.Schema)
	default:
		panic("something unimplemented")
	}
}
