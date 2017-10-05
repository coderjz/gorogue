package main

import (
	"log"
	"os"
	"time"

	"github.com/coderjz/gorogue/game"
	"github.com/nsf/termbox-go"
)

const animationSpeed = 10 * time.Millisecond

const (
	StateIntroScrolling = iota
	StateIntroScrolled
	StateInstructions
	StateMainGame
	StateGameOverStarting
	StateGameOverMenuDisplayed
)

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

	state := StateIntroScrolling

	intro := game.NewIntro()
	intro.Render() //This is asynchronous
	intro.ScrollCompleted = func() {
		state = StateIntroScrolled
	}

	instructions := game.Instructions{}
	gameover := game.NewGameOver()
	gameover.MenuDisplayed = func() {
		state = StateGameOverMenuDisplayed
	}
	var mainGame *game.Game

	for {
		ev := <-eventQueue
		switch state {
		case StateIntroScrolling:
			if ev.Key == termbox.KeyEsc {
				return
			}
			intro.CompleteScrolling()
		case StateIntroScrolled:
			if ev.Key == termbox.KeyEsc {
				return
			}

			if ev.Key == termbox.KeyArrowUp || ev.Ch == 'k' {
				intro.SelectPrevChoice()
			} else if ev.Key == termbox.KeyArrowDown || ev.Ch == 'j' {
				intro.SelectNextChoice()
			} else if ev.Key == termbox.KeySpace || ev.Key == termbox.KeyEnter {
				switch intro.GetSelectedChoice() {
				case 0: //Start game
					mainGame = game.NewGame()
					mainGame.UpdateFOV()
					mainGame.Render()
					state = StateMainGame
				case 1: //Instructions
					instructions.Render()
					state = StateInstructions
				case 2: //Exit
					return
				}
			}
		case StateInstructions:
			if ev.Key == termbox.KeyEsc || ev.Key == termbox.KeySpace || ev.Key == termbox.KeyEnter {
				state = StateIntroScrolled
				intro.RenderScrolled()
			}
		case StateMainGame:
			playerActed := false
			//Copy logic here from main game logic
			switch {
			case ev.Key == termbox.KeyArrowUp || ev.Ch == 'k':
				mainGame.ClearMessages()
				playerActed = mainGame.MovePlayer(game.UP)
			case ev.Key == termbox.KeyArrowDown || ev.Ch == 'j':
				mainGame.ClearMessages()
				playerActed = mainGame.MovePlayer(game.DOWN)
			case ev.Key == termbox.KeyArrowLeft || ev.Ch == 'h':
				mainGame.ClearMessages()
				playerActed = mainGame.MovePlayer(game.LEFT)
			case ev.Key == termbox.KeyArrowRight || ev.Ch == 'l':
				mainGame.ClearMessages()
				playerActed = mainGame.MovePlayer(game.RIGHT)
			case ev.Key == termbox.KeySpace:
				playerActed = mainGame.ChangeFloor()
			case ev.Key == termbox.KeyEsc:
				return
			}

			if !playerActed {
				time.Sleep(animationSpeed)
				continue
			}

			mainGame.HealPlayerFromActions()
			mainGame.UpdateFOV()
			mainGame.UpdateMonsters()
			if mainGame.IsGameOver {
				mainGame.ClearMessages()
				mainGame.StopMessageChan()
				state = StateGameOverStarting
				go gameover.Render(mainGame.GetPlayerPos())
				continue
			}
			mainGame.Render()
			time.Sleep(animationSpeed)
		case StateGameOverStarting:
			if ev.Key == termbox.KeyEsc {
				return
			}
		case StateGameOverMenuDisplayed:
			if ev.Key == termbox.KeyEsc {
				return
			}

			if ev.Key == termbox.KeyArrowUp || ev.Ch == 'k' {
				gameover.SelectPrevChoice()
			} else if ev.Key == termbox.KeyArrowDown || ev.Ch == 'j' {
				gameover.SelectNextChoice()
			} else if ev.Key == termbox.KeySpace || ev.Key == termbox.KeyEnter {
				switch gameover.GetSelectedChoice() {
				case 0: //Start game
					mainGame = game.NewGame()
					mainGame.UpdateFOV()
					mainGame.Render()
					state = StateMainGame
				case 1: //Exit
					return
				}
			}
		}
	}
}
