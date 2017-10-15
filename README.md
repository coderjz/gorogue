# gorogue
Simple roguelike game in Go

![Introduction](assets/doc_1.png?raw=true)
![Gameplay](assets/doc_2.png?raw=true)


## Play
Requires go environment to be setup.

```
go get github.com/coderjz/gorogue
cd $GOPATH/src/github.com/coderjz/gorogue
make build 
./gorogue
```

## Purpose

It was recommended to make a simple text-based roguelike as a way to learn a new programming language. Decided to try this for Go and found it is a very effective and fun approach.

The code is intentionally light on comments and does not include unit tests, as I did not intend for this to become a large project requiring maintenance.  It is also more fun to test the simple logic by playing the game.

## Features

* Randomly generated dungeon layouts
* Multiple dungeon floors
* Simple monster AI
* Field of view based on rooms visited
* Player experience points and level ups
* Different monster types
* Animated introduction and game over screen
