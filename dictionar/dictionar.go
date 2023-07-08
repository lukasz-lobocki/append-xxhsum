package dictionar

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"xxhsum/globals"
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
			loadLine(line, `^XXH64\s*\((?P<fileName>.*)\)\s*=\s*(?P<hashValue>\w+)$`, data)
			/*
				^ asserts the start of the line.
				(\w+) captures one or more word characters as the algorithm name.
				\s* matches zero or more whitespace characters.
				\( matches the opening parenthesis.
				(.*) captures any character (greedy) until the last occurrence of a closing parenthesis.
					This ensures that the match group captures the text between the opening parenthesis and the last closing parenthesis in the file name.
					Should work correctly even when the file name contains nested parentheses.
				\) matches the last closing parenthesis.
				\s* matches zero or more whitespace characters.
				= matches the equals sign.
				\s* matches zero or more whitespace characters.
				(\w+) captures one or more word characters as the hash value.
				$ asserts the end of the line.
			*/
		} else {
			// Load GNU-style line
			loadLine(line, `^(?P<hashValue>.*?)  (?P<fileName>.*)$`, data)
			/*
				^ asserts the start of the line.
				(.*?) captures any characters (except newline characters) lazily in the first group. The ? makes the * quantifier non-greedy
					so that it matches as few characters as possible.
				'  ' matches the double space between the two groups.
				(.*) captures any remaining characters (except newline characters) greedily in the second group.
				$ asserts the end of the line.
			*/
		}
	}

	// Check for any scanning errors
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error scanning: %s; %w", inputFile, err)
	}

	return data, nil
}

func loadLine(line string, pattern string, data map[string]string) {

	regex := regexp.MustCompile(pattern)
	matches := regex.FindStringSubmatch(line)

	if len(matches) == 3 {
		groupNames := regex.SubexpNames()
		result := make(map[string]string)

		for i, match := range matches {
			result[groupNames[i]] = match
		}
		data[result["fileName"]] = result["hashValue"]
	}
}

func DumpXXHSumDict(inputData map[string]string) {
	// Print the dictionary contents
	for key, value := range inputData {
		log.Printf(globals.BLUE+"DUMP"+globals.RESET+" %s  %s\n", value, key)
	}
}
