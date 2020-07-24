package sim

import (
	"math/rand"

	"github.com/kupospelov/gol/sim/event"
)

type Simulation struct {
	w, h  int
	field [][]bool
}

func New(w, h int) Simulation {
	var r Simulation
	r.Resize(w, h)
	return r
}

func (g *Simulation) Reset() {
	for x := 0; x < g.w; x++ {
		for y := 0; y < g.h; y++ {
			g.field[x][y] = false
		}
	}
}

func (g *Simulation) Set(events []event.Cell) {
	for _, e := range events {
		g.field[e.X][e.Y] = e.Alive
	}
}

func (g *Simulation) Resize(w, h int) {
	field := make([][]bool, w)
	for x := 0; x < w; x++ {
		field[x] = make([]bool, h)
		if x < g.w {
			copy(field[x], g.field[x])
		}
	}

	g.w, g.h = w, h
	g.field = field
}

func (g *Simulation) Randomize(p int) {
	n := g.h * g.w * p / 100
	for i := 0; i < n; i++ {
		x := rand.Intn(g.w)
		y := rand.Intn(g.h)
		g.field[x][y] = true
	}

	// Remove cells that are going to die right away
	g.Run()
}

func (g *Simulation) Dump() []event.Cell {
	r := make([]event.Cell, 0)
	for x := 0; x < g.w; x++ {
		for y := 0; y < g.h; y++ {
			r = append(r, event.Cell{x, y, g.field[x][y]})
		}
	}
	return r
}

func (g *Simulation) Run() []event.Cell {
	r := make([]event.Cell, 0)
	for x := 0; x < g.w; x++ {
		for y := 0; y < g.h; y++ {
			n := g.neighbours(x, y)
			if g.field[x][y] {
				if n < 2 || n > 3 {
					r = append(r, event.Cell{x, y, false})
				}
			} else {
				if n == 3 {
					r = append(r, event.Cell{x, y, true})
				}
			}
		}
	}

	for _, e := range r {
		g.field[e.X][e.Y] = e.Alive
	}

	return r
}

func (g *Simulation) neighbours(x, y int) int {
	return g.cell(x-1, y-1) + g.cell(x-1, y) + g.cell(x-1, y+1) + g.cell(x, y-1) + g.cell(x, y+1) + g.cell(x+1, y-1) + g.cell(x+1, y) + g.cell(x+1, y+1)
}

func (g *Simulation) cell(x, y int) int {
	if g.field[(g.w+x)%g.w][(g.h+y)%g.h] {
		return 1
	}
	return 0
}
