Remove changes in your git patch files if they match a regular expression

Example usage

```console
foo@bar:~$ cat file.diff
diff --git a/file.txt b/file.txt
index abcdef1..1234567 100644
--- a/file.txt
+++ b/file.txt
@@ -1,4 +1,2 @@
 Some input
-Very long input
 Another line
-Line 4
foo@bar:~$ cat file.diff | ./patchmatch "long"
diff --git a/file.txt b/file.txt
index abcdef1..1234567 100644
--- a/file.txt
+++ b/file.txt
@@ -1,4 +1,3 @@
 Some input
 Very long input
 Another line
-Line 4
```

