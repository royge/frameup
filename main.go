package main

import (
	"flag"
	"fmt"
	"log"
	"sync"
	"time"
)

var (
	overlay = "./assets/overlay.png"
	delay   = 100 * time.Millisecond
)

type dimension struct {
	Width  int
	Height int
}

func main() {
	src := flag.String("s", "", "Source directory. (Ex. -s=/Users/roye/Desktop/mypics)")
	dst := flag.String("o", "", "Output directory. (Ex. -o=/Users/roye/Desktop/output)")
	ext := flag.String("x", ".jpg", "Picture files extensions.")

	flag.Parse()

	var (
		inDirChan  = make(chan string, 1)
		outDirChan = make(chan string, 1)
		dirWg      sync.WaitGroup
		fileWg     sync.WaitGroup
		dimensions = []dimension{
			dimension{500, 600},
			dimension{300, 400},
			dimension{400, 100},
		}
	)

	go func() {
		defer close(inDirChan)
		if err := scanDir(&dirWg, *src, inDirChan); err != nil {
			log.Fatalf("error scanning %s directory: %v", *src, err)
		}
	}()

	for v := range inDirChan {
		c := make(chan string, 1)
		go func(dir string) {
			defer close(c)
			if err := scan(&fileWg, dir, c, *ext); err != nil {
				log.Fatalf("error scanning %s directory: %v", dir, err)
			}
		}(v)

		for w := range c {
			go func(file string) {

				for _, d := range dimensions {
					func(d dimension) {
						if err := crop(file, *dst, d.Width, d.Height); err != nil {
							fmt.Printf("error cropping picture file %s: %v", file, err)
						}
					}(d)
				}

				fileWg.Done()
			}(w)
		}

		fileWg.Wait()

		// Done scanning 1 directory.
		dirWg.Done()
	}

	dirWg.Wait()

	go func() {
		defer close(outDirChan)
		if err := scanDir(&dirWg, *dst, outDirChan); err != nil {
			log.Fatalf("error scanning %s directory: %v", *src, err)
		}
	}()

	for v := range outDirChan {
		c := make(chan string, 1)
		go func(dir string) {
			defer close(c)
			if err := scan(&fileWg, dir, c, *ext); err != nil {
				log.Fatalf("error scanning %s directory: %v", dir, err)
			}
		}(v)

		for w := range c {
			go func(file string) {
				fileWg.Done()
				fmt.Println("file:", file)
			}(w)
		}

		fileWg.Wait()

		// Done scanning 1 directory.
		dirWg.Done()
	}

	dirWg.Wait()
}
