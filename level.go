package main

import (
	"bufio"
	"strings"
)

//Level represents one level of the dungeon
type Level struct {
	cells  Cells
	startX int
	startY int
}

//Cells is a type for a double array of cells
type Cells [][]Cell

func (c Cells) get(x, y int) Cell {
	return c[y][x]
}

//Cell represents a single cell or tile in the level
type Cell struct {
	content rune
	visible bool
}

//Recommended max size is 80x24. Allow bottom row for text and one spacer.
const maxX int = 80
const maxY int = 22

//Temporary for testing.. want to do this differently
const testLayout string = `
#######                   ######
#.....#   #############   #....#
#.....#   #...........#   #....#
#.....#   #.....#####.#   #....#
#.....#   #####.#   #.#####....#
###.###       #.#   #..........#
###.###       #.#   #######....#
#.....#     ###.##        ######
#.....#     #....#
#.....#######....#
#................#
##################
`

// NewLevel generates the level
func NewLevel() *Level {
	return &Level{
		cells:  buildLevel(),
		startX: 3,
		startY: 2,
	}
}

func buildLevel() [][]Cell {
	cells := make([][]Cell, maxY)
	for i := range cells {
		cells[i] = make([]Cell, maxX)
	}

	scanner := bufio.NewScanner(strings.NewReader(testLayout))
	y := 0
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) >= maxX {
			panic("Attempt to build level with longer line than is supported")
		}

		for x, ch := range line {
			cells[y][x] = Cell{
				content: ch,
				visible: true,
			}
		}
		y++
		if y == maxY {
			panic("Attempt to build level with more lines than is supported.")
		}
	}
	return cells
}
