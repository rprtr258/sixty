package sixty

import "testing"

// The isSchemaValue methods are unexported interface markers; they can only be
// executed through the SchemaValue interface from inside the package.
func TestSchemaValueMarkerMethods(t *testing.T) {
	t.Parallel()

	values := []SchemaValue{
		BitValue(true),
		IntValue(-1),
		UintValue(1),
		RuneValue('a'),
		EnumValue(0),
		EnumIntValue(0),
		ProductValue{},
		CoproductValue{},
		ListValue{},
		ArrayValue{},
		OptionalValue{},
		StringValue(""),
	}
	for _, v := range values {
		v.isSchemaValue()
	}
	// direct calls too: the marker bodies are empty, so the coverage counter
	// sits on the function entry and interface dispatch may not register it
	BitValue(true).isSchemaValue()
	IntValue(-1).isSchemaValue()
	UintValue(1).isSchemaValue()
	RuneValue('a').isSchemaValue()
	EnumValue(0).isSchemaValue()
	EnumIntValue(0).isSchemaValue()
	ProductValue{}.isSchemaValue()
	CoproductValue{}.isSchemaValue()
	ListValue{}.isSchemaValue()
	ArrayValue{}.isSchemaValue()
	OptionalValue{}.isSchemaValue()
	StringValue("").isSchemaValue()
}
