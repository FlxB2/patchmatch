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
	input, err := io.ReadAll(os.Stdin)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error reading stdin:", err)
		return
	}

	// normalize new lines
	stdIn := strings.ReplaceAll(string(input), "\r\n", "\n")

	hunks := SplitHunks(stdIn)

	for _, h := range hunks {
		converted, err := ConvertHunk(h, "abc")

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Println(converted.Str())
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

func ConvertHunk(h hunk, contains string) (hunk, error) {
	// indexes are just preimage:
	// linesIncluded = <content length> - <"+" lines>
	// postimage: linesIncluded: <content length> - <"-" lines>

	regex, err := regexp.Compile(contains)

	if err != nil {
		return hunk{}, fmt.Errorf("could not compile %s", contains)
	}

	newContent := ""
	for _, line := range strings.Split(h.content, "\n") {
		if regex.FindString(line) != "" {
			if line[0] == '-' {
				h.postImage.linesIncluded += 1
			} else if line[0] == '+' {
				h.postImage.linesIncluded -= 1
			}

			after := " " + line[1:]
			newContent += "\n" + after
		} else {
			newContent += "\n" + line
		}
	}

	h.content = strings.TrimPrefix(newContent, "\n")
	return h, nil
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
