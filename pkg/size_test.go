package sixty_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	// TODO: move everything to root
	sixty "github.com/rprtr258/sixty/pkg"
	"github.com/rprtr258/sixty/tests"
)

func TestBttvDB(t *testing.T) {
	assert.Equal(t, sixty.Inf, sixty.SchemaSize(tests.BttvEmotesDBSchema))
}
