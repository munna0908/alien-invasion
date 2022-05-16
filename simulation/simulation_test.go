package simulation

import (
	"errors"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/munna0908/alien-invasion/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInitAliens_InvalidCities(t *testing.T) {
	tests := []struct {
		name       string
		shouldFail bool
		err        error
		cityCount  int
		alienCount int
	}{
		{
			name:       "City count less than zero",
			shouldFail: true,
			err:        ErrInvalidCityCount,
			cityCount:  -1,
			alienCount: 20,
		},
		{
			name:       "City count less equal to zero",
			shouldFail: true,
			err:        ErrInvalidCityCount,
			cityCount:  0,
			alienCount: 20,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test cities
			world, cities := createTestWorld(tt.cityCount)
			// Create test simulation instance
			simulation := createTestSimulation(world, tt.alienCount)
			// Assign aliens to cities
			err := simulation.InitAliens(cities, tt.alienCount)
			if tt.shouldFail {
				assert.ErrorIs(t, err, tt.err, "Invalid city count error expected")
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
func TestInitAliens_AlienCount(t *testing.T) {
	tests := []struct {
		name       string
		shouldFail bool
		err        error
		alienCount int
		cityCount  int
	}{
		{
			name:       "Aliens count greater than twice the city count",
			shouldFail: true,
			err:        ErrInvalidAliensCount,
			alienCount: 21,
			cityCount:  10,
		},
		{
			name:       "Alien count is zero",
			shouldFail: true,
			err:        ErrInvalidAliensCount,
			alienCount: 0,
			cityCount:  10,
		},
		{
			name:       "Alien count is less than zero",
			shouldFail: true,
			err:        ErrInvalidAliensCount,
			alienCount: -1,
			cityCount:  10,
		},
		{
			name:       "Valid alien count",
			shouldFail: false,
			err:        ErrInvalidAliensCount,
			alienCount: 20,
			cityCount:  10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test cities
			world, cities := createTestWorld(tt.cityCount)
			// Create test simulation instance
			simulation := createTestSimulation(world, tt.alienCount)
			// Assign aliens to cities
			err := simulation.InitAliens(cities, tt.alienCount)
			if tt.shouldFail {
				assert.ErrorIs(t, err, tt.err, "Invalid aliens count error expected")
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
func TestInitAliens_MaxOccupancy(t *testing.T) {
	testWorld, cities := createTestWorld(20)
	simulation, err := NewSimulation(testWorld, 40, 0)
	require.NoError(t, err)
	// initilize the aliens
	err = simulation.InitAliens(cities, 40)
	require.NoError(t, err)

	// cities occupancy should be <=2
	for _, city := range testWorld {
		if len(city.OccupiedAliens) > 2 {
			require.Less(t, len(city.OccupiedAliens), 3, "Max occupancy per city is expected to be <=2")
		}
	}
}

func TestCleanupAliens(t *testing.T) {
	testWorld, cities, err := createTestWorldWithNeighbours(2, [][]int{
		{1},
		{0},
	})
	if err != nil {
		t.Fatalf("Error creating test cities %s", err)
	}
	// create simulation instance
	simulation, err := NewSimulation(testWorld, 2, 0)
	require.NoError(t, err)
	// create test aliens
	testAliens := createTestAliens(2, cities)
	alienIds := make(map[int]interface{})

	for id, city := range testAliens {
		simulation.aliens.AddAlien(id, city)

		alienIds[id] = nil
	}
	// cleanup aliens
	simulation.cleanupAliens(alienIds)
	// verify
	for id := range testAliens {
		require.Nil(t, simulation.aliens.GetAlien(id))
	}
}

func TestCleanupRoads(t *testing.T) {
	testWorld, cities, err := createTestWorldWithNeighbours(4, [][]int{
		{1, 2, 3},
		{0, 2, 3},
		{0, 1, 3},
		{0},
	})
	if err != nil {
		t.Fatalf("Error creating test cities %s", err)
	}

	simulation, err := NewSimulation(testWorld, 4, 0)
	require.NoError(t, err)

	randomCity := pickRandomCity(cities)
	simulation.cleanupRoads(randomCity)

	for _, city := range testWorld {
		for _, neighbourCity := range city.Neighbours {
			require.NotEqual(t, neighbourCity, randomCity, "Destroyed up city still exits")
		}
	}
}

func createTestWorldWithNeighbours(citiesCount int, neighbours [][]int) (types.World, []*types.City, error) {
	cities := make([]*types.City, 0, citiesCount)
	world := types.NewWorldMap()

	if len(neighbours) < citiesCount {
		return nil, nil, errors.New("wrong number of neighbours")
	}

	for i := 0; i < citiesCount; i++ {
		city := &types.City{
			Name:       fmt.Sprintf("testCity_%d", i),
			Neighbours: make(map[types.Direction]*types.City, len(neighbours[i])),
		}
		cities = append(cities, city)
	}

	for j := 0; j < citiesCount; j++ {
		neighbour := neighbours[j]
		for k := 0; k < len(neighbour); k++ {
			cities[j].Neighbours[types.Direction(k)] = cities[neighbour[k]]
		}
		world.AddCity(cities[j]) //nolint
	}

	return world, cities, nil
}

func createTestAliens(count int, city []*types.City) map[int]*types.City {
	rand.Seed(time.Now().UnixNano())

	aliens := make(map[int]*types.City, count)

	for i := 0; i < count; i++ {
		aliens[i] = city[rand.Intn(len(city))] //nolint
	}

	return aliens
}

func createTestWorld(citiesCount int) (types.World, []*types.City) {
	cities := make([]*types.City, 0)
	world := types.NewWorldMap()

	for i := 0; i < citiesCount; i++ {
		city := &types.City{
			Name:       fmt.Sprintf("testCity_%d", i),
			Neighbours: make(map[types.Direction]*types.City),
		}
		cities = append(cities, city)
		world.AddCity(city) //nolint
	}

	return world, cities
}

func createTestSimulation(world types.World, aliensCount int) *Simulation {
	return &Simulation{
		worldMap: world,
		aliens:   make(map[int]*types.City, aliensCount),
	}
}
