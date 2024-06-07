package utils

import (
	"fmt"
	"sort"
	"strings"
)

func PrintMapAsTable(hdr1 string, hdr2 string, data map[string]any) {
	// Find the maximum key length for formatting
	maxKeyLength := len(hdr1)
	for key := range data {
		if len(key) > maxKeyLength {
			maxKeyLength = len(key)
		}
	}

	// Print the header
	fmt.Println(strings.Repeat("-", maxKeyLength+len(hdr2)+5))
	fmt.Printf("%-*s | %s\n", maxKeyLength, hdr1, hdr2)
	fmt.Println(strings.Repeat("-", maxKeyLength+len(hdr2)+5))

	// Sort keys for consistent output
	keys := make([]string, 0, len(data))
	for key := range data {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	// Print each key-value pair
	for _, key := range keys {
		fmt.Printf("%-*s | %d\n", maxKeyLength, key, data[key])
	}
	fmt.Println(strings.Repeat("-", maxKeyLength+len(hdr2)+5))
}
