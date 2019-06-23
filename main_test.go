package main

import (
	"log"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestScanDir(t *testing.T) {
	var (
		wg       sync.WaitGroup
		inputDir = "./test/input/"
		c        = make(chan string, 2)
		expected = []string{
			"test/input/folder1",
			"test/input/folder2",
			"test/input/folder3",
		}
		actual = []string{}
		r      = require.New(t)
	)

	go func() {
		defer close(c)
		err := scanDir(&wg, inputDir, c)
		if err != nil {
			log.Fatalf("error scanning directory: %v", err)
		}
	}()

	for v := range c {
		go func(s string) {
			actual = append(actual, s)
			wg.Done()
		}(v)
	}

	wg.Wait()

	r.Equal(expected, actual)
}

func TestScan(t *testing.T) {
	var (
		wg       sync.WaitGroup
		dir      = "./test/input/folder1/"
		c        = make(chan string, 2)
		expected = []string{
			"test/input/folder1/file01.jpg",
			"test/input/folder1/file02.jpg",
			"test/input/folder1/file03.jpg",
		}
		actual = []string{}
		r      = require.New(t)
	)

	go func() {
		defer close(c)
		err := scan(&wg, dir, c, ".jpg")
		if err != nil {
			log.Fatalf("error scanning files: %v", err)
		}
	}()

	for v := range c {
		go func(s string) {
			actual = append(actual, s)
			wg.Done()
		}(v)
	}

	wg.Wait()

	r.Equal(expected, actual)
}
