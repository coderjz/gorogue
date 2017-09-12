package main

type Monster struct {
	x      int
	y      int
	hp     int
	name   string
	symbol rune //What is displayed in the cell
	active bool
}

var monsterList []*Monster

type MonsterType int

const (
	Page MonsterType = iota
	Squire
	Knight
	Commander
)

func init() {
	monsterList = []*Monster{
		{
			hp:     10,
			name:   "Page",
			symbol: 'P',
		},
		{
			hp:     25,
			name:   "Squire",
			symbol: 'S',
		},
		{
			hp:     47,
			name:   "Knight",
			symbol: 'K',
		},
		{
			hp:     100,
			name:   "Commander",
			symbol: 'C',
		},
	}

	for _, m := range monsterList {
		m.active = false
	}
}

//NewMonster creates a new monster instance based on the monster type
func NewMonster(monsterType MonsterType) *Monster {
	monster := &Monster{}
	switch monsterType {
	case Page:
		*monster = *monsterList[0]
	case Squire:
		*monster = *monsterList[1]
	case Knight:
		*monster = *monsterList[2]
	case Commander:
		*monster = *monsterList[3]
	}
	return monster
}
