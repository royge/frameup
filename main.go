package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	bg      = "./assets/bg.jpg"
	overlay = "./assets/overlay.png"
	delay   = 100 * time.Millisecond
	dims    = []string{
		"500x600",
		"467x467",
		"555x555",
	}
)

type dimension struct {
	Width  int
	Height int
}

type dimensions []dimension

func (d dimensions) Keys() []string {
	s := make([]string, len(d))
	for i, v := range d {
		s[i] = fmt.Sprintf("%dx%d", v.Width, v.Height)
	}

	return s
}

func parseDimension(s string) (*dimension, error) {
	p := strings.Split(s, "x")
	if len(p) < 2 {
		return nil, errors.New("invalid dimension")
	}

	d := dimension{}
	var err error

	d.Width, err = strconv.Atoi(p[0])
	if err != nil {
		return nil, err
	}
	d.Height, err = strconv.Atoi(p[1])
	if err != nil {
		return nil, err
	}

	return &d, nil
}

func classify(files []string, keys []string) map[string]string {
	m := map[string]string{}
	for _, v := range keys {
		m[v] = ""
	}

	for _, v := range files {
		func(file string) {
			for i := range m {
				if strings.Contains(file, i) {
					m[i] = file
				}
			}
		}(v)
	}

	return m
}

func main() {
	cmd := flag.String("c", "crop", "Action")
	src := flag.String("s", "", "Source directory. (Ex. -s=/Users/roye/Desktop/mypics)")
	dst := flag.String("o", "", "Output directory. (Ex. -o=/Users/roye/Desktop/output)")
	ext := flag.String("x", ".jpg", "Picture files extensions.")

	flag.Parse()

	var (
		inDirChan  = make(chan string, 4)
		outDirChan = make(chan string, 4)
		dirWg      sync.WaitGroup
		fileWg     sync.WaitGroup
	)

	switch *cmd {
	case "crop":
		go func() {
			defer close(inDirChan)
			if err := scanDir(&dirWg, *src, inDirChan); err != nil {
				log.Fatalf("error scanning %s directory: %v", *src, err)
			}
		}()

		for v := range inDirChan {
			c := make(chan string, 4)
			go func(dir string) {
				defer close(c)
				if err := scan(&fileWg, dir, c, *ext); err != nil {
					log.Fatalf("error scanning %s directory: %v", dir, err)
				}
			}(v)

			for w := range c {
				go func(file string) {

					for _, v := range dims {
						d, _ := parseDimension(v)
						func(d dimension) {
							if err := crop(file, *dst, d.Width, d.Height); err != nil {
								fmt.Printf("error cropping picture file %s: %v", file, err)
							}
						}(*d)
					}

					fileWg.Done()
				}(w)
			}

			fileWg.Wait()

			// Done scanning 1 directory.
			dirWg.Done()
		}

		dirWg.Wait()
	case "frame":
		go func() {
			defer close(outDirChan)
			if err := scanDir(&dirWg, *dst, outDirChan); err != nil {
				log.Fatalf("error scanning %s directory: %v", *src, err)
			}
		}()

		for v := range outDirChan {
			c := make(chan string, 1)
			files := []string{}
			mu := sync.Mutex{}

			go func(dir string) {
				defer close(c)
				if err := scan(&fileWg, dir, c, *ext); err != nil {
					log.Fatalf("error scanning %s directory: %v", dir, err)
				}
			}(v)

			for f := range c {
				mu.Lock()
				files = append(files, f)
				mu.Unlock()
				fileWg.Done()
			}

			m := classify(files, dims)
			err := frame(m, *dst, bg, overlay)
			if err != nil {
				log.Fatalf("error creating frame: %v", err)
			}

			fileWg.Wait()

			// Done scanning 1 directory.
			dirWg.Done()
		}

		dirWg.Wait()
	}
}
