package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"
	"xxhsum/arg_handling"
	"xxhsum/dictionar"

	"github.com/briandowns/spinner"
	"github.com/cespare/xxhash/v2"
)

const (
	RESET  string = "\033[0m"
	RED    string = "\033[31m"
	GREEN  string = "\033[32m"
	YELLOW string = "\033[33m"
	Blue   string = "\033[34m"
	Purple string = "\033[35m"
	Cyan   string = "\033[36m"
	Gray   string = "\033[37m"
	White  string = "\033[97m"
)

func search_dir(root string, dict map[string]string, xxhsum_filepath string, bsd_style bool, verbose bool) {
	var (
		line string
	)

	err := filepath.WalkDir(root, func(path string, di fs.DirEntry, err error) error {
		if err != nil {
			log.Printf("error accessing path %s; skipping %v\n", path, err)
			return nil
		}

		if di.IsDir() {
			// Skip directories
			return nil
		}

		if fileInfo, err := di.Info(); err != nil {
			log.Printf("error accessing file %s; skipping %v\n", di.Name(), err)
			return nil
		} else {
			if fileInfo.Mode()&os.ModeSymlink != 0 {
				return nil // Skip symbolic links
			}
		}

		if rel_path, err := filepath.Rel(filepath.Dir(xxhsum_filepath), path); err != nil {
			log.Printf("error resolving relative path; skipping %v\n", err)
		} else {
			rel_path = "./" + rel_path

			if _, ok := dict[rel_path]; ok {
				// Found
				if verbose {
					log.Printf(GREEN+"INFO"+RESET+" %s exists; skipping\n", rel_path)
				}
			} else {
				// Not found
				if checksum, err := calculateXXHash(path); err != nil {
					log.Printf("error calculating xxHash: %v\n", err)
				} else {
					// Calculate line
					if bsd_style {
						line = fmt.Sprintf("XXH64 (%s) = %s\n", rel_path, checksum)
					} else {
						line = fmt.Sprintf("%s  %s\n", checksum, rel_path)
					}

					// Emit line
					if err := append_to_file(xxhsum_filepath, line); err != nil {
						log.Printf("error appending to file %s; skipping %v\n", xxhsum_filepath, err)
					}
					if verbose {
						fmt.Print(line)
					}
				}
			}
		}
		return nil
	})

	if err != nil {
		log.Fatalf("Error walking the path %s: %v\n", root, err)
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

func append_to_file(filename string, content string) error {
	// Open the file in append mode, create it if it doesn't exist
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return err
	}

	// Create a buffered writer for efficient writing
	writer := bufio.NewWriter(file)

	// Write the content to the file
	_, err = writer.WriteString(content)
	if err != nil {
		return err
	}

	// Flush the buffer to ensure the content is written
	err = writer.Flush()
	if err != nil {
		return err
	}
	return file.Close()
}

func init() {
	log.SetPrefix(filepath.Base(os.Args[0] + " "))
	log.SetFlags(0)
	flag.Usage = func() { fmt.Printf(arg_handling.Usage, filepath.Base(os.Args[0])) }
}

func main() {

	var (
		verbose            bool              = false
		bsd_style          bool              = false
		xxhsum_filepath    string            = ""
		given_path         string            = ""
		parent_path        string            = ""
		dict               map[string]string = nil
		xxhsum_file_exists bool              = false
		err                error             = nil
	)

	/*
		Parsing input
	*/
	flag.BoolVar(&verbose, "verbose", false, "increase the verbosity.")
	flag.BoolVar(&verbose, "v", false, "increase the verbosity.")
	flag.BoolVar(&bsd_style, "bsd-style", false, "BSD-style checksum lines.")
	flag.BoolVar(&bsd_style, "b", false, "BSD-style checksum lines.")
	flag.StringVar(&xxhsum_filepath, "xxhsum-filepath", "", "FILEPATH to file to append to.")
	flag.StringVar(&xxhsum_filepath, "x", "", "FILEPATH to file to append to.")

	flag.Parse()

	/*
		Parsing PATH argument for given_path
	*/
	if flag.NArg() != 1 {
		log.Fatalln(RED + "PATH agrument missing." + RESET)
	}

	given_path, err = arg_handling.Arg_parse(flag.Arg(0), verbose)
	if err != nil {
		log.Fatalf(RED+"%s"+RESET, err)
	}

	parent_path = filepath.Dir(given_path)

	/*
		Parsing parameter xxhsum-filepath
	*/
	if xxhsum_filepath == "" {
		xxhsum_filepath = given_path + ".xxhsum"
		if verbose {
			log.Printf("--xxhsum-filepath defaulted to %s\n", xxhsum_filepath)
		}
	}

	xxhsum_filepath, xxhsum_file_exists, err = arg_handling.Param_parse(xxhsum_filepath, verbose)
	if err != nil {
		log.Fatalf(RED+"%s"+RESET, err)
	}

	/*
		Doing the do
	*/
	if verbose {
		log.Printf(YELLOW+"DEBUG"+RESET+" given_path=%v\n", given_path)
		log.Printf(YELLOW+"DEBUG"+RESET+" parent_dir=%v\n", parent_path)
		log.Printf(YELLOW+"DEBUG"+RESET+" xxhsum-path=%v\n", xxhsum_filepath)
		log.Printf(YELLOW+"DEBUG"+RESET+" xxhsum-path exists=%t\n", xxhsum_file_exists)
	}

	if xxhsum_file_exists {
		s := spinner.New(spinner.CharSets[14], 100*time.Millisecond, spinner.WithWriter(os.Stderr))
		s.Prefix = "Loading: "
		s.FinalMSG = "Loading complete\n"
		s.Start()

		// Load xxhsum_file to dictionary
		dict, err = dictionar.Load_xxhsum_file(xxhsum_filepath, bsd_style)

		s.Stop()

		if err != nil {
			log.Fatalf(RED+"%s"+RESET, err)
		}

		if verbose {
			dictionar.Dump_xxhsum_dict(dict)
		}
	}

	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond, spinner.WithWriter(os.Stderr))
	if !verbose {
		s.Prefix = "Searching: "
		s.FinalMSG = "Searching complete\n"
		s.Start()
	}

	// Search given_path against dictionary
	search_dir(given_path, dict, xxhsum_filepath, bsd_style, verbose)

	if !verbose {
		s.Stop()
	}

	dict = nil
}
