package sixty_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/rprtr258/sixty"
	"github.com/rprtr258/sixty/internal"
)

func TestBttvDB(t *testing.T) {
	got := sixty.SchemaSize(internal.BttvEmotesDBSchema)
	assert.Equal(t, sixty.Inf, got)
}
