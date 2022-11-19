package utils

import "strings"

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
