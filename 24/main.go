package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	txt, err := os.ReadFile("input.txt")
	if err != nil {
		panic(err)
	}
	g := newGame(string(txt))
	fmt.Println("game loaded")
	r := g.simulate()
	fmt.Printf("got result %d\n", r)
}

func newGame(txt string) *game {
	start := cord{x: 1, y: 0}
	blizzards := []blizzard{}

	lines := strings.Split(txt, "\n")
	for y, line := range lines {
		for x, c := range line {
			switch c {
			case '#':
			case '.':
			case '>':
				blizzards = append(blizzards, blizzard{c: cord{x: x, y: y}, d: right})
			case '<':
				blizzards = append(blizzards, blizzard{c: cord{x: x, y: y}, d: left})
			case '^':
				blizzards = append(blizzards, blizzard{c: cord{x: x, y: y}, d: up})
			case 'v':
				blizzards = append(blizzards, blizzard{c: cord{x: x, y: y}, d: down})
			}
		}
	}
	endY := len(lines) - 1
	endX := len(lines[0]) - 1

	return &game{
		blizzards:   blizzards,
		expeditions: []expedition{{c: start, state: first}},
		start:       start,
		endCorner:   cord{x: endX, y: endY},
		end:         cord{x: endX - 1, y: endY},
	}
}

const (
	up direction = iota
	right
	down
	left
)

type direction int
type cord struct {
	x int
	y int
}

func (c cord) move(d direction) cord {
	switch d {
	case up:
		return cord{c.x, c.y - 1}
	case down:
		return cord{c.x, c.y + 1}
	case left:
		return cord{c.x - 1, c.y}
	case right:
		return cord{c.x + 1, c.y}
	}
	panic(1)
}

type blizzard struct {
	c cord
	d direction
}
type state int

const (
	first state = iota
	second
	third
)

type expedition struct {
	c     cord
	state state
}
type game struct {
	blizzards       []blizzard
	blizzardsLookup map[cord]struct{}
	expeditions     []expedition
	endCorner       cord
	end             cord
	start           cord
}

func (g *game) simulate() int {
	time := 0
	for {
		time++
		fmt.Printf("time %d expeditions %d\n", time, len(g.expeditions))
		g.moveBlizzards()
		newExpeditinsSet := make(map[expedition]struct{})
		for _, exp := range g.expeditions {

			if g.canMove(exp.c) {
				newExpeditinsSet[exp] = struct{}{}
			}

			cDown := exp.c.move(down)
			if g.canMove(cDown) {
				if cDown == g.end && exp.state == third {
					return time
				}
				if cDown == g.end && exp.state == first {
					newExpeditinsSet[expedition{c: cDown, state: second}] = struct{}{}
				} else {
					newExpeditinsSet[expedition{c: cDown, state: exp.state}] = struct{}{}
				}
			}

			cUp := exp.c.move(up)
			if g.canMove(cUp) {
				if cUp == g.start && exp.state == second {
					newExpeditinsSet[expedition{c: cUp, state: third}] = struct{}{}
				} else {
					newExpeditinsSet[expedition{c: cUp, state: exp.state}] = struct{}{}
				}
			}

			cLeft := exp.c.move(left)
			if g.canMove(cLeft) {
				newExpeditinsSet[expedition{c: cLeft, state: exp.state}] = struct{}{}
			}

			cRight := exp.c.move(right)
			if g.canMove(cRight) {
				newExpeditinsSet[expedition{c: cRight, state: exp.state}] = struct{}{}
			}
		}
		newExpeditions := make([]expedition, len(newExpeditinsSet))
		i := 0
		for e := range newExpeditinsSet {
			newExpeditions[i] = e
			i++
		}
		g.expeditions = newExpeditions
	}
}

func (g *game) moveBlizzards() {
	nextBlizzards := make([]blizzard, len(g.blizzards))
	for i, blizz := range g.blizzards {
		newC := blizz.c.move(blizz.d)
		if newC.y >= g.endCorner.y {
			newC.y = 1
		}
		if newC.y <= 0 {
			newC.y = g.endCorner.y - 1
		}
		if newC.x >= g.endCorner.x {
			newC.x = 1
		}
		if newC.x <= 0 {
			newC.x = g.endCorner.x - 1
		}
		nextBlizzards[i] = blizzard{c: newC, d: blizz.d}
	}
	nextLookup := make(map[cord]struct{})
	for _, blizz := range nextBlizzards {
		nextLookup[blizz.c] = struct{}{}
	}
	g.blizzards = nextBlizzards
	g.blizzardsLookup = nextLookup
}

func (g *game) canMove(c cord) bool {
	_, ok := g.blizzardsLookup[c]
	if ok {
		return false
	}
	if c == g.start || c == g.end {
		return true
	}
	if c.x <= 0 || c.y <= 0 || c.x >= g.endCorner.x || c.y >= g.endCorner.y {
		return false
	}
	return true
}
