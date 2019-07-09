package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCombine(t *testing.T) {
	r := require.New(t)
	count := 3
	expected := []int{
		0, 0, 0,
		0, 0, 1,
		0, 1, 1,
		1, 1, 1,
		1, 1, 0,
		1, 0, 0,
		2, 0, 0,
		2, 0, 1,
		2, 1, 1,
		2, 1, 0,
		2, 2, 0,
		2, 2, 1,
		2, 2, 2,
	}
	actual := combine(count)
	r.Equal(expected, actual)
}
