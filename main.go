package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

//TODO: function which handles printing if verbose is active

func main() {
	// Main should only handle filenames and arguments, then starting each function
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
			writeToFile(lines1, ".todo.txt.mergebase")
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
		writeToFile(lines0, ".todo.txt.backup")
		writeToFile(mergedLines, ".todo.txt.mergebase")
		writeToFile(mergedLines, fileNames[0])
		// TODO: delete old files
	}
	os.Exit(0)
}
