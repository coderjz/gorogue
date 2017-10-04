package game

import (
	"time"

	termbox "github.com/nsf/termbox-go"
)

var header = [][]rune{
	[]rune{' ', ' ', '_', '_', '_', '_', '_', '_', '_', '_', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' '},
	[]rune{' ', '/', ' ', ' ', '_', '_', '_', '_', '_', '/', ' ', ' ', '_', '_', '_', '_', ' ', ' ', ' ', ' ', ' ', '_', '_', '_', '_', '_', '_', '_', ' ', ' ', '_', '_', '_', '_', ' ', ' ', ' ', '_', '_', '_', '_', ' ', ' ', '_', '_', ' ', '_', '_', ' ', ' ', ' ', '_', '_', '_', '_', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' '},
	[]rune{'/', ' ', ' ', ' ', '\\', ' ', ' ', '_', '_', '_', ' ', '/', ' ', ' ', '_', ' ', '\\', ' ', ' ', ' ', ' ', '\\', '_', ' ', ' ', '_', '_', ' ', '\\', '/', ' ', ' ', '_', ' ', '\\', ' ', '/', ' ', '_', '_', '_', '\\', '|', ' ', ' ', '|', ' ', ' ', '\\', '_', '/', ' ', '_', '_', ' ', '\\', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' '},
	[]rune{'\\', ' ', ' ', ' ', ' ', '\\', '_', '\\', ' ', ' ', '(', ' ', ' ', '<', '_', '>', ' ', ')', ' ', ' ', ' ', ' ', '|', ' ', ' ', '|', ' ', '\\', '(', ' ', ' ', '<', '_', '>', ' ', ')', ' ', '/', '_', '/', ' ', ' ', '>', ' ', ' ', '|', ' ', ' ', '/', '\\', ' ', ' ', '_', '_', '_', '/', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' '},
	[]rune{' ', '\\', '_', '_', '_', '_', '_', '_', ' ', ' ', '/', '\\', '_', '_', '_', '_', '/', ' ', ' ', ' ', ' ', ' ', '|', '_', '_', '|', ' ', ' ', ' ', '\\', '_', '_', '_', '_', '/', '\\', '_', '_', '_', ' ', ' ', '/', '|', '_', '_', '_', '_', '/', ' ', ' ', '\\', '_', '_', '_', ' ', ' ', '>', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' '},
	[]rune{' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', '\\', '/', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', '/', '_', '_', '_', '_', '_', '/', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', '\\', '/', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' '},
	[]rune{' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' '},
	[]rune{' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', '.', '_', '_', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' '},
	[]rune{' ', ' ', '_', '_', '_', '_', '_', '_', '_', '_', '_', '_', '_', ' ', ' ', ' ', ' ', ' ', ' ', ' ', '_', '_', '_', '_', ' ', ' ', ' ', '_', '_', '_', '_', ' ', ' ', ' ', ' ', ' ', '|', ' ', ' ', '|', '_', '_', ' ', ' ', ' ', '_', '_', '_', '_', ' ', ' ', ' ', '_', '_', '_', '_', '_', ' ', ' ', ' ', '_', '_', '_', '_', ' ', ' '},
	[]rune{' ', '/', ' ', ' ', '_', ' ', '\\', '_', ' ', ' ', '_', '_', ' ', '\\', ' ', ' ', ' ', ' ', ' ', '/', ' ', '_', '_', '_', '\\', ' ', '/', ' ', ' ', '_', ' ', '\\', ' ', ' ', ' ', ' ', '|', ' ', ' ', '|', ' ', ' ', '\\', ' ', '/', ' ', ' ', '_', ' ', '\\', ' ', '/', ' ', ' ', ' ', ' ', ' ', '\\', '_', '/', ' ', '_', '_', ' ', '\\', ' '},
	[]rune{'(', ' ', ' ', '<', '_', '>', ' ', ')', ' ', ' ', '|', ' ', '\\', '/', ' ', ' ', ' ', ' ', '/', ' ', '/', '_', '/', ' ', ' ', '>', ' ', ' ', '<', '_', '>', ' ', ')', ' ', ' ', ' ', '|', ' ', ' ', ' ', 'Y', ' ', ' ', '(', ' ', ' ', '<', '_', '>', ' ', ')', ' ', ' ', 'Y', ' ', 'Y', ' ', ' ', '\\', ' ', ' ', '_', '_', '_', '/', ' '},
	[]rune{' ', '\\', '_', '_', '_', '_', '/', '|', '_', '_', '|', ' ', ' ', ' ', ' ', ' ', ' ', ' ', '\\', '_', '_', '_', ' ', ' ', '/', ' ', '\\', '_', '_', '_', '_', '/', ' ', ' ', ' ', ' ', '|', '_', '_', '_', '|', ' ', ' ', '/', '\\', '_', '_', '_', '_', '/', '|', '_', '_', '|', '_', '|', ' ', ' ', '/', '\\', '_', '_', '_', ' ', ' ', '>'},
	[]rune{' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', '/', '_', '_', '_', '_', '_', '/', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', '\\', '/', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ', '\\', '/', ' ', ' ', ' ', ' ', ' ', '\\', '/', ' '},
}

type Intro struct {
	isScrolling     bool
	startX          int
	startY          int
	selectedChoice  int
	choices         []string
	ScrollCompleted func()
}

func NewIntro() *Intro {
	termWidth, termHeight := termbox.Size()
	introWidth := len(header[0])
	introHeight := len(header) + 8 //+ 8 for the menu choices being added

	return &Intro{
		isScrolling:    true,
		startX:         (termWidth - introWidth) / 2,
		startY:         (termHeight - introHeight) / 4,
		selectedChoice: 0,
		choices:        []string{"New Game", "Instructions", "Exit"},
	}
}

func (i *Intro) Render() {
	_, termHeight := termbox.Size()

	go func() {
		for y := termHeight - 1; y > i.startY; y-- {
			i.renderScrolling(i.startX, y)
			time.Sleep(20 * time.Millisecond)
			if !i.isScrolling {
				return
			}
		}
		i.CompleteScrolling()
	}()
}

func (i *Intro) CompleteScrolling() {
	i.isScrolling = false
	if i.ScrollCompleted != nil {
		i.ScrollCompleted()
	}
	i.RenderScrolled()
}

func (i *Intro) renderScrolling(startX, startY int) {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	for y, row := range header {
		for x, c := range row {
			termbox.SetCell(startX+x, startY+y, c, termbox.ColorWhite, termbox.ColorBlack)
		}
	}
	termbox.Flush()
}

func (i *Intro) RenderScrolled() {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	for y, row := range header {
		for x, c := range row {
			termbox.SetCell(i.startX+x, i.startY+y, c, termbox.ColorWhite, termbox.ColorBlack)
		}
	}

	y := len(header) + i.startY + 3
	for choiceIndex, choice := range i.choices {
		foreColor := termbox.ColorWhite
		if i.selectedChoice == choiceIndex {
			foreColor = termbox.ColorYellow
			termbox.SetCell(i.startX, y, '*', foreColor, termbox.ColorBlack)
		}
		for j, c := range choice {
			termbox.SetCell(i.startX+j+2, y, c, foreColor, termbox.ColorBlack)
		}
		y += 2
	}
	termbox.Flush()
}

func (i *Intro) SelectPrevChoice() {
	i.selectedChoice = (i.selectedChoice - 1 + len(i.choices)) % len(i.choices)
	i.RenderScrolled()
}

func (i *Intro) SelectNextChoice() {
	i.selectedChoice = (i.selectedChoice + 1) % len(i.choices)
	i.RenderScrolled()
}

func (i *Intro) GetSelectedChoice() int {
	return i.selectedChoice
}
