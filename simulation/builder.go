package simulation

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/munna0908/alien-invasion/types"
	"github.com/pkg/errors"
)

var (
	ErrInvalidLine      = errors.New("invalid line")
	ErrInvalidNeighbour = errors.New("invalid neighbour")
	ErrNoNeighbours     = errors.New("no neighbours")
	ErrEmptyFile        = errors.New("empty file")
)

// BuildMap reads the input file and create a map of cities
func BuildMap(filePath string) (types.World, []*types.City, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, nil, errors.Wrap(err, "error reading file")
	}

	st, err := file.Stat()
	if err != nil {
		return nil, nil, errors.Wrap(err, "error reading file")
	}

	if st.Size() == 0 {
		return nil, nil, ErrEmptyFile
	}
	// Create a world map and cities instance
	worldMap := types.NewWorldMap()
	cities := make([]*types.City, 0)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line = strings.TrimRight(line, " "); len(line) == 0 {
			continue
		}

		tokens := strings.Split(line, " ")
		if len(tokens) == 1 {
			// City should have aleast one neighbour
			return nil, nil, ErrNoNeighbours
		}

		city := worldMap.GetCity(tokens[0])
		if city == nil {
			// If city doesnt exists create one
			city = types.NewCity(tokens[0], len(tokens[1:]))
			if err := worldMap.AddCity(city); err != nil {
				return nil, nil, err
			}

			cities = append(cities, city)
		}

		for _, links := range tokens[1:] {
			neighbours := strings.Split(links, "=")
			if len(neighbours) != 2 || neighbours[1] == "" {
				return nil, nil, ErrInvalidNeighbour
			}
			// Parse the neighbours and create the cities if required
			neighbourCity := worldMap.GetCity(neighbours[1])
			if neighbourCity == nil {
				neighbourCity = types.NewCity(neighbours[1], 0)
				if err := worldMap.AddCity(neighbourCity); err != nil {
					return nil, nil, err
				}

				cities = append(cities, neighbourCity)
			}
			// Add neighbours to the respective city
			if err := city.AddNeighbour(neighbours[0], neighbourCity); err != nil {
				return nil, nil, errors.Wrap(err, "error adding neighbour")
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, nil, err
	}

	return worldMap, cities, nil
}

// PrintMap prints the leftout cities
func PrintMap(worldMap types.World) {
	fmt.Println("*****************************************")
	fmt.Println("Cities Left After Invasion")
	fmt.Println("*****************************************")

	for _, city := range worldMap {
		if city != nil {
			fmt.Println(city.String())
		}
	}

	fmt.Println()
}
