package main

import (
	"fmt"
	"math"
	"strconv"

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
	level    *Level
	player   *Player
	messages []string
}

func NewGame() *Game {
	level := NewLevel()

	return &Game{
		level:    level,
		player:   NewPlayer(level.startX, level.startY),
		messages: make([]string, 0, 2),
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

	//Render menu
	menuRow := 20
	nextLevelExp := strconv.Itoa(g.player.nextLevelEXP)
	if g.player.nextLevelEXP == math.MaxInt32 {
		nextLevelExp = "MAX"
	}
	menuString := fmt.Sprintf("HP: %d/%d EXP: %d/%s LVL: %d: ", g.player.hp, g.player.maxHP, g.player.exp, nextLevelExp, g.player.level)
	for i := 0; i < len(menuString); i++ {
		termbox.SetCell(i, menuRow, rune(menuString[i]), foregroundColor, backgroundColor)
	}

	//TODO: A better approach to render messages is to show each message for say a second and a half or so.
	//And if there are more messages left to then show those ones are shown for another second and a half or so.
	//But if player first clears messages (by say moving) then we would not show those messages.
	//Recommend to make a separate "GameMessages" struct that would take responsibility for this logic

	//Render the message logs
	messageRow := 22
	for i, m := range g.messages {
		//We only support showing up to two messages
		if i >= 2 {
			break
		}

		//Make sure string can't go outside of our range
		if len(m) > maxX+1 {
			m = m[:maxX+1]
		}
		for j, c := range m {
			termbox.SetCell(j, messageRow, c, foregroundColor, backgroundColor)
		}
		messageRow++
	}
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

	for i, m := range g.level.monsters {
		if m.x == newX && m.y == newY {
			//TODO: Better damage formula
			damage := 3
			m.hp -= damage
			if m.hp <= 0 {
				g.messages = append(g.messages, fmt.Sprintf("You killed %s.", m.name))

				//Just remove the monster from the array, it shouldn't re-render
				g.level.monsters = append(g.level.monsters[:i], g.level.monsters[i+1:]...)

				g.player.exp += m.exp
				if g.player.ProcessLevelUp() {
					g.messages = append(g.messages, fmt.Sprintf("You reached level %d!", g.player.level))
				}
			} else {
				g.messages = append(g.messages, fmt.Sprintf("You attacked %s for %d damage.", m.name, damage))
			}
			return true
		}
	}

	g.player.x = newX
	g.player.y = newY

	return true
}

func (g *Game) moveMonster(m *Monster) {
	x, y := g.determineMonsterMoveNewPos(m.x, m.y)

	if g.player.x == x && g.player.y == y {
		//TODO: Better calculate damage
		damage := 1
		g.player.hp -= damage
		if g.player.hp <= 0 {
			//TODO: Game over
		} else {
			g.messages = append(g.messages, fmt.Sprintf("%s has attacked you for %d damage.", m.name, damage))
		}
	} else {
		m.x = x
		m.y = y
	}
}

func (g *Game) updateMonsters() {
	for _, m := range g.level.monsters {
		if g.level.cells.get(m.x, m.y).visible {
			m.active = true
		}

		if m.active {
			g.moveMonster(m)
		}
	}
}

//Given the input (x, y), what is the best way to move towards the player
func (g *Game) determineMonsterMoveNewPos(x, y int) (int, int) {
	//For now, we'll just do it based off a naive check of player coordinates vs. monster coordinates
	//TODO: A* Search ?
	if g.player.x < x && g.monsterCanMoveTo(x-1, y) {
		return x - 1, y
	} else if g.player.x > x && g.monsterCanMoveTo(x+1, y) {
		return x + 1, y
	} else if g.player.y < y && g.monsterCanMoveTo(x, y-1) {
		return x, y - 1
	} else if g.player.y > y && g.monsterCanMoveTo(x, y+1) {
		return x, y + 1
	} else {
		return x, y
	}
}

func (g *Game) monsterCanMoveTo(x, y int) bool {
	c := g.level.cells.get(x, y)
	if c.content == WALL {
		return false
	}

	for _, m := range g.level.monsters {
		if m.x == x && m.y == y {
			return false
		}
	}
	return true
}

func (g *Game) clearMessages() {
	g.messages = g.messages[:0]
}
