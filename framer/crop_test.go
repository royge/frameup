package framer

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCrop(t *testing.T) {
	var (
		dir      = "../testdata/output/"
		s        = "../testdata/sub/sample.jpg"
		r        = require.New(t)
		w        = 200
		h        = 150
		expected = fmt.Sprintf("%ssub/%dx%d-sample.jpg", dir, w, h)
	)

	if err := Crop(s, dir, w, h); err != nil {
		t.Fatalf("error cropping image %s: %v", s, err)
	}

	_, err := os.Stat(expected)
	r.Empty(err)

	os.Remove(expected)
}
