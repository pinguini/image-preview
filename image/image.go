package image

import (
	"bytes"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"

	"io"

	"github.com/nfnt/resize"
)

func init() {
	image.RegisterFormat("gif", "GIF8?a", gif.Decode, gif.DecodeConfig)
	image.RegisterFormat("jpeg", "ÿØÿá", jpeg.Decode, jpeg.DecodeConfig)
	image.RegisterFormat("png", "\x89PNG\r\n\x1a\n", png.Decode, png.DecodeConfig)
}

func ResizeImage(r io.Reader, newWidth, newHeight int) ([]byte, error) {
	img, _, err := image.Decode(r)
	if err != nil {
		return nil, fmt.Errorf("failed to decode image: %v", err)
	}

	resizedImg := resize.Thumbnail(uint(newWidth), uint(newHeight), img, resize.Lanczos3)
	buf := new(bytes.Buffer)
	err = jpeg.Encode(buf, resizedImg, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to encode image: %v", err)
	}

	return buf.Bytes(), nil
}
