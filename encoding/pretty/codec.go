package pretty

import "github.com/rprtr258/sixty"

// Codec round-trips schema values through the pretty-printed text format.
var Codec = sixty.CodecFunc(encode, decode)
