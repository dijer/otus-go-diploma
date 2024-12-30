package server

import (
	"errors"
	"io"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/dijer/otus-go-diploma/internal/cache"
	"github.com/dijer/otus-go-diploma/internal/config"
	"github.com/dijer/otus-go-diploma/internal/resizer"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type HTTPClientInterface interface {
	Do(req *http.Request) (*http.Response, error)
}

type MockResizer struct {
	mock.Mock
	client HTTPClientInterface
}

func (m *MockResizer) ResizeImg(path string, headers http.Header) (*string, bool, error) {
	args := m.Called(path, headers)
	resizedImgPath := args.String(0)
	return &resizedImgPath, args.Bool(1), args.Error(2)
}

func (m *MockResizer) CreateHash(width, height int64, url string) string {
	args := m.Called(width, height, url)
	return args.String(0)
}

func (m *MockResizer) LoadImg(url, dest string, headers http.Header) error {
	req, err := http.NewRequest("GET", "http://"+url, nil)
	if err != nil {
		return err
	}

	for key, values := range headers {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}

	resp, err := m.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New("failed load img")
	}

	file, err := os.Create(dest)
	if err != nil {
		return errors.New(dest)
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	return err
}

func TestNewServer_Success(t *testing.T) {
	os.Setenv("PORT", "8080")
	os.Setenv("HOST", "localhost")
	os.Setenv("CACHE_DIR", "tmp")
	os.Setenv("CACHE_SIZE", "1024")

	config, err := config.NewConfig()
	require.NoError(t, err)

	imageCache := cache.NewCache(int(config.Cache.Size))
	resizer := resizer.NewResizer(imageCache, config.Cache.Dir)

	server := NewServer(config.Server, resizer)

	require.NotNil(t, server)
}

func TestServer_Start(t *testing.T) {
	mockResizer := new(MockResizer)

	config := config.ServerConfig{
		Port: 8080,
		Host: "localhost",
	}

	server := NewServer(config, mockResizer)

	require.NotNil(t, server)

	mockResizer.On("ResizeImg", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil, false, nil)

	go func() {
		err := server.Start()
		require.NoError(t, err)
	}()

	time.Sleep(500 * time.Millisecond)
	require.True(t, true)
}
