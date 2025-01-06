package resizer

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"image"
	"net/http"
	"path/filepath"

	"github.com/dijer/otus-go-diploma/internal/cache"
	"github.com/disintegration/imaging"
)

type Resizer interface {
	GetCached(hash string) (*string, bool)
	ResizeImg(img image.Image, width, height int64, hash string) (image.Image, error)
	LoadImg(ctx context.Context, url string, header *http.Header) (image.Image, error)
	CreateHash(width, height int64, url string) string
}

type resizer struct {
	cache    cache.Cache
	CacheDir string
}

func New(cache cache.Cache, cacheDir string) Resizer {
	return &resizer{
		cache:    cache,
		CacheDir: cacheDir,
	}
}

func (r *resizer) GetCached(hash string) (*string, bool) {
	if val, ok := r.cache.Get(cache.Key(hash)); ok {
		return val.(*string), true
	}

	return nil, false
}

func (r *resizer) ResizeImg(img image.Image, width, height int64, hash string) (image.Image, error) {
	dest := filepath.Join(r.CacheDir, hash+".jpg")
	resizedImg := imaging.Fill(img, int(width), int(height), imaging.Center, imaging.Lanczos)
	r.cache.Set(cache.Key(hash), &dest)
	err := imaging.Save(resizedImg, dest)
	if err != nil {
		return nil, err
	}

	return resizedImg, nil
}

func (r *resizer) LoadImg(ctx context.Context, url string, header *http.Header) (image.Image, error) {
	request, err := http.NewRequestWithContext(ctx, "GET", "http://"+url, nil)
	if err != nil {
		return nil, err
	}

	for key, values := range *header {
		for _, value := range values {
			request.Header.Add(key, value)
		}
	}

	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("failed load img")
	}

	img, _, err := image.Decode(resp.Body)
	if err != nil {
		return nil, err
	}

	return img, nil
}

func (r *resizer) CreateHash(width, height int64, url string) string {
	hash := sha256.New()
	hash.Write([]byte(fmt.Sprintf("%s_%d_%d", url, width, height)))
	return hex.EncodeToString(hash.Sum(nil))
}
