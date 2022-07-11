package sixty

import (
	"fmt"
	"regexp"
)

var IdentifierRegex = regexp.MustCompile("^[a-zA-Z]+[a-zA-Z.-0-9]*$")

type Validator struct {
	Errs []error
}

func (v *Validator) appendError(format string, args ...any) {
	v.Errs = append(v.Errs, fmt.Errorf(format, args...))
}

func (e *Validator) Write(t Schema) {
	switch tt := t.(type) {
	case String, Bit, Int, Uint, Rune:
	case Enum:
		if len(tt.Values) < 2 {
			e.appendError("wrong enum Values count, must be at least 2: %v", tt.Values)
		}
	case EnumInt:
		if tt.Max <= tt.Min {
			e.appendError("incorrect enum int, max must be greater than min: min=%d, max=%d", tt.Min, tt.Max)
		}
	case Product:
		if len(tt) < 2 {
			e.appendError("wrong product field count, must be at least 2: %d", len(tt))
		}
		for _, field := range tt {
			e.Write(field)
		}
	case Coproduct:
		if len(tt) < 2 {
			e.appendError("wrong coproduct field count, must be at least 2: %d", len(tt))
		}
		for _, variant := range tt {
			e.Write(variant)
		}
	case List:
		e.Write(tt.Schema)
	case Array:
		if tt.Count < 2 {
			e.appendError("wrong array length, must be at least 2: %d", tt.Count)
		}
		e.Write(tt.Schema)
	case Optional:
		e.Write(tt.Schema)
	case Alias:
		if !IdentifierRegex.MatchString(tt.Name) {
			e.appendError("invalid alias name: %q", tt.Name)
		}
		e.Write(tt.Schema)
	default:
		panic("something unimplemented")
	}
}
