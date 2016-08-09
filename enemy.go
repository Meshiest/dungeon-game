package main

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/meshiest/go-dungeon/dungeon"
	"math"
)

type Enemy struct {
	X, Y, Size, DPS float64
	Health          int
}

func (e *Enemy) CollideWithEnemy(other *Enemy) {
	vector := mgl32.Vec2{float32(e.X - other.X), float32(e.Y - other.Y)}
	minDist := float32(e.Size/2.0 + other.Size/2.0)
	if vector.Len() < minDist {
		pushDist := float64((vector.Len() - minDist) / 2.0)
		angle := math.Atan2(float64(vector.Y()), float64(vector.X()))
		e.X += math.Cos(angle+math.Pi) * pushDist
		e.Y += math.Sin(angle+math.Pi) * pushDist
		other.X += math.Cos(angle) * pushDist
		other.Y += math.Sin(angle) * pushDist
	}
}

func (e *Enemy) CollideWithPlayer(player *Player) bool {
	vector := mgl32.Vec2{float32(e.X - player.X), float32(e.Y - player.Y)}
	minDist := float32(e.Size/2.0 + player.Size/2.0)
	if vector.Len() < minDist {
		pushDist := float64((vector.Len() - minDist) / 2.0)
		angle := math.Atan2(float64(vector.Y()), float64(vector.X()))
		e.X += math.Cos(angle+math.Pi) * pushDist
		e.Y += math.Sin(angle+math.Pi) * pushDist
		player.X += math.Cos(angle) * pushDist
		player.Y += math.Sin(angle) * pushDist
		return true
	}
	return false
}

func (e *Enemy) CollideWithDungeon(dungeon *dungeon.Dungeon) {
	dist := e.Size / 2.0

	x := int(math.Floor(e.X + 0.5))
	y := int(math.Floor(e.Y + 0.5))

	if x >= 0 && y >= 0 && x < len(dungeon.Grid) && y < len(dungeon.Grid) {
		if (y == 0 || dungeon.Grid[y-1][x] == 0) && e.Y-float64(y) < -dist {
			e.Y = float64(y) - dist
		}
		if (y == len(dungeon.Grid)-1 || dungeon.Grid[y+1][x] == 0) && e.Y-float64(y) > dist {
			e.Y = float64(y) + dist
		}
		if (x == 0 || dungeon.Grid[y][x-1] == 0) && e.X-float64(x) < -dist {
			e.X = float64(x) - dist
		}
		if (x == len(dungeon.Grid)-1 || dungeon.Grid[y][x+1] == 0) && e.X-float64(x) > dist {
			e.X = float64(x) + dist
		}
	}
}

type EnemyRender struct {
	Dist  float32
	Enemy Enemy
}

type EnemyRenderOrder []*EnemyRender

func (c EnemyRenderOrder) Len() int           { return len(c) }
func (c EnemyRenderOrder) Swap(i, j int)      { c[i], c[j] = c[j], c[i] }
func (c EnemyRenderOrder) Less(i, j int) bool { return c[i].Dist > c[j].Dist }
