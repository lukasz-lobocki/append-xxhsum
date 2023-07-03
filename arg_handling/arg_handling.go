package arg_handling

import (
	"errors"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

const Usage = `
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

func Arg_parse(arg string, verbose bool) string {
	// PATH

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

func Param_parse(param string, verbose bool) string {
	// xxhsum-filepath

	var (
		err       error
		file_path string
	)

	if strings.HasPrefix(param, "~/") {
		param = expand_tilde(param)
	}
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

func expand_tilde(path string) string {
	if usr, err := user.Current(); err != nil {
		log.Fatalln("No current user:", err)
	} else {
		if dir, err := filepath.Abs(usr.HomeDir); err != nil {
			log.Fatalln("No homedir:", err)
		} else {
			if path == "~" {
				// In case of "~", which won't be caught by the "else if"
				path = dir
			} else if strings.HasPrefix(path, "~/") {
				// Use strings.HasPrefix so we don't match paths like
				// "/something/~/something/"
				path = filepath.Join(dir, path[2:])
			}
		}
	}
	return path
}
