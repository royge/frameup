package main

import (
	"flag"
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"
	"image/png"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/muesli/smartcrop"
	"github.com/muesli/smartcrop/nfnt"
	"github.com/nfnt/resize"
)

var (
	overlay = "./assets/overlay.png"
	delay   = 100 * time.Millisecond
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

func scan(wg *sync.WaitGroup, dir string, c chan string, exts string) error {
	err := filepath.Walk(dir, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.HasSuffix(info.Name(), exts) {
			p := path.Join(dir, info.Name())
			if ok, _ := exists(p); ok {
				wg.Add(1)
				c <- path.Join(dir, info.Name())
			}
			time.Sleep(delay)
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

// exists returns whether the given file or directory exists
func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

func scanDir(wg *sync.WaitGroup, dir string, c chan string) error {
	err := filepath.Walk(dir, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			p := path.Join(dir, info.Name())
			if ok, _ := exists(p); ok {
				wg.Add(1)
				c <- path.Join(dir, info.Name())
			}
			time.Sleep(delay)
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func crop(wg *sync.WaitGroup, s string, dir string) error {
	defer wg.Done()
	if s != "" {
		r, err := os.Open(s)
		if err != nil {
			return fmt.Errorf("unable to open raw image, %v", err)
		}
		defer r.Close()

		name := filepath.Base(s)

		p := path.Join(
			dir,
			name,
		)

		if err := resizeAndCrop(r, 1200, 1800, p); err != nil {
			return fmt.Errorf("unable to crop image, %v", err)
		}
	}

	return nil
}

func frame(wg *sync.WaitGroup, s string, dir string) error {
	defer wg.Done()
	if s != "" {
		f, err := os.Open(s)
		if err != nil {
			return fmt.Errorf("unable to open raw image, %v", err)
		}
		defer f.Close()

		name := filepath.Base(s)

		p := path.Join(
			dir,
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

func main() {
	s := flag.String("s", "", "Source directory. (Ex. -s=/Users/roye/Desktop/mypics)")
	o := flag.String("o", "", "Output directory. (Ex. -o=/Users/roye/Desktop/output)")
	x := flag.String("x", ".jpg", "Picture files extensions.")

	flag.Parse()

	c := make(chan string, 4)

	var wg sync.WaitGroup

	go func() {
		if err := scan(&wg, *s, c, *x); err != nil {
			log.Fatalf("error scanning directory: %v", err)
		}

		wg.Wait()
		close(c)
	}()

	// for v := range c {
	// 	go func(s string) {
	// 		fmt.Printf("\ncropping %s...", s)
	// 		if err := crop(&wg, s, *o); err != nil {
	// 			fmt.Printf("\nerror cropping picture: %v", err)
	// 		} else {
	// 			fmt.Printf("\n%s Done!", s)
	// 		}
	// 	}(v)
	// }

	for v := range c {
		go func(s string) {
			fmt.Printf("\nframing %s...", s)
			if err := frame(&wg, s, *o); err != nil {
				fmt.Printf("\nerror framing picture: %v", err)
			} else {
				fmt.Printf("\n%s Done!", s)
			}
		}(v)
	}

	fmt.Print("\n")
}
