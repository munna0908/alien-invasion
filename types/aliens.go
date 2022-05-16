package types

// Aliens map alien-id to city
type Aliens map[int]*City

// AddAlien creates a new entry in the map
func (a Aliens) AddAlien(id int, city *City) {
	a[id] = city
}

// GetAlien returns the city associated with the given alien Id
func (a Aliens) GetAlien(id int) *City {
	return a[id]
}

// DeleteAlien removes the entry from the map based on alien Id
func (a Aliens) DeleteAlien(id int) {
	delete(a, id)
}
