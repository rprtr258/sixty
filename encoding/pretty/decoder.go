package pretty

import (
	"github.com/rprtr258/sixty"

	"github.com/rprtr258/sixty/encoding/text"
)

// Decode parses data written by pretty.Encode, which is the text encoding with
// insignificant whitespace and trailing commas — both accepted by text.Decode.
func Decode(t sixty.Schema, data string) (sixty.SchemaValue, error) {
	// the text parse error is the most precise one for pretty data as well
	return text.Decode(t, data) //nolint:wrapcheck
}
