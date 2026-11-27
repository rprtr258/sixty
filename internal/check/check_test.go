package check_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/rprtr258/sixty"
	"github.com/rprtr258/sixty/internal/check"
)

func TestAssertSchemas(t *testing.T) {
	t.Parallel()

	tt, vv, ok := check.AssertSchemas[sixty.Int, sixty.IntValue](sixty.Int{}, sixty.IntValue(3))
	require.True(t, ok)
	require.Equal(t, sixty.Int{}, tt)
	require.Equal(t, sixty.IntValue(3), vv)

	_, _, ok = check.AssertSchemas[sixty.Int, sixty.IntValue](sixty.Bit{}, sixty.IntValue(3))
	require.False(t, ok)

	_, _, ok = check.AssertSchemas[sixty.Int, sixty.IntValue](sixty.Int{}, sixty.UintValue(3))
	require.False(t, ok)

	_, _, ok = check.AssertSchemas[sixty.Int, sixty.IntValue](sixty.Bit{}, sixty.UintValue(3))
	require.False(t, ok)
}
