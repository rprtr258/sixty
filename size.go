package main

import "math"

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

// TODO: remake to some accum visitor
type SizeVisitor struct {
	Size
}

func (sv *SizeVisitor) VisitBit(Bit)   { sv.Size.Add(1) }
func (sv *SizeVisitor) VisitInt(Int)   { sv.Size = Inf }
func (sv *SizeVisitor) VisitUint(Uint) { sv.Size = Inf }
func (sv *SizeVisitor) VisitRune(Rune) { sv.Size = Inf }
func (sv *SizeVisitor) VisitProduct(product Product) {
	productSize := Size(1)
	for _, t := range product {
		sv2 := &SizeVisitor{}
		t.Visit(sv2)
		productSize.Add(sv2.Size)
	}
	sv.Size.Add(productSize)
}
func (sv *SizeVisitor) VisitCoproduct(coproduct Coproduct) {
	coproductSize := Size(0)
	for _, t := range coproduct {
		sv2 := &SizeVisitor{}
		t.Visit(sv2)
		if coproductSize < sv2.Size {
			coproductSize.Add(sv2.Size)
		}
	}
	// TODO: own Log2 on int
	sv.Size.Add(coproductSize).Add(Size(math.Ceil(math.Log2(float64(len(coproduct))))))
}
func (sv *SizeVisitor) VisitList(List) { sv.Size = Inf }
func (sv *SizeVisitor) VisitArray(arr Array) {
	sv2 := &SizeVisitor{}
	arr.Type.Visit(sv2)
	sv.Size.Add(sv2.Size.Multiply(int(arr.Count)))
}
func (sv *SizeVisitor) VisitOptional(optional Optional) {
	optional.Type.Visit(sv)
	sv.Size.Add(1)
}
func (sv *SizeVisitor) VisitEnum(enum Enum) {
	sv.Size.Add(Size(math.Ceil(math.Log2(float64(len(enum.Values))))))
}
func (sv *SizeVisitor) VisitEnumInt(enumInt EnumInt) {
	sv.Size.Add(Size(math.Ceil(math.Log2(float64(enumInt.Max - enumInt.Min)))))
}
func (sv *SizeVisitor) VisitAlias(alias Alias) { alias.Type.Visit(sv) }
