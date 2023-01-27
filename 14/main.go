package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

const (
	X          = 1000
	Y          = 500
	sandStartX = 500
	sandStartY = 0
)

func main() {
	f, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	g := newGame()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		readLine(line, g)
	}
	s := 0
	for {
		stop := g.addSand()
		if stop {
			break
		}
		s++
	}
	g.draw()
	fmt.Println(s)

}

func readLine(line string, g *game) {
	segmentsString := strings.Split(line, " -> ")
	segmentsCoords := make([][2]int, len(segmentsString))
	for i, s := range segmentsString {
		xy := strings.Split(s, ",")
		x, err := strconv.Atoi(xy[0])
		if err != nil {
			panic(err)
		}
		y, err := strconv.Atoi(xy[1])
		if err != nil {
			panic(err)
		}
		segmentsCoords[i] = [2]int{x, y}
	}
	nPath := len(segmentsCoords) - 1
	for i := 0; i < nPath; i++ {
		start := segmentsCoords[i]
		end := segmentsCoords[i+1]
		g.addRockPath(start[0], start[1], end[0], end[1])
	}

}

type game struct {
	cave       [X][Y]bool
	maxY       int
	minX, maxX int
}

func newGame() *game {
	return &game{
		minX: X,
	}
}
func (g *game) draw() {
	for j := 0; j <= g.maxY; j++ {
		for i := g.minX; i <= g.maxX; i++ {
			t := g.cave[i][j]
			if t {
				fmt.Print(" # ")
			} else {
				fmt.Print(" . ")
			}
		}
		fmt.Println()
	}
}

func (g *game) addSand() bool {
	sandX := sandStartX
	sandY := sandStartY
	for {
		newSandY := sandY + 1
		isTouchingFloor := newSandY == g.maxY+2
		occupiedMiddle := g.cave[sandX][newSandY] || isTouchingFloor
		if !occupiedMiddle {
			sandY = newSandY
			continue
		}
		newSandXLeft := sandX - 1
		occupiedLeft := g.cave[newSandXLeft][newSandY] || isTouchingFloor
		if !occupiedLeft {
			sandY = newSandY
			sandX = newSandXLeft
			continue
		}
		newSandXRight := sandX + 1
		occupieRight := g.cave[newSandXRight][newSandY] || isTouchingFloor
		if !occupieRight {
			sandY = newSandY
			sandX = newSandXRight
			continue
		}
		g.cave[sandX][sandY] = true
		break
	}
	return sandY == 0
}

func (g *game) addRockPath(x1, y1, x2, y2 int) {
	if g.maxY < y1 {
		g.maxY = y1
	}
	if g.maxY < y2 {
		g.maxY = y2
	}
	if g.minX > x1 {
		g.minX = x1
	}
	if g.minX > x2 {
		g.minX = x2
	}
	if g.maxX < x1 {
		g.maxX = x1
	}
	if g.maxX < x2 {
		g.maxX = x2
	}
	horizontal := y1 == y2
	if horizontal {
		start := x1
		end := x2
		if x2 < x1 {
			start = x2
			end = x1
		}
		for i := start; i <= end; i++ {
			g.cave[i][y1] = true
		}
		return
	}
	if x1 != x2 {
		fmt.Println(x1, y1, x2, y2)
		panic("line is diagonal")
	}
	start := y1
	end := y2
	if y2 < y1 {
		start = y2
		end = y1
	}
	for i := start; i <= end; i++ {
		g.cave[x1][i] = true
	}
}
