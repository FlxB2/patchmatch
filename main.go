package main

import (
	"fmt"
	"io"
	"os"
)

func main() {
	input, err := io.ReadAll(os.Stdin)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error reading stdin:", err)
		return
	}

	fmt.Println("Read all input:", ParsePatchFile(string(input)))
}

func ParsePatchFile(content string) string {
	return content
}
