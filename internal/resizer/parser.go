package resizer

import (
	"errors"
	"strconv"
	"strings"
)

type parseData struct {
	Width, Height int64
	URL           string
}

var (
	ErrNotSetDimensions = errors.New("not set width, height or url")
)

func parseURL(path string) (data *parseData, err error) {
	trimmed := strings.ReplaceAll(path, "/fill/", "")
	parts := strings.SplitN(trimmed, "/", 3)

	if len(parts) < 3 {
		err = ErrNotSetDimensions
		return
	}

	width, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return
	}

	height, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return
	}

	url := parts[2]

	data = &parseData{
		Width:  width,
		Height: height,
		URL:    url,
	}
	return
}
