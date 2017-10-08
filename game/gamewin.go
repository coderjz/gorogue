package game

import (
	"strings"

	termbox "github.com/nsf/termbox-go"
)

type GameWin struct {
	selectedChoice int
	choices        []string
	MenuDisplayed  func()
}

func NewGameWin() *GameWin {
	return &GameWin{
		selectedChoice: 0,
		choices:        []string{"Play again", "Exit"},
	}
}

func (g *GameWin) Render() {
	message := `_____.___.                __      __.__
\__  |   | ____  __ __   /  \    /  \__| ____
 /   |   |/  _ \|  |  \  \   \/\/   /  |/    \
 \____   (  <_> )  |  /   \        /|  |   |  \
 / ______|\____/|____/     \__/\  / |__|___|  /
 \/                             \/          \/ `

	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

	lines := strings.Split(message, "\n")

	termWidth, termHeight := termbox.Size()
	msgX := len(lines[1])
	msgY := len(lines)
	startX := (termWidth - msgX) / 2
	startY := (termHeight - msgY) / 4

	for y, line := range lines {
		for x, c := range line {
			termbox.SetCell(startX+x, startY+y, c, termbox.ColorWhite, termbox.ColorBlack)
		}
	}

	y := len(lines) + startY + 3
	for choiceIndex, choice := range g.choices {
		foreColor := termbox.ColorWhite
		if g.selectedChoice == choiceIndex {
			foreColor = termbox.ColorYellow
			termbox.SetCell(startX, y, '*', foreColor, termbox.ColorBlack)
		}
		for j, c := range choice {
			termbox.SetCell(startX+j+2, y, c, foreColor, termbox.ColorBlack)
		}
		y += 2
	}

	termbox.Flush()
}

func (g *GameWin) SelectPrevChoice() {
	g.selectedChoice = (g.selectedChoice - 1 + len(g.choices)) % len(g.choices)
	g.Render()
}

func (g *GameWin) SelectNextChoice() {
	g.selectedChoice = (g.selectedChoice + 1) % len(g.choices)
	g.Render()
}

func (g *GameWin) GetSelectedChoice() int {
	return g.selectedChoice
}
