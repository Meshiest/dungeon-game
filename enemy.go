package main

import (
    "github.com/meshiest/go-dungeon/dungeon"
    "math"
)

type Enemy struct {
  X, Y, Size float64
  Health int
}

func (e *Enemy) CollideWithDungeon(dungeon *dungeon.Dungeon) {
  dist := e.Size/2.0

  x := int(math.Floor(e.X+0.5))
  y := int(math.Floor(e.Y+0.5))

  if x >= 0 && y >= 0 && x < len(dungeon.Grid) && y < len(dungeon.Grid) {
    if (y == 0 || dungeon.Grid[y-1][x] == 0) && e.Y - float64(y) < -dist {
      e.Y = float64(y)-dist
    }
    if (y == len(dungeon.Grid) - 1 || dungeon.Grid[y+1][x] == 0) && e.Y - float64(y) > dist {
      e.Y = float64(y)+dist
    }
    if (x == 0 || dungeon.Grid[y][x-1] == 0) && e.X - float64(x) < -dist {
      e.X = float64(x)-dist
    }
    if (x == len(dungeon.Grid) - 1 || dungeon.Grid[y][x+1] == 0) && e.X - float64(x) > dist {
      e.X = float64(x)+dist
    }
  }
}