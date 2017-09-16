package main

//PlayerLevel stores the information related to each level of the player
type playerLevel struct {
	level     int
	neededEXP int //Total for this level
	maxHP     int
}

var playerLevels = map[int]playerLevel{
	1: {
		level:     1,
		neededEXP: 0,
		maxHP:     20,
	},
	2: {
		level:     2,
		neededEXP: 8,
		maxHP:     32,
	},
	3: {
		level:     3,
		neededEXP: 20,
		maxHP:     43,
	},
	4: {
		level:     4,
		neededEXP: 45,
		maxHP:     58,
	},
	5: {
		level:     5,
		neededEXP: 80,
		maxHP:     71,
	},
	6: {
		level:     6,
		neededEXP: 125,
		maxHP:     90,
	},
}
