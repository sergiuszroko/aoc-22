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
	elvs := []cord{}
	lines := strings.Split(string(txt), "\n")
	y := len(lines)
	for _, line := range lines {
		y--
		for x, c := range line {
			if c == '#' {
				elvs = append(elvs, cord{x: x, y: y})
				continue
			} else if c == '.' {
				continue
			}
			fmt.Println(c)
			fmt.Println(string(c))
			panic("dk")
		}
	}
	g := newGame(elvs)
	g.draw()
	i := 1
	for {
		move := g.doOneRound()
		if !move {
			break
		}
		i++
	}
	g.draw()
	fmt.Println(i)
}

func newGame(elves []cord) *game {
	elvesMap := make(map[cord]struct{}, len(elves))
	for _, el := range elves {
		elvesMap[el] = struct{}{}
	}

	return &game{
		elves: elvesMap,
		rules: getRules(),
	}
}

type cord struct {
	x int
	y int
}

func (c cord) getAdjacent() []cord {
	return []cord{{c.x + 1, c.y + 1}, {c.x + 1, c.y}, {c.x + 1, c.y - 1}, {c.x, c.y + 1}, {c.x, c.y - 1}, {c.x - 1, c.y + 1}, {c.x - 1, c.y}, {c.x - 1, c.y - 1}}
}

type game struct {
	rules     []func(cord, map[cord]struct{}) (cord, bool)
	elves     map[cord]struct{}
	startRule int
}

func (g *game) doOneRound() bool {
	move := false
	new := make(map[cord]struct{})
	proposedCount := make(map[cord]int)
	moves := make(map[cord]cord)
	for el := range g.elves {
		if ok := isAdjacentFree(el, g.elves); ok {
			new[el] = struct{}{}
			continue // doesnt move
		}
		newCord, ok := g.getNewCord(el, g.elves)
		if !ok {
			new[el] = struct{}{}
			continue // doesnt move
		}
		moves[el] = newCord
		proposedCount[newCord] = proposedCount[newCord] + 1
	}
	for oldCord, newCord := range moves {
		if proposedCount[newCord] == 1 {
			new[newCord] = struct{}{}
			move = true
			continue
		}
		new[oldCord] = struct{}{}
	}
	g.elves = new
	g.startRule++
	return move
}
func (g *game) draw() {
	var minX, maxX, minY, maxY int
	for el := range g.elves {
		minX = el.x
		maxX = el.x
		minY = el.y
		maxY = el.y
		break
	}
	for el := range g.elves {
		if minX > el.x {
			minX = el.x
		}
		if maxX < el.x {
			maxX = el.x
		}
		if minY > el.y {
			minY = el.y
		}
		if maxY < el.y {
			maxY = el.y
		}
	}
	fmt.Println(cord{maxX, maxY})
	b := strings.Builder{}
	for j := maxY; j >= minY; j-- {
		for i := minX; i <= maxX; i++ {
			if _, ok := g.elves[cord{x: i, y: j}]; ok {
				b.WriteRune('#')
				continue
			}
			b.WriteRune('.')
		}
		b.WriteRune('\n')
	}
	fmt.Println(b.String())
}

func (g *game) countEmpty() int {
	var minX, maxX, minY, maxY int
	for el := range g.elves {
		minX = el.x
		maxX = el.x
		minY = el.y
		maxY = el.y
		break
	}
	for el := range g.elves {
		if minX > el.x {
			minX = el.x
		}
		if maxX < el.x {
			maxX = el.x
		}
		if minY > el.y {
			minY = el.y
		}
		if maxY < el.y {
			maxY = el.y
		}
	}
	count := 0
	for i := minX; i <= maxX; i++ {
		for j := minY; j <= maxY; j++ {
			if _, ok := g.elves[cord{x: i, y: j}]; ok {
				continue
			}
			count++
		}
	}
	return count
}

func (g *game) getNewCord(el cord, elves map[cord]struct{}) (cord, bool) {
	i := g.startRule
	for {
		if i-g.startRule >= 4 {
			break
		}
		rule := g.rules[i%4]
		newCord, ok := rule(el, elves)
		if ok {
			return newCord, true
		}
		i++
	}
	return el, false
}

func isAdjacentFree(el cord, elves map[cord]struct{}) bool {
	for _, n := range el.getAdjacent() {
		if _, ok := elves[n]; ok {
			return false
		}
	}
	return true
}

func getRules() []func(cord, map[cord]struct{}) (cord, bool) {
	return []func(cord, map[cord]struct{}) (cord, bool){
		func(c cord, m map[cord]struct{}) (cord, bool) {
			_, nw := m[cord{c.x + 1, c.y + 1}]
			_, n := m[cord{c.x, c.y + 1}]
			_, ne := m[cord{c.x - 1, c.y + 1}]
			if nw || n || ne {
				return cord{}, false
			}
			return cord{c.x, c.y + 1}, true
		},
		func(c cord, m map[cord]struct{}) (cord, bool) {
			_, sw := m[cord{c.x + 1, c.y - 1}]
			_, s := m[cord{c.x, c.y - 1}]
			_, se := m[cord{c.x - 1, c.y - 1}]
			if sw || s || se {
				return cord{}, false
			}
			return cord{c.x, c.y - 1}, true
		},
		func(c cord, m map[cord]struct{}) (cord, bool) {
			_, nw := m[cord{c.x - 1, c.y + 1}]
			_, w := m[cord{c.x - 1, c.y}]
			_, sw := m[cord{c.x - 1, c.y - 1}]
			if nw || w || sw {
				return cord{}, false
			}
			return cord{c.x - 1, c.y}, true
		},
		func(c cord, m map[cord]struct{}) (cord, bool) {
			_, ne := m[cord{c.x + 1, c.y + 1}]
			_, e := m[cord{c.x + 1, c.y}]
			_, se := m[cord{c.x + 1, c.y - 1}]
			if ne || e || se {
				return cord{}, false
			}
			return cord{c.x + 1, c.y}, true
		},
	}
}
