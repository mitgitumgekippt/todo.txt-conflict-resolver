package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
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
			fmt.Printf("Issue!! - In line %d, you need choose between '%s' and '%s'! Resolve issue and run again.\n", i, lines1[i], lines2[i])
			os.Exit(1)
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

//TODO: function which handles printing if verbose is active

func findFiles(baseFileName string) []string {
	pattern := baseFileName + "*conflict*"
	files, err := filepath.Glob(pattern)
	if err != nil {
		fmt.Printf("There was an error while looking for files: %s \n", err)
		os.Exit(1)
	}
	return files
}

func main() {
	args := os.Args
	args = args[1:] // stops deletion of binary
	var run_dry bool = false
	var run_verbose bool = false
	//var run_force bool = false
	var run_check bool = false
	var run_init bool = false
	var fileNames []string
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
			//vllt skippen von fragen bei init und double check TODO (soft force and hard force)
			//run_force = true
		case "--verbose":
			//show what is done (show proposed changes)
			run_verbose = true

		default:
			fileNames = append(fileNames, argument)
			//run_dry = false
			//run_force = false
			//run_verbose = false
			//run_check = false
		}
	}

	if len(fileNames) == 0 {
		fmt.Println("Error: Please specify a base file name. Consult --help if instructions are unclear.")
		os.Exit(1)
	} else if len(fileNames) == 1 {
		fmt.Println("Searching for filenames")
		fileNames = findFiles(fileNames[0])
		fmt.Println(fileNames)
		// TODO: change variable names according to naming convention in go
	}
	//TODO use fileNames
	fileNames = []string{"todo.txt", "todo.txt.conflict"}

	lines0 := openFile(".todo.txt.mergebackup") // TODO wenn kein Init: fragen ob trotzdem gemerged werden soll
	lines1 := openFile(fileNames[0])
	lines2 := openFile(fileNames[1])

	if run_init {
		if run_verbose {
			fmt.Println("Start init...")
			fmt.Println("-- todo.txt: ")
			printTxt(lines1)
			fmt.Println("Writing init-file")
		}
		if !run_dry {
			writeToFile(lines1, ".todo.txt.mergebackup")
		}
		os.Exit(0)
	}

	fmt.Println("Start merging...")
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
		for _, filename := range fileNames {
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
		writeToFile(mergedLines, fileNames[0])
	}
	os.Exit(0)
}
