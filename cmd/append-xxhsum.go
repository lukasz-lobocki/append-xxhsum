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
	"time"

	"github.com/briandowns/spinner"
	"github.com/cespare/xxhash/v2"
	"github.com/lukasz-lobocki/append-xxhsum/pkg/utils"
)

// To be updated with -ldflags during compilation.
var (
	// Version numer shown in help message.
	version string = "development"
	// Git last commit ID shown in help message.
	gitCommit string = ""
	// Time of build shown in help message.
	buildTime string = ""
)

// Walks the directory and adds hashes, missing in the map, to the file.
func searchDir(root string, dict map[string]string, xxhsumFilepath string, bsdStyle bool, verbose bool) int {

	var (
		line string
		i    int
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
					log.Printf(utils.GREEN+"INFO"+utils.RESET+" %s exists; skipping\n", rel_path)
				}
			} else {
				// Not found
				if checksum, err := calculateXXHash(path); err != nil {
					log.Printf("error calculating xxHash: %v\n", err)
				} else {
					// Calculate line
					line = calculateLine(bsdStyle, rel_path, checksum)

					// Emit line
					i = i + emitLine(xxhsumFilepath, line, verbose)
				}
			}
		}
		return nil
	})

	if err != nil {
		log.Fatalf("Error walking the path %s: %v\n", root, err)
	}

	return i
}

// Formats output string according to BSD or default specifiaction.
func calculateLine(bsdStyle bool, relPath string, checksum string) string {
	if bsdStyle {
		return fmt.Sprintf("XXH64 (%s) = %s\n", relPath, checksum)
	}
	return fmt.Sprintf("%s *%s\n", checksum, relPath)
}

// Outputs a line.
func emitLine(xxhsumFilepath string, line string, verbose bool) (linesEmited int) {
	if verbose {
		fmt.Print(line)
	}
	if err := appendToFile(xxhsumFilepath, line); err != nil {
		log.Printf("error appending to file %s; skipping %v\n", xxhsumFilepath, err)
	} else {
		linesEmited = 1
	}
	return
}

// Outputs if parameter is directory or symbolic-link.
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

// Outputs XXH64 hash for the file.
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

	return fmt.Sprintf("%0*x", 16, hash.Sum64()), nil
}

// Appends a string to the file.
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

// Prints some DEBUG info.
func debugVariables(verbose bool, givenPath string, xxhsumFilepath string, xxhsumFileExists bool) {
	log.Printf(utils.YELLOW+"DEBUG"+utils.RESET+" given_path=%v\n", givenPath)
	log.Printf(utils.YELLOW+"DEBUG"+utils.RESET+" xxhsum-path=%v\n", xxhsumFilepath)
	log.Printf(utils.YELLOW+"DEBUG"+utils.RESET+" xxhsum-path exists=%t\n", xxhsumFileExists)
}

// Sets up initial the program environment.
func init() {
	log.SetFlags(0)
	flag.Usage = func() { fmt.Printf(utils.Usage, filepath.Base(os.Args[0]), version, gitCommit, buildTime) }
}

// Main routine.
func main() {

	var (
		verbose          bool              = false
		debug            bool              = false
		bsdStyle         bool              = false
		xxhsumFileExists bool              = false
		xxhsumFilepath   string            = ""
		givenPath        string            = ""
		dict             map[string]string = nil
		err              error             = nil
		s                *spinner.Spinner  = nil
		i                int               = 0
	)

	defer func() { dict = nil }()

	/*
		Parsing input
	*/
	flag.BoolVar(&verbose, "verbose", false, "increase the verbosity.")
	flag.BoolVar(&verbose, "v", false, "increase the verbosity.")
	flag.BoolVar(&debug, "debug", false, "show debug information.")
	flag.BoolVar(&debug, "d", false, "show debug information.")
	flag.BoolVar(&bsdStyle, "bsd-style", false, "BSD-style checksum lines.")
	flag.BoolVar(&bsdStyle, "b", false, "BSD-style checksum lines.")
	flag.StringVar(&xxhsumFilepath, "xxhsum-filepath", "", "FILEPATH to file to append to.")
	flag.StringVar(&xxhsumFilepath, "x", "", "FILEPATH to file to append to.")
	flag.Parse()

	/*
		Parsing PATH argument for given_path
	*/
	if flag.NArg() != 1 {
		log.Fatalln(utils.RED + "PATH agrument missing" + utils.RESET)
	}

	givenPath, err = utils.ArgParse(flag.Arg(0), verbose)
	if err != nil {
		log.Fatalf(utils.RED+"%s"+utils.RESET, err)
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

	xxhsumFilepath, xxhsumFileExists, err = utils.ParamParse(xxhsumFilepath, verbose)
	if err != nil {
		log.Fatalf(utils.RED+"%s"+utils.RESET, err)
	}

	/*
		Doing the do
	*/
	if debug {
		debugVariables(verbose, givenPath, xxhsumFilepath, xxhsumFileExists)
	}

	if xxhsumFileExists {
		/*
			Load xxhsum_file to dictionary
		*/
		s = spinner.New(spinner.CharSets[14], 1000*time.Millisecond, spinner.WithWriter(os.Stderr),
			spinner.WithSuffix(" Loading existing xxhsum file"),
			spinner.WithFinalMSG(fmt.Sprintf("Loading existing %s xxhsum file complete\n", xxhsumFilepath)))
		s.Start()
		dict, err = utils.LoadXXHSumFile(xxhsumFilepath, bsdStyle)
		s.Stop()

		if err != nil {
			log.Fatalf(utils.RED+"%s"+utils.RESET, err)
		}

		if verbose {
			/*
				Dump xxhsum_file dictionary
			*/
			utils.DumpXXHSumDict(dict)
		}
	} else {
		// Create a GNU-style file with heading comment.
		if !bsdStyle {
			appendToFile(xxhsumFilepath, "# XXH64 hashes https://xxhash.com/\n")
			appendToFile(xxhsumFilepath, "# To verify use xxhsum --check --quiet -\n")
		}
	}

	/*
	   Search given_path against dictionary
	*/
	s = spinner.New(spinner.CharSets[14], 1000*time.Millisecond, spinner.WithWriter(os.Stderr),
		spinner.WithSuffix(" Searching and appending new xxhashes to xxhsum file"),
		spinner.WithFinalMSG(fmt.Sprintf("Searching %s and appending new xxhashes to %s xxhsum file complete\n", givenPath, xxhsumFilepath)))
	if !verbose {
		s.Start()
	}
	i = searchDir(givenPath, dict, xxhsumFilepath, bsdStyle, verbose)
	if !verbose {
		s.Stop()
	}

	log.Printf("Hashes added %d\n", i)
}
