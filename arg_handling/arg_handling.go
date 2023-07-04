package arg_handling

import (
	"errors"
	"fmt"
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

func Arg_parse(arg string, verbose bool) (string, error) {
	// PATH

	var (
		err      error  = nil
		dir_path string = ""
	)

	if dir_path, err = filepath.Abs(arg); err != nil {
		return "", errors.Join(fmt.Errorf("error resolving filepath: %s", arg), err)
	}

	if file_info, err := os.Stat(dir_path); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return "", errors.Join(fmt.Errorf("does not exist: %s", dir_path), err)
		} else {
			return "", errors.Join(fmt.Errorf("error accessing file: %s", dir_path), err)
		}
	} else {
		if !file_info.Mode().IsDir() {
			return "", fmt.Errorf("exists but is not a directory: %s", dir_path)
		} else if verbose {
			log.Printf("%s is a directory.\n", dir_path)
		}
	}
	return dir_path, nil
}

func Param_parse(param string, verbose bool) (string, bool, error) {
	// xxhsum-filepath

	var (
		err       error  = nil
		file_path string = ""
	)

	file_path, err = expand_tilde(param)
	if err != nil {
		return "", false, err
	}

	if file_info, err := os.Stat(file_path); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return "", false, errors.Join(fmt.Errorf("error accessing file: %s", file_path), err)
		} else {
			return file_path, false, nil
		}
	} else {
		if !file_info.Mode().IsRegular() {
			return "", true, fmt.Errorf("exists but is not a file: %s", file_path)
		} else {
			if verbose {
				log.Printf("%s is a file.\n", file_path)
			}
			return file_path, true, nil
		}
	}
}

func expand_tilde(in_path string) (string, error) {

	if !strings.HasPrefix(in_path, "~") {
		return in_path, nil
	}

	if usr, err := user.Current(); err != nil {
		return "", errors.Join(fmt.Errorf("no current user"), err)
	} else {
		if dir, err := filepath.Abs(usr.HomeDir); err != nil {
			return "", errors.Join(fmt.Errorf("no homedir"), err)
		} else {
			switch true {
			case !strings.HasPrefix(in_path, "~"):
				{
					return "", errors.Join(fmt.Errorf("no current user"), err)
				}
			case in_path == "~":
				{
					return dir, nil
				}
			case strings.HasPrefix(in_path, "~/"):
				{
					return filepath.Join(dir, in_path[2:]), nil
				}
			default:
				{
					return "", fmt.Errorf("error resolving filepath: %s", in_path)
				}
			}
		}
	}
}
