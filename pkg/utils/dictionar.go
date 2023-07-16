package utils

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
)

// Loads xxhsum_file to the map.
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
			loadLine(line, `^XXH64 \((?P<fileName>.*)\) = (?P<hashValue>\w+)$`, data)
			/*
				^ asserts the start of the line.
				XXH64 matches exact characters as the algorithm name.
				' ' matches space between groups.
				\( matches the opening parenthesis.
				(.*) captures any character (greedy) until the last occurrence of a closing parenthesis.
					This ensures that the match group captures the text between the opening parenthesis and the last closing parenthesis in the file name.
					Should work correctly even when the file name contains nested parentheses.
				\) matches the last closing parenthesis.
				' ' matches space between groups.
				= matches the equal sign.
				' ' matches space between groups.
				(\w+) captures one or more word characters as the hash value.
				$ asserts the end of the line.
			*/
		} else {
			// Load GNU-style line
			loadLine(line, `^(?P<hashValue>\w+) [ \*](?P<fileName>.*)$`, data)
			/*
				^ asserts the start of the line.
				(\w+) captures one or more word characters as the hash value.
				' ' matches single space between the two groups.
				'[ \*]' matches either single space or single asterisk.
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

// Outputs xxhsum_file map.
func DumpXXHSumDict(inputData map[string]string) {
	// Print the dictionary contents
	for key, value := range inputData {
		log.Printf(BLUE+"DUMP"+RESET+" %s  %s\n", value, key)
	}
}
