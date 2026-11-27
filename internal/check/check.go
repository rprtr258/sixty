// Package check provides schema/value type assertion helpers for codecs.
package check

import "github.com/rprtr258/sixty"

// AssertSchemas checks that the schema is of type S and the value is of type V.
func AssertSchemas[S sixty.Schema, V any](t sixty.Schema, v any) (tt S, vv V, ok bool) {
	tt, okSchema := t.(S)
	vv, okValue := v.(V)
	ok = okSchema && okValue
	return
}
