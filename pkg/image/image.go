package image

import (
	"bytes"
	"errors"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"net/http"

	"github.com/nfnt/resize"
)

func Compress(data []byte, width, height uint) ([]byte, error) {
	contentType := http.DetectContentType(data)
	reader := bytes.NewReader(data)

	var (
		img image.Image
		err error
	)

	switch contentType {
	case "image/png":
		img, err = png.Decode(reader)
	case "image/jpeg":
		img, err = jpeg.Decode(reader)
	case "image/gif":
		img, err = gif.Decode(reader)
	default:
		return nil, errors.New("unsupport file type:" + contentType)
	}

	if err != nil {
		return nil, err
	}

	m := resize.Resize(width, height, img, resize.Lanczos3)
	buffer := bytes.NewBuffer([]byte{})

	err = jpeg.Encode(buffer, m, nil)
	if err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}
