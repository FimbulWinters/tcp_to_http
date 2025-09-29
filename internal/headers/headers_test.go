package headers

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidSingleHeader(t *testing.T) {

	headers := NewHeaders()
	data := []byte("Host: localhost:42069\r\n\r\n")
	n, done, err := headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	require.Equal(t, "localhost:42069", headers["host"])
	require.Equal(t, 23, n)
	require.False(t, done)
}

func TestInValidSpacingHeader(t *testing.T) {

	headers := NewHeaders()
	data := []byte("                   Host :  localhost:42069                  \r\n\r\n")
	n, done, err := headers.Parse(data)
	require.Error(t, err)
	require.Equal(t, 0, n)
	require.False(t, done)
}
