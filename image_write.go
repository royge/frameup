package main

import (
	"image"
	"image/jpeg"
	"os"
	"path/filepath"
)

func writeImage(img image.Image, name string) error {
	dir := filepath.Dir(name)
	os.MkdirAll(dir, os.ModePerm)

	fso, err := os.Create(name)
	if err != nil {
		return err
	}
	defer fso.Close()

	return jpeg.Encode(fso, img, &jpeg.Options{Quality: 100})
}
