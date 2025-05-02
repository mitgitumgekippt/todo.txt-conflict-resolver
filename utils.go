package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
)

func openFile(filename string) []string {
	// Open the file1
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error opening file:", err)
		if err.Error() == "open .todo.txt.mergebackup: no such file or directory" {
			fmt.Println("You should run ./merge --init FILE1 first. For help, run with --help .")
		}
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
			fmt.Printf("Error writing in line %d to file: %s", index, err)
			os.Exit(1)
		}
	}

	// Close the file
	defer file.Close()
}

func findFiles(baseFileName string) []string {
	pattern := baseFileName + "*conflict*"
	files, err := filepath.Glob(pattern)
	if err != nil {
		fmt.Printf("There was an error while looking for files: %s \n", err)
		os.Exit(1)
	}
	return files
}
