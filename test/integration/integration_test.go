//go:build integration
// +build integration

package main

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCache(t *testing.T) {
	req1, err := http.NewRequest("GET", "http://0.0.0.0:8081/fill/30/30/imgserver/images/_gopher_original_1024x504.jpg", nil)
	require.NoError(t, err)

	req1.Header.Set("Cache-Control", "no-cache, no-store, must-revalidate")
	req1.Header.Set("Pragma", "no-cache")
	req1.Header.Set("Expires", "0")

	client := &http.Client{}
	resp1, err := client.Do(req1)
	require.NoError(t, err)
	defer resp1.Body.Close()

	require.Equal(t, http.StatusOK, resp1.StatusCode)

	req2, err := http.NewRequest("GET", "http://0.0.0.0:8081/fill/30/30/imgserver/images/_gopher_original_1024x504.jpg", nil)
	require.NoError(t, err)

	req2.Header.Set("Cache-Control", "no-cache, no-store, must-revalidate")
	req2.Header.Set("Pragma", "no-cache")
	req2.Header.Set("Expires", "0")

	resp2, err := client.Do(req2)
	require.NoError(t, err)
	defer resp2.Body.Close()

	require.Equal(t, http.StatusOK, resp2.StatusCode)

	require.Equal(t, resp2.Header.Get("X-Cache"), "HIT")
}

func TestServerNotExists(t *testing.T) {
	resp, _ := http.Get("http://0.0.0.0:8081/fill/30/30/something/images/_gopher_original_1024x504.jpg")

	require.Equal(t, resp.StatusCode, http.StatusNotFound)
}

func TestImgNotExists(t *testing.T) {
	resp, _ := http.Get("http://0.0.0.0:8081/fill/30/30/imgserver/images/something.jpg")

	require.Equal(t, resp.StatusCode, http.StatusNotFound)
}

func TestImgNotImage(t *testing.T) {
	resp, _ := http.Get("http://0.0.0.0:8081/fill/30/30/imgserver/images/test.txt")

	require.Equal(t, resp.StatusCode, http.StatusNotFound)
}

func TestServerErrParse(t *testing.T) {
	resp, _ := http.Get("http://0.0.0.0:8081////imgserver/images/_gopher_original_1024x504.jpg")

	require.Equal(t, resp.StatusCode, http.StatusNotFound)
}

func TestServerReturnsImg(t *testing.T) {
	resp, _ := http.Get("http://0.0.0.0:8081/fill/30/30/imgserver/images/_gopher_original_1024x504.jpg")

	require.Equal(t, resp.StatusCode, http.StatusOK)
	require.Equal(t, resp.Header.Get("Content-Type"), "image/jpeg")
}

func TestOriginalImageLessResizeImage(t *testing.T) {
	resp, _ := http.Get("http://0.0.0.0:8081/fill/1050/550/imgserver/images/_gopher_original_1024x504.jpg")

	require.Equal(t, resp.StatusCode, http.StatusOK)
	require.Equal(t, resp.Header.Get("Content-Type"), "image/jpeg")
}
