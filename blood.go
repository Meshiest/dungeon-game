package main

import (
	"math"
	"math/rand"
)

type Blood struct {
	X, Y, Angle float64
}

func NewBlood(x, y float64) *Blood {
	return &Blood{
		X:     x,
		Y:     y,
		Angle: rand.Float64() * math.Pi * 2.0,
	}
}
