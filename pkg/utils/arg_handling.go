package utils

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

// Text of help.
const Usage string = `
Usage: %s [--xxhsum-filepath FILEPATH] [--bsd-style] [--verbose] [--help] PATH

Recursively adds missing xxhsum (XXH64) hashes from PATH to --xxhsum-filepath.

Arguments:
PATH                            PATH to analyze

Parameters:
-x, --xxhsum-filepath           FILEPATH of file to append to. Defaults to PATH\..\DIRNAME.xxhsum
-b, --bsd-style                 BSD-style checksum lines. Defaults to GNU-style
-v, --verbose                   increase the verbosity
-d, --debug                     show debug information
-h, --help                      show this help message and exit

To verify use xxhsum --check --quiet FILEPATH

version: %s commit: %s built: %s
`

// Parses PATH part of program invocation.
func ArgParse(arg string, verbose bool) (string, error) {

	var (
		err      error  = nil
		dir_path string = ""
	)

	if dir_path, err = filepath.Abs(filepath.Clean(arg)); err != nil {
		return "", fmt.Errorf("error resolving filepath: %s; %w", arg, err)
	}

	if file_info, err := os.Stat(dir_path); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return "", fmt.Errorf("does not exist: %s; %w", dir_path, err)
		} else {
			return "", fmt.Errorf("error accessing file: %s; %w", dir_path, err)
		}
	} else {
		if !file_info.Mode().IsDir() {
			return "", fmt.Errorf("exists but is not a directory: %s", dir_path)
		} else if verbose {
			log.Printf("%s is a directory\n", dir_path)
		}
	}
	return dir_path, nil
}

// Pareses xxhsum-filepath part of program invocation.
func ParamParse(param string, verbose bool) (string, bool, error) {

	var (
		err       error  = nil
		file_path string = ""
	)

	param, err = expandTilde(param)
	if err != nil {
		return "", false, err
	}
	if file_path, err = filepath.Abs(filepath.Clean(param)); err != nil {
		return "", false, fmt.Errorf("error resolving filepath: %s; %w", param, err)
	}

	if file_info, err := os.Stat(file_path); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return file_path, false, nil
		} else {
			return "", false, fmt.Errorf("error accessing file: %s; %w", file_path, err)
		}
	} else {
		if file_info.Mode().IsRegular() {
			if verbose {
				log.Printf("%s is a file\n", file_path)
			}
			return file_path, true, nil
		} else {
			return "", true, fmt.Errorf("exists but is not a file: %s", file_path)
		}
	}
}

// Replaces ~ tilde from inputPath with current user's homedir.
func expandTilde(inputPath string) (string, error) {

	if !strings.HasPrefix(inputPath, "~") {
		return filepath.Clean(inputPath), nil
	}

	if usr, err := user.Current(); err != nil {
		return "", fmt.Errorf("no current user; %w", err)
	} else {
		if dir, err := filepath.Abs(usr.HomeDir); err != nil {
			return "", fmt.Errorf("no homedir; %w", err)
		} else {
			switch true {
			case inputPath == "~":
				{
					return dir, nil
				}
			case strings.HasPrefix(inputPath, "~/"):
				{
					return filepath.Join(dir, inputPath[2:]), nil
				}
			default:
				{
					return "", fmt.Errorf("error resolving filepath: %s", inputPath)
				}
			}
		}
	}
}
