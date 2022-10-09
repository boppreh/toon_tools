package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
)

func die(status int) {
	fmt.Println("Press any key to exit...")
	fmt.Scanln()
	os.Exit(status)
}

func regexGlobOrDie(path string, patternString string) (matches []string) {
	entries, err := ioutil.ReadDir(path)
	if err != nil {
		fmt.Printf("ERROR: Could not read files in %s because of: %s\n", path, err)
		die(2)
	}

	pattern := regexp.MustCompile("^" + patternString + "$")
	for _, entry := range entries {
		if !entry.IsDir() && pattern.MatchString(entry.Name()) {
			matches = append(matches, filepath.Join(path, entry.Name()))
		}
	}

	return matches
}

func main() {
	fmt.Println("LT_Tom, welcome to Filename Restorer!")
	fmt.Println("This program takes two folders, and renames the files in one to match the names in the other.")
	fmt.Println()

	if len(os.Args) != 3 {
		fmt.Println("ERROR")
		fmt.Println("Please run this program by dragging and dropping exactly two folders on it.")
		die(1)
	}

	// Path of the folders with the files to rename.
	path1 := os.Args[1]
	path2 := os.Args[2]

	// Find all files in the given folder. Lazy solution by matching all names with a dot in it.
	numberedFiles1 := regexGlobOrDie(path1, `\d+\..+`)
	numberedFiles2 := regexGlobOrDie(path2, `\d+\..+`)

	if len(numberedFiles1) == 0 && len(numberedFiles2) == 0 {
		fmt.Println("ERROR: Neither folder contains files with numeric names. Don't know which one to rename.")
		die(3)
	}

	var sourceFiles, destFiles []string
	if len(numberedFiles1) > len(numberedFiles2) {
		destFiles = numberedFiles1
		sourceFiles = regexGlobOrDie(path2, `.+\..+`)
	} else {
		destFiles = numberedFiles2
		sourceFiles = regexGlobOrDie(path1, `.+\..+`)
	}

	if len(sourceFiles) != len(destFiles) {
		fmt.Println("ERROR: the number of file in the folders doesn't match.")
		die(3)
	}

	// Rename the files in the folder.
	for i := 0; i < len(destFiles); i++ {
		sourceFile := sourceFiles[i]
		destFile := destFiles[i]

		newDestFile := filepath.Join(filepath.Dir(destFile), filepath.Base(sourceFile))
		fmt.Printf("mv %s %s\n", destFile, newDestFile)
		err := os.Rename(destFile, newDestFile)
		if err != nil {
			fmt.Printf("ERROR: failed to rename %s to %s because of: %s. Continuing with the rest of the files...\n", destFile, newDestFile, err)
		}
	}

	fmt.Println("\n\nAll done!\n")
	die(0)
}
