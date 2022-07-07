package sixty

import (
	"math"
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

func (x Size) Multiply(n int) Size {
	if x == Inf {
		return Inf
	}
	return Size(int(x) * n)
}

func (x Size) String() string {
	if x == Inf {
		return "Inf"
	}
	return strconv.Itoa(int(x))
}

// TODO: remake to some accum visitor
type SizeVisitor struct{}

func (sv SizeVisitor) VisitBit(Bit) Size   { return 1 }
func (sv SizeVisitor) VisitInt(Int) Size   { return Inf }
func (sv SizeVisitor) VisitUint(Uint) Size { return Inf }
func (sv SizeVisitor) VisitRune(Rune) Size { return Inf }
func (sv SizeVisitor) VisitEnum(enum Enum) Size {
	return Size(math.Ceil(math.Log2(float64(len(enum.Values)))))
}
func (sv SizeVisitor) VisitEnumInt(enumInt EnumInt) Size {
	return Size(math.Ceil(math.Log2(float64(enumInt.Max - enumInt.Min))))
}
func (sv SizeVisitor) VisitProduct(product []Size) (s Size) {
	for _, x := range product {
		s.Add(x)
	}
	return
}
func (sv SizeVisitor) VisitCoproduct(coproduct []Size) (s Size) {
	for _, x := range coproduct {
		if s < x {
			s = x
		}
	}
	// TODO: own Log2 on int
	s.Add(Size(math.Ceil(math.Log2(float64(len(coproduct))))))
	return
}
func (sv SizeVisitor) VisitList(Size) Size                    { return Inf }
func (sv SizeVisitor) VisitArray(array Size, count uint) Size { return array.Multiply(int(count)) }
func (sv SizeVisitor) VisitOptional(optional Size) Size       { return optional + 1 }
func (sv SizeVisitor) VisitAlias(alias Size, _ string) Size   { return alias }
