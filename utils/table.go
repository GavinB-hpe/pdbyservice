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

	// Print each key-value pair
	for _, key := range keys {
		fmt.Printf("%-*s | %d\n", maxKeyLength, key, data[key])
	}
	fmt.Println(strings.Repeat("-", maxKeyLength+len(hdr2)+5))
}

func PrintMapIntAsSortedTable(hdr1 string, hdr2 string, data map[string]int) {
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

	// Sort keys by data value
	keys := make([]string, 0, len(data))
	for key := range data {
		keys = append(keys, key)
	}
	// sort.Strings(keys)
	sort.Slice(keys, func(i, j int) bool { return data[keys[i]] > data[keys[j]] })

	// Print each key-value pair
	for _, key := range keys {
		fmt.Printf("%-*s | %d\n", maxKeyLength, key, data[key])
	}
	fmt.Println(strings.Repeat("-", maxKeyLength+len(hdr2)+5))
}
