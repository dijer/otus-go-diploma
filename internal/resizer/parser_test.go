package resizer

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParser_Success(t *testing.T) {
	data, err := parseURL("/fill/100/200/localhost/img.jpg")

	require.NoError(t, err)
	require.Equal(t, data, &parseData{
		Width:  100,
		Height: 200,
		URL:    "localhost/img.jpg",
	})
}

func TestParser_ErrorDimensions(t *testing.T) {
	data, err := parseURL("/fill/100/")

	require.Error(t, err)
	require.Empty(t, data)
}

func TestParser_IncorrectWidth(t *testing.T) {
	data, err := parseURL("/fill/abc/200/localhost/img.jpg")

	require.Error(t, err)
	require.Empty(t, data)
}

func TestParser_IncorrectHeight(t *testing.T) {
	data, err := parseURL("/fill/100/abc/localhost/img.jpg")

	require.Error(t, err)
	require.Empty(t, data)
}
