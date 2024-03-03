package main

import (
	"fmt"
	"strings"
)

type Dimensions struct {
	Length int
	Width  int
}

func NewDimensions(str string) Dimensions {
	str = strings.Trim(str, "()")

	var length, width int
	fmt.Sscanf(str, "%d,%d", &length, &width)

	return Dimensions{length, width}
}
