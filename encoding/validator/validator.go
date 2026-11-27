package validator

import (
	"fmt"
	"regexp"

	"github.com/rprtr258/sixty"
)

// InvalidSchemaError - schema definition rejected by Validator, Reason describes why.
type InvalidSchemaError struct {
	Reason string
}

func (e InvalidSchemaError) Error() string {
	return e.Reason
}

var IdentifierRegex = regexp.MustCompile("^[a-zA-Z]+[a-zA-Z.-0-9]*$")

type Validator struct {
	Errs []error
}

func (v *Validator) appendErrorf(format string, args ...any) {
	v.Errs = append(v.Errs, InvalidSchemaError{Reason: fmt.Sprintf(format, args...)})
}

func (e *Validator) Write(t sixty.Schema) {
	switch tt := t.(type) {
	case sixty.String, sixty.Bit, sixty.Int, sixty.Uint, sixty.Rune:
	case sixty.Enum:
		if len(tt.Values) < 2 {
			e.appendErrorf("wrong enum Values count, must be at least 2: %v", tt.Values)
		}
	case sixty.EnumInt:
		if tt.Max <= tt.Min {
			e.appendErrorf("incorrect enum int, max must be greater than min: min=%d, max=%d", tt.Min, tt.Max)
		}
	case sixty.Product:
		if len(tt) < 2 {
			e.appendErrorf("wrong product field count, must be at least 2: %d", len(tt))
		}
		for _, field := range tt {
			e.Write(field.Schema)
		}
	case sixty.Coproduct:
		if len(tt) < 2 {
			e.appendErrorf("wrong coproduct field count, must be at least 2: %d", len(tt))
		}
		for _, variant := range tt {
			e.Write(variant.Schema)
		}
	case sixty.List:
		e.Write(tt.Schema)
	case sixty.Array:
		if tt.Count < 2 {
			e.appendErrorf("wrong array length, must be at least 2: %d", tt.Count)
		}
		e.Write(tt.Schema)
	case sixty.Optional:
		e.Write(tt.Schema)
	default:
		panic("something unimplemented")
	}
}
