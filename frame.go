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

func frame(m map[string]string, outDir string, bg string, overlay string) error {
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

	iom := map[string]io.Reader{}

	for i, v := range m {
		o, err := os.Open(v)
		if err != nil {
			return fmt.Errorf("unable to open overlay image, %v", err)
		}
		defer o.Close()
		iom[i] = o
	}

	err = addFrame(iom, f, o, w)
	if err != nil {
		return fmt.Errorf("expected no error, got %v", err)
	}

	return nil
}

func addFrame(m map[string]io.Reader, bg, ol io.Reader, w io.Writer) error {
	img, err := jpeg.Decode(bg)
	if err != nil {
		return fmt.Errorf("unable to decode jpeg image: %v", err)
	}

	overlay, err := png.Decode(ol)
	if err != nil {
		return fmt.Errorf("unable to decode png image: %v", err)
	}

	b := img.Bounds()
	ib := image.NewRGBA(b)

	draw.Draw(ib, b, img, image.ZP, draw.Src)

	pt := image.ZP

	for k, v := range m {
		img, err := jpeg.Decode(v)
		if err != nil {
			return fmt.Errorf("unable to decode jpeg image: %v", err)
		}

		fmt.Println(pt)
		draw.Draw(ib, b, img, pt, draw.Over)

		dim, err := parseDimension(k)
		if err != nil {
			return fmt.Errorf("unable to parse dimension %s: %v", k, err)
		}

		fmt.Println(dim)

		pt = image.Pt(0, pt.Y-dim.Height)
	}

	draw.Draw(ib, b, overlay, image.ZP, draw.Over)

	jpeg.Encode(w, ib, &jpeg.Options{Quality: 100})

	return nil
}
