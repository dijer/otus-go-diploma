package server

import (
	"fmt"
	"net"
	"net/http"
	"strconv"

	"github.com/dijer/otus-go-diploma/internal/config"
	"github.com/dijer/otus-go-diploma/internal/resizer"
	"github.com/stretchr/testify/mock"
)

type Server interface {
	Start() error
	ResizerHandler(w http.ResponseWriter, r *http.Request)
}

type HTTPServer struct {
	address string
	resizer resizer.Resizer
}

func NewServer(config config.ServerConfig, resizer resizer.Resizer) Server {
	address := net.JoinHostPort(config.Host, strconv.Itoa(int(config.Port)))

	return &HTTPServer{
		address: address,
		resizer: resizer,
	}
}

type MockServer struct {
	mock.Mock
}

func (m *MockServer) Start() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockServer) ResizerHandler(http.ResponseWriter, *http.Request) {
}

func (s *HTTPServer) Start() error {
	fmt.Println("Server run! at", s.address)

	mux := http.NewServeMux()

	mux.HandleFunc("/fill/", s.ResizerHandler)
	return http.ListenAndServe(s.address, mux)
}

func (s *HTTPServer) ResizerHandler(w http.ResponseWriter, r *http.Request) {
	resizedImgPath, isCached, err := s.resizer.ResizeImg(r.URL.Path, r.Header)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if isCached {
		w.Header().Set("X-Cache", "HIT")
	}

	http.ServeFile(w, r, *resizedImgPath)
}
