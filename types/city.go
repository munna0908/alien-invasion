package types

import (
	"errors"
	"fmt"
	"math/rand"
	"strings"
)

var (
	ErrInvalidDirection = errors.New("invalid direction")
	ErrNoNeighbours     = errors.New("no neighbour")
)

//Direction in an integer represetation of a real world direction
type Direction int

const (
	North Direction = iota
	South
	East
	West
)

// City maintains the links to the neighbouring cities and alien occupancy
type City struct {
	Name           string
	Neighbours     map[Direction]*City
	OccupiedAliens map[int]interface{}
}

func NewCity(name string, neighboursCount int) *City {
	return &City{
		Name:       name,
		Neighbours: make(map[Direction]*City, neighboursCount),
	}
}

// GetNeighbours returns a non nil random neighbour
func (c *City) PickRandomNeighbours() (*City, error) {
	validNeigbours := make([]*City, 0)

	for _, city := range c.Neighbours {
		if city != nil {
			validNeigbours = append(validNeigbours, city)
		}
	}

	if len(validNeigbours) > 0 {
		return validNeigbours[rand.Intn(len(validNeigbours))], nil
	}

	return nil, ErrNoNeighbours
}

// AddAlien adds the alien to the city
func (c *City) AddAlien(alien int) {
	if c.OccupiedAliens == nil {
		c.OccupiedAliens = make(map[int]interface{})
	}

	c.OccupiedAliens[alien] = nil
}

// AddNeighbour adds the given city as neighbour if the direction is valid
func (c *City) AddNeighbour(direction string, city *City) error {
	if city.Neighbours == nil {
		city.Neighbours = make(map[Direction]*City)
	}

	switch strings.ToLower(direction) {
	case "north":
		c.Neighbours[North] = city
	case "south":
		c.Neighbours[South] = city
	case "east":
		c.Neighbours[East] = city
	case "west":
		c.Neighbours[West] = city
	default:
		return ErrInvalidDirection
	}

	return nil
}

// GetDirection returns the string representation of the given direction
func GetDirection(direction Direction) string {
	switch direction {
	case North:
		return "north"
	case South:
		return "south"
	case East:
		return "east"
	case West:
		return "west"
	}

	return "invalid direction"
}

// String implements the stringer interface
func (c *City) String() string {
	neighbours := ""

	for direction, city := range c.Neighbours {
		if city != nil {
			neighbours += fmt.Sprintf("%s=%s ", GetDirection(direction), city.Name)
		}
	}

	return fmt.Sprintf("%s %s", c.Name, neighbours)
}
