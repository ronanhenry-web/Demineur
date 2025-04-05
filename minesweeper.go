package main

import (
	"math/rand"
	"net/http"
	_ "net/http/pprof"
	"strings"
	"time"
	"sync"
)

type Tile struct {
	isBomb      bool
	isUncovered bool
	isFlagged   bool
	nearbyBombs int
	x           int
	y           int
}

var RandomGenerator = rand.New(rand.NewSource(3))

func main() {

	go func() {
		println("pprof en écoute sur http://localhost:6060")
		http.ListenAndServe("localhost:6060", nil)
	}()
	
	size := 25
	bombCount := 75
	grid := generateGrid(size, bombCount)
	//displayGrid(grid, true)
	solve(grid, bombCount)

	time.Sleep(30 * time.Second)
}

func generateGrid(size int, bombCount int) [][]Tile {
	if size*size < bombCount {
		return nil
	}

	grid := make([][]Tile, size)

	for x := range size {
		grid[x] = make([]Tile, size)
		for y := range size {
			tile := Tile{}
			tile.x = x
			tile.y = y
			grid[x][y] = tile
		}
	}

	bombsPlaced := 0
	for bombsPlaced < bombCount {
		x := RandomGenerator.Intn(size)
		y := RandomGenerator.Intn(size)
		if !grid[x][y].isBomb {
			grid[x][y].isBomb = true
			forEachNeighbour(grid, x, y, func(tile Tile) {
				grid[tile.x][tile.y].nearbyBombs++
			})
			bombsPlaced++
		}
	}

	return grid
}

func solve(grid [][]Tile, bombCount int) {
	gridSize := len(grid)
	tilesCount := gridSize * gridSize
	emptyTilesCount := tilesCount - bombCount

	hasFailed := false
	x := RandomGenerator.Intn(gridSize - 1)
	y := RandomGenerator.Intn(gridSize - 1)

	for {
		//println("Played :", x, y)
		if grid[x][y].isBomb {
			hasFailed = true
			break
		}

		grid = uncoverTile(grid, x, y)
		grid = flagTiles(grid)
		//displayGrid(grid, false)
		var uncoveredTilesCount = 0
		for _, row := range grid {
			for _, tile := range row {
				if tile.isUncovered {
					uncoveredTilesCount++
				}
			}
		}
		if uncoveredTilesCount == emptyTilesCount {
			break
		}

		nextTile := getFirstSafeTile(grid)
		if nextTile == nil {
			hasFailed = true
			break
		}
		x = nextTile.x
		y = nextTile.y
	}

	//displayGrid(grid, false)

	if hasFailed {
		println("*BOOM*")
	} else {
		println("SUCCESS")
	}
}

func displayGrid(grid [][]Tile, showAll bool) {
	size := len(grid)

	for _, row := range grid {
		print(strings.Repeat("-", size*4+1))
		print("\n|")
		for _, tile := range row {
			if !showAll && tile.isFlagged {
				print(" ⚑ ")
			} else if !showAll && !tile.isUncovered {
				print("███")
			} else if tile.isBomb {
				print(" X ")
			} else if tile.nearbyBombs != 0 {
				print(" ")
				print(tile.nearbyBombs)
				print(" ")
			} else {
				print("   ")
			}
			print("|")
		}
		print("\n")
	}
	print(strings.Repeat("-", size*4+1))
	print("\n")
}

func forEachNeighbour(grid [][]Tile, x, y int, action func(tile Tile)) {
	size := len(grid)

	directions := [8][2]int{
		{-1, -1}, {-1, 0}, {-1, 1},
		{0, -1}, {0, 1},
		{1, -1}, {1, 0}, {1, 1},
	}

	for _, dir := range directions {
		nx, ny := x+dir[0], y+dir[1]
		if nx >= 0 && nx < size && ny >= 0 && ny < size {
			action(grid[nx][ny])
		}
	}
}

func getNeighboursLeft(grid [][]Tile, tile Tile) []Tile {
	var neighboursLeft []Tile
	forEachNeighbour(grid, tile.x, tile.y, func(tile Tile) {
		if !tile.isUncovered {
			neighboursLeft = append(neighboursLeft, tile)
		}
	})
	return neighboursLeft
}

func getNearbyFlaggedBombsCount(grid [][]Tile, tile Tile) int {
	flaggedBombsCount := 0
	forEachNeighbour(grid, tile.x, tile.y, func(tile Tile) {
		if tile.isFlagged {
			flaggedBombsCount++
		}
	})
	return flaggedBombsCount
}

func flagTiles(grid [][]Tile) [][]Tile {
	var wg sync.WaitGroup
	newFlag := true
	for newFlag {
		newFlag = false
		for _, row := range grid {
			for _, tile := range row {
				wg.Add(1)
				go func() {
					defer wg.Done()
					if !tile.isUncovered || tile.nearbyBombs == 0 {
						return
					}

					neighboursLeft := getNeighboursLeft(grid, tile)

					if len(neighboursLeft) == tile.nearbyBombs {
						for _, neighbour := range neighboursLeft {
							if !neighbour.isFlagged {
								newFlag = true
								grid[neighbour.x][neighbour.y].isFlagged = true
							}
						}
					}
				}()
			}
		}
		wg.Wait()
	}
	return grid
}

func getFirstSafeTile(grid [][]Tile) *Tile {
	for _, row := range grid {
		for _, tile := range row {
			if !tile.isUncovered || tile.nearbyBombs == 0 {
				continue
			}

			neighboursLeft := getNeighboursLeft(grid, tile)
			nearbyFlaggedBombsCount := getNearbyFlaggedBombsCount(grid, tile)
			if tile.nearbyBombs == nearbyFlaggedBombsCount &&
				len(neighboursLeft) > nearbyFlaggedBombsCount {
				for _, neighbour := range neighboursLeft {
					if !neighbour.isFlagged {
						return &grid[neighbour.x][neighbour.y]
					}
				}
			}
		}
	}
	return nil
}

func uncoverTile(grid [][]Tile, x int, y int) [][]Tile {
	grid[x][y].isUncovered = true
	tile := grid[x][y]

	if tile.nearbyBombs > 0 {
		return grid
	}

	forEachNeighbour(grid, x, y, func(tile Tile) {
		if !tile.isUncovered {
			grid = uncoverTile(grid, tile.x, tile.y)
		}
	})

	return grid
}
