package main

import (
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"
	"image/png"
	"io"
	"os"
	"path"
	"path/filepath"
)

func combine(count int) []int {
	return nil
}

func frame(m map[string][]string, outDir string, bg string, overlay string) error {
	if bg != "" {
		f, err := os.Open(bg)
		if err != nil {
			return fmt.Errorf("unable to open raw image, %v", err)
		}
		defer f.Close()

		name := filepath.Base(bg)

		p := path.Join(
			outDir,
			name,
		)

		o, err := os.Open(overlay)
		if err != nil {
			return fmt.Errorf("unable to open overlay image, %v", err)
		}
		defer o.Close()

		w, err := os.Create(p)
		if err != nil {
			return fmt.Errorf("unable to create watermarked file, %v", err)
		}
		defer w.Close()

		err = addFrame(f, o, w)
		if err != nil {
			return fmt.Errorf("expected no error, got %v", err)
		}
	}

	return nil
}

func addFrame(r, o io.Reader, w io.Writer) error {
	img, err := jpeg.Decode(r)
	if err != nil {
		return fmt.Errorf("unable to decode jpeg image: %v", err)
	}

	overlay, err := png.Decode(o)
	if err != nil {
		return fmt.Errorf("unable to decode png image: %v", err)
	}

	// offset := getOffset(img, overlay)

	b := img.Bounds()
	m := image.NewRGBA(b)

	draw.Draw(m, b, img, image.ZP, draw.Src)
	draw.Draw(
		m,
		b,
		overlay,
		image.ZP,
		draw.Over,
	)

	jpeg.Encode(w, m, &jpeg.Options{Quality: 100})

	return nil
}

func getOffset(img, overlay image.Image) image.Point {
	b := img.Bounds().Max
	o := overlay.Bounds().Max

	top := float64(b.Y)*.5 - float64(o.Y)*.5
	left := float64(b.X)*.5 - float64(o.X)*.5
	return image.Pt(int(left), int(top))
}

func join(r, o io.Reader, w io.Writer) error {
	img, err := jpeg.Decode(r)
	if err != nil {
		return fmt.Errorf("unable to decode jpeg image: %v", err)
	}

	overlay, err := png.Decode(o)
	if err != nil {
		return fmt.Errorf("unable to decode png image: %v", err)
	}

	// offset := getOffset(img, overlay)

	b := img.Bounds()
	m := image.NewRGBA(b)

	draw.Draw(m, b, img, image.ZP, draw.Src)
	draw.Draw(
		m,
		b,
		overlay,
		image.ZP,
		draw.Over,
	)

	jpeg.Encode(w, m, &jpeg.Options{Quality: 100})

	return nil
}
