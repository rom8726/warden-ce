package jobs

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestThroughOneStrategy_Present(t *testing.T) {
	str := newThroughOneStrategy(6)

	require.False(t, str.Present(0))
	require.True(t, str.Present(1))
	require.False(t, str.Present(2))
	require.True(t, str.Present(3))
	require.False(t, str.Present(4))
	require.True(t, str.Present(5))
	require.False(t, str.Present(6))
	require.False(t, str.Present(7))
	require.False(t, str.Present(8))
	require.False(t, str.Present(9))
}
