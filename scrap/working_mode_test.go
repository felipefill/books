package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWorkingModeFromString(t *testing.T) {
	scrapOnlyModeString := "sCrAp_OnLy"
	scrapAndStoreModeString := "scrap_and_store"
	retrieveAllModeString := "RETRIEVE_ALL"
	unknownModeString := "i_like_dogs"
	emptyString := ""

	assert.Equal(t, ScrapOnly, WorkingModeFromString(scrapOnlyModeString))
	assert.Equal(t, ScrapAndStore, WorkingModeFromString(scrapAndStoreModeString))
	assert.Equal(t, RetrieveAll, WorkingModeFromString(retrieveAllModeString))
	assert.Equal(t, RetrieveAll, WorkingModeFromString(unknownModeString))
	assert.Equal(t, RetrieveAll, WorkingModeFromString(emptyString))
}
