package bot

import (
	"bytes"
	"fmt"
)

func splitStringIfNeeded(text string) []string {
	// Use a byte slice for efficient and UTF-8 safe operations
	textBytes := []byte(text)
	totalSize := len(textBytes)
	fmt.Printf("Total string size: %d bytes\n", totalSize)

	// If the size is within the limit, return it as a single element slice
	if totalSize <= 2000 {
		return []string{text}
	}

	var parts []string
	remainingBytes := textBytes

	const limit = 2000

	for len(remainingBytes) > limit {
		// Define the search area for the comma (the first 2000 bytes)
		searchArea := remainingBytes[:limit]

		// Find the last comma in the search area
		splitPos := bytes.LastIndex(searchArea, []byte(","))

		// If no comma is found, we must split at the 2000 byte mark to avoid an infinite loop.
		// This ensures progress even with unstructured data.
		if splitPos == -1 {
			splitPos = limit
		}

		// Add the part before the split position to our slice of strings
		parts = append(parts, string(remainingBytes[:splitPos]))

		// Update the remaining bytes, skipping the comma itself
		// (+1 to move past the comma)
		remainingBytes = remainingBytes[splitPos+1:]
	}

	// Add the final remaining part to the slice
	if len(remainingBytes) > 0 {
		parts = append(parts, string(remainingBytes))
	}

	return parts
}
