package framer

import (
	"fmt"
	"image"
	"io"
	"os"
	"path"
	"path/filepath"

	"github.com/muesli/smartcrop"
	"github.com/muesli/smartcrop/nfnt"
	"github.com/nfnt/resize"
)

// Crop pictures inside a directory.
func Crop(s string, dir string, w, h int) error {
	if s != "" {
		r, err := os.Open(s)
		if err != nil {
			return fmt.Errorf("unable to open %s from %s, %v", s, dir, err)
		}
		defer r.Close()

		name := filepath.Base(s)
		fdir := filepath.Base(filepath.Dir(s))
		p := path.Join(
			dir,
			fdir,
			fmt.Sprintf("%dx%d-%s", w, h, name),
		)

		if err := resizeAndCrop(r, w, h, p); err != nil {
			return fmt.Errorf("unable to crop image, %v", err)
		}
	}

	return nil
}

func resizeAndCrop(f io.Reader, width, height int, name string) error {
	img, _, err := image.Decode(f)
	if err != nil {
		return fmt.Errorf("failed to decode image, %v", err)
	}

	analyzer := smartcrop.NewAnalyzer(
		nfnt.NewDefaultResizer(),
	)
	topCrop, err := analyzer.FindBestCrop(img, width, height)
	if err != nil {
		return fmt.Errorf("failed to find best crop, %v", err)
	}

	type SubImager interface {
		SubImage(r image.Rectangle) image.Image
	}

	croppedImg := img.(SubImager).SubImage(topCrop)
	croppedImg = resize.Resize(uint(width), uint(height), croppedImg, resize.Bicubic)

	return writeImage(croppedImg, name)
}
