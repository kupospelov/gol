package sim

import (
	"sort"
	"testing"

	"github.com/kupospelov/gol/sim/event"
)

func TestRun(t *testing.T) {
	tests := []struct {
		name   string
		size   int
		input  []event.Cell
		output []event.Cell
	}{
		{
			"still life",
			3,
			[]event.Cell{
				event.Cell{0, 0, true},
				event.Cell{1, 0, true},
				event.Cell{0, 1, true},
				event.Cell{1, 1, true},
			},
			[]event.Cell{},
		},
		{
			"oscillator",
			5,
			[]event.Cell{
				event.Cell{2, 1, true},
				event.Cell{2, 2, true},
				event.Cell{2, 3, true},
			},
			[]event.Cell{
				event.Cell{2, 1, false},
				event.Cell{2, 3, false},
				event.Cell{1, 2, true},
				event.Cell{3, 2, true},
			},
		},
		{
			"oscillator on edge",
			5,
			[]event.Cell{
				event.Cell{0, 4, true},
				event.Cell{0, 0, true},
				event.Cell{0, 1, true},
			},
			[]event.Cell{
				event.Cell{0, 4, false},
				event.Cell{0, 1, false},
				event.Cell{1, 0, true},
				event.Cell{4, 0, true},
			},
		},
	}

	for _, test := range tests {
		g := New(test.size, test.size)
		g.Set(test.input)

		actual := g.Run()
		expected := test.output

		sort.Sort(event.Cells(actual))
		sort.Sort(event.Cells(expected))

		equal := len(actual) == len(expected)
		if equal {
			for i, e := range actual {
				if e != expected[i] {
					equal = false
					break
				}
			}
		}

		if !equal {
			t.Errorf("%s: Got the following events: %v, expected: %v", test.name, actual, expected)
		}
	}
}

func BenchmarkRun(b *testing.B) {
	g := New(100, 100)
	g.Set([]event.Cell{
		event.Cell{1, 0, true},
		event.Cell{1, 1, true},
		event.Cell{1, 2, true},
	})

	for i := 0; i < b.N; i++ {
		g.Run()
	}
}
