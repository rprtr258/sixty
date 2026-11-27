package sixty

// Codec is a paired encoder and decoder sharing one data representation T.
type Codec[T any] interface {
	Encode(Schema, SchemaValue) (T, error)
	Decode(Schema, T) (SchemaValue, error)
}

type codecFunc[T any] struct {
	encode func(Schema, SchemaValue) (T, error)
	decode func(Schema, T) (SchemaValue, error)
}

func (c codecFunc[T]) Encode(s Schema, v SchemaValue) (T, error) { return c.encode(s, v) }
func (c codecFunc[T]) Decode(s Schema, t T) (SchemaValue, error) { return c.decode(s, t) }

func CodecFunc[T any](
	encode func(Schema, SchemaValue) (T, error),
	decode func(Schema, T) (SchemaValue, error),
) Codec[T] {
	return codecFunc[T]{encode, decode}
}
