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
	"xxhsum/globals"

	"github.com/briandowns/spinner"
	"github.com/cespare/xxhash/v2"
)

func searchDir(root string, dict map[string]string, xxhsumFilepath string, bsdStyle bool, verbose bool) {

	var (
		line string
	)

	err := filepath.WalkDir(root, func(path string, di fs.DirEntry, err error) error {
		if err != nil {
			log.Printf("error accessing path %s; skipping %v\n", path, err)
			return nil
		}

		// Skip directories and symbolic links
		shouldReturn, returnValue := skipDirs(di)
		if shouldReturn {
			return returnValue
		}

		if rel_path, err := filepath.Rel(filepath.Dir(xxhsumFilepath), path); err != nil {
			log.Printf("error resolving relative path; skipping %v\n", err)
		} else {
			rel_path = "./" + rel_path

			if _, ok := dict[rel_path]; ok {
				// Found
				if verbose {
					log.Printf(globals.GREEN+"INFO"+globals.RESET+" %s exists; skipping\n", rel_path)
				}
			} else {
				// Not found
				if checksum, err := calculateXXHash(path); err != nil {
					log.Printf("error calculating xxHash: %v\n", err)
				} else {
					// Calculate line
					line = calculateLine(bsdStyle, rel_path, checksum)

					// Emit line
					emitLine(xxhsumFilepath, line, verbose)
				}
			}
		}
		return nil
	})

	if err != nil {
		log.Fatalf("Error walking the path %s: %v\n", root, err)
	}
}

func calculateLine(bsdStyle bool, relPath string, checksum string) string {
	if bsdStyle {
		return fmt.Sprintf("XXH64 (%s) = %s\n", relPath, checksum)
	} else {
		return fmt.Sprintf("%s  %s\n", checksum, relPath)
	}
}

func emitLine(xxhsumFilepath string, line string, verbose bool) {
	if err := appendToFile(xxhsumFilepath, line); err != nil {
		log.Printf("error appending to file %s; skipping %v\n", xxhsumFilepath, err)
	}
	if verbose {
		fmt.Print(line)
	}
}

func skipDirs(di fs.DirEntry) (bool, error) {
	if di.IsDir() {
		return true, nil
	}

	if fileInfo, err := di.Info(); err != nil {
		log.Printf("error accessing file %s; skipping %v\n", di.Name(), err)
		return true, nil
	} else {
		if fileInfo.Mode()&os.ModeSymlink != 0 {
			return true, nil
		}
	}
	return false, nil
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

func appendToFile(filename string, content string) error {
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

func debugVariables(verbose bool, givenPath string, xxhsumFilepath string, xxhsumFileExists bool) {
	log.Printf(globals.YELLOW+"DEBUG"+globals.RESET+" given_path=%v\n", givenPath)
	log.Printf(globals.YELLOW+"DEBUG"+globals.RESET+" xxhsum-path=%v\n", xxhsumFilepath)
	log.Printf(globals.YELLOW+"DEBUG"+globals.RESET+" xxhsum-path exists=%t\n", xxhsumFileExists)
}

func init() {
	log.SetPrefix(filepath.Base(os.Args[0] + " "))
	log.SetFlags(0)
	flag.Usage = func() { fmt.Printf(arg_handling.Usage, filepath.Base(os.Args[0])) }
}

type configStru struct {
	verbose  bool
	bsdStyle bool
}

type pathsStru struct {
	xxhsumFileExists bool
	xxhsumFilepath   string
	givenPath        string
}

func main() {

	var (
		verbose          bool              = false
		bsdStyle         bool              = false
		xxhsumFileExists bool              = false
		xxhsumFilepath   string            = ""
		givenPath        string            = ""
		dict             map[string]string = nil
		err              error             = nil
		s                *spinner.Spinner  = nil
	)

	defer func() { dict = nil }()

	config := configStru{verbose: false, bsdStyle: false}
	log.Printf("DiUPA %v %v", &config.verbose, &config.bsdStyle)

	/*
		Parsing input
	*/
	flag.BoolVar(&verbose, "verbose", false, "increase the verbosity.")
	flag.BoolVar(&verbose, "v", false, "increase the verbosity.")
	flag.BoolVar(&bsdStyle, "bsd-style", false, "BSD-style checksum lines.")
	flag.BoolVar(&bsdStyle, "b", false, "BSD-style checksum lines.")
	flag.StringVar(&xxhsumFilepath, "xxhsum-filepath", "", "FILEPATH to file to append to.")
	flag.StringVar(&xxhsumFilepath, "x", "", "FILEPATH to file to append to.")
	flag.Parse()

	/*
		Parsing PATH argument for given_path
	*/
	if flag.NArg() != 1 {
		log.Fatalln(globals.RED + "PATH agrument missing." + globals.RESET)
	}

	givenPath, err = arg_handling.ArgParse(flag.Arg(0), verbose)
	if err != nil {
		log.Fatalf(globals.RED+"%s"+globals.RESET, err)
	}

	/*
		Parsing parameter xxhsum-filepath
	*/
	if xxhsumFilepath == "" {
		xxhsumFilepath = givenPath + ".xxhsum"
		if verbose {
			log.Printf("--xxhsum-filepath defaulted to %s\n", xxhsumFilepath)
		}
	}

	xxhsumFilepath, xxhsumFileExists, err = arg_handling.ParamParse(xxhsumFilepath, verbose)
	if err != nil {
		log.Fatalf(globals.RED+"%s"+globals.RESET, err)
	}

	/*
		Doing the do
	*/
	if verbose {
		debugVariables(verbose, givenPath, xxhsumFilepath, xxhsumFileExists)
	}

	if xxhsumFileExists {
		/*
			Load xxhsum_file to dictionary
		*/
		s = spinner.New(spinner.CharSets[14], 100*time.Millisecond, spinner.WithWriter(os.Stderr),
			spinner.WithSuffix(" Loading"), spinner.WithFinalMSG("Loading complete\n"))
		s.Start()
		dict, err = dictionar.LoadXXHSumFile(xxhsumFilepath, bsdStyle)
		s.Stop()

		if err != nil {
			log.Fatalf(globals.RED+"%s"+globals.RESET, err)
		}

		if verbose {
			/*
				Dump xxhsum_file dictionary
			*/
			dictionar.DumpXXHSumDict(dict)
		}
	}

	/*
	   Search given_path against dictionary
	*/
	s = spinner.New(spinner.CharSets[14], 100*time.Millisecond, spinner.WithWriter(os.Stderr),
		spinner.WithSuffix(" Searching"), spinner.WithFinalMSG("Searching complete\n"))
	if !verbose {
		s.Start()
	}
	searchDir(givenPath, dict, xxhsumFilepath, bsdStyle, verbose)
	if !verbose {
		s.Stop()
	}
}
