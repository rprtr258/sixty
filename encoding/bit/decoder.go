package bit

import (
	"bytes"
	"fmt"
	"math/bits"
	"unicode/utf8"

	"github.com/icza/bitio"
	"go.uber.org/multierr"

	"github.com/rprtr258/sixty"
	"github.com/rprtr258/sixty/internal/zigzag"
)

type decoder struct {
	buffer *bitio.CountReader
	Error  error
}

// DecodeError - data decoding failure, Bit is the number of bits read before the failure.
type DecodeError struct {
	Bit    int64
	Reason string
}

func (e DecodeError) Error() string {
	return fmt.Sprintf("error at bit %d: %s", e.Bit, e.Reason)
}

func (e *decoder) appendErrorf(format string, args ...any) {
	multierr.AppendInto(&e.Error, DecodeError{
		Bit:    e.buffer.BitsCount,
		Reason: fmt.Sprintf(format, args...),
	})
}

// failed reports whether decoding must stop because of a decode or read error.
func (e *decoder) failed() bool {
	return e.Error != nil || e.buffer.TryError != nil
}

// readUint reads an Elias-gamma-coded uint: value+1 as (n-1) zero bits
// followed by the n-bit binary of value+1, so 0 takes a single bit.
func (e *decoder) readUint() uint64 {
	zeros := 0
	for !e.buffer.TryReadBool() {
		zeros++
		if e.buffer.TryError != nil {
			return 0
		}
	}
	if zeros >= 64 {
		e.appendErrorf("uint is too long: %d bits", zeros+1)
		return 0
	}
	// the 1-bit consumed above is the leading bit of the (zeros+1)-bit value
	rest := e.buffer.TryReadBits(uint8(zeros))
	if e.buffer.TryError != nil {
		return 0
	}
	return (uint64(1)<<uint(zeros) | rest) - 1
}

func (e *decoder) readInt() int64 {
	return zigzag.Decode(e.readUint())
}

// readRune reads a UTF-8 encoded rune from the bit stream.
func (e *decoder) readRune() rune {
	buf := [4]byte{e.buffer.TryReadByte()}
	var size int
	switch lead := buf[0]; {
	case lead < 0x80:
		size = 1
	case lead&0xE0 == 0xC0:
		size = 2
	case lead&0xF0 == 0xE0:
		size = 3
	case lead&0xF8 == 0xF0:
		size = 4
	default:
		e.appendErrorf("invalid utf-8 lead byte %#x", lead)
		return utf8.RuneError
	}
	for i := 1; i < size; i++ {
		buf[i] = e.buffer.TryReadByte()
	}
	r, n := utf8.DecodeRune(buf[:size])
	if r == utf8.RuneError && n != size {
		e.appendErrorf("invalid utf-8 rune %#x", buf[:size])
		return utf8.RuneError
	}
	return r
}

// maxPrealloc caps size hints taken from the data so that malformed
// lengths cannot crash make with an absurd capacity.
const maxPrealloc = 1 << 10

func (e *decoder) readString() sixty.StringValue {
	stringLength := e.readUint()
	runes := make([]rune, 0, min(stringLength, maxPrealloc))
	for range stringLength {
		runes = append(runes, e.readRune())
		if e.failed() {
			return ""
		}
	}
	return sixty.StringValue(runes)
}

func (e *decoder) readBit() sixty.BitValue {
	if e.buffer.TryReadBool() {
		return true
	}
	return false
}

func (e *decoder) readEnum(tt sixty.Enum) (sixty.EnumValue, error) {
	n := e.buffer.TryReadBits(uint8(bits.Len(uint(len(tt.Values)))))
	if n < uint64(len(tt.Values)) {
		return sixty.EnumValue(n), nil
	}
	e.appendErrorf("found %d'th enum value while there are only %d of them", n, len(tt.Values))
	return 0, e.Error
}

func (e *decoder) readEnumInt(tt sixty.EnumInt) (sixty.EnumIntValue, error) {
	n := e.buffer.TryReadBits(uint8(bits.Len64(uint64(tt.Max - tt.Min))))
	if int64(n) <= tt.Max-tt.Min {
		return sixty.EnumIntValue(int64(n) + tt.Min), nil
	}
	e.appendErrorf("enum int value %d out of range [%d..%d]", int64(n)+tt.Min, tt.Min, tt.Max)
	return 0, e.Error
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

func (e *decoder) readCoproduct(tt sixty.Coproduct) (sixty.CoproductValue, error) {
	variant := e.buffer.TryReadBits(uint8(bits.Len(uint(len(tt)))))
	if variant < uint64(len(tt)) {
		return sixty.CoproductValue{
			SchemaValue: e.Read(tt[variant].Schema),
			Variant:     uint(variant),
		}, nil
	}
	e.appendErrorf("found variant %d while coproduct has only %d of them", variant, len(tt))
	return sixty.CoproductValue{}, e.Error
}

func (e *decoder) readList(tt sixty.List) sixty.ListValue {
	listLen := e.readUint()
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

func (e *decoder) readOptional(tt sixty.Optional) (sixty.OptionalValue, error) {
	switch isPresent := e.readUint(); isPresent {
	case 0:
		return sixty.OptionalValue{IsValid: false}, nil
	case 1:
		return sixty.OptionalValue{
			SchemaValue: e.Read(tt.Schema),
			IsValid:     true,
		}, nil
	default:
		e.appendErrorf("option variant must be either 0 or 1, but it is %d", isPresent)
		return sixty.OptionalValue{}, e.Error
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
		return sixty.IntValue(e.readInt())
	case sixty.Uint:
		return sixty.UintValue(e.readUint())
	case sixty.Rune:
		return sixty.RuneValue(e.readRune())
	case sixty.Enum:
		v, err := e.readEnum(tt)
		if err != nil {
			return nil
		}
		return v
	case sixty.EnumInt:
		v, err := e.readEnumInt(tt)
		if err != nil {
			return nil
		}
		return v
	case sixty.Product:
		return e.readProduct(tt)
	case sixty.Coproduct:
		v, err := e.readCoproduct(tt)
		if err != nil {
			return nil
		}
		return v
	case sixty.List:
		return e.readList(tt)
	case sixty.Array:
		return e.readArray(tt)
	case sixty.Optional:
		v, err := e.readOptional(tt)
		if err != nil {
			return nil
		}
		return v
	default:
		e.appendErrorf("unsupported schema type %T", t)
		return nil
	}
}

func decode(t sixty.Schema, data []byte) (sixty.SchemaValue, error) {
	buffer := bytes.NewBuffer(data)
	reader := bitio.NewCountReader(buffer)
	decoder := decoder{buffer: reader}
	res := decoder.Read(t)
	if decoder.Error != nil {
		return res, decoder.Error
	}
	return res, reader.TryError
}
