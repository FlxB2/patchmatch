#!/bin/sh

apply_and_check_diff() {
  if [ "$#" -ne 2 ]; then
    echo "Usage: apply_and_check_diff <diff_filename> <regex>"
    return 1
  fi

  file="$1"
  diff_file="${1}.diff"
  regex="$2"

  # overwrite w/ files from ../diff
  cp -r ../modified/* .
  # save git diff
  git diff ${file} > basic.diff
  # stash changes to be in baseline state
  git stash >/dev/null 2>&1

  new_diff_file="${file}_new.diff"

  # Filter the diff file using patchmatch and apply the changes
  cat "$diff_file" | ../../patchmatch "$regex" > "$new_diff_file"
  git apply "$new_diff_file"

  # Check if git apply was successful
  if [ $? -ne 0 ]; then
    echo "❌ ${file} failure"
    echo "Git apply failed"
    echo "Before:"
    cat "$diff_file"
    echo "After:"
    cat "$new_diff_file"
    return 1
  fi

  # Compare the resulting file with the expected file
  if diff -q basic ../expected/basic >/dev/null; then
    echo "✅ ${file} success"
  else
    echo "❌ ${file} failure"
    return 1
  fi

  git stash > /dev/null 2>&1
}

# build patchmatch
go build ../

# init git 
git config --global user.email "you@example.com"
git config --global user.name "Your Name"
git config --global init.defaultBranch main

# setup git repo
cd testFiles/baseline
git init > /dev/null 2>&1
git add -A > /dev/null 2>&1
git commit -m "test" > /dev/null 2>&1

failure=0
echo "-- Starting Test --"
apply_and_check_diff basic "content" || failure=1

if [ "$failure" -ne 0 ]; then
  echo "One or more tests failed."
  exit 1
fi

echo "-- End Tests --"