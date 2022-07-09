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
	// TODO: cache idents
	e.builder.WriteString(strings.Repeat(e.ident, int(e.identSize)))
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
		// TODO: iter while
		e.builder.WriteString(strings.Join(tt.Values, ", "))
		e.builder.WriteRune('}')
	case sixty.EnumInt:
		e.builder.WriteRune('(')
		e.builder.WriteString(strconv.Itoa(tt.Min))
		e.builder.WriteString("..")
		e.builder.WriteString(strconv.Itoa(tt.Max))
		e.builder.WriteRune(')')
	case sixty.Product:
		e.builder.WriteRune('(')
		e.builder.WriteString(e.prefix)
		e.identSize++
		// TODO: iter-while
		for i, field := range tt {
			e.Write(field)
			if i < len(tt)-1 {
				e.builder.WriteRune(',')
				e.builder.WriteString(e.prefix)
			}
		}
		if e.printTrailingComma {
			e.builder.WriteRune(',')
			e.builder.WriteString(e.prefix)
		}
		e.identSize--
		e.builder.WriteRune(')')
	case sixty.Coproduct:
		e.builder.WriteRune('{')
		e.builder.WriteString(e.prefix)
		e.identSize++
		// TODO: iter-while
		for i, variant := range tt {
			e.Write(variant)
			if i < len(tt)-1 {
				e.builder.WriteRune(',')
				e.builder.WriteString(e.prefix)
			}
		}
		if e.printTrailingComma {
			e.builder.WriteRune(',')
			e.builder.WriteString(e.prefix)
		}
		e.identSize--
		e.builder.WriteRune('}')
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
	case sixty.Alias:
		e.builder.WriteString(tt.Name)
		if _, ok := e.aliases[tt.Name]; ok {
			return
		}
		pp := encoder{aliases: e.aliases}
		pp.Write(tt.Schema)
		e.aliases[tt.Name] = pp.String()
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
