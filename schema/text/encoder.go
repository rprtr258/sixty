package text

import (
	"strconv"
	"strings"

	"github.com/rprtr258/sixty"
)

type encoder struct {
	builder            strings.Builder
	identSize          uint
	aliases            map[string]string
	ident              string
	disableIndent      bool
	prefix             string
	printTrailingComma bool
}

func NewPrettyEncoder(builder strings.Builder) encoder {
	return encoder{
		builder:            builder,
		aliases:            make(map[string]string),
		identSize:          0,
		ident:              "  ",
		prefix:             "\n",
		printTrailingComma: true,
	}
}

func NewEncoder(builder strings.Builder) encoder {
	return encoder{
		builder:            builder,
		aliases:            make(map[string]string),
		identSize:          0,
		ident:              "",
		prefix:             "",
		printTrailingComma: false,
	}
}

func (e *encoder) Write(t sixty.Schema) {
	if !e.disableIndent {
		e.builder.WriteString(strings.Repeat(e.ident, int(e.identSize)))
	} else {
		e.disableIndent = true
	}
	switch tt := t.(type) {
	case sixty.String:
		e.builder.WriteString(string(sixty.StringKind))
	case sixty.Bit:
		e.builder.WriteString("bit")
	case sixty.Int:
		e.builder.WriteString("int")
	case sixty.Uint:
		e.builder.WriteString("uint")
	case sixty.Rune:
		e.builder.WriteString("rune")
	case sixty.Enum:
		e.builder.WriteRune('{')
		i := 0
	iter_enum:
		e.builder.WriteString(tt.Values[i])
		if i < len(tt.Values)-1 {
			i++
			e.builder.WriteString(", ")
			goto iter_enum
		}
		e.builder.WriteRune('}')
	case sixty.EnumInt:
		e.builder.WriteRune('[')
		e.builder.WriteString(strconv.FormatInt(tt.Min, 10))
		e.builder.WriteString("..")
		e.builder.WriteString(strconv.FormatInt(tt.Max, 10))
		e.builder.WriteRune(']')
	case sixty.Product:
		e.builder.WriteRune('(')
		e.builder.WriteString(e.prefix)
		e.identSize++
		i := 0
	iter_product:
		e.builder.WriteString(strings.Repeat(e.ident, int(e.identSize)))
		e.builder.WriteString(tt[i].Name)
		e.builder.WriteString(" = ")
		e.disableIndent = true
		e.Write(tt[i].Schema)
		if i < len(tt)-1 {
			i++
			e.builder.WriteString(",")
			e.builder.WriteString(e.prefix)
			goto iter_product
		}
		if e.printTrailingComma {
			e.builder.WriteRune(',')
			e.builder.WriteString(e.prefix)
		}
		e.identSize--
		e.builder.WriteRune(')')
	case sixty.Coproduct:
		e.builder.WriteRune('[')
		e.builder.WriteString(e.prefix)
		e.identSize++
		i := 0
	iter_coproduct:
		e.builder.WriteString(strings.Repeat(e.ident, int(e.identSize)))
		e.builder.WriteString(tt[i].Name)
		e.builder.WriteString(" = ")
		e.disableIndent = true
		e.Write(tt[i].Schema)
		if i < len(tt)-1 {
			i++
			e.builder.WriteString(",")
			e.builder.WriteString(e.prefix)
			goto iter_coproduct
		}
		if e.printTrailingComma {
			e.builder.WriteRune(',')
			e.builder.WriteString(e.prefix)
		}
		e.identSize--
		e.builder.WriteRune(']')
	case sixty.List:
		e.builder.WriteString("[]")
		tmp := e.identSize
		e.identSize = 0
		e.Write(tt.Schema)
		e.identSize = tmp
	case sixty.Array:
		e.builder.WriteRune('[')
		e.builder.WriteString(strconv.FormatUint(uint64(tt.Count), 10))
		e.builder.WriteRune(']')
		tmp := e.identSize
		e.identSize = 0
		e.Write(tt.Schema)
		e.identSize = tmp
	case sixty.Optional:
		e.builder.WriteRune('?')
		tmp := e.identSize
		e.identSize = 0
		e.Write(tt.Schema)
		e.identSize = tmp
	// TODO: remove
	// case sixty.Alias:
	// 	e.builder.WriteString(tt.Name)
	// 	if _, ok := e.aliases[tt.Name]; ok {
	// 		return
	// 	}
	// 	pp := encoder{aliases: e.aliases}
	// 	pp.Write(tt.Schema)
	// 	e.aliases[tt.Name] = pp.String()
	default:
		panic("something unimplemented")
	}
}

func (e *encoder) String() string {
	var res strings.Builder
	for k, v := range e.aliases {
		res.WriteString(k)
		res.WriteString(" = ")
		res.WriteString(v)
		res.WriteString(e.prefix)
		res.WriteString(e.prefix)
	}
	res.WriteString(e.builder.String())
	return res.String()
}

func PrettyEncode(t sixty.Schema) string {
	tsv := NewPrettyEncoder(strings.Builder{})
	tsv.Write(t)
	return tsv.String()
}

func Encode(t sixty.Schema) string {
	tsv := NewEncoder(strings.Builder{})
	tsv.Write(t)
	return tsv.String()
}
