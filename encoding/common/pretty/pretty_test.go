package pretty_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/rprtr258/sixty/encoding/common/pretty"
	"github.com/rprtr258/sixty/internal"
)

func TestBttvDBEncode(t *testing.T) {
	t.Parallel()

	got, err := pretty.Encode(internal.BttvEmotesDBSchema, internal.BttvEmotesDBData)
	assert.NoError(t, err)
	assert.Equal(t, `[
  {
    "Kappa",
    [
      "KappaOwner1",
      "KappaOwner2",
    ],
  },
  {
    "PogChamp",
    [
      "PogChampOwner1",
      "PogChampOwner2",
    ],
  },
]`, got)
}
