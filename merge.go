package main

import (
	"bufio"
	"fmt"
	"os"
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

}

func main() {
	args := os.Args
	var run_dry bool = false
	var run_verbose bool = false
	var run_force bool = false
	for _, argument := range args {
		switch argument {
		case "--dry":
			// If this option is selected, no file is altered
			run_dry = true
		case "--init":
			//Initializes the common base-point
			// runInit()
		case "--verbose":
			//show what is done (show proposed changes)
			run_verbose = true
		case "--force":
			//force changes when it is unclear (multiple changes)
			run_force = true
		case "--help":
			printHelp()
		default:
			run_dry = false
			run_force = false
			run_verbose = false
		}
	}

	fmt.Println("Start merging...")

	lines0 := openFile(".todo.mergebackup.txt") // wenn kein Init: fragen ob trotzdem gemerged werden soll
	lines1 := openFile("todo.txt")
	lines2 := openFile("todo.conflict.txt")

	fmt.Println("-- todo.backup: ")
	printTxt(lines0)
	fmt.Println("-- todo.txt: ")
	printTxt(lines1)
	fmt.Println("-- todo.conflict: ")
	printTxt(lines2)

	mergedLines, proposedChanges := mergeFiles(lines0, lines1, lines2)

	fmt.Println("-- Proposed Changes:")
	printTxt(proposedChanges)

	fmt.Println("File would look like this: ", mergedLines)

	os.Exit(0)
}
