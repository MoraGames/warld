package main

import (
	"bufio"
	"os"
	"strings"
)

var (
	Scanner = bufio.NewReader(os.Stdin)
)

func Scan(variable *string) error {
	// Read the input line using Scanner
	str, err := Scanner.ReadString('\n')
	if err != nil {
		return err
	}

	// Trim the line ending from the string and assign it to the variable
	*variable = strings.Trim(str, "\r\n")
	return nil
}
