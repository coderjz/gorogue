package game

import (
	"math"
	"math/rand"
)

//Returns a random number in the range [min, max)
func random(min int, max int, rand *rand.Rand) int {
	return rand.Intn(max-min) + min
}

func min(a, b int) int {
	if a <= b {
		return a
	}
	return b
}

func round(val float64) int {
	if val < 0 {
		return int(val - 0.5)
	}
	return int(val + 0.5)
}

func distance(x1, y1, x2, y2 int) float64 {
	return math.Sqrt(float64((x2-x1)*(x2-x1) + (y2-y1)*(y2-y1)))
}
