package main

import (
	"math"
)

//Player is the main character
//question - do we generalize this to 'creature' or 'object'?
type Player struct {
	content        rune
	x              int
	y              int
	level          int
	strength       int
	defense        int
	maxHP          int
	hp             int
	exp            int
	nextLevelEXP   int
	healingActions int
}

const neededHealingActions int = 8

//NewPlayer creates a player struct
func NewPlayer(startX, startY int) *Player {
	level1 := playerLevels[1]
	level2 := playerLevels[2]
	return &Player{
		content:      '@',
		x:            startX,
		y:            startY,
		level:        1,
		maxHP:        level1.maxHP,
		hp:           level1.maxHP,
		strength:     level1.strength,
		defense:      level1.defense,
		exp:          0,
		nextLevelEXP: level2.neededEXP,
	}
}

//HealFromActions handles the automatic healing in the game that happens after doing enough actions
func (p *Player) HealFromActions() {
	if p.hp == p.maxHP {
		p.healingActions = 0
		return
	}

	p.healingActions++
	if p.healingActions == neededHealingActions {
		p.hp = min(p.hp+p.level, p.maxHP)
		p.healingActions = 0
	}
}

//ProcessLevelUp checks if the player can level up, and if so, updates the appropriate player stats
//Returns if there was a level up or not
func (p *Player) ProcessLevelUp() bool {
	if p.exp < p.nextLevelEXP {
		return false
	}

	newLevel := playerLevels[p.level+1]
	p.level = newLevel.level
	p.maxHP = newLevel.maxHP
	p.hp = newLevel.maxHP
	p.strength = newLevel.strength
	p.defense = newLevel.defense

	afterNewLevel, ok := playerLevels[p.level+1]
	if !ok {
		//We are at max level
		p.nextLevelEXP = math.MaxInt32
		return true
	}

	p.nextLevelEXP = afterNewLevel.neededEXP

	//In case more than one level up happens from one exp increase
	p.ProcessLevelUp()
	return true
}
