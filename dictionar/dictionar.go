package dictionar

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
)

func LoadXXHSumFile(inputFile string, bsdStyle bool) (map[string]string, error) {

	var (
		file    *os.File          = nil
		scanner *bufio.Scanner    = nil
		err     error             = nil
		data    map[string]string = nil
	)

	// Open the text file
	if file, err = os.Open(inputFile); err != nil {
		return nil, fmt.Errorf("error opening file: %s; %w", inputFile, err)
	}
	defer file.Close()

	// Create a dictionary (map) to store the data
	data = make(map[string]string)

	// Read the file line by line
	scanner = bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		if bsdStyle {
			// Load BSD-style line
			loadBSDStyleLine(line, data)
		} else {
			// Load GNU-style line
			loadGNUStyleLine(line, data)
		}
	}

	// Check for any scanning errors
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error scanning: %s; %w", inputFile, err)
	}

	return data, nil
}

func loadGNUStyleLine(line string, data map[string]string) {
	parts := strings.Split(line, "  ")
	if len(parts) == 2 {
		fileName := strings.TrimSpace(parts[1])
		hashValue := strings.TrimSpace(parts[0])
		data[fileName] = hashValue
	}
}

func loadBSDStyleLine(line string, data map[string]string) {
	pattern := `^(\w+)\s*\((.*?)\)\s*=\s*(\w+)$`
	/*
		^ asserts the start of the line.
		(\w+) captures one or more word characters as the algorithm name.
		\s* matches zero or more whitespace characters.
		\( matches the opening parenthesis.
		(.*?) captures any character (non-greedy) until the first occurrence of a closing parenthesis.
			This ensures that the match group captures the text between the opening parenthesis and the last closing parenthesis in the file name.
			Should work correctly even when the file name contains nested parentheses.
		\) matches the last closing parenthesis.
		\s* matches zero or more whitespace characters.
		= matches the equals sign.
		\s* matches zero or more whitespace characters.
		(\w+) captures one or more word characters as the hash value.
		$ asserts the end of the line.
	*/

	re := regexp.MustCompile(pattern)
	match := re.FindStringSubmatch(line)

	if len(match) == 4 {
		algoName := match[1]
		fileName := match[2]
		hashValue := match[3]
		if algoName == "XXH64" {
			data[fileName] = hashValue
		}
	}
}

func DumpXXHSumDict(inputData map[string]string) {
	// Print the dictionary contents
	for key, value := range inputData {
		log.Printf("DUMP %s  %s\n", value, key)
	}
}
