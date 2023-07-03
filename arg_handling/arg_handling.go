package arg_handling

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
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
