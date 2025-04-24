package server

import (
	"os"
	"testing"

	"github.com/dijer/otus-go-diploma/internal/cache"
	"github.com/dijer/otus-go-diploma/internal/config"
	"github.com/dijer/otus-go-diploma/internal/resizer"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestNewServer_Success(t *testing.T) {
	os.Setenv("PORT", "8080")
	os.Setenv("HOST", "localhost")
	os.Setenv("CACHE_DIR", "tmp")
	os.Setenv("CACHE_SIZE", "1024")

	config, err := config.New()
	require.NoError(t, err)

	imageCache := cache.New(config.Cache.Size)
	resizer := resizer.New(imageCache, config.Cache.Dir)

	server := New(config.Server, resizer)

	require.NotNil(t, server)
}

func TestServer_Start(t *testing.T) {
	MockImage := new(resizer.MockResizer)

	config := config.ServerConfig{
		Port: 8085,
		Host: "localhost",
	}

	server := New(config, MockImage)

	require.NotNil(t, server)

	MockImage.On("ResizeImg", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return("dest", nil)

	go func() {
		err := server.Start()
		require.NoError(t, err)
	}()

	<-server.Started()
	require.True(t, true)
}
