package types

import (
	"errors"
)

var (
	ErrCityExists = errors.New("city already exists")
)

// World is mapping of city name to City instance
type World map[string]*City

func NewWorldMap() World {
	return make(map[string]*City)
}

// DeleteCity removes the city entry from the map
func (w World) DeleteCity(city string) {
	delete(w, city)
}

// AddCity adds new city to the map
func (w World) AddCity(city *City) error {
	if _, ok := w[city.Name]; ok {
		return ErrCityExists
	}

	w[city.Name] = city

	return nil
}

// GetCity returns the city associated with the city name
func (w World) GetCity(name string) *City {
	return w[name]
}
