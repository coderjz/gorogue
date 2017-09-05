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
	level *Level
	x     int
	y     int
}

func NewGame() *Game {
	level := NewLevel()

	return &Game{
		level: level,
		x:     level.startX,
		y:     level.startY,
	}
}

func (g *Game) render() {
	termbox.Clear(backgroundColor, backgroundColor)

	for y, line := range g.level.cells {
		for x, cell := range line {
			if cell.visible {
				termbox.SetCell(x, y, cell.content, foregroundColor, backgroundColor)
			}
		}
	}

	termbox.SetCell(g.x, g.y, '@', foregroundColor, backgroundColor)
	termbox.Flush()
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
	newX := g.x + x
	newY := g.y + y

	if g.level.cells.get(newX, newY).content == WALL {
		return
	}

	g.x = newX
	g.y = newY
}
