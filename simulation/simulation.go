package simulation

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"os"

	"github.com/munna0908/alien-invasion/types"
)

var (
	ErrInvalidCityCount   = errors.New("cities count too low")
	ErrInvalidAliensCount = errors.New("invalid aliens count")
)

// Simulation simulates the alien invasion on the given cities
type Simulation struct {
	count         int
	maxIterations int
	worldMap      types.World
	aliens        types.Aliens
}

func NewSimulation(worldMap types.World, aliensCount, maxIterations int) (*Simulation, error) {
	if len(worldMap) <= 0 {
		return nil, ErrInvalidCityCount
	}

	// Assumption: 0 < Aliens_count <= 2*cities_count
	if aliensCount <= 0 || aliensCount > 2*len(worldMap) {
		return nil, ErrInvalidAliensCount
	}

	return &Simulation{
		count:         0,
		worldMap:      worldMap,
		maxIterations: maxIterations,
		aliens:        make(map[int]*types.City, aliensCount),
	}, nil
}

// InitAliens allocates the aliens to random cities
func (s *Simulation) InitAliens(cities []*types.City, aliensCount int) error {
	if len(cities) <= 0 {
		return ErrInvalidCityCount
	}

	// Assumption: 0 < Aliens_count <= 2*cities_count
	if aliensCount <= 0 || aliensCount > 2*len(cities) {
		return ErrInvalidAliensCount
	}

	for alienID := 0; alienID < aliensCount; {
		city := pickRandomCity(cities)
		if city == nil {
			continue
		}

		if city.OccupiedAliens == nil {
			city.OccupiedAliens = make(map[int]interface{})
		}

		// Only two aliens are allocated per city.
		if len(city.OccupiedAliens) >= 2 {
			continue
		}

		city.AddAlien(alienID)
		s.aliens.AddAlien(alienID, city)
		alienID++
	}

	return nil
}

// CanContinue checks
func (s *Simulation) CanContinue() bool {
	if s.count >= s.maxIterations || len(s.aliens) == 0 || len(s.worldMap) == 0 {
		return false
	}

	return true
}

// Run starts the alien invasion
func (s *Simulation) Run(closeCh chan os.Signal) {
	s.checkForFight()

	defer func() { fmt.Println("Aliens left", len(s.aliens)) }()

	for s.CanContinue() {
		select {
		case <-closeCh:
			fmt.Println("*****************************************")
			fmt.Println("Stopping the invasion")
			fmt.Println("*****************************************")

			return
		default:
			s.moveAliens()
			s.count++
		}
	}
}

// cleanupAliens removes the aliens from the alien map
func (s *Simulation) cleanupAliens(aliens map[int]interface{}) []int {
	keys := make([]int, 0, len(aliens))
	for alien := range aliens {
		keys = append(keys, alien)
		s.aliens.DeleteAlien(alien)
	}

	return keys
}

// checkForFight checks for a fight between aliens, in case of a fight the city will be destroyed
func (s *Simulation) checkForFight() {
	for _, currentCity := range s.aliens {
		if len(currentCity.OccupiedAliens) > 1 {
			s.distroyCity(currentCity)

			continue
		}
	}
}

// moveAliens picks a random neighbour and moves the alien, in case of a fight the city is destroyed
func (s *Simulation) moveAliens() {
	for alien, currentCity := range s.aliens {
		// Get random neighbour
		newCity, err := currentCity.PickRandomNeighbours()
		if err != nil {
			// Alien is trapped
			continue
		}

		// Move the alien and update occupancy
		newCity.AddAlien(alien)
		s.aliens.AddAlien(alien, newCity)
		delete(s.worldMap[currentCity.Name].OccupiedAliens, alien)

		// If an alien exists in the chosen city, than city can be destroyed in the same iteration.
		if len(newCity.OccupiedAliens) > 1 {
			s.distroyCity(newCity)

			continue
		}
	}
}

// distroyCity deletes the city and associated roads,aliens
func (s *Simulation) distroyCity(city *types.City) {
	// Cleanup the linking roads
	s.cleanupRoads(city)
	// Delete the aliens
	aliens := s.cleanupAliens(city.OccupiedAliens)
	// Delete the city from world map
	s.worldMap.DeleteCity(city.Name)
	fmt.Printf("%s has been destroyed by alien %d and alien %d ! \n", city.Name, aliens[0], aliens[1])
}

// cleanupRoads removes all the inward/outward links
func (s *Simulation) cleanupRoads(c *types.City) {
	c.Neighbours = nil

	for _, city := range s.worldMap {
		for direction, neighbourCity := range city.Neighbours {
			if neighbourCity != nil && neighbourCity.Name == c.Name {
				city.Neighbours[direction] = nil
			}
		}
	}
}

// pickRandomCity chooses a random city from the given set of cities
func pickRandomCity(cities []*types.City) *types.City {
	randomBigInt, err := rand.Int(rand.Reader, big.NewInt(int64(len(cities))))
	if err != nil {
		return nil
	}

	return cities[randomBigInt.Int64()]
}
