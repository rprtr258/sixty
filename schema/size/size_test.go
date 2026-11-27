package size_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/rprtr258/sixty/internal"
	"github.com/rprtr258/sixty/schema/size"
)

func TestBttvDB(t *testing.T) {
	t.Parallel()

	got := size.SchemaSize(internal.BttvEmotesDBSchema)
	assert.Equal(t, size.Inf, got)
}
