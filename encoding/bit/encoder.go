package bit

import (
	"bytes"
	"math/bits"
	"unicode/utf8"

	"github.com/icza/bitio"

	"github.com/rprtr258/sixty"
)

type Encoder struct {
	buffer *bitio.Writer
}

func (e *Encoder) writeUint(x uint64) {
	n := bits.Len64(x)
	for i := 1; i < n; i++ {
		e.buffer.TryWriteBool(true)
	}
	e.buffer.TryWriteBool(false)
	e.buffer.TryWriteBits(x, uint8(n))
}
func (e *Encoder) writeInt(x int64) {
	e.writeUint(sixty.ToZigZag(x))
}

func (e *Encoder) Write(t sixty.Schema, v sixty.SchemaValue) {
	if _, vv, ok := sixty.AssertSchemas[sixty.String, sixty.StringValue](t, v); ok {
		e.writeUint(uint64(len(vv)))
		e.buffer.TryWrite([]byte(vv))
	} else if _, vv, ok := sixty.AssertSchemas[sixty.Bit, sixty.BitValue](t, v); ok {
		e.buffer.TryWriteBool(bool(vv))
	} else if _, vv, ok := sixty.AssertSchemas[sixty.Int, sixty.IntValue](t, v); ok {
		e.writeInt(int64(vv))
	} else if _, vv, ok := sixty.AssertSchemas[sixty.Uint, sixty.UintValue](t, v); ok {
		e.writeUint(uint64(vv))
	} else if _, vv, ok := sixty.AssertSchemas[sixty.Rune, sixty.RuneValue](t, v); ok {
		runeBuffer := make([]byte, 4)
		n := utf8.EncodeRune(runeBuffer, rune(vv))
		e.buffer.TryWrite(runeBuffer[:n])
	} else if _, vv, ok := sixty.AssertSchemas[sixty.Enum, sixty.EnumValue](t, v); ok {
		u64 := uint64(vv)
		e.buffer.TryWriteBits(u64, uint8(bits.Len64(u64)))
	} else if tt, vv, ok := sixty.AssertSchemas[sixty.EnumInt, sixty.EnumIntValue](t, v); ok {
		u64 := uint64(int64(vv) - tt.Min)
		e.buffer.TryWriteBits(u64, uint8(bits.Len64(u64)))
	} else if tt, vv, ok := sixty.AssertSchemas[sixty.Product, sixty.ProductValue](t, v); ok {
		for i, x := range vv {
			e.Write(tt[i], x)
		}
	} else if tt, vv, ok := sixty.AssertSchemas[sixty.Coproduct, sixty.CoproductValue](t, v); ok {
		u64 := uint64(len(tt))
		e.buffer.TryWriteBits(u64, uint8(bits.Len64(u64)))
		e.Write(tt[vv.Variant], vv.SchemaValue)
	} else if tt, vv, ok := sixty.AssertSchemas[sixty.List, sixty.ListValue](t, v); ok {
		e.writeUint(uint64(len(vv)))
		for _, x := range vv {
			e.Write(tt.Schema, x)
		}
	} else if tt, vv, ok := sixty.AssertSchemas[sixty.Array, sixty.ArrayValue](t, v); ok {
		for _, x := range vv {
			e.Write(tt.Schema, x)
		}
	} else if tt, vv, ok := sixty.AssertSchemas[sixty.Optional, sixty.OptionalValue](t, v); ok {
		if !vv.IsValid {
			e.writeUint(0)
		} else {
			e.writeUint(1)
			e.Write(tt.Schema, vv.SchemaValue)
		}
	} else if tt, vv, ok := sixty.AssertSchemas[sixty.Alias, sixty.AliasValue](t, v); ok {
		e.Write(tt.Schema, vv.Value)
	} else {
		panic("something unimplemented")
	}
}

func Encode(t sixty.Schema, v sixty.SchemaValue) ([]byte, error) {
	var buffer bytes.Buffer
	writer := bitio.NewWriter(&buffer)
	encoder := Encoder{
		buffer: writer,
	}
	encoder.Write(t, v)
	writer.TryAlign()
	return buffer.Bytes(), writer.TryError
}
