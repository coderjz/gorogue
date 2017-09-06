package main

import (
	"math/rand"
	"time"
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

//Recommended max size is 80x24. 0-based indexing. Allow bottom row for text and one spacer.
const maxX int = 79
const maxY int = 21

const minNumRooms = 6
const maxNumRooms = 10
const minRoomSize = 4
const maxRoomSize = 10

var levelRand *rand.Rand

// NewLevel generates the level
func NewLevel() *Level {
	levelRand = rand.New(rand.NewSource(time.Now().UnixNano()))

	cells := make(Cells, maxY+1) //+1 needed because maxX, maxY are 0-based
	for i := range cells {
		cells[i] = make([]Cell, maxX+1)
	}

	rooms := generateRooms()
	convertRoomsToCells(rooms, &cells)

	startX, startY := rooms[0].getCenter()

	return &Level{
		cells:  cells,
		startX: startX,
		startY: startY,
	}
}

//Room represents a single room within a level
type Room struct {
	x1 int
	y1 int
	x2 int
	y2 int
}

func (r *Room) getCenter() (int, int) {
	return (r.x1 + r.x2) / 2,
		(r.y1 + r.y2) / 2
}

func (r *Room) intersects(r2 *Room) bool {
	return (r.x1 <= r2.x2 && r.x2 >= r2.x1 &&
		r.y1 <= r2.y2 && r.y2 >= r2.y1)
}

func convertRoomsToCells(rooms []*Room, cells *Cells) {
	for _, room := range rooms {
		for x := room.x1; x <= room.x2; x++ {
			for y := room.y1; y <= room.y2; y++ {
				if x == room.x1 || x == room.x2 ||
					y == room.y1 || y == room.y2 {
					//Wall
					(*cells)[y][x] = Cell{
						content: WALL,
						visible: true,
					}
				} else {
					//Floor
					(*cells)[y][x] = Cell{
						content: FLOOR,
						visible: true,
					}
				}
			}
		}
	}
}

func random(min int, max int, rand *rand.Rand) int {
	return rand.Intn(max-min) + min
}

func min(a, b int) int {
	if a <= b {
		return a
	}
	return b
}

func generateRooms() []*Room {
	numRooms := random(minNumRooms, maxNumRooms, levelRand)
	rooms := make([]*Room, numRooms)
	for i := 0; i < numRooms; i++ {
		x1 := random(0, maxX-minRoomSize, levelRand)
		y1 := random(0, maxY-minRoomSize, levelRand)
		x2 := min(maxX, x1+random(minRoomSize, maxRoomSize, levelRand))
		y2 := min(maxY, y1+random(minRoomSize, maxRoomSize, levelRand))
		r := &Room{
			x1: x1,
			y1: y1,
			x2: x2,
			y2: y2,
		}
		//Check if intersects any other room
		intersects := false
		for j := 0; j < i; j++ {
			if r.intersects(rooms[j]) {
				intersects = true
				break
			}
		}

		if intersects {
			i--
			continue
		}

		rooms[i] = r
	}
	return rooms
}
