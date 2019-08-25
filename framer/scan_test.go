package framer

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestScanDir(t *testing.T) {
	var (
		wg       sync.WaitGroup
		mu       sync.Mutex
		inputDir = "../testdata/input/"
		c        = make(chan string, 1)
		expected = []string{
			"../testdata/input/folder1",
			"../testdata/input/folder2",
			"../testdata/input/folder3",
			"../testdata/input/folder4",
			"../testdata/input/folder5",
		}
		actual  = []string{}
		r       = require.New(t)
		scanner = Scanner{}
	)

	go func() {
		defer close(c)
		err := scanner.ScanDir(&wg, inputDir, c)
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
		dir      = "../testdata/input/folder1/"
		c        = make(chan string, 1)
		expected = []string{
			"../testdata/input/folder1/file01.jpg",
			"../testdata/input/folder1/file02.jpg",
			"../testdata/input/folder1/file03.jpg",
		}
		actual  = []string{}
		r       = require.New(t)
		scanner = Scanner{}
	)

	go func() {
		defer close(c)
		err := scanner.Scan(&wg, dir, c, ".jpg")
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
