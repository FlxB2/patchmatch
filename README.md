# Patchmatch
Removes changes in your git diffs if they match a regular expression


```console
foo@bar:~$ git diff | ./patchmatch "long"
```

<table>
<tr>
<tr>
<th>Before</th>
<th>After</th>
</tr>
<td>
 
```diff
diff --git a/file.txt b/file.txt
index abcdef1..1234567 100644
--- a/file.txt
+++ b/file.txt
@@ -1,4 +1,2 @@
 Some input
-Very long input
 Another line
-Fourth line
```
</td>
<td>

```diff
diff --git a/file.txt b/file.txt
index abcdef1..1234567 100644
--- a/file.txt
+++ b/file.txt
@@ -1,4 +1,3 @@
 Some input
 Very long input
 Another line
-Fourth line
```

</td>
</tr>
</table>


## Usage
```console
./patchmatch [-k] <regex>
```
- `<regex>`: The regular expression pattern, following [RE2 Syntax](https://github.com/google/re2/wiki/Syntax)
- `-k`: An optional flag that, if provided, keeps the changes that match the regex. Everything else is removed.

## Installation
For OSX and Linux, simply run `make install`. This will add the binary to `/usr/local/bin`.  
You can compile it yourself using `make build`.  
Run tests using `make tests`.

## Which changes are removed?
Because splitting changed lines is hard, Patchmatch removes blocks of changes entirely if they match the regular expression.   
Similar to existing visual tools, e.g. VSCodes diff view.  


For example:

```console
foo@bar:~$ git diff | patchmatch "((m|M)odified|(r|R)ewritten|(u|U)pdated)"
```

<table>
<thead>
<tr>
<th>Before</th>
<th>After</th>
</thead>
<tbody>
<td>

```diff
diff --git a/complex b/complex
index 04ab2d9..9701e89 100644
--- a/complex
+++ b/complex
@@ -1,27 +1,26 @@
 # Introduction
-This is the initial content of the test file.
+This is the modified content of the test file.
 It includes multiple sections and lines of text.
-All sections follow a consistent structure.
+Some lines were changed.
 
 # Section 1
-The quick brown fox jumps over the lazy dog.
-This sentence contains all letters of the alphabet.
-Numbers follow: 1234567890.
-Special characters: !@#$%^&*()_+.
+The quick brown fox jumps over the lazy cat.
+This sentence contains all letters of the alphabet.
+Numbers have been updated: 9876543210.
+Special characters removed.
 
 # Section 2
 More content in section 2.
-This section has fewer sentences.
-Yet, it is important for the test.
+This section has been shortened.
 End of section 2.
 
 # Section 3
-Section 3 has some lines that will be changed.
-Here we go with the modifications.
-The quick brown fox jumps over something new.
+Section 3 has been completely rewritten.
+These lines are entirely different from the original.
+Expect significant changes here.
 End of section 3.
 
 # Conclusion
-This is the end of the file.
-It will be modified in the new version.
-End of baseline version.
+This is the conclusion of the modified file.
+It has some additional notes.
+End of modified version.
```

</td>
<td valign="top">

```diff
diff --git a/complex b/complex
index 04ab2d9..9701e89 100644
--- a/complex
+++ b/complex
@@ -1,27 +1,26 @@
 # Introduction
 This is the initial content of the test file.
 It includes multiple sections and lines of text.
-All sections follow a consistent structure.
+Some lines were changed.
 
 # Section 1
 The quick brown fox jumps over the lazy dog.
 This sentence contains all letters of the alphabet.
 Numbers follow: 1234567890.
 Special characters: !@#$%^&*()_+.
 
 # Section 2
 More content in section 2.
-This section has fewer sentences.
-Yet, it is important for the test.
+This section has been shortened.
 End of section 2.
 
 # Section 3
 Section 3 has some lines that will be changed.
 Here we go with the modifications.
 The quick brown fox jumps over something new.
 End of section 3.
 
 # Conclusion
 This is the end of the file.
 It will be modified in the new version.
 End of baseline version.
```

</td>
</tbody>
</table>




