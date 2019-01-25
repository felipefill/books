package main

import "strings"

// WorkingMode how the handler should work
type WorkingMode int

const (
	// RetrieveAll will only get currently stored books and return them
	RetrieveAll WorkingMode = 0

	// ScrapOnly will scrap Kotlin books from website and return them
	ScrapOnly WorkingMode = 1

	// ScrapAndStore will scarp Kotlin books, store them and then retrieve all stored books and return
	ScrapAndStore WorkingMode = 2
)

// WorkingModeFromString receives a string and returns equivalent WorkingMode
func WorkingModeFromString(mode string) WorkingMode {
	normalizedMode := strings.ToUpper(mode)

	if normalizedMode == "SCRAP_ONLY" {
		return ScrapOnly
	}

	if normalizedMode == "SCRAP_AND_STORE" {
		return ScrapAndStore
	}

	return RetrieveAll
}
