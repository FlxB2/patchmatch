package main

import (
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type hunk struct {
	preImage  lineIndex
	postImage lineIndex

	// e.g.
	// diff --git a/file.txt b/file.txt
	// index abcdef1..1234567 100644
	// --- a/file.txt
	// +++ b/file.txt
	header string

	// e.g.
	// @@ -1,4 +1,4 @@
	//  Line 1
	// +Line 2
	// -Line 3
	content string
}

type lineIndex struct {
	start         int
	linesIncluded int
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Error: Nothing to match against. Usage: ./patchmatch <regex>")
		os.Exit(1)
	}

	regexPattern := os.Args[1]
	regex, err := regexp.Compile(regexPattern)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error compiling regex:", err)
		os.Exit(1)
	}

	input, err := io.ReadAll(os.Stdin)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error reading stdin:", err)
		os.Exit(1)
	}
	stdIn := strings.ReplaceAll(string(input), "\r\n", "\n")

	hunks := SplitHunks(stdIn)

	for _, h := range hunks {
		converted := ConvertHunk(h, regex)
		if converted != nil {
			fmt.Println(converted.Str())
		}
	}
}

func SplitHunks(content string) []hunk {
	content = strings.Trim(content, "\n ")
	idx := strings.Index(content, "@@")

	if idx == -1 {
		return []hunk{}
	}

	hunks := []hunk{}
	currentHunkHeader := ""
	inHunk := false
	splitted := strings.Split(content, "\n")
	for i, line := range splitted {
		if i != len(splitted)-1 {
			line += "\n"
		}

		newHunk := false
		if strings.HasPrefix(line, "@@") {
			newHunk = true
			inHunk = true

			re := regexp.MustCompile("[0-9]+")
			digits := re.FindAllString(line, -1)

			intSlice := make([]int, len(digits))
			for i, digit := range digits {
				intSlice[i], _ = strconv.Atoi(digit)
			}

			hunks = append(hunks, hunk{
				header: strings.TrimSuffix(currentHunkHeader, "\n"),
				preImage: lineIndex{
					start:         intSlice[0],
					linesIncluded: intSlice[1],
				},
				postImage: lineIndex{
					start:         intSlice[2],
					linesIncluded: intSlice[3],
				},
				content: "",
			})
			currentHunkHeader = ""
		} else if !inHunk {
			currentHunkHeader += line
		}

		if inHunk && !newHunk {
			hunks[len(hunks)-1].content += line
		}
	}

	return hunks
}

func ConvertHunk(h hunk, regex *regexp.Regexp) *hunk {
	// naively search the whole diff, if we don't match
	// anything right away just return
	// maybe remove this? might be unnecessary
	if regex.FindString(h.content) == "" {
		return &h
	}

	resultingLines := []string{}
	lastDiffStart := -1
	inRemovingDiff := false
	for _, line := range strings.Split(h.content, "\n") {
		inDiff := isChangeLine(line)
		if inDiff {
			if lastDiffStart == -1 {
				lastDiffStart = len(resultingLines) - 1
			}

			if !inRemovingDiff && regex.FindString(line) != "" {
				// This is way too drastic and will remove too many
				// lines :) - might fix
				// especially for cases with lines like: + + - + - - etc
				// we just remove all connected differences w/o space in between
				for j := len(resultingLines) - 1; j >= lastDiffStart; j-- {
					if resultingLines[j][0] == '-' {
						h.postImage.linesIncluded += 1
						resultingLines[j] = " " + resultingLines[j][1:]
					} else if resultingLines[j][0] == '+' {
						h.postImage.linesIncluded -= 1
						resultingLines = resultingLines[:len(resultingLines) - 1]
					}
				}

				inRemovingDiff = true
			}

			if inRemovingDiff {
				if line[0] == '-' {
					// keep - lines to restore previous line
					h.postImage.linesIncluded += 1
					newLine := " " + line[1:]
					resultingLines = append(resultingLines, newLine)
				} else if line[0] == '+' {
					// remove + lines
					h.postImage.linesIncluded -= 1
				}
			} else {
				resultingLines = append(resultingLines, line)
			}

		} else {
			lastDiffStart = -1
			inRemovingDiff = false
			resultingLines = append(resultingLines, line)
		}
	}

	h.content = strings.Join(resultingLines, "\n")

	if h.empty() {
		return nil
	}

	return &h
}

func (h hunk) Str() string {
	hunk := fmt.Sprintf(
		"@@ -%d,%d +%d,%d @@\n%s",
		h.preImage.start, h.preImage.linesIncluded,
		h.postImage.start, h.postImage.linesIncluded,
		h.content)
	if h.header != "" {
		hunk = fmt.Sprintf("%s\n%s", h.header, hunk)
	}
	return hunk
}

func (h hunk) empty() bool {
	for _, line := range strings.Split(h.content, "\n") {
		if isChangeLine(line) {
			return false
		}
	}
	return true
}

func isChangeLine(s string) bool {
	return len(s) > 0 && (s[0] == '-' || s[0] == '+')
}
