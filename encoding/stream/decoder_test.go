package stream_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/rprtr258/sixty"
	"github.com/rprtr258/sixty/encoding/stream"
	"github.com/rprtr258/sixty/internal"
)

func TestDecode_BttvDB(t *testing.T) {
	t.Parallel()

	type entry struct {
		Emote string
		Ids   []string
	}

	type BttvEmotesDB []entry

	var got BttvEmotesDB
	err := stream.Decode(internal.BttvEmotesDBSchema, internal.BttvEmotesDBStream, &got)
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

func TestDecode_Coproduct(t *testing.T) {
	t.Parallel()

	type s any
	var got s
	err := stream.Decode(sixty.Coproduct{
		{"int", sixty.Int{}},
		{"string", sixty.String{}},
	}, []byte{0x00, 0b01010100}, &got)
	assert.NoError(t, err)
	assert.Equal(t, s(int64(42)), got)
}
