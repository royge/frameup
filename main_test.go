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
		"qwerty/12345/100x200-test-2.jpg",
		"qwerty/12345/140x205-test-2.jpg",
		"qwerty/12345/150x105-test-2.jpg",
		"qwerty/12345/100x200-test-3.jpg",
		"qwerty/12345/140x205-test-3.jpg",
		"qwerty/12345/150x105-test-3.jpg",
	}
	keys := []string{
		"100x200",
		"140x205",
		"150x105",
	}
	expected := map[string][]string{
		"100x200": []string{
			"qwerty/12345/100x200-test.jpg",
			"qwerty/12345/100x200-test-2.jpg",
			"qwerty/12345/100x200-test-3.jpg",
		},
		"140x205": []string{
			"qwerty/12345/140x205-test.jpg",
			"qwerty/12345/140x205-test-2.jpg",
			"qwerty/12345/140x205-test-3.jpg",
		},
		"150x105": []string{
			"qwerty/12345/150x105-test.jpg",
			"qwerty/12345/150x105-test-2.jpg",
			"qwerty/12345/150x105-test-3.jpg",
		},
	}

	actual := classify(input, keys)
	r.Equal(expected, actual)
}
