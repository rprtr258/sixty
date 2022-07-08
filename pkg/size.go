package sixty

import (
	"strconv"
)

// Size is either size of struct value (>= 0) or infinity (really -1)
type Size int

const Inf Size = -1

func (x *Size) Add(y Size) *Size {
	if *x != Inf && y != Inf {
		*x += y
	}
	return x
}

func (x Size) Multiply(n uint) Size {
	if x == Inf {
		return Inf
	}
	return Size(uint(x) * n)
}

func (x Size) String() string {
	if x == Inf {
		return "Inf"
	}
	return strconv.Itoa(int(x))
}

func SchemaSize(t Schema) Size {
	switch tt := t.(type) {
	case String:
		return Inf
	case Bit:
		return 1
	case Int:
		return Inf
	case Uint:
		return Inf
	case Rune:
		return Inf
	case Enum:
		return Size(Log2(len(tt.Values)))
	case EnumInt:
		return Size(Log2((tt.Max - tt.Min)))
	case Product:
		sum := Size(0)
		for _, x := range tt {
			sum.Add(SchemaSize(x))
		}
		return sum
	case Coproduct:
		max := Size(0)
		for _, x := range tt {
			xSize := SchemaSize(x)
			if max < xSize {
				max = xSize
			}
		}
		max.Add(Size(Log2(len(tt))))
		return max
	case List:
		return Inf
	case Array:
		return SchemaSize(tt.Schema).Multiply(tt.Count)
	case Optional:
		return 1 + SchemaSize(tt.Schema)
	case Alias:
		return SchemaSize(tt.Schema)
	default:
		panic("something unimplemented")
	}
}
