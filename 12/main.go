package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	f, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	if err != nil {
		log.Fatal(err)
	}
	g := newGame()
	i := 0
	for scanner.Scan() {
		line := scanner.Text()
		g.readLine(line, i)
		i++
	}
	g.maxY = i
	// g.printMap()
	g.g()
	fmt.Println(g.findShortestA())
}

type coord struct {
	x, y int
}

type game struct {
	start, end coord
	elevation  map[coord]int
	fewestStep map[coord]int
	maxX, maxY int
}

func newGame() *game {
	g := &game{}
	g.elevation = make(map[coord]int)
	g.fewestStep = make(map[coord]int)
	return g
}

func (g *game) printMap() {
	b := strings.Builder{}
	for i := 0; i < g.maxY; i++ {
		for j := 0; j < g.maxX; j++ {
			b.WriteString(fmt.Sprintf(" %02d ", g.elevation[coord{j, i}]))
		}
		b.WriteString("\n")
	}
	fmt.Println(b.String())
}

func (g *game) readLine(line string, n int) {
	g.maxX = len(line)
	for i, c := range line {
		if c == 'S' {
			g.start = coord{i, n}
			g.elevation[coord{i, n}] = 0
			continue
		}
		if c == 'E' {
			g.end = coord{i, n}
			g.elevation[coord{i, n}] = 27
			continue
		}
		g.elevation[coord{i, n}] = int(c - 96)
	}
}

func (g *game) g() {
	g.findPaths(g.end)
}

func (g *game) findPaths(c coord) {
	next := g.getAdjacent(c)
	currentElevation := g.elevation[c]
	currentSteps := g.fewestStep[c]
	for _, n := range next {
		nextElevation := g.elevation[n]
		if nextElevation < currentElevation-1 {
			continue
		}
		nextSteps, ok := g.fewestStep[n]
		if ok && currentSteps+1 >= nextSteps {
			continue
		}
		g.fewestStep[n] = currentSteps + 1
		g.findPaths(n)
	}
}

func (g *game) getAdjacent(c coord) []coord {
	possible := []coord{{-1, 0}, {0, 1}, {1, 0}, {0, -1}}
	var result []coord
	for _, p := range possible {
		r := coord{c.x + p.x, c.y + p.y}
		if r.x < 0 || r.x > g.maxX || r.y < 0 || r.y > g.maxY {
			continue
		}
		result = append(result, r)
	}
	return result
}
func (g *game) findShortestA() int {
	min := g.fewestStep[g.start]
	fmt.Println(min)
	for k, v := range g.elevation {
		if v != 1 {
			continue
		}
		current := g.fewestStep[k]
		if current == 0 {
			continue
		}
		if min > current {
			min = current
		}
	}
	return min
}
