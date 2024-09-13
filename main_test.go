package main

import (
	"testing"
)

func TestParsePatchFile(t *testing.T) {
	input := "Test"
	result := ParsePatchFile(input)

	if result != input {
		t.Errorf("wrong output %s for output %s", result, input)
	}
}
