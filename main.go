package main

import (
	"fmt"
	"github.com/muesli/smartcrop"
	"github.com/muesli/smartcrop/nfnt"
	"github.com/nfnt/resize"
	"image"
	"image/draw"
	"image/jpeg"
	"image/png"
	"io"
	"log"
	"os"
)

func writeImage(img image.Image, name string) error {
	fso, err := os.Create(name)
	if err != nil {
		return err
	}
	defer fso.Close()

	return jpeg.Encode(fso, img, &jpeg.Options{Quality: 100})
}

func getOffset(img, overlay image.Image) image.Point {
	b := img.Bounds().Max
	o := overlay.Bounds().Max

	top := float64(b.Y)*.5 - float64(o.Y)*.5
	left := float64(b.X)*.5 - float64(o.X)*.5
	return image.Pt(int(left), int(top))
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

func crop(f io.Reader, width, height int, name string) error {
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

var (
	overlay = "./assets/overlay.png"
	sample  = "./test/dummy.jpg"
	framed  = "./test/framed.jpg"
	raw     = "./test/raw.jpg"
	cropped = "./test/cropped.jpg"
)

func main() {
	f, err := os.Open(sample)
	if err != nil {
		log.Fatalf("unable to open test image, %v", err)
	}
	defer f.Close()

	o, err := os.Open(overlay)
	if err != nil {
		log.Fatalf("unable to open overlay image, %v", err)
	}
	defer o.Close()

	w, err := os.Create(framed)
	if err != nil {
		log.Fatalf("unable to create watermarked file, %v", err)
	}
	defer w.Close()

	err = addFrame(f, o, w)
	if err != nil {
		log.Fatalf("expected no error, got %v", err)
	}
}
