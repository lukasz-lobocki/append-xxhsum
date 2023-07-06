package dictionar

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
)

func Load_xxhsum_file(in_file string, bsd_style bool) (map[string]string, error) {

	var (
		file *os.File          = nil
		err  error             = nil
		data map[string]string = nil
	)

	// Open the text file
	if file, err = os.Open(in_file); err != nil {
		return nil, fmt.Errorf("error opening file: %s; %w", in_file, err)
	}
	defer file.Close()

	// Create a dictionary (map) to store the data
	data = make(map[string]string)

	// Read the file line by line
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		if bsd_style {
			// Load BSD-style line
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
		} else {
			// Load GNU-style line
			parts := strings.Split(line, "  ")
			if len(parts) == 2 {
				fileName := strings.TrimSpace(parts[1])
				hashValue := strings.TrimSpace(parts[0])
				data[fileName] = hashValue
			}
		}
	}

	// Check for any scanning errors
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error scanning: %s; %w", in_file, err)
	}

	return data, nil
}

func Dump_xxhsum_dict(in_data map[string]string) {
	// Print the dictionary contents
	for key, value := range in_data {
		log.Printf("DUMP %s  %s\n", value, key)
	}
}
