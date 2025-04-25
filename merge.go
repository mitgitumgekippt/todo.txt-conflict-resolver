package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func openFile(filename string) []string {
	// Open the file1
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error opening file:", err)
		os.Exit(1)
	}
	defer file.Close() // Ensure the file is closed when the function exits

	// Create Scanner
	scanner := bufio.NewScanner(file)

	// Create a slice to hold the lines
	var lines []string

	// Read each line and append to the slice
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	// Check for errors during scanning
	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
	}

	return lines
}

func mergeFiles(lines0 []string, lines1 []string, lines2 []string) ([]string, []string) {
	// Naive algorithm, which assumes the presence of the backupfile and that no entry was deleted

	// Cases if content differs
	minsize := min(len(lines0), len(lines1), len(lines2))
	var mergedLines []string
	var proposedChanges []string
	for i := range minsize {
		if lines0[i] == lines1[i] && lines0[i] == lines2[i] {
			// Case 0: Content is the same
			mergedLines = append(mergedLines, lines0[i])
		} else if lines0[i] != lines1[i] && lines0[i] == lines2[i] {
			// Case 1: Content in File 1 differs
			mergedLines = append(mergedLines, lines1[i])
			change := "Change: " + lines0[i] + " --> " + lines1[i]
			proposedChanges = append(proposedChanges, change)
		} else if lines0[i] == lines1[i] && lines0[i] != lines2[i] {
			// Case 2: Content in File 2 differs
			mergedLines = append(mergedLines, lines2[i])
			change := "change: " + lines0[i] + " --> " + lines2[i]
			proposedChanges = append(proposedChanges, change)
		} else if lines0[i] != lines1[i] && lines0[i] != lines2[i] {
			// Case 3: Content in File 3 differs
			// TODO: force change when option is choosen
			fmt.Println("Issue!! - In line %d, you need choose between '%s' and '%s'", i, lines1[i], lines2[i])
		}
	}

	// Case if new tasks were appended
	additionalLines := lines1[(minsize):]
	//fmt.Println("additional lines 1:", additionalLines)
	mergedLines = append(mergedLines, additionalLines...)
	for _, value := range additionalLines {
		change := "new: " + value
		proposedChanges = append(proposedChanges, change)
	}

	additionalLines = lines2[(minsize):]
	//fmt.Println("additional lines 2:", additionalLines)
	mergedLines = append(mergedLines, additionalLines...)
	for _, value := range additionalLines {
		change := "new: " + value
		proposedChanges = append(proposedChanges, change)
	}

	return mergedLines, proposedChanges
}

func printTxt(lines []string) {
	for _, value := range lines {
		fmt.Println(value)
	}

}

func printHelp() {
	fmt.Println(`usage: merge-todoconflict [OPTIONS] FILE1 FILE2 [FILE3 ...]

Merges two or more todo.txt files into a single unified version.

Positional arguments:
  FILE               Paths to the conflict files to be merged (at least two)

Options:
  --double-check 	 Confirm the changes before writing them. Implies --verbose. Overrides --dry.
  --dry              Perform a dry run without altering any file
  --force            Overwrite conflicts with the most recent entry (aka if every file has a different entry for the same todo)
  --init			 Initializes the common base point. If multiple files names are specified, then the first one is taken.
  --verbose          Shows the file contents and file changes
  --help             Show this help message and exit`)
}

func writeToFile(content []string, fileName string) {
	// Create the file
	file, err := os.Create(fileName)
	if err != nil {
		fmt.Println("Error creating file:", err)
		os.Exit(1)
	}

	// Write content to the file
	for index, line := range content {
		_, err = file.WriteString(line + "\n")
		if err != nil {
			fmt.Println("Error writing in line %d to file: %s", index, err)
			os.Exit(1)
		}
	}

	// Close the file
	defer file.Close()
}

func main() {
	args := os.Args
	args = args[1:] // stops deletion of binary
	var run_dry bool = false
	var run_verbose bool = false
	//var run_force bool = false
	var run_check bool = false
	var run_init bool = false
	var file_names []string
	for _, argument := range args {
		switch argument {
		case "--double-check":
			// Show changes and confirm before writing the changes
			run_verbose = true
			run_check = true
		case "--dry":
			// If this option is selected, no file is altered
			run_dry = true
		case "--help":
			printHelp()
			os.Exit(0) // runInit()
		case "--init":
			//Initializes the common base-point
			run_init = true
		case "--force":
			//force changes when it is unclear (multiple changes)
			//vllt skippen von fragen bei init und double check TODO
			//run_force = true
		case "--verbose":
			//show what is done (show proposed changes)
			run_verbose = true

		default:
			file_names = append(file_names, argument)
			//run_dry = false
			//run_force = false
			//run_verbose = false
			//run_check = false
		}
	}

	//TODO use file_names
	file_names = []string{"todo.txt", "todo.txt.conflict"}

	if run_init {
		// TODO
		// check for existing backup and ask
		// get content
		// if verbose ...
		// if dry ...
		// write to file
	}

	fmt.Println("Start merging...")

	lines0 := openFile(".todo.txt.mergebackup") // TODO wenn kein Init: fragen ob trotzdem gemerged werden soll
	lines1 := openFile(file_names[0])
	lines2 := openFile(file_names[1])

	mergedLines, proposedChanges := mergeFiles(lines0, lines1, lines2)

	if run_verbose {
		fmt.Println("-- todo.backup: ")
		printTxt(lines0)
		fmt.Println("-- todo.txt: ")
		printTxt(lines1)
		fmt.Println("-- todo.conflict: ")
		printTxt(lines2)
		fmt.Println("-- Proposed Changes:")
		printTxt(proposedChanges)
		fmt.Println("File would look like this: ", mergedLines)

		if run_check {
			fmt.Println("Do you want confirm the changes? (y/n)")
			reader := bufio.NewReader(os.Stdin)
			input, _ := reader.ReadString('\n')
			input = strings.TrimSpace(strings.ToLower(input))

			if input == "y" || input == "yes" {
				fmt.Println("Continuing...")
				run_dry = false
			} else {
				fmt.Println("Exiting.")
				run_dry = true
			}
		}
	}

	if !run_dry {
		// delete old files
		for _, filename := range file_names {
			os.Remove(filename)
		}
		if run_verbose {
			fmt.Println("Delete old files...")
		}

		// Write new files
		if run_verbose {
			fmt.Println("Write to file ...")
		}
		writeToFile(mergedLines, ".todo.txt.mergebackup")
		writeToFile(mergedLines, file_names[0])
	}
	os.Exit(0)
}
