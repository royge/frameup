package main

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

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
				fmt.Println(info.Name())
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
