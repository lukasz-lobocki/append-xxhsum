package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const usage = `
Usage: %s [--xxhsum-filepath FILEPATH] [--verbose] [--help] PATH

Recursively adds missing xxhsum hashes from PATH to --xxhsum-filepath.

Arguments:
PATH                          PATH to analyze. Must exist and be readable (+r) and browsable/executable (+x).

Parameters:
-x, --xxhsum-filepath         FILEPATH to file to append to. Defaults to PATH\..\DIRNAME.xxhsum. Must be readable (+r)
                              and writable (+w).
-v, --verbose                 increase the verbosity of the bash script.
-h, --help                    show this help message and exit.
`

func init() {
	log.SetPrefix(filepath.Base(os.Args[0] + `: `))
	log.SetFlags(0)
	flag.Usage = func() { fmt.Printf(usage, filepath.Base(os.Args[0])) }
}

func main() {

	const DEBUG = true

	var (
		verbose         bool
		xxhsum_filepath string
		given_path      string
		parent_dir      string
		dict            map[string]string
	)

	flag.BoolVar(&verbose, "verbose", false, "increase the verbosity.")
	flag.BoolVar(&verbose, "v", false, "increase the verbosity.")
	flag.StringVar(&xxhsum_filepath, "xxhsum-filepath", "", "FILEPATH to file to append to.")
	flag.StringVar(&xxhsum_filepath, "x", "", "FILEPATH to file to append to.")

	flag.Parse()

	// Diagnosing argument for given_path
	if flag.NArg() != 1 {
		log.Fatalf("PATH agrument missing.\n")
	}
	given_path = agr_parse(flag.Arg(0), verbose)
	parent_dir = filepath.Dir(given_path)

	//Diagnosing parameter xxhsum-filepath
	if xxhsum_filepath == "" {
		xxhsum_filepath = given_path + ".xxhsum"
		if verbose {
			log.Printf("--xxhsum-filepath defaulted to %s\n", xxhsum_filepath)
		}
	}
	xxhsum_filepath = param_parse(xxhsum_filepath, verbose)

	if DEBUG {
		log.Printf("DEBUG given_path=%v\n", given_path)
		log.Printf("DEBUG parent_dir=%v\n", parent_dir)
		log.Printf("DEBUG xxhsum-path=%v\n", xxhsum_filepath)
	}

	dict = load_xxhsum_file(xxhsum_filepath)

	if verbose {
		dump_xxhsum_dict(dict)
	}

	search_dir(given_path)

	dict = nil
	if DEBUG {
		log.Fatalln("DUPA")
	}

	rel2, _ := filepath.Rel(filepath.Dir(xxhsum_filepath), "/home/lukasz/Documents/aabbaa/fajle.txt")
	fmt.Printf("rel_parent_fajle: %s\n", `./`+rel2)
}

func agr_parse(arg string, verbose bool) string {

	var (
		err      error
		dir_path string
	)

	if dir_path, err = filepath.Abs(arg); err != nil {
		log.Fatalln("Error resolving filepath:", err)
	}

	if file_info, err := os.Stat(dir_path); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			log.Fatalf("%s does not exist.\n", dir_path)
		} else {
			log.Fatalln("Error accessing file:", err)
		}
	} else {
		if !file_info.Mode().IsDir() {
			log.Fatalf("%s exists but is not a directory.\n", dir_path)
		} else if verbose {
			log.Printf("%s is a directory.\n", dir_path)
		}
	}

	return dir_path
}

func param_parse(param string, verbose bool) string {

	var (
		err       error
		file_path string
	)

	if file_path, err = filepath.Abs(param); err != nil {
		log.Fatalln("Error resolving filepath:", err)
	}

	if file_info, err := os.Stat(file_path); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			if file, err := os.Create(file_path); err != nil {
				defer file.Close()
				log.Fatalln("Error creating file:", err)
			}
			log.Printf("%s created.\n", file_path)
		} else {
			log.Fatalln("Error accessing file:", err)
		}
	} else {
		if !file_info.Mode().IsRegular() {
			log.Fatalf("%s exists but is not a file.\n", file_path)
		} else if verbose {
			log.Printf("%s is a file.\n", file_path)
		}
	}

	return file_path
}

func load_xxhsum_file(in_file string) map[string]string {

	var (
		file *os.File
		err  error
		data map[string]string
	)

	// Open the text file
	file, err = os.Open(in_file)
	if err != nil {
		defer file.Close()
		log.Fatalln("Error opening file:", err)
	}

	// Create a dictionary (map) to store the data
	data = make(map[string]string)

	// Read the file line by line
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, "  ")
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[1])
			value := strings.TrimSpace(parts[0])
			data[key] = value
		}
	}

	// Check for any scanning errors
	if err := scanner.Err(); err != nil {
		log.Fatalln("Error scanning file:", err)
		return nil
	}

	return data
}

func dump_xxhsum_dict(in_data map[string]string) {
	// Print the dictionary contents
	for key, value := range in_data {
		fmt.Printf("DUMP: %s {%s}\n", key, value)
	}
}

func search_dir(dir string) {
	// Specify the root directory
	root := dir

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("Error accessing path %s: %v\n", path, err)
			return nil
		}

		if info.IsDir() {
			// Skip directories
			return nil
		}

		// Print the file path
		fmt.Printf("SEARCH: %s\n", path)

		return nil
	})

	if err != nil {
		fmt.Printf("Error walking the path %s: %v\n", root, err)
		return
	}
}
