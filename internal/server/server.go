package server

import (
	"fmt"
	"image/jpeg"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/dijer/otus-go-diploma/internal/config"
	"github.com/dijer/otus-go-diploma/internal/resizer"
)

type Server interface {
	Start() error
	ResizerHandler(w http.ResponseWriter, r *http.Request)
	Started() <-chan struct{}
}

type HTTPServer struct {
	config      config.ServerConfig
	address     string
	resizer     resizer.Resizer
	StartedChan chan struct{}
}

func New(config config.ServerConfig, resizer resizer.Resizer) Server {
	address := net.JoinHostPort(config.Host, strconv.Itoa(config.Port))

	return &HTTPServer{
		address:     address,
		resizer:     resizer,
		StartedChan: make(chan struct{}),
		config:      config,
	}
}

func (s *HTTPServer) Start() error {
	fmt.Println("Server run! at", s.address)

	mux := http.NewServeMux()
	mux.HandleFunc("/fill/{width}/{height}/{url...}", s.ResizerHandler)

	server := &http.Server{
		Addr:         s.address,
		Handler:      mux,
		ReadTimeout:  20 * time.Second,
		WriteTimeout: 20 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		close(s.StartedChan)
	}()

	return server.ListenAndServe()
}

func (s *HTTPServer) ResizerHandler(w http.ResponseWriter, r *http.Request) {
	width, err := strconv.ParseInt(r.PathValue("width"), 10, 64)
	if err != nil {
		http.Error(w, "err parse width in pathname", http.StatusNotFound)
		return
	}

	height, err := strconv.ParseInt(r.PathValue("height"), 10, 64)
	if err != nil {
		http.Error(w, "err parse height in pathname", http.StatusNotFound)
		return
	}

	url := r.PathValue("url")

	hash := s.resizer.CreateHash(width, height, url)

	if val, ok := s.resizer.GetCached(hash); ok {
		w.Header().Set("X-Cache", "HIT")
		http.ServeFile(w, r, *val)
		return
	}

	img, err := s.resizer.LoadImg(r.Context(), url, &r.Header)
	if err != nil {
		http.Error(w, "err load image", http.StatusNotFound)
		return
	}

	resizedImg, err := s.resizer.ResizeImg(img, width, height, hash)
	if err != nil {
		http.Error(w, "failed to resize image", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "image/jpeg")
	err = jpeg.Encode(w, resizedImg, nil)
	if err != nil {
		http.Error(w, "failed to encode image", http.StatusNotFound)
	}
}

func (s *HTTPServer) Started() <-chan struct{} {
	return s.StartedChan
}
