package size

import (
	"math/bits"
	"strconv"

	"github.com/rprtr258/sixty"
)

// Size is either size of struct value (>= 0) or infinity (really -1).
type Size int

const Inf Size = -1

func log2(x uint64) Size {
	return Size(bits.Len64(x))
}

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

func SchemaSize(t sixty.Schema) Size {
	switch tt := t.(type) {
	case sixty.String:
		return Inf
	case sixty.Bit:
		return 1
	case sixty.Int:
		return Inf
	case sixty.Uint:
		return Inf
	case sixty.Rune:
		return Inf
	case sixty.Enum:
		return log2(uint64(len(tt.Values)))
	case sixty.EnumInt:
		return log2(uint64(tt.Max - tt.Min))
	case sixty.Product:
		sum := Size(0)
		for _, x := range tt {
			sum.Add(SchemaSize(x.Schema))
		}
		return sum
	case sixty.Coproduct:
		max := Size(0)
		for _, x := range tt {
			xSize := SchemaSize(x.Schema)
			if max < xSize {
				max = xSize
			}
		}
		max.Add(log2(uint64(len(tt))))
		return max
	case sixty.List:
		return Inf
	case sixty.Array:
		return SchemaSize(tt.Schema).Multiply(tt.Count)
	case sixty.Optional:
		return 1 + SchemaSize(tt.Schema)
	default:
		panic("something unimplemented")
	}
}
