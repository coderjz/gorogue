package game

import (
	"math/rand"
	"time"
)

var damageRand = rand.New(rand.NewSource(time.Now().UnixNano()))

//Calcuates how much damage will be done by given strength to a target with given defense
func calculateDamage(strength, defense int) int {
	//If strength not up to defense, 20% chance of doing 1 point of damage
	if strength <= defense {
		num := random(0, 10, damageRand)
		if num < 2 {
			return 1
		}
		return 0
	}

	//Base damage is just straight difference
	damage := float64(strength - defense)

	//We roll two dice together and add them to get a more "centarlized" random number (80% of numbers in [4, 14], 10% in [0, 3], 10% in [15, 18])
	diceRoll := float64(random(0, 10, damageRand) + random(0, 10, damageRand))
	normalizedDiceRoll := diceRoll / 18.0

	//Give us a result between [0.5 * damage, 1.5 * damage]
	return round((0.5 * damage) + (0.5 * normalizedDiceRoll * damage))
}
