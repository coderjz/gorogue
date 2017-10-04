package game

import (
	"strings"

	termbox "github.com/nsf/termbox-go"
)

type Instructions struct{}

func (i *Instructions) Render() {
	instructions := `                                  INSTRUCTIONS

Objective: Obtain the chalice of riches and get out alive

Symbol          Name
P               Page
S               Squire
K               Knight
C               Commander
*               Chalice of riches
<               Stairs for next floor
>               Stairs for previous floor

Key              Action
Up OR k          Move up
Down OR j        Move down
Left OR h        Move left
Right OR l       Move right
Escape           Exit game
Space            Use stairs, pick up chalice of riches
`

	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	for y, line := range strings.Split(instructions, "\n") {
		for x, c := range line {
			foreColor := termbox.ColorWhite
			if y == 9 && x == 0 {
				foreColor = termbox.ColorYellow
			}
			if y >= 5 && y <= 8 && x == 0 {
				foreColor = termbox.ColorRed
			}
			termbox.SetCell(x, y, c, foreColor, termbox.ColorBlack)
		}
	}
	termbox.Flush()

}
