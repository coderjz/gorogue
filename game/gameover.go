package game

import (
	"strings"
	"time"

	termbox "github.com/nsf/termbox-go"
)

type GameOver struct {
	selectedChoice int
	choices        []string
	MenuDisplayed  func()
}

func NewGameOver() *GameOver {
	return &GameOver{
		selectedChoice: 0,
		choices:        []string{"Play again", "Exit"},
	}
}

func (g *GameOver) Render(finalX, finalY int) {
	g.colorGraduallyRed(finalX, finalY)

	time.Sleep(1000 * time.Millisecond)
	g.selectedChoice = 0
	g.renderMenu()

	if g.MenuDisplayed != nil {
		g.MenuDisplayed()
	}
}

func (g *GameOver) colorGraduallyRed(finalX, finalY int) {
	ch1, count1 := g.generateFromUpperLeft(finalX, finalY)
	ch2, count2 := g.generateFromLowerRight(finalX, finalY)

	//Avoid division by 0
	if count1 == 0 {
		count1 = 1
	}
	if count2 == 0 {
		count2 = 1
	}

	//We want to render the red over 50 frames
	numCh1 := 1 + (count1 / 50)
	numCh2 := 1 + (count2 / 50)

	for {
		var p1 *point
		var p2 *point
		for i := 0; i < numCh1; i++ {
			p1 = <-ch1
			if p1 == nil {
				break
			}
			termbox.SetCell(p1.x, p1.y, ' ', termbox.ColorWhite, termbox.ColorRed)
		}

		for i := 0; i < numCh2; i++ {
			p2 = <-ch2
			if p2 == nil {
				break
			}
			termbox.SetCell(p2.x, p2.y, ' ', termbox.ColorWhite, termbox.ColorRed)
		}

		if p1 == nil && p2 == nil {
			break
		}
		termbox.Flush()
		time.Sleep(10 * time.Millisecond)
	}
	termbox.Flush()

	//Last cell, slow down for emphasis
	time.Sleep(500 * time.Millisecond)
	termbox.SetCell(finalX, finalY, ' ', termbox.ColorWhite, termbox.ColorRed)
	termbox.Flush()
}

type point struct {
	x int
	y int
}

func (g *GameOver) generateFromUpperLeft(finalX, finalY int) (ch chan *point, count int) {
	ch = make(chan *point)
	_, termHeight := termbox.Size()
	go func() {
		defer close(ch)
		for x := 0; x < finalX; x++ {
			for y := 0; y < termHeight; y++ {
				ch <- &point{x: x, y: y}
			}
		}

		for y := 0; y < finalY; y++ {
			ch <- &point{x: finalX, y: y}
		}
	}()
	count = (termHeight * finalX) + finalY
	return ch, count
}

func (g *GameOver) generateFromLowerRight(finalX, finalY int) (ch chan *point, count int) {
	ch = make(chan *point)
	termWidth, termHeight := termbox.Size()
	go func() {
		defer close(ch)
		for x := termWidth - 1; x > finalX; x-- {
			for y := termHeight - 1; y >= 0; y-- {
				ch <- &point{x: x, y: y}
			}
		}

		for y := termHeight - 1; y > finalY; y-- {
			ch <- &point{x: finalX, y: y}
		}
	}()
	count = (termHeight*termWidth - finalX - 1) + termHeight - finalY - 1
	return ch, count
}

func (g *GameOver) renderMenu() {
	message := `  ________                         ________
 /  _____/_____    _____   ____    \__  __ \___  __ ___________
/   \  ___\__  \  /     \_/ __ \    /  | |  \  \/ // __ \_  __ \
\    \_\  \/ __ \|  Y Y  \  ___/   /   |_|   \   /\  ___/|  | \/
 \______  (____  /__|_|  /\___  >  \_______  /\_/  \___  >__|
        \/     \/      \/     \/           \/          \/`

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

func (g *GameOver) SelectPrevChoice() {
	g.selectedChoice = (g.selectedChoice - 1 + len(g.choices)) % len(g.choices)
	g.renderMenu()
}

func (g *GameOver) SelectNextChoice() {
	g.selectedChoice = (g.selectedChoice + 1) % len(g.choices)
	g.renderMenu()
}

func (g *GameOver) GetSelectedChoice() int {
	return g.selectedChoice
}
