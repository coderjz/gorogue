package main

import (
	"log"
	"os"
	"time"

	"github.com/coderjz/gorogue/game"
	"github.com/nsf/termbox-go"
)

const animationSpeed = 10 * time.Millisecond

func main() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}

	defer termbox.Close()

	eventQueue := make(chan termbox.Event)
	go func() {
		for {
			eventQueue <- termbox.PollEvent()
		}
	}()

	log.SetFlags(log.Lshortfile)
	logFileName := "./logs"
	f, err := os.Create(logFileName)
	if err != nil {
		panic("Cannot make log file")
	}
	log.SetOutput(f)

	log.Printf("\n\nStarting game at %s", time.Now().Format("010206_030405"))

	//TODO: Check terminal size (termbox.Size(), if not big enough and output error message)
	//Maybe do that in the render itself or do it here with a check in the game loop?

	g := game.NewGame()
	g.UpdateFOV()
	g.Render()

	for {
		ev := <-eventQueue
		//TODO: If ev.Type is resize, render error message

		playerActed := false

		if ev.Type == termbox.EventKey {
			switch {
			case ev.Key == termbox.KeyArrowUp || ev.Ch == 'k':
				g.ClearMessages()
				playerActed = g.MovePlayer(game.UP)
			case ev.Key == termbox.KeyArrowDown || ev.Ch == 'j':
				g.ClearMessages()
				playerActed = g.MovePlayer(game.DOWN)
			case ev.Key == termbox.KeyArrowLeft || ev.Ch == 'h':
				g.ClearMessages()
				playerActed = g.MovePlayer(game.LEFT)
			case ev.Key == termbox.KeyArrowRight || ev.Ch == 'l':
				g.ClearMessages()
				playerActed = g.MovePlayer(game.RIGHT)
			case ev.Key == termbox.KeySpace:
				playerActed = g.ChangeFloor()
			case ev.Key == termbox.KeyEsc:
				return
			}
		}

		if !playerActed {
			time.Sleep(animationSpeed)
			continue
		}

		g.HealPlayerFromActions()
		g.UpdateFOV()
		g.UpdateMonsters()
		g.Render()
		time.Sleep(animationSpeed)
	}
}
