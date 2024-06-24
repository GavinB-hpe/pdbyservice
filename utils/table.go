package utils

import (
	"fmt"
	"sort"
	"strings"
)

const (
	EDGECHAR = "|"
)

func clip(mx int, s string) string {
	if len(s) > mx {
		return s[:mx]
	}
	return s
}

func Print2DArrayAsTable(mxc int, headers []string, data [][]string) {
	// if headers not specified, use first line of data
	if headers == nil || len(headers) < 1 {
		headers = data[0]
		data = data[1:]
	}
	// get the column widths
	lengths := make([]int, len(headers))
	for i, h := range headers {
		lengths[i] = len(h)
	}
	for _, line := range data {
		for i, _ := range headers {
			if len(line[i]) > lengths[i] {
				lengths[i] = len(line[i])
			}
		}
	}
	// clip col widths for large cols
	for i, l := range lengths {
		if l > mxc {
			lengths[i] = mxc
		}
	}
	// add extra for edges with spaces
	for i, _ := range headers {
		lengths[i] += 4 // need extra for ++ and __
	}
	// draw top line
	for i, _ := range headers {
		fmt.Print(EDGECHAR, strings.Repeat("=", lengths[i]-1), EDGECHAR)
	}
	fmt.Println()
	// Header row
	for i, _ := range headers {
		fmt.Print(EDGECHAR, " ", headers[i], strings.Repeat(" ", lengths[i]-len(headers[i])-3), " ", EDGECHAR)
	}
	fmt.Println()
	// draw bottom line
	for i, _ := range headers {
		fmt.Print(EDGECHAR, strings.Repeat("=", lengths[i]-1), EDGECHAR)
	}
	fmt.Println()
	// draw data
	for _, line := range data {
		for i, _ := range headers {
			fmt.Print(EDGECHAR)
			fmt.Printf(" %-*s", lengths[i]-2, clip(mxc, line[i]))
			fmt.Print(EDGECHAR)
		}
		fmt.Println()
	}

}

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
