package main

//Player is the main character
//question - do we generalize this to 'creature' or 'object'?
type Player struct {
	content        rune
	x              int
	y              int
	level          int
	maxHP          int
	hp             int
	healingActions int
}

const neededHealingActions int = 8

//NewPlayer creates a player struct
func NewPlayer(startX, startY int) *Player {
	return &Player{
		content: '@',
		x:       startX,
		y:       startY,
		level:   1,
		maxHP:   20,
		hp:      20,
	}
}

func (p *Player) HealFromActions() {
	p.healingActions++
	if p.healingActions == neededHealingActions {
		p.hp = min(p.hp+p.level, p.maxHP)
		p.healingActions = 0
	}
}
