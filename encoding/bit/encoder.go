package bit

import (
	"bytes"
	"math/bits"
	"unicode/utf8"

	"github.com/icza/bitio"

	"github.com/rprtr258/sixty"
)

type encoder struct {
	buffer *bitio.Writer
}

func (e *encoder) writeUint(x uint64) {
	n := bits.Len64(x)
	for i := 1; i < n; i++ {
		e.buffer.TryWriteBool(true)
	}
	e.buffer.TryWriteBool(false)
	e.buffer.TryWriteBits(x, uint8(n))
}

func (e *encoder) writeInt(x int64) {
	e.writeUint(sixty.ToZigZag(x))
}

func (e *encoder) Write(t sixty.Schema, v sixty.SchemaValue) {
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
	} else if tt, vv, ok := sixty.AssertSchemas[sixty.Enum, sixty.EnumValue](t, v); ok {
		// width must match the decoder's bits.Len(len(values)) read
		e.buffer.TryWriteBits(uint64(vv), uint8(bits.Len64(uint64(len(tt.Values)))))
	} else if tt, vv, ok := sixty.AssertSchemas[sixty.EnumInt, sixty.EnumIntValue](t, v); ok {
		// value stored as offset from Min, at the width of the full [Min..Max] span
		e.buffer.TryWriteBits(uint64(int64(vv)-tt.Min), uint8(bits.Len64(uint64(tt.Max-tt.Min))))
	} else if tt, vv, ok := sixty.AssertSchemas[sixty.Product, sixty.ProductValue](t, v); ok {
		for i, x := range vv {
			e.Write(tt[i].Schema, x)
		}
	} else if tt, vv, ok := sixty.AssertSchemas[sixty.Coproduct, sixty.CoproductValue](t, v); ok {
		e.buffer.TryWriteBits(uint64(vv.Variant), uint8(bits.Len64(uint64(len(tt)))))
		e.Write(tt[vv.Variant].Schema, vv.SchemaValue)
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
	} else {
		panic("something unimplemented")
	}
}

func encode(t sixty.Schema, v sixty.SchemaValue) ([]byte, error) {
	var buffer bytes.Buffer
	writer := bitio.NewWriter(&buffer)
	encoder := encoder{writer}
	encoder.Write(t, v)
	writer.TryAlign()
	return buffer.Bytes(), writer.TryError
}
