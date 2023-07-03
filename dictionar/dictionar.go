package dictionar

import (
	"bufio"
	"log"
	"os"
	"strings"
)

func Load_xxhsum_file(in_file string) map[string]string {

	var (
		file *os.File          = nil
		err  error             = nil
		data map[string]string = nil
	)

	// Open the text file
	file, err = os.Open(in_file)
	if err != nil {
		defer file.Close()
		log.Fatalln("Error opening file:", err)
	}

	// Create a dictionary (map) to store the data
	data = make(map[string]string)

	// Read the file line by line
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, "  ")
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[1])
			value := strings.TrimSpace(parts[0])
			data[key] = value
		}
	}

	// Check for any scanning errors
	if err := scanner.Err(); err != nil {
		log.Fatalln("Error scanning file:", err)
		return nil
	}

	return data
}

func Dump_xxhsum_dict(in_data map[string]string) {
	// Print the dictionary contents
	for key, value := range in_data {
		log.Printf("DUMP %s  %s\n", value, key)
	}
}
