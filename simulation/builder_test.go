package simulation

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/munna0908/alien-invasion/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBuildMap(t *testing.T) {
	testWorldMap, testCities, err := createTestWorldWithNeighbours(2, [][]int{
		{1},
		{0},
	})
	if err != nil {
		t.Fatalf("Error creating test cities %s", err)
	}
	// Create temp file
	file, fileName := createTempFile(t)
	// Write cities to file
	for _, v := range testCities {
		_, err = file.WriteString(v.String() + "\n")
		require.NoError(t, err)
	}
	// Read file and build the map
	worldMap, _, err := BuildMap(fileName)
	require.NoError(t, err)
	// Verify the all cities are present
	for cityName := range testWorldMap {
		require.Contains(t, worldMap, cityName)
	}
}
func TestBuildMapInvalidFilePath(t *testing.T) {
	path := "invalidtestfile.txt"
	_, _, err := BuildMap(path)
	require.Error(t, err)
	require.Contains(t, err.Error(), "error reading file")
}
func TestBuildMapEmptyFile(t *testing.T) {
	// Create temp file
	_, fileName := createTempFile(t)
	// Build the world map
	_, _, err := BuildMap(fileName)
	// Assert the error
	require.ErrorIs(t, err, ErrEmptyFile, "ErrEmptyFile is expected")
}

func TestBuildMapInvalidLine(t *testing.T) {
	tests := []struct {
		name       string
		shouldFail bool
		err        error
		data       string
	}{
		{
			name:       "No Neighbours",
			shouldFail: true,
			err:        ErrNoNeighbours,
			data:       "Banglore",
		},
		{
			name:       "Invalid Neighbours",
			shouldFail: true,
			err:        ErrInvalidNeighbour,
			data:       "Banglore north",
		},
		{
			name:       "Invalid Neighbours",
			shouldFail: true,
			err:        ErrInvalidNeighbour,
			data:       "Banglore north=",
		},
		{
			name:       "Invalid Direction",
			shouldFail: true,
			err:        types.ErrInvalidDirection,
			data:       "Banglore northe=hyderabad",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temp file
			_, fileName := createTempFile(t)
			// Write data to file
			err := os.WriteFile(fileName, []byte(tt.data), 0600)
			require.NoError(t, err)
			// Build map
			_, _, err = BuildMap(fileName)
			if tt.shouldFail {
				assert.ErrorIs(t, err, tt.err, "Expecting error")
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// createTempFile creates a temporary file and removes the file after test execution
func createTempFile(t *testing.T) (*os.File, string) {
	t.Helper()

	file, err := ioutil.TempFile("./", "tempFile*.txt")
	require.NoError(t, err)

	t.Cleanup(func() {
		os.Remove(file.Name())
	})

	return file, file.Name()
}
