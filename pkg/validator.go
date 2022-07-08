package sixty

import "fmt"

type Validator struct {
	Errs []error
}

func (v *Validator) appendError(format string, args ...any) {
	v.Errs = append(v.Errs, fmt.Errorf(format, args...))
}

func (v *Validator) VisitBit(Bit) {}

func (v *Validator) VisitInt(Int) {}

func (v *Validator) VisitUint(Uint) {}

func (v *Validator) VisitRune(Rune) {}

func (v *Validator) VisitProduct(product Product) {
	if len(product) < 2 {
		v.appendError("wrong product field count, must be at least 2: %d", len(product))
	}
	for _, field := range product {
		field.Visit(v)
	}
}

func (v *Validator) VisitCoproduct(coproduct Coproduct) {
	if len(coproduct) < 2 {
		v.appendError("wrong coproduct field count, must be at least 2: %d", len(coproduct))
	}
	for _, field := range coproduct {
		field.Visit(v)
	}
}

func (v *Validator) VisitList(list List) {
	list.Schema.Visit(v)
}

func (v *Validator) VisitArray(array Array) {
	if array.Count < 2 {
		v.appendError("wrong array length, must be at least 2: %d", array.Count)
	}
	array.Schema.Visit(v)
}

func (v *Validator) VisitOptional(optional Optional) {
	optional.Schema.Visit(v)
}

func (v *Validator) VisitEnum(enum Enum) {
}

func (v *Validator) VisitEnumInt(enumInt EnumInt) {
	if enumInt.Max <= enumInt.Min {
		v.appendError("incorrect enum int, max must be greater than min: min=%d, max=%d", enumInt.Min, enumInt.Max)
	}
}

func (v *Validator) VisitAlias(alias Alias) {
	// TODO: check format
	if alias.Name == "" {
		v.appendError("wrong alias name: %q", alias.Name)
	}
	alias.Schema.Visit(v)
}

func (v *Validator) VisitString(String) {}
