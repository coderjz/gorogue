package main

import (
	"log"
	"math/rand"
	"time"
)

//Level represents one level of the dungeon
type Level struct {
	cells  Cells
	rooms  []*Room
	startX int
	startY int
}

//Cells is a type for a double array of cells
type Cells [][]Cell

func (c Cells) get(x, y int) Cell {
	return c[y][x]
}

func (c Cells) set(x, y int, cell Cell) {
	c[y][x] = cell
}

//Cell represents a single cell or tile in the level
type Cell struct {
	content rune
	visible bool
}

//Recommended max size is 80x24. 0-based indexing. Allow bottom row for text and one spacer.
const maxX int = 79
const maxY int = 21

const minRoomSize = 4
const maxRoomSize = 10
const numAttemptRooms = 20

var levelRand *rand.Rand

// NewLevel generates the level
func NewLevel() *Level {
	levelRand = rand.New(rand.NewSource(time.Now().UnixNano()))

	cells := make(Cells, maxY+1) //+1 needed because maxX, maxY are 0-based
	for i := range cells {
		cells[i] = make([]Cell, maxX+1)
		for j := range cells[i] {
			cells[i][j] = Cell{
				content: WALL,
				visible: false,
			}
		}
	}

	rooms := generateRooms()
	for i, r := range rooms {
		log.Printf("Generated room #%d: %+v\n", i, r)
	}
	convertRoomsToCells(rooms, &cells)

	for i := 0; i < len(rooms)-1; i++ {
		r1x, r1y := rooms[i].getCenter()
		r2x, r2y := rooms[i+1].getCenter()
		generateTunnel(r1x, r1y, r2x, r2y, &cells)
	}

	startX, startY := rooms[0].getCenter()

	return &Level{
		cells:  cells,
		rooms:  rooms,
		startX: startX,
		startY: startY,
	}
}

func (l *Level) roomContainsPoint(x, y int) *Room {
	for _, r := range l.rooms {
		if r.pointIntersects(x, y) {
			return r
		}
	}
	return nil
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

func (r *Room) pointIntersects(x, y int) bool {
	//Do not check equals because we don't want to intersect when we are along a wall
	return x > r.x1 && x < r.x2 && y > r.y1 && y < r.y2
}

func convertRoomsToCells(rooms []*Room, cells *Cells) {
	for _, room := range rooms {
		for x := room.x1 + 1; x < room.x2; x++ {
			for y := room.y1 + 1; y < room.y2; y++ {
				c := Cell{
					content: FLOOR,
					visible: false,
				}

				cells.set(x, y, c)
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
	rooms := make([]*Room, 0)
	for i := 0; i < numAttemptRooms; i++ {
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
		//Check if intersects any already added room
		intersects := false
		for _, existingRoom := range rooms {
			if r.intersects(existingRoom) {
				intersects = true
				continue
			}
		}

		if intersects {
			continue
		}
		rooms = append(rooms, r)
	}
	return rooms
}

func generateTunnel(x1, y1, x2, y2 int, cells *Cells) {
	if x1 == x2 {
		log.Printf("Generating only vertical tunnel from (%d, %d) to (%d, %d)\n", x1, y1, x2, y2)
		generateVertTunnel(x1, y1, y2, cells)
		return
	}

	if y1 == y2 {
		log.Printf("Generating only horizontal tunnel from (%d, %d) to (%d, %d)\n", x1, y1, x2, y2)
		generateHorizTunnel(x1, x2, y1, cells)
		return
	}

	dir := levelRand.Intn(2)
	if dir == 0 {
		log.Printf("Generating first horizontal tunnel from (%d, %d) to (%d, %d); then vertical from (%d, %d) to (%d, %d)\n", x1, y1, x2, y1, x2, y1, x2, y2)
		generateHorizTunnel(x1, x2, y1, cells)
		generateVertTunnel(x2, y1, y2, cells)
	} else {
		log.Printf("Generating first vertical tunnel from (%d, %d) to (%d, %d); then horizontal from (%d, %d) to (%d, %d)\n", x1, y1, x1, y2, x1, y2, x2, y2)
		generateVertTunnel(x1, y1, y2, cells)
		generateHorizTunnel(x1, x2, y2, cells)
	}
}

func generateHorizTunnel(x1, x2, y int, cells *Cells) {
	if x1 > x2 {
		x1, x2 = x2, x1
	}
	for i := x1; i <= x2; i++ {
		cells.set(i, y, Cell{
			content: FLOOR,
			visible: false,
		})
	}
}

func generateVertTunnel(x, y1, y2 int, cells *Cells) {
	if y1 > y2 {
		y1, y2 = y2, y1
	}
	for i := y1; i <= y2; i++ {
		cells.set(x, i, Cell{
			content: FLOOR,
			visible: false,
		})
	}
}
