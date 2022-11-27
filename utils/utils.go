package utils

import (
	"encoding/csv"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
	"unsafe"

	"github.com/seehuhn/mt19937"
)

func ContainsAny(str string, subs ...string) bool {
	for _, sub := range subs {
		fmt.Println("Str: ",str," Sub: ",sub)
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

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

var src = rand.NewSource(time.Now().UnixNano())

func RandStringBytesMaskImprSrcUnsafe(n int) string {
	b := make([]byte, n)
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return *(*string)(unsafe.Pointer(&b))
}

func GeneratePlaylistName(nouns []string, adjectives []string) string {
	rng := rand.New(mt19937.New())
	rng.Seed(time.Now().UnixNano())

	noun := nouns[rng.Int63()%int64(len(nouns))]
	adj := adjectives[rng.Int63()%int64(len(adjectives))]
	return fmt.Sprintf("%s %s", adj, noun)
}