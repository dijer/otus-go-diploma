package resizer

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateHash(t *testing.T) {
	r := NewResizer(nil, "")
	hash := r.CreateHash(200, 300, "example.com/img.jpg")
	expected := "997d93937e9e9a13db08ef8b9c068c28"
	require.Equal(t, expected, hash)
}
