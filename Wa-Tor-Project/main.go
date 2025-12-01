/*!
 * \file main.go
 * \brief Wa-Tor Simulation in Go
 *
 * This file implements the Wa-Tor predator-prey simulation.
 * The simulation contains fish and sharks on a toroidal grid.
 * Each chronon (time step) updates the world according to the rules:
 * - Fish move and reproduce
 * - Sharks move, hunt fish, reproduce, and starve
 */

package main

import (
	"fmt"
	"math/rand"
	"time"
)

/*!
 * \brief Species type for creatures.
 * Used to identify whether a cell is empty, a fish, or a shark.
 */
type Species int

const (
	Empty Species = iota ///< Empty cell
	Fish                 ///< Fish creature
	Shark                ///< Shark creature
)

/*!
 * \brief Represents an individual fish or shark.
 */
type Creature struct {
	Species   Species ///< Type of creature
	Age       int     ///< Age in chronons
	Energy    int     ///< Remaining energy (only for sharks)
	LastBreed int     ///< Chronons since last reproduction
}

/*!
 * \brief Represents the Wa-Tor simulation world.
 */
type World struct {
	Grid       [][]*Creature ///< 2D grid of creatures
	Size       int           ///< Width/Height of the square grid
	FishBreed  int           ///< Chronons needed for a fish to reproduce
	SharkBreed int           ///< Chronons needed for a shark to reproduce
	Starve     int           ///< Shark energy before starvation
}

/*!
 * \brief Main function to run the simulation.
 *
 * It initializes the world, places fish and sharks,
 * and iteratively processes chronons, printing the grid and population.
 */
func main() {
	fmt.Println("Wa-Tor Simulation:")

	// Simulation parameters
	params := struct {
		NumShark   int ///< Initial number of sharks
		NumFish    int ///< Initial number of fish
		FishBreed  int ///< Fish reproduction rate
		SharkBreed int ///< Shark reproduction rate
		Starve     int ///< Shark starvation time
		GridSize   int ///< Size of the square grid
	}{
		NumShark:   100,
		NumFish:    300,
		FishBreed:  3,
		SharkBreed: 10,
		Starve:     5,
		GridSize:   50,
	}

	rand.Seed(time.Now().UnixNano())

	// Create and initialize world
	world := createWorld(params.GridSize)
	initializeWorld(world, params)

	// Run simulation
	for chronon := 0; chronon < 10000; chronon++ {
		world = processChronon(world, params)

		// Count populations
		fishCount, sharkCount := countPopulation(world)

		// Print population and grid
		fmt.Printf("Chronon %d | Fish=%d | Sharks=%d\n", chronon, fishCount, sharkCount)
		printWorld(world)

		// Stop if all life extinct
		if fishCount == 0 && sharkCount == 0 {
			fmt.Println("All life extinct!")
			break
		}

		time.Sleep(100 * time.Millisecond)
	}
}

/*!
 * \brief Print the current state of the world grid.
 * \param world Pointer to the World to print.
 *
 * Symbols:
 * - '.' = empty cell
 * - 'F' = fish
 * - 'S' = shark
 */
func printWorld(world *World) {
	for y := 0; y < world.Size; y++ {
		for x := 0; x < world.Size; x++ {
			c := world.Grid[x][y]
			if c == nil {
				fmt.Print(". ")
			} else if c.Species == Fish {
				fmt.Print("F ")
			} else {
				fmt.Print("S ")
			}
		}
		fmt.Println()
	}
	fmt.Println()
}

/*!
 * \brief Create a new empty world of given size.
 * \param size Width/Height of the square grid.
 * \return Pointer to the newly created World.
 */
func createWorld(size int) *World {
	grid := make([][]*Creature, size)
	for i := range grid {
		grid[i] = make([]*Creature, size)
	}
	return &World{
		Grid: grid,
		Size: size,
	}
}

/*!
 * \brief Initialize the world with sharks and fish placed randomly.
 * \param world Pointer to the World to initialize.
 * \param params Simulation parameters.
 */
func initializeWorld(world *World, params struct {
	NumShark, NumFish, FishBreed, SharkBreed, Starve, GridSize int
}) {
	// Place sharks
	for i := 0; i < params.NumShark; i++ {
		for {
			x, y := rand.Intn(world.Size), rand.Intn(world.Size)
			if world.Grid[x][y] == nil {
				world.Grid[x][y] = &Creature{
					Species:   Shark,
					Energy:    params.Starve,
					LastBreed: 0,
				}
				break
			}
		}
	}

	// Place fish
	for i := 0; i < params.NumFish; i++ {
		for {
			x, y := rand.Intn(world.Size), rand.Intn(world.Size)
			if world.Grid[x][y] == nil {
				world.Grid[x][y] = &Creature{
					Species:   Fish,
					LastBreed: 0,
				}
				break
			}
		}
	}

	world.FishBreed = params.FishBreed
	world.SharkBreed = params.SharkBreed
	world.Starve = params.Starve
}

/*!
 * \brief Process one chronon (time step) for the world.
 * \param oldWorld Current state of the world.
 * \param params Simulation parameters.
 * \return Pointer to the new World state after processing.
 */
func processChronon(oldWorld *World, params struct {
	NumShark, NumFish, FishBreed, SharkBreed, Starve, GridSize int
}) *World {
	newWorld := createWorld(oldWorld.Size)
	newWorld.FishBreed = oldWorld.FishBreed
	newWorld.SharkBreed = oldWorld.SharkBreed
	newWorld.Starve = oldWorld.Starve

	for x := 0; x < oldWorld.Size; x++ {
		for y := 0; y < oldWorld.Size; y++ {
			creature := oldWorld.Grid[x][y]
			if creature == nil {
				continue
			}

			// Skip if already moved
			if newWorld.Grid[x][y] != nil {
				continue
			}

			creature.Age++
			creature.LastBreed++

			switch creature.Species {
			case Fish:
				processFish(oldWorld, newWorld, x, y, creature)
			case Shark:
				processShark(oldWorld, newWorld, x, y, creature)
			}
		}
	}

	return newWorld
}

/*!
 * \brief Process movement and reproduction of a fish.
 * \param oldWorld Current world state.
 * \param newWorld Next world state.
 * \param x X position of the fish.
 * \param y Y position of the fish.
 * \param fish Pointer to the fish Creature.
 */
func processFish(oldWorld, newWorld *World, x, y int, fish *Creature) {
	adjacent := getAdjacentPositions(x, y, oldWorld.Size)

	emptyCells := [][2]int{}
	for _, pos := range adjacent {
		if oldWorld.Grid[pos[0]][pos[1]] == nil &&
			newWorld.Grid[pos[0]][pos[1]] == nil {
			emptyCells = append(emptyCells, pos)
		}
	}

	if len(emptyCells) == 0 {
		newWorld.Grid[x][y] = fish
		return
	}

	newPos := emptyCells[rand.Intn(len(emptyCells))]
	newX, newY := newPos[0], newPos[1]

	if fish.LastBreed >= oldWorld.FishBreed {
		newWorld.Grid[x][y] = &Creature{
			Species:   Fish,
			LastBreed: 0,
		}
		newWorld.Grid[newX][newY] = fish
		fish.LastBreed = 0
	} else {
		newWorld.Grid[newX][newY] = fish
	}
}

/*!
 * \brief Process movement, hunting, and reproduction of a shark.
 * \param oldWorld Current world state.
 * \param newWorld Next world state.
 * \param x X position of the shark.
 * \param y Y position of the shark.
 * \param shark Pointer to the shark Creature.
 */
func processShark(oldWorld, newWorld *World, x, y int, shark *Creature) {
	shark.Energy--

	if shark.Energy <= 0 {
		return
	}

	adjacent := getAdjacentPositions(x, y, oldWorld.Size)

	// Look for fish to eat
	fishCells := [][2]int{}
	for _, pos := range adjacent {
		if oldWorld.Grid[pos[0]][pos[1]] != nil &&
			oldWorld.Grid[pos[0]][pos[1]].Species == Fish &&
			newWorld.Grid[pos[0]][pos[1]] == nil {
			fishCells = append(fishCells, pos)
		}
	}

	if len(fishCells) > 0 {
		newPos := fishCells[rand.Intn(len(fishCells))]
		newX, newY := newPos[0], newPos[1]

		shark.Energy = oldWorld.Starve

		if shark.LastBreed >= oldWorld.SharkBreed {
			newWorld.Grid[x][y] = &Creature{
				Species:   Shark,
				Energy:    oldWorld.Starve,
				LastBreed: 0,
			}
			newWorld.Grid[newX][newY] = shark
			shark.LastBreed = 0
		} else {
			newWorld.Grid[newX][newY] = shark
		}
		return
	}

	// Move to empty adjacent cell if no fish
	emptyCells := [][2]int{}
	for _, pos := range adjacent {
		if oldWorld.Grid[pos[0]][pos[1]] == nil &&
			newWorld.Grid[pos[0]][pos[1]] == nil {
			emptyCells = append(emptyCells, pos)
		}
	}

	if len(emptyCells) == 0 {
		newWorld.Grid[x][y] = shark
		return
	}

	newPos := emptyCells[rand.Intn(len(emptyCells))]
	newX, newY := newPos[0], newPos[1]

	if shark.LastBreed >= oldWorld.SharkBreed {
		newWorld.Grid[x][y] = &Creature{
			Species:   Shark,
			Energy:    oldWorld.Starve,
			LastBreed: 0,
		}
		newWorld.Grid[newX][newY] = shark
		shark.LastBreed = 0
	} else {
		newWorld.Grid[newX][newY] = shark
	}
}

/*!
 * \brief Get 4 adjacent positions with wrapping around edges.
 * \param x X coordinate.
 * \param y Y coordinate.
 * \param size Grid size.
 * \return Slice of 4 [x,y] coordinates.
 */
func getAdjacentPositions(x, y, size int) [][2]int {
	return [][2]int{
		{(x - 1 + size) % size, y}, // West
		{(x + 1) % size, y},        // East
		{x, (y - 1 + size) % size}, // North
		{x, (y + 1) % size},        // South
	}
}

/*!
 * \brief Count number of fish and sharks in the world.
 * \param world Pointer to the World.
 * \return fishCount Number of fish.
 * \return sharkCount Number of sharks.
 */
func countPopulation(world *World) (int, int) {
	fish, sharks := 0, 0
	for x := 0; x < world.Size; x++ {
		for y := 0; y < world.Size; y++ {
			if world.Grid[x][y] != nil {
				if world.Grid[x][y].Species == Fish {
					fish++
				} else {
					sharks++
				}
			}
		}
	}
	return fish, sharks
}
