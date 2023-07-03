package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"xxhsum/arg_handling"
	"xxhsum/dictionar"

	"github.com/cespare/xxhash/v2"
)

func search_dir(root string, dict map[string]string, xxhsum_filepath string, verbose bool) {
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Printf("Error accessing path %s: %v\n", path, err)
			return nil
		}

		if info.IsDir() {
			// Skip directories
			return nil
		}

		if rel_path, err := filepath.Rel(filepath.Dir(xxhsum_filepath), path); err != nil {
			log.Printf("Error resolving filepath: %v", err)
		} else {
			rel_path = "./" + rel_path

			if _, ok := dict[rel_path]; ok {
				if verbose {
					log.Printf("%s exists.\n", rel_path)
				}
				return nil
			} else {
				if checksum, err := calculateXXHash(path); err != nil {
					log.Printf("Error calculating xxHash: %v", err)
				} else {
					output := fmt.Sprintf("%s  %s\n", checksum, rel_path)
					fmt.Print(output)
					append_to_file(xxhsum_filepath, output)
				}
			}
		}

		return nil
	})

	if err != nil {
		log.Fatalf("Error walking the path %s: %v\n", root, err)
		return
	}
}

func calculateXXHash(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := xxhash.New()

	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return strconv.FormatUint(hash.Sum64(), 16), nil
}

func append_to_file(filename string, content string) {
	// Open the file in append mode, create it if it doesn't exist
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Create a buffered writer for efficient writing
	writer := bufio.NewWriter(file)

	// Write the content to the file
	_, err = writer.WriteString(content)
	if err != nil {
		log.Fatal(err)
	}

	// Flush the buffer to ensure the content is written
	err = writer.Flush()
	if err != nil {
		log.Fatal(err)
	}
}

func init() {
	log.SetPrefix(filepath.Base(os.Args[0] + `: `))
	log.SetFlags(0)
	flag.Usage = func() { fmt.Printf(arg_handling.Usage, filepath.Base(os.Args[0])) }
}

func main() {

	const DEBUG = true

	var (
		verbose         bool              = false
		xxhsum_filepath string            = ""
		given_path      string            = ""
		parent_path     string            = ""
		dict            map[string]string = nil
		exists          bool              = false
	)

	/*
		Parsing input
	*/
	flag.BoolVar(&verbose, "verbose", false, "increase the verbosity.")
	flag.BoolVar(&verbose, "v", false, "increase the verbosity.")
	flag.StringVar(&xxhsum_filepath, "xxhsum-filepath", "", "FILEPATH to file to append to.")
	flag.StringVar(&xxhsum_filepath, "x", "", "FILEPATH to file to append to.")

	flag.Parse()

	// Parsing PATH argument for given_path
	if flag.NArg() != 1 {
		log.Fatalf("PATH agrument missing.\n")
	}
	given_path = arg_handling.Arg_parse(flag.Arg(0), verbose)
	parent_path = filepath.Dir(given_path)

	// Parsing parameter xxhsum-filepath
	if xxhsum_filepath == "" {
		xxhsum_filepath = given_path + ".xxhsum"
		if verbose {
			log.Printf("--xxhsum-filepath defaulted to %s\n", xxhsum_filepath)
		}
	}

	xxhsum_filepath, exists = arg_handling.Param_parse(xxhsum_filepath, verbose)

	if DEBUG {
		log.Printf("DEBUG given_path=%v\n", given_path)
		log.Printf("DEBUG parent_dir=%v\n", parent_path)
		log.Printf("DEBUG xxhsum-path=%v\n", xxhsum_filepath)
		log.Printf("DEBUG xxhsum-path exists=%t\n", exists)
	}

	/*
		Doing the do
	*/
	if exists {
		dict = dictionar.Load_xxhsum_file(xxhsum_filepath)

		if verbose {
			dictionar.Dump_xxhsum_dict(dict)
		}
	}

	search_dir(given_path, dict, xxhsum_filepath, verbose)

	if DEBUG {
		log.Fatalln("DUPA")
	}

	dict = nil
}
