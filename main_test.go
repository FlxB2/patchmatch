package main

import (
	"strings"
	"testing"
)

func TestHunkToString(t *testing.T) {
	tests := []struct {
		name     string
		hunk     hunk
		expected string
	}{
		{
			name: "Simple Hunk",
			hunk: hunk{
				preImage: lineIndex{
					start:         10,
					linesIncluded: 5,
				},
				postImage: lineIndex{
					start:         20,
					linesIncluded: 5,
				},
				content: "-old line\n+new line",
			},
			expected: "@@ -10,5 +20,5 @@\n-old line\n+new line",
		},
		{
			name: "Empty Content",
			hunk: hunk{
				preImage: lineIndex{
					start:         30,
					linesIncluded: 0,
				},
				postImage: lineIndex{
					start:         40,
					linesIncluded: 0,
				},
				content: "",
			},
			expected: "@@ -30,0 +40,0 @@\n",
		},
		{
			name: "Multi-line Content",
			hunk: hunk{
				preImage: lineIndex{
					start:         1,
					linesIncluded: 3,
				},
				postImage: lineIndex{
					start:         4,
					linesIncluded: 3,
				},
				content: "-line 1\n-line 2\n+line 3\n+line 4",
			},
			expected: "@@ -1,3 +4,3 @@\n-line 1\n-line 2\n+line 3\n+line 4",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.hunk.Str()
			if result != tt.expected {
				t.Errorf("Hunk.ToString() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestSplitHunks(t *testing.T) {
	input := `diff --git a/file.txt b/file.txt
index abcdef1..1234567 100644
--- a/file.txt
+++ b/file.txt
@@ -1,4 +1,4 @@
 Line 1
 Line 2
 Line 3
-Line 4
+Line 4 modified
`

	expectedResult := []hunk{{
		preImage: lineIndex{
			start:         1,
			linesIncluded: 4,
		},
		postImage: lineIndex{
			start:         1,
			linesIncluded: 4,
		},
		header: `diff --git a/file.txt b/file.txt
index abcdef1..1234567 100644
--- a/file.txt
+++ b/file.txt`,
		content: ` Line 1
 Line 2
 Line 3
-Line 4
+Line 4 modified`},
	}

	result := SplitHunks(input)

	if len(result) != len(expectedResult) {
		t.Errorf("wrong result %+v for output %+v", result, input)
	}
	for i, s := range expectedResult {
		if s != result[i] {
			t.Errorf("wrong output at index %d was \n %+v but should be \n %+v",
				i, result[i], s)
		}
	}

}

func DebugString(content string) string {
	spaces := strings.ReplaceAll(content, " ", "_")
	res := strings.ReplaceAll(spaces, "\n", "\\n")
	return res
}
