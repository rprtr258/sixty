package bytes

import (
	"bytes"

	"github.com/rprtr258/sixty"
)

// TODO: bit encoder/decoder
// TODO: streaming encoder/decoder, with size of element for lists (if it is not known)
type Encoder struct {
	// TODO: make unexported
	Buffer sixty.Writer
}

func (v *Encoder) writeUint(x uint64) {
	// TODO: implement
	switch {
	case x < 128:
		v.Buffer.WriteByte(byte(x))
	}
}
func (v *Encoder) writeInt(x int64) {
	// TODO: test overflow
	v.writeUint(uint64((x >> 63) ^ (x << 1)))
}

func (e *Encoder) Write(t sixty.Schema, v sixty.SchemaValue) {
	if _, vv, ok := sixty.AssertSchemas[sixty.String, sixty.StringValue](t, v); ok {
		e.writeUint(uint64(len(vv)))
		e.Buffer.WriteString(string(vv))
	} else if _, vv, ok := sixty.AssertSchemas[sixty.Bit, sixty.BitValue](t, v); ok {
		if vv {
			e.Buffer.WriteByte(1)
		} else {
			e.Buffer.WriteByte(0)
		}
	} else if _, vv, ok := sixty.AssertSchemas[sixty.Int, sixty.IntValue](t, v); ok {
		e.writeInt(int64(vv))
	} else if _, vv, ok := sixty.AssertSchemas[sixty.Uint, sixty.UintValue](t, v); ok {
		e.writeUint(uint64(vv))
	} else if _, vv, ok := sixty.AssertSchemas[sixty.Rune, sixty.RuneValue](t, v); ok {
		e.Buffer.WriteRune(rune(vv))
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
	encoder := Encoder{Buffer: sixty.NewWriter(&buffer)}
	encoder.Write(t, v)
	return buffer.Bytes(), encoder.Buffer.Error
}
