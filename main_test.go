package main

import (
	"fmt"
	"os"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestScanDir(t *testing.T) {
	var (
		wg       sync.WaitGroup
		mu       sync.Mutex
		inputDir = "./testdata/input/"
		c        = make(chan string, 2)
		expected = []string{
			"testdata/input/folder1",
			"testdata/input/folder2",
			"testdata/input/folder3",
		}
		actual = []string{}
		r      = require.New(t)
	)

	go func() {
		defer close(c)
		err := scanDir(&wg, inputDir, c)
		if err != nil {
			t.Fatalf("error scanning directory: %v", err)
		}
	}()

	for v := range c {
		go func(s string) {
			mu.Lock()
			actual = append(actual, s)
			mu.Unlock()
			wg.Done()
		}(v)
	}

	wg.Wait()

	r.Equal(expected, actual)
}

func TestScan(t *testing.T) {
	var (
		wg       sync.WaitGroup
		mu       sync.Mutex
		dir      = "./testdata/input/folder1/"
		c        = make(chan string, 2)
		expected = []string{
			"testdata/input/folder1/file01.jpg",
			"testdata/input/folder1/file02.jpg",
			"testdata/input/folder1/file03.jpg",
		}
		actual = []string{}
		r      = require.New(t)
	)

	go func() {
		defer close(c)
		err := scan(&wg, dir, c, ".jpg")
		if err != nil {
			t.Fatalf("error scanning files: %v", err)
		}
	}()

	for v := range c {
		go func(s string) {
			mu.Lock()
			actual = append(actual, s)
			mu.Unlock()
			wg.Done()
		}(v)
	}

	wg.Wait()

	r.Equal(expected, actual)
}

func TestCrop(t *testing.T) {
	var (
		dir      = "./testdata/output/"
		s        = "./testdata/sub/sample.jpg"
		r        = require.New(t)
		w        = 200
		h        = 150
		expected = fmt.Sprintf("%ssub/%dx%d-sample.jpg", dir, w, h)
	)

	if err := crop(s, dir, w, h); err != nil {
		t.Fatalf("error cropping image %s: %v", s, err)
	}

	_, err := os.Stat(expected)
	r.Empty(err)

	os.Remove(expected)
}
