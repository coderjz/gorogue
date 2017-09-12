package main

import (
	termbox "github.com/nsf/termbox-go"
)

type Direction int

const (
	UP Direction = iota
	DOWN
	LEFT
	RIGHT
)

const backgroundColor = termbox.ColorBlack
const foregroundColor = termbox.ColorWhite
const monsterForegroundColor = termbox.ColorRed

const WALL rune = '#'
const GUY rune = '@'
const FLOOR rune = '.'
const ENEMY rune = 'x'

type Game struct {
	level  *Level
	player *Player
}

func NewGame() *Game {
	level := NewLevel()

	return &Game{
		level:  level,
		player: NewPlayer(level.startX, level.startY),
	}
}

func (g *Game) render() {
	termbox.Clear(backgroundColor, backgroundColor)

	//Display dungeon tiles
	for y, line := range g.level.cells {
		for x, cell := range line {
			if cell.visible {
				termbox.SetCell(x, y, cell.content, foregroundColor, backgroundColor)
			}
		}
	}

	//Add monsters on top of cells
	for _, m := range g.level.monsters {
		c := g.level.cells.get(m.x, m.y)
		if c.visible {
			termbox.SetCell(m.x, m.y, m.symbol, monsterForegroundColor, backgroundColor)
		}
	}

	termbox.SetCell(g.player.x, g.player.y, g.player.content, foregroundColor, backgroundColor)
	termbox.Flush()
}

func (g *Game) updateFOV() {
	x, y := g.player.x, g.player.y
	r := g.level.roomContainsPoint(x, y)
	if r != nil {
		for x = r.x1; x <= r.x2; x++ {
			for y = r.y1; y <= r.y2; y++ {
				g.level.cells[y][x].visible = true
			}
		}
	} else {
		g.level.cells[y][x].visible = true
		g.level.cells[y][x+1].visible = true
		g.level.cells[y][x-1].visible = true
		g.level.cells[y+1][x].visible = true
		g.level.cells[y-1][x].visible = true
	}
}

//Return value is if the move requested counts as a player action.
//Moving into a wall does not count as an action
func (g *Game) movePlayer(dir Direction) bool {
	x := 0
	y := 0
	switch dir {
	case UP:
		y = -1
	case DOWN:
		y = 1
	case LEFT:
		x = -1
	case RIGHT:
		x = 1
	}
	newX := g.player.x + x
	newY := g.player.y + y

	if g.level.cells.get(newX, newY).content == WALL {
		return false
	}

	g.player.x = newX
	g.player.y = newY

	return true
}

func (g *Game) moveMonster(m *Monster, dir Direction) {
	x := 0
	y := 0
	switch dir {
	case UP:
		y = -1
	case DOWN:
		y = 1
	case LEFT:
		x = -1
	case RIGHT:
		x = 1
	}
	newX := m.x + x
	newY := m.y + y

	if g.level.cells.get(newX, newY).content == WALL {
		return
	}

	m.x = newX
	m.y = newY
}

func (g *Game) updateMonsters() {
	for _, m := range g.level.monsters {
		if g.level.cells.get(m.x, m.y).visible {
			m.active = true
		}

		if m.active {
			g.moveMonster(m, LEFT)
		}
	}

}
