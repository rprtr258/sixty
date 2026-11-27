package json

import "github.com/rprtr258/sixty"

// Codec round-trips schema values through the json format.
var Codec = sixty.CodecFunc(encode, decode)
