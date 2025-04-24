package resizer

import (
	"bytes"
	"context"
	"image"
	"image/color"
	"image/jpeg"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dijer/otus-go-diploma/internal/cache"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateHash(t *testing.T) {
	r := New(nil, "")
	hash := r.CreateHash(200, 300, "example.com/img.jpg")
	expected := "5adfc8c6c371135f84981d2bf1f774dd29ddf2c7e5825c0513cdb4617dbb36ba"
	require.Equal(t, expected, hash)
}

func TestGetCached(t *testing.T) {
	testCache := cache.New(10)
	r := New(testCache, "/tmp")

	testKey := "key"
	testValue := "cached_image_path"
	testCache.Set(cache.Key(testKey), &testValue)

	val, ok := r.GetCached(testKey)
	require.True(t, ok)
	assert.Equal(t, testValue, *val)

	val, ok = r.GetCached("bad_key")
	require.False(t, ok)
	assert.Nil(t, val)
}

func TestResizeImg(t *testing.T) {
	mockResizer := new(MockResizer)

	testImg := image.NewRGBA(image.Rect(0, 0, 200, 200))
	mockHash := "hash"
	mockWidth := int64(100)
	mockHeight := int64(100)

	expectedImg := image.NewRGBA(image.Rect(0, 0, 100, 100))
	expectedImg.Set(0, 0, color.RGBA{255, 255, 0, 255})

	mockResizer.On("ResizeImg", testImg, mockWidth, mockHeight, mockHash).Return(expectedImg, nil)

	resizedImg, err := mockResizer.ResizeImg(testImg, mockWidth, mockHeight, mockHash)

	require.NoError(t, err)
	require.NotNil(t, resizedImg)
	require.Equal(t, expectedImg.Bounds(), resizedImg.Bounds())
	require.Equal(t, expectedImg.At(0, 0), resizedImg.At(0, 0))
}

func TestLoadImg(t *testing.T) {
	testImg := createTestImage(100, 100)
	buf := new(bytes.Buffer)
	err := jpeg.Encode(buf, testImg, nil)
	require.NoError(t, err)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(buf.Bytes())
	}))
	defer server.Close()

	testCache := cache.New(10)
	r := New(testCache, "/tmp")

	ctx := context.Background()
	url := server.URL[len("http://"):]
	headers := http.Header{
		"User-Agent": {"Test-Agent"},
	}

	loadedImg, err := r.LoadImg(ctx, url, &headers)

	require.NoError(t, err)
	assert.NotNil(t, loadedImg)

	_, err = r.LoadImg(ctx, "bad.url", &headers)
	require.Error(t, err)
}

func createTestImage(width, height int) image.Image {
	img := image.NewNRGBA(image.Rect(0, 0, width, height))
	white := color.NRGBA{255, 255, 255, 255}
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, white)
		}
	}
	return img
}
