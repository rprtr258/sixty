package bytes

import (
	"bytes"
	"encoding/binary"

	"github.com/rprtr258/sixty"
)

// TODO: streaming encoder/decoder, with size of element for lists (if it is not known)
type Encoder struct {
	buffer     sixty.Writer
	uintBuffer []byte
}

func (v *Encoder) writeUint(x uint64) {
	n := binary.PutUvarint(v.uintBuffer, x)
	for i := 0; i < n; i++ {
		v.buffer.WriteByte(v.uintBuffer[i])
	}
}
func (v *Encoder) writeInt(x int64) {
	v.writeUint(sixty.ToZigZag(x))
}

func (e *Encoder) Write(t sixty.Schema, v sixty.SchemaValue) {
	if _, vv, ok := sixty.AssertSchemas[sixty.String, sixty.StringValue](t, v); ok {
		e.writeUint(uint64(len(vv)))
		e.buffer.WriteString(string(vv))
	} else if _, vv, ok := sixty.AssertSchemas[sixty.Bit, sixty.BitValue](t, v); ok {
		if vv {
			e.buffer.WriteByte(1)
		} else {
			e.buffer.WriteByte(0)
		}
	} else if _, vv, ok := sixty.AssertSchemas[sixty.Int, sixty.IntValue](t, v); ok {
		e.writeInt(int64(vv))
	} else if _, vv, ok := sixty.AssertSchemas[sixty.Uint, sixty.UintValue](t, v); ok {
		e.writeUint(uint64(vv))
	} else if _, vv, ok := sixty.AssertSchemas[sixty.Rune, sixty.RuneValue](t, v); ok {
		e.buffer.WriteRune(rune(vv))
	} else if _, vv, ok := sixty.AssertSchemas[sixty.Enum, sixty.EnumValue](t, v); ok {
		e.writeUint(uint64(vv))
	} else if tt, vv, ok := sixty.AssertSchemas[sixty.EnumInt, sixty.EnumIntValue](t, v); ok {
		e.writeUint(uint64(int64(vv) - tt.Min))
	} else if tt, vv, ok := sixty.AssertSchemas[sixty.Product, sixty.ProductValue](t, v); ok {
		for i, x := range vv {
			e.Write(tt[i], x)
		}
	} else if tt, vv, ok := sixty.AssertSchemas[sixty.Coproduct, sixty.CoproductValue](t, v); ok {
		e.writeUint(uint64(vv.Variant))
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
	encoder := Encoder{
		buffer:     sixty.NewWriter(&buffer),
		uintBuffer: make([]byte, 10),
	}
	encoder.Write(t, v)
	return buffer.Bytes(), encoder.buffer.Error
}
