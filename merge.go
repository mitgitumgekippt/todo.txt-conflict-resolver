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

func mergeFiles(lines0 []string, lines1 []string, lines2 []string) []string {
	// Naive algorithm, which assumes the presence of the backupfile and that no entry was deleted

	// Cases if content differs
	minsize := min(len(lines0), len(lines1), len(lines2))
	var mergedLines []string
	for i := range minsize {
		if lines0[i] == lines1[i] && lines0[i] == lines2[i] {
			// Case 0: Content is the same
			mergedLines = append(mergedLines, lines0[i])
		} else if lines0[i] != lines1[i] && lines0[i] == lines2[i] {
			// Case 1: Content in File 1 differs
			mergedLines = append(mergedLines, lines1[i])
		} else if lines0[i] == lines1[i] && lines0[i] != lines2[i] {
			// Case 2: Content in File 2 differs
			mergedLines = append(mergedLines, lines2[i])
		} else if lines0[i] != lines1[i] && lines0[i] != lines2[i] {
			// Case 3: Content in File 3 differs
			fmt.Println("Issue!! - In line %d, you need choose between '%s' and '%s'", i, lines1[i], lines2[i])
		}
	}

	// Case if new tasks were appended
	additionalLines := lines1[(minsize):]
	fmt.Println("Test additional lines 1:", additionalLines)
	mergedLines = append(mergedLines, additionalLines...)

	additionalLines = lines2[(minsize):]
	fmt.Println("Test additional lines 2:", additionalLines)
	mergedLines = append(mergedLines, additionalLines...)

	return mergedLines
}

func main() {
	fmt.Println("Start merging...")

	lines0 := openFile(".todo.mergebackup.txt")
	lines1 := openFile("todo.txt")
	lines2 := openFile("todo.conflict.txt")

	fmt.Println("todo.backup: ", lines0)
	fmt.Println("todo.txt: ", lines1)
	fmt.Println("todo.conflict: ", lines2)

	fmt.Println("merged: ", mergeFiles(lines0, lines1, lines2))

	os.Exit(0)
}
