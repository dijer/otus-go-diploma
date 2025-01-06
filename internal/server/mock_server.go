package server

import (
	"net/http"

	"github.com/stretchr/testify/mock"
)

type MockServer struct {
	mock.Mock
	startedChan chan struct{}
}

func NewMockServer() *MockServer {
	return &MockServer{
		startedChan: make(chan struct{}),
	}
}

func (m *MockServer) Start() error {
	args := m.Called()

	close(m.startedChan)

	return args.Error(0)
}

func (m *MockServer) ResizerHandler(http.ResponseWriter, *http.Request) {
}

func (m *MockServer) Started() <-chan struct{} {
	return m.startedChan
}
