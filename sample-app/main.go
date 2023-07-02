package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

const usage = `
Usage: %s [--xxhsum-filepath=FILEPATH] [--verbose] [--help] PATH

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

	var (
		verbose         bool
		xxhsum_filepath string
		given_path      string
		parent_dir      string
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

	if verbose {
		log.Printf("DEBUG given_path=%v\n", given_path)
		log.Printf("DEBUG parent_dir=%v\n", parent_dir)
		log.Printf("DEBUG xxhsum-path=%v\n", xxhsum_filepath)
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
		log.Fatalln(err)
	}

	if file_info, err := os.Stat(dir_path); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			log.Fatalf("%s does not exist.\n", dir_path)
		} else {
			log.Fatalln(err)
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
		log.Fatalln(err)
	}

	if file_info, err := os.Stat(file_path); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			if file, err := os.Create(file_path); err != nil {
				defer file.Close()
				log.Fatalln(err)
			}
			log.Printf("%s created.\n", file_path)
		} else {
			log.Fatalln(err)
		}
	} else {
		if !file_info.Mode().IsRegular() {
			log.Fatalf("%s is not a file.\n", file_path)
		} else if verbose {
			log.Printf("%s is a file.\n", file_path)
		}
	}

	return file_path
}
