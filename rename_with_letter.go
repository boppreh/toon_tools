package main

import (
    "bufio"
    "fmt"
    "os"
    "strings"
    "path/filepath"
    "sort"
)

func main() {
	fmt.Println("LT_Tom, welcome to File Renamer!")
	fmt.Println("This program adds a letter (or more) to all files in a given folder.")
	fmt.Println()

	if len(os.Args) <= 1 {
		fmt.Println("ERROR")
		fmt.Println("Please run this program by dragging and dropping a folder onto it.")
		fmt.Println("Press any key to exit...")
		fmt.Scanln()
		os.Exit(1)
	}

	// Path of the folder with the files to rename.
	path := os.Args[1]

    // Find all files in the given folder. Lazy solution by matching all names with a dot in it.
    files, err := filepath.Glob(filepath.Join(path, "*.*"))
    if err != nil {
    	fmt.Printf("ERROR: Could not read files in %s because of: %s\n", path, err)
    	fmt.Printf("Press any key to exit...")
    	fmt.Scanln()
    	os.Exit(2)
    }
    fmt.Printf("Will rename the following files: %s\n\n", files)

	// Characters that must not appear in a Windows filename.
	BAD_CHARACTERS := ".\\/:<>\"|?*"

	// Ask for the characters to add, in a command prompt window.
	// If there's a problem, keep trying until a valid answer is given.
    var str string
    reader := bufio.NewReader(os.Stdin)
    for {
        fmt.Printf("What letter to add to the files in %s? ", path)
        str, _ = reader.ReadString('\n')
        if str == "" {
        	fmt.Printf("ERROR: Cannot be empty.\n\n")
        } else if strings.ContainsAny(str, BAD_CHARACTERS) {
        	fmt.Printf("ERROR: Windows filenames cannot contain any of the following characters: %s\n\n", BAD_CHARACTERS)
        } else {
        	break
        }
    }
    prefix := strings.TrimSpace(str)

    // Sort list of files by length, longest first. This guarantees that adding a prefix will not replace an existing file.
	sort.Slice(files, func(i, j int) bool {
        return len(files[i]) > len(files[j])
    })

    // Rename the files in the folder.
    for _, file := range files {
    	newFile := filepath.Join(filepath.Dir(file), prefix + filepath.Base(file))
    	fmt.Printf("mv %s %s\n", file, newFile)
    	err = os.Rename(file, newFile)
    	if err != nil {
    		fmt.Printf("ERROR: failed to rename %s to %s because of: %s. Continuing with the rest of the files...\n", file, newFile, err)
    	}
    }

	fmt.Println("Press any key to exit...")
	fmt.Scanln()
	os.Exit(0)
}
