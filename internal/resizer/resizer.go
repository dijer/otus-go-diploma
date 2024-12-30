package resizer

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/dijer/otus-go-diploma/internal/cache"

	"github.com/disintegration/imaging"
)

type Resizer interface {
	ResizeImg(path string, header http.Header) (imgURL *string, isCached bool, err error)
	LoadImg(url, dest string, headers http.Header) error
	CreateHash(width, height int64, url string) string
}

type resizer struct {
	cache    cache.Cache
	cacheDir string
}

func NewResizer(cache cache.Cache, cacheDir string) Resizer {
	return &resizer{
		cache:    cache,
		cacheDir: cacheDir,
	}
}

func (r *resizer) ResizeImg(path string, header http.Header) (imgURL *string, isCached bool, err error) {
	isCached = false

	data, err := parseURL(path)
	if err != nil {
		return
	}

	width := data.Width
	height := data.Height
	url := data.URL

	if width == 0 {
		err = errors.New("zero width")
		return
	}

	if height == 0 {
		err = errors.New("zero height not allowed")
		return
	}

	if url == "" {
		err = errors.New("empty url not allowed")
		return
	}

	hash := r.CreateHash(width, height, url)

	if val, ok := r.cache.Get(cache.Key(hash)); ok {
		isCached = true
		imgURL = val.(*string)
		return
	}

	filePath := filepath.Join(r.cacheDir, hash+".jpg")
	err = r.LoadImg(url, filePath, header)
	if err != nil {
		return
	}

	img, err := imaging.Open(filePath)
	if err != nil {
		return
	}

	resizedImg := imaging.Fill(img, int(width), int(height), imaging.Center, imaging.Lanczos)
	err = imaging.Save(resizedImg, filePath)
	if err != nil {
		return
	}

	r.cache.Set(cache.Key(hash), &filePath)
	imgURL = &filePath
	return
}

func (r *resizer) LoadImg(url, dest string, headers http.Header) error {
	req, err := http.NewRequest("GET", "http://"+url, nil)
	if err != nil {
		return err
	}

	for key, values := range headers {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}

	client := &http.Client{}
	resp, err := client.Do(req)
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

func (r *resizer) CreateHash(width, height int64, url string) string {
	hash := md5.New()
	hash.Write([]byte(fmt.Sprintf("%s_%d_%d", url, width, height)))
	return hex.EncodeToString(hash.Sum(nil))
}
