package inflect

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGoReceiverName(t *testing.T) {
	require.Equal(t, "u", GoReceiver("User"))
	require.Equal(t, "u", GoReceiver("UserPost"))
}
