package game

import (
	"fmt"
	"math"
	"strconv"
	"time"

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
const FLOOR_PREV rune = '>'
const FLOOR_NEXT rune = '<'

const numLevels int = 3

type Game struct {
	levels           []*Level
	currLevel        *Level
	currLevelPos     int
	player           *Player
	messages         []string
	clearMessageChan chan struct{}
	renderedOnce     bool
}

func NewGame() *Game {
	levels := make([]*Level, numLevels)
	for i := 0; i < numLevels; i++ {
		levels[i] = NewLevel()
	}

	return &Game{
		levels:           levels,
		currLevelPos:     0,
		currLevel:        levels[0],
		player:           NewPlayer(levels[0].prevFloorX, levels[0].prevFloorY),
		messages:         make([]string, 0, 2),
		clearMessageChan: make(chan struct{}),
		renderedOnce:     false,
	}
}

func (g *Game) Render() {
	if g.renderedOnce {
		g.clearMessageChan <- struct{}{} //Clear the message chan from last time
	}
	termbox.Clear(backgroundColor, backgroundColor)

	//Display dungeon tiles
	for y, line := range g.currLevel.cells {
		for x, cell := range line {
			if cell.visible {
				termbox.SetCell(x, y, cell.content, foregroundColor, backgroundColor)
			}
		}
	}

	//Add monsters on top of cells
	for _, m := range g.currLevel.monsters {
		c := g.currLevel.cells.get(m.x, m.y)
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

	//Render the message logs
	messageRow := 22
	nextMessageToShow := 0
	g.renderedOnce = true

	displayMessages := func() (done bool) {
		//First clear any old messages
		for i := 0; i < 2; i++ {
			for j := 0; j <= maxX; j++ {
				termbox.SetCell(j, messageRow+i, ' ', foregroundColor, backgroundColor)
			}
		}

		//We only support showing up to two messages at a time
		for i := 0; i < 2; i++ {
			if nextMessageToShow >= len(g.messages) {
				if i == 0 { //If i > 0, then we need to come back in again to clear the last message
					return true
				}
				return false
			}
			for j, c := range g.messages[nextMessageToShow] {
				termbox.SetCell(j, messageRow+i, c, foregroundColor, backgroundColor)
			}
			nextMessageToShow++
		}
		return false
	}

	displayMessages()
	termbox.Flush()
	go func() {
		//Ticker displays next messages after a delay.
		messageTicker := time.NewTicker(time.Millisecond * 1500)
		for {
			select {
			case <-messageTicker.C:
				if displayMessages() {
					messageTicker.Stop()
				}
				termbox.Flush()
			case <-g.clearMessageChan:
				//Stop showing more messages and clean up resources (by returning out of go routine)
				//This channel is sent data each time we render the next frame
				messageTicker.Stop()
				return
			}
		}
	}()
}

func (g *Game) UpdateFOV() {
	x, y := g.player.x, g.player.y
	r := g.currLevel.roomContainsPoint(x, y)
	if r != nil {
		for x = r.x1; x <= r.x2; x++ {
			for y = r.y1; y <= r.y2; y++ {
				g.currLevel.cells[y][x].visible = true
			}
		}
	} else {
		g.currLevel.cells[y][x].visible = true
		g.currLevel.cells[y][x+1].visible = true
		g.currLevel.cells[y][x-1].visible = true
		g.currLevel.cells[y+1][x].visible = true
		g.currLevel.cells[y-1][x].visible = true
	}
}

//Return value is if the move requested counts as a player action.
//Moving into a wall does not count as an action
func (g *Game) MovePlayer(dir Direction) bool {
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

	if g.currLevel.cells.get(newX, newY).content == WALL {
		return false
	}

	for i, m := range g.currLevel.monsters {
		if m.x == newX && m.y == newY {
			damage := calculateDamage(g.player.strength, m.defense)
			m.hp -= damage
			if m.hp <= 0 {
				g.messages = append(g.messages, fmt.Sprintf("You killed %s.", m.name))

				//Just remove the monster from the array, it shouldn't re-render
				g.currLevel.monsters = append(g.currLevel.monsters[:i], g.currLevel.monsters[i+1:]...)

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
		damage := calculateDamage(m.strength, g.player.defense)
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

func (g *Game) ChangeFloor() bool {
	if g.player.x == g.currLevel.prevFloorX && g.player.y == g.currLevel.prevFloorY {
		if g.currLevelPos <= 0 {
			//TODO - show dialog asking if player wants to leave. Special case to handle this (custom return value enum?)
			return false
		}
		g.currLevelPos--
		g.currLevel = g.levels[g.currLevelPos]
		g.player.x = g.currLevel.nextFloorX
		g.player.y = g.currLevel.nextFloorY
		return true
	} else if g.player.x == g.currLevel.nextFloorX && g.player.y == g.currLevel.nextFloorY {
		if g.currLevelPos >= len(g.levels)-1 {
			//Should never happen
			//TODO: do not render the next floor on the top most floor
			return false
		}
		g.currLevelPos++
		g.currLevel = g.levels[g.currLevelPos]
		g.player.x = g.currLevel.prevFloorX
		g.player.y = g.currLevel.prevFloorY
		return true
	}
	return false
}

func (g *Game) HealPlayerFromActions() {
	g.player.HealFromActions()
}

func (g *Game) UpdateMonsters() {
	for _, m := range g.currLevel.monsters {
		if g.currLevel.cells.get(m.x, m.y).visible {
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
	c := g.currLevel.cells.get(x, y)
	if c.content == WALL {
		return false
	}

	for _, m := range g.currLevel.monsters {
		if m.x == x && m.y == y {
			return false
		}
	}
	return true
}

func (g *Game) ClearMessages() {
	g.messages = g.messages[:0]
}
