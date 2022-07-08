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

type sizeVisitor struct{}

var SizeVisitor = DFSFactory((AccumVisitor[Size])(sizeVisitor{}))

func (sv sizeVisitor) VisitBit(Bit) Size   { return 1 }
func (sv sizeVisitor) VisitInt(Int) Size   { return Inf }
func (sv sizeVisitor) VisitUint(Uint) Size { return Inf }
func (sv sizeVisitor) VisitRune(Rune) Size { return Inf }
func (sv sizeVisitor) VisitEnum(enum Enum) Size {
	return Size(Log2(len(enum.Values)))
}
func (sv sizeVisitor) VisitEnumInt(enumInt EnumInt) Size {
	return Size(Log2((enumInt.Max - enumInt.Min)))
}
func (sv sizeVisitor) VisitProduct(product []Size) (s Size) {
	for _, x := range product {
		s.Add(x)
	}
	return
}
func (sv sizeVisitor) VisitCoproduct(coproduct []Size) (s Size) {
	for _, x := range coproduct {
		if s < x {
			s = x
		}
	}
	s.Add(Size(Log2(len(coproduct))))
	return
}
func (sv sizeVisitor) VisitList(Size) Size                    { return Inf }
func (sv sizeVisitor) VisitArray(array Size, count uint) Size { return array.Multiply(int(count)) }
func (sv sizeVisitor) VisitOptional(optional Size) Size       { return optional + 1 }
func (sv sizeVisitor) VisitAlias(alias Size, _ string) Size   { return alias }
