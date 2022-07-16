package stringify

import "github.com/rprtr258/sixty"

func BoolToText(b sixty.BitValue) string {
	if b {
		return "true"
	} else {
		return "false"
	}
}
