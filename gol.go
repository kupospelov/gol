package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/gdamore/tcell"
	"github.com/kupospelov/gol/sim"
	"github.com/kupospelov/gol/sim/event"
)

const (
	EventRefresh = iota
	EventPause
	EventResize
	EventRestart
)

var names []rune

func init() {
	ranges := [...]struct{ from, to rune }{
		{'a', 'z'},
		{'A', 'Z'},
		{'0', '9'},
	}

	for _, r := range ranges {
		for c := r.from; c <= r.to; c++ {
			names = append(names, c)
		}
	}
}

func printEvents(s tcell.Screen, events []event.Cell) {
	active := tcell.StyleDefault.Foreground(tcell.ColorGreen)
	inactive := tcell.StyleDefault
	for _, e := range events {
		if e.Alive {
			bold := rand.Intn(2) > 0
			ch := names[rand.Intn(len(names))]

			s.SetContent(e.X, e.Y, ch, nil, active.Bold(bold))
		} else {
			s.SetContent(e.X, e.Y, ' ', nil, inactive)
		}
	}
	s.Show()
}

func eventHandler(s tcell.Screen, events chan int) {
	defer close(events)
	for {
		ev := s.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventKey:
			switch ev.Key() {
			case tcell.KeyEscape, tcell.KeyEnter:
				return
			case tcell.KeyRune:
				switch ev.Rune() {
				case 'a':
					events <- EventRestart
				case 'r':
					events <- EventRefresh
				case 'q':
					return
				case ' ':
					events <- EventPause
				}
			}
		case *tcell.EventResize:
			events <- EventResize
			s.Sync()
		}
	}
}

func run(s tcell.Screen, g *sim.Simulation, c *uint64) {
	events := g.Run()
	printEvents(s, events)
	*c++
}

func main() {
	init := flag.Int("init", 50, "The percentage of cells to populate on start")
	flag.Parse()

	if *init < 0 || *init > 100 {
		fmt.Println("init: must be between 0 and 100")
		os.Exit(1)
	}

	tcell.SetEncodingFallback(tcell.EncodingFallbackASCII)
	s, e := tcell.NewScreen()
	if e != nil {
		fmt.Fprintf(os.Stderr, "%v\n", e)
		os.Exit(1)
	}
	if e = s.Init(); e != nil {
		fmt.Fprintf(os.Stderr, "%v\n", e)
		os.Exit(1)
	}

	s.Clear()

	events := make(chan int)
	go eventHandler(s, events)

	w, h := s.Size()
	g := sim.New(w, h)
	g.Randomize(*init)
	printEvents(s, g.Dump())

	var paused bool
	var counter uint64
	ticks := time.Tick(time.Millisecond * 100)

loop:
	for {
		select {
		case event, ok := <-events:
			if !ok {
				break loop
			}
			switch event {
			case EventRefresh:
				run(s, &g, &counter)
			case EventPause:
				paused = !paused
			case EventResize:
				g.Resize(s.Size())
			case EventRestart:
				s.Clear()
				g.Reset()
				g.Randomize(*init)
			}
		case <-ticks:
			if !paused {
				run(s, &g, &counter)
			}
		}
	}

	s.Fini()
	fmt.Printf("Finished %d iterations.\n", counter)
}
