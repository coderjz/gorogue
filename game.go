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
	g.updateFOV()

	termbox.Clear(backgroundColor, backgroundColor)

	for y, line := range g.level.cells {
		for x, cell := range line {
			if cell.visible {
				termbox.SetCell(x, y, cell.content, foregroundColor, backgroundColor)
			}
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

func (g *Game) move(dir Direction) {
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
		return
	}

	g.player.x = newX
	g.player.y = newY
}
