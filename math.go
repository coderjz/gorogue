package main

import "math/rand"

func random(min int, max int, rand *rand.Rand) int {
	return rand.Intn(max-min) + min
}

func min(a, b int) int {
	if a <= b {
		return a
	}
	return b
}
