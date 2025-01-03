package app

import (
	"errors"
	"net/http"
	"testing"

	"github.com/dijer/otus-go-diploma/internal/config"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockServer struct {
	mock.Mock
}

func (m *MockServer) Start() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockServer) ResizerHandler(http.ResponseWriter, *http.Request) {
}

func TestNewApp_Success(t *testing.T) {
	cfg := config.Config{
		Cache: config.CacheConfig{
			Size: 1024,
			Dir:  "tmp",
		},
		Server: config.ServerConfig{
			Port: 8080,
			Host: "localhost",
		},
	}

	app := NewApp(cfg)

	require.NotNil(t, app)
	require.NotNil(t, app.cache)
	require.NotNil(t, app.server)
}

func TestApp_RunError(t *testing.T) {
	cfg := config.Config{
		Cache: config.CacheConfig{
			Size: 1024,
			Dir:  "tmp",
		},
		Server: config.ServerConfig{
			Port: 8080,
			Host: "localhost",
		},
	}

	mockServer := new(MockServer)
	mockServer.On("Start").Return(errors.New("server err"))

	app := NewApp(cfg)
	app.server = mockServer

	err := app.Run()

	require.Error(t, err)
	require.Equal(t, "server err", err.Error())
}
