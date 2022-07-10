package bytes_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/rprtr258/sixty/encoding/bytes"
	"github.com/rprtr258/sixty/internal"
)

type entry struct {
	Emote string
	Ids   []string
}

type BttvEmotesDB []entry

func TestBttvDBDecode(t *testing.T) {
	var got BttvEmotesDB
	err := bytes.Decode(internal.BttvEmotesDBSchema, internal.BttvEmotesDBBytes, &got)
	assert.NoError(t, err)
	assert.Equal(t,
		BttvEmotesDB([]entry{{
			Emote: "Kappa",
			Ids:   []string{"KappaOwner1", "KappaOwner2"},
		}, {
			Emote: "PogChamp",
			Ids:   []string{"PogChampOwner1", "PogChampOwner2"},
		}}),
		got,
	)
}
