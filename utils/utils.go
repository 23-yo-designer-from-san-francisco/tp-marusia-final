package utils

import (
	"encoding/csv"
	"log"
	"os"
	"strings"
)

func ContainsAny(str string, subs ...string) bool {
	for _, sub := range subs {
		if strings.Contains(strings.ToLower(str), strings.ToLower(sub)) {
			return true
		}
	}
	return false
}

func ContainsAll(str string, subs ...string) bool {
	count := 0
	for _, sub := range subs {
		if strings.Contains(str, sub) {
			count += 1
		}
	}
	return count == len(subs)
}

func ReadCsvFile(filePath string) []string {
	f, err := os.Open(filePath)
	if err != nil {
		log.Fatal("Unable to read input file "+filePath, err)
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	records, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal("Unable to parse file as CSV for "+filePath, err)
	}

	return records[0]
}
