package main

import (
	"github.com/meshiest/go-dungeon/dungeon"
	"math"
)

type Player struct {
	X, Y, Yaw, Pitch, Speed, Size, Health float64
}

func (p *Player) CollideWithDungeon(dungeon *dungeon.Dungeon) {
	dist := p.Size / 2.0

	x := int(math.Floor(p.X + 0.5))
	y := int(math.Floor(p.Y + 0.5))

	if x >= 0 && y >= 0 && x < len(dungeon.Grid) && y < len(dungeon.Grid) {
		if (y == 0 || dungeon.Grid[y-1][x] == 0) && p.Y-float64(y) < -dist {
			p.Y = float64(y) - dist
		}
		if (y == len(dungeon.Grid)-1 || dungeon.Grid[y+1][x] == 0) && p.Y-float64(y) > dist {
			p.Y = float64(y) + dist
		}
		if (x == 0 || dungeon.Grid[y][x-1] == 0) && p.X-float64(x) < -dist {
			p.X = float64(x) - dist
		}
		if (x == len(dungeon.Grid)-1 || dungeon.Grid[y][x+1] == 0) && p.X-float64(x) > dist {
			p.X = float64(x) + dist
		}
	}
}
