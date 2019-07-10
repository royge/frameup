package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestClassify(t *testing.T) {
	r := require.New(t)

	input := []string{
		"qwerty/12345/100x200-test.jpg",
		"qwerty/12345/140x205-test.jpg",
		"qwerty/12345/150x105-test.jpg",
	}
	keys := []string{
		"100x200",
		"140x205",
		"150x105",
	}
	expected := map[string]string{
		"100x200": "qwerty/12345/100x200-test.jpg",
		"140x205": "qwerty/12345/140x205-test.jpg",
		"150x105": "qwerty/12345/150x105-test.jpg",
	}

	actual := classify(input, keys)
	r.Equal(expected, actual)
}

func TestDimensionsKeys(t *testing.T) {
	r := require.New(t)

	dims := dimensions([]dimension{
		dimension{500, 600},
		dimension{300, 400},
		dimension{400, 100},
	})

	expected := []string{
		"500x600",
		"300x400",
		"400x100",
	}

	actual := dims.Keys()
	r.Equal(expected, actual)
}

func TestParseDimension(t *testing.T) {
	r := require.New(t)

	dim, err := parseDimension("100x200")
	r.Empty(err)
	r.Equal(100, dim.Width)
	r.Equal(200, dim.Height)

	dim, err = parseDimension("100by200")
	r.NotEmpty(err)

	dim, err = parseDimension("100-200")
	r.NotEmpty(err)

	dim, err = parseDimension("100x")
	r.NotEmpty(err)

	dim, err = parseDimension("x100")
	r.NotEmpty(err)
}
