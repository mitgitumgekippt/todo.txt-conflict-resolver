package main

import "fmt"

func printTxt(lines []string) {
	for _, value := range lines {
		fmt.Println(value)
	}
}

func printHelp() {
	fmt.Println(`usage: merge-todoconflict [OPTIONS] FILE1 FILE2 [FILE3 ...]

Merges two or more todo.txt files into a single unified version. Needed to be run with --init, before any merging can be done.
The merging algorythm is similar to the one of git, but with some assumptions:
- No line was deleted
If only one filename is given, it searches for possible conflicting files by itself.

Positional arguments:
  FILE               Paths to the conflict files to be merged (at least two)

Options:
  --double-check 	 Confirm the changes before writing them. Implies --verbose. Overrides --dry.
  --dry              Perform a dry run without altering any file
  --force            Overwrite conflicts with the most recent entry (aka if every file has a different entry for the same todo)
  --init			 Initializes the common base point. If multiple files names are specified, then the first one is taken. Stops after.
  --verbose          Shows the file contents and file changes
  --help             Show this help message and exit`)
}
