package framer

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// Dimension of the picture.
type Dimension struct {
	Width  int
	Height int
}

// Dimensions is a list of Dimension.
type Dimensions []Dimension

// Keys of dimensions.
func (d Dimensions) Keys() []string {
	s := make([]string, len(d))
	for i, v := range d {
		s[i] = fmt.Sprintf("%dx%d", v.Width, v.Height)
	}

	return s
}

// ParseDimension from a string input.
func ParseDimension(s string) (*Dimension, error) {
	p := strings.Split(s, "x")
	if len(p) < 2 {
		return nil, errors.New("invalid dimension")
	}

	d := Dimension{}
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

// Classify files by keys.
func Classify(files []string, keys []string) map[string]string {
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
