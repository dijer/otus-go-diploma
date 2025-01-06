package resizer

import (
	"context"
	"image"
	"net/http"

	"github.com/stretchr/testify/mock"
)

type MockResizer struct {
	mock.Mock
}

func (m *MockResizer) GetCached(hash string) (*string, bool) {
	args := m.Called(hash)
	if args.Get(0) == nil {
		return nil, args.Bool(1)
	}
	val := args.String(0)
	return &val, args.Bool(1)
}

func (m *MockResizer) ResizeImg(img image.Image, width, height int64, hash string) (image.Image, error) {
	args := m.Called(img, width, height, hash)
	resizedImg, _ := args.Get(0).(image.Image)
	return resizedImg, args.Error(1)
}

func (m *MockResizer) LoadImg(ctx context.Context, url string, header *http.Header) (image.Image, error) {
	args := m.Called(ctx, url, header)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(image.Image), args.Error(1)
}

func (m *MockResizer) CreateHash(width, height int64, url string) string {
	args := m.Called(width, height, url)
	return args.String(0)
}
