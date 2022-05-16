package cli

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/munna0908/alien-invasion/simulation"
)

const (
	// DefaultIteration is 10000 as specified in the requirement document
	DefaultIterations = 10000
	DefaultAliens     = 0
)

var (
	maxIterations int
	alientsCount  int
	worldFilePath string
)

func init() {
	flag.IntVar(&maxIterations, "iterations", DefaultIterations, "Number of iterations")
	flag.IntVar(&alientsCount, "aliens", 0, "Number of aliens")
	flag.StringVar(&worldFilePath, "input-file", "", "Location of input world file")
	flag.Parse()
}

func validateFlags() error {
	if maxIterations <= 0 {
		return errors.New("invalid iterations")
	}

	if alientsCount <= 0 {
		return errors.New("invalid aliens count")
	}

	if len(worldFilePath) == 0 {
		return errors.New("invalid file path")
	}

	if _, err := os.Stat(worldFilePath); os.IsNotExist(err) {
		return errors.New("world file not found")
	}

	return nil
}

func Execute(closeCh chan os.Signal) {
	// Validate the flags
	if err := validateFlags(); err != nil {
		log.Printf("Error validating flags err=%s \n", err.Error())
		flag.Usage()

		return
	}

	// Build the world map
	worldMap, cities, err := simulation.BuildMap(worldFilePath)
	if err != nil {
		log.Printf("Error building world map err=%s \n", err.Error())

		return
	}
	// Assumption: Aliens_count <= 2*Cities_count
	if alientsCount > 2*len(cities) {
		log.Printf("Error invalid aliens count")

		return
	}
	// Create Simulation instance
	simulator, err := simulation.NewSimulation(worldMap, alientsCount, maxIterations)
	if err != nil {
		log.Printf("Error creating Simulation instance err=%s \n", err.Error())

		return
	}
	// Allocate aliens to the cities
	if err = simulator.InitAliens(cities, alientsCount); err != nil {
		log.Printf("Error initiating aliens err=%s \n", err.Error())

		return
	}

	fmt.Println("*****************************************")
	fmt.Println("Aliens Started Invasion ...!!! ")
	fmt.Println("*****************************************")

	// Start the simulation
	simulator.Run(closeCh)
	//Print the left over cities
	simulation.PrintMap(worldMap)
}
