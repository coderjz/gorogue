package game

//PlayerLevel stores the information related to each level of the player
type playerLevel struct {
	level     int
	neededEXP int //Total for this level
	maxHP     int
	strength  int
	defense   int
}

var playerLevels = map[int]playerLevel{
	1: {
		level:     1,
		neededEXP: 0,
		maxHP:     20,
		strength:  10,
		defense:   8,
	},
	2: {
		level:     2,
		neededEXP: 8,
		maxHP:     32,
		strength:  13,
		defense:   10,
	},
	3: {
		level:     3,
		neededEXP: 20,
		maxHP:     43,
		strength:  16,
		defense:   12,
	},
	4: {
		level:     4,
		neededEXP: 45,
		maxHP:     58,
		strength:  19,
		defense:   14,
	},
	5: {
		level:     5,
		neededEXP: 80,
		maxHP:     71,
		strength:  22,
		defense:   16,
	},
	6: {
		level:     6,
		neededEXP: 125,
		maxHP:     90,
		strength:  25,
		defense:   18,
	},
}
