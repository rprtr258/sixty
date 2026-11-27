package bytes

import (
	"bytes"
	"fmt"

	"go.uber.org/multierr"

	"github.com/rprtr258/sixty"
)

type decoder struct {
	reader sixty.Reader
	Error  error
}

// DecodeError - data decoding failure, Pos counts scanned bytes.
type DecodeError struct {
	Pos    uint
	Reason string
}

func (e DecodeError) Error() string {
	return fmt.Sprintf("error at position %d: %s", e.Pos, e.Reason)
}

func (e *decoder) appendErrorf(format string, args ...any) {
	multierr.AppendInto(&e.Error, DecodeError{
		Pos:    e.reader.Pos(),
		Reason: fmt.Sprintf(format, args...),
	})
}

// failed reports whether decoding must stop because of a decode or read error.
func (e *decoder) failed() bool {
	return e.Error != nil || e.reader.Error != nil
}

// maxPrealloc caps size hints taken from the data so that malformed
// lengths cannot crash make with an absurd capacity.
const maxPrealloc = 1 << 10

func (e *decoder) readString() sixty.StringValue {
	stringLength := e.reader.ReadUint64()
	sb := bytes.Buffer{}
	sb.Grow(int(min(stringLength, maxPrealloc)))
	for range stringLength {
		sb.WriteRune(e.reader.ReadRune())
		if e.failed() {
			return ""
		}
	}
	return sixty.StringValue(sb.String())
}

func (e *decoder) readBit() sixty.BitValue {
	switch b := e.reader.ReadByte(); b {
	case 0:
		return false
	case 1:
		return true
	default:
		e.appendErrorf("expected bit which is either 0 or 1, found: %d", b)
		return false
	}
}

func (e *decoder) readEnum(tt sixty.Enum) sixty.EnumValue {
	n := e.reader.ReadUint64()
	if n < uint64(len(tt.Values)) {
		return sixty.EnumValue(n)
	}
	e.appendErrorf("found %d'th enum value while there are only %d of them", n, len(tt.Values))
	return 0
}

func (e *decoder) readEnumInt(tt sixty.EnumInt) sixty.EnumIntValue {
	n := e.reader.ReadUint64()
	if int64(n) > tt.Max-tt.Min {
		e.appendErrorf("enum int value %d out of range [%d..%d]", int64(n)+tt.Min, tt.Min, tt.Max)
		return 0
	}
	return sixty.EnumIntValue(int64(n) + tt.Min)
}

func (e *decoder) readProduct(tt sixty.Product) sixty.ProductValue {
	res := make(sixty.ProductValue, 0, len(tt))
	for _, x := range tt {
		res = append(res, e.Read(x.Schema))
		if e.failed() {
			return nil
		}
	}
	return res
}

func (e *decoder) readCoproduct(tt sixty.Coproduct) sixty.CoproductValue {
	variant := e.reader.ReadUint64()
	if variant >= uint64(len(tt)) {
		e.appendErrorf("found variant %d while coproduct has only %d of them", variant, len(tt))
		return sixty.CoproductValue{}
	}
	return sixty.CoproductValue{
		SchemaValue: e.Read(tt[variant].Schema),
		Variant:     uint(variant),
	}
}

func (e *decoder) readList(tt sixty.List) sixty.ListValue {
	listLen := e.reader.ReadUint64()
	res := make(sixty.ListValue, 0, min(listLen, maxPrealloc))
	for range listLen {
		res = append(res, e.Read(tt.Schema))
		if e.failed() {
			return nil
		}
	}
	return res
}

func (e *decoder) readArray(tt sixty.Array) sixty.ArrayValue {
	res := make(sixty.ArrayValue, 0, tt.Count)
	for i := uint(0); i < tt.Count; i++ {
		res = append(res, e.Read(tt.Schema))
		if e.failed() {
			return nil
		}
	}
	return res
}

func (e *decoder) readOptional(tt sixty.Optional) sixty.OptionalValue {
	isPresent := e.reader.ReadByte()
	switch isPresent {
	case 0:
		return sixty.OptionalValue{IsValid: false}
	case 1:
		return sixty.OptionalValue{
			SchemaValue: e.Read(tt.Schema),
			IsValid:     true,
		}
	default:
		e.appendErrorf("option variant must be either 0 or 1, but it is %d", isPresent)
		return sixty.OptionalValue{}
	}
}

func (e *decoder) Read(t sixty.Schema) sixty.SchemaValue {
	if e.failed() {
		return nil
	}
	switch tt := t.(type) {
	case sixty.String:
		return e.readString()
	case sixty.Bit:
		return e.readBit()
	case sixty.Int:
		return sixty.IntValue(e.reader.ReadInt64())
	case sixty.Uint:
		return sixty.UintValue(e.reader.ReadUint64())
	case sixty.Rune:
		return sixty.RuneValue(e.reader.ReadRune())
	case sixty.Enum:
		return e.readEnum(tt)
	case sixty.EnumInt:
		return e.readEnumInt(tt)
	case sixty.Product:
		return e.readProduct(tt)
	case sixty.Coproduct:
		return e.readCoproduct(tt)
	case sixty.List:
		return e.readList(tt)
	case sixty.Array:
		return e.readArray(tt)
	case sixty.Optional:
		return e.readOptional(tt)
	default:
		e.appendErrorf("unsupported schema type %T", t)
		return nil
	}
}

func decode(t sixty.Schema, data []byte) (sixty.SchemaValue, error) {
	buffer := bytes.NewBuffer(data)
	decoder := decoder{reader: sixty.NewReader(buffer)}
	res := decoder.Read(t)
	if decoder.Error != nil {
		return res, decoder.Error
	}
	return res, decoder.reader.Error
}
