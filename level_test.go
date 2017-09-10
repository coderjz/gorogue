package main

import (
	"math/rand"
	"testing"
	"time"
)

func TestGetPointInRoom(t *testing.T) {
	levelRand = rand.New(rand.NewSource(time.Now().UnixNano()))
	rooms := []Room{
		{
			x1: 10,
			x2: 14,
			y1: 10,
			y2: 14,
		},
		{
			x1: 10,
			x2: 20,
			y1: 30,
			y2: 34,
		},
		{
			x1: 0,
			x2: 4,
			y1: 10,
			y2: 20,
		},
	}

	for _, r := range rooms {
		for i := 0; i < 1000; i++ {
			x, y := r.getPointInRoom()

			if x <= r.x1 || x >= r.x2 || y <= r.y1 || y >= r.y2 {
				t.Fatalf("Invalid point (%d, %d) generated that is not within room walls. Room: %+v", x, y, r)
			}
		}
	}
}
