package game

import (
	termbox "github.com/nsf/termbox-go"
)

type LeaveDialog struct {
	selectedChoice int
	numChoices     int
}

func NewLeaveDialog() *LeaveDialog {
	return &LeaveDialog{
		selectedChoice: 0,
		numChoices:     2,
	}
}

//Render displays the leave dialog on top of the existing content
//Unlike some of the other screen renders, we do NOT want to clear the whole terminal
func (d *LeaveDialog) Render() {
	lines := []string{
		"                       ",
		" Give up and go home?  ",
		"                       ",
		"   Stay                ",
		"   Go home             ",
		"                       ",
	}

	startX := 28
	startY := 6

	//Outside box. Need to do this because specifying corners in strings and putting in loop added extra spaces
	width := len(lines[0])
	height := len(lines)

	termbox.SetCell(startX, startY, '┌', termbox.ColorWhite, termbox.ColorBlack)
	termbox.SetCell(startX+width+1, startY, '┐', termbox.ColorWhite, termbox.ColorBlack)
	termbox.SetCell(startX, startY+height+1, '└', termbox.ColorWhite, termbox.ColorBlack)
	termbox.SetCell(startX+width+1, startY+height+1, '┘', termbox.ColorWhite, termbox.ColorBlack)
	for i := 1; i <= width; i++ {
		termbox.SetCell(startX+i, startY, '─', termbox.ColorWhite, termbox.ColorBlack)
		termbox.SetCell(startX+i, startY+height+1, '─', termbox.ColorWhite, termbox.ColorBlack)
	}
	for i := 1; i <= height; i++ {
		termbox.SetCell(startX, startY+i, '│', termbox.ColorWhite, termbox.ColorBlack)
		termbox.SetCell(startX+width+1, startY+i, '│', termbox.ColorWhite, termbox.ColorBlack)
	}

	//Content within box
	startX = startX + 1
	startY = startY + 1
	for y, line := range lines {
		for x, c := range line {
			foreColor := termbox.ColorWhite

			if y == (3 + d.selectedChoice) {
				if x > 0 && x < len(line)-1 {
					foreColor = termbox.ColorYellow
				}
				if x == 1 {
					c = '*'
				}
			}

			termbox.SetCell(startX+x, startY+y, c, foreColor, termbox.ColorBlack)
		}
	}
	termbox.Flush()
}

func (d *LeaveDialog) SelectPrevChoice() {
	d.selectedChoice = (d.selectedChoice - 1 + d.numChoices) % d.numChoices
	d.Render()
}

func (d *LeaveDialog) SelectNextChoice() {
	d.selectedChoice = (d.selectedChoice + 1) % d.numChoices
	d.Render()
}

func (d *LeaveDialog) GetSelectedChoice() int {
	return d.selectedChoice
}
