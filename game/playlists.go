package game

import (
	"fmt"
	"math/rand"
)

func GeneratePlaylistName(nouns []string, adjectives []string, rng *rand.Rand) string {
	noun := nouns[rng.Int63()%int64(len(nouns))]
	adj := adjectives[rng.Int63()%int64(len(adjectives))]
	return fmt.Sprintf("%s %s", adj, noun)
}
