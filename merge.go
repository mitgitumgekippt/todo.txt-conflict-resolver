package main

import (
	"fmt"
	"os"
)

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
