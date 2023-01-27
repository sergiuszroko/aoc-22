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
	scanner.Scan()
	line := scanner.Text()
	wind := readLine(line)
	game := newGame(wind)
	// deltas := []int{}
	// last := 0
	for i := 0; i < 997406; i++ {
		game.fallNextRock()
		// now := len(game.tower) - 1
		// deltas = append(deltas, now-last)
		// last = now
	}
	// fmt.Println(deltas)
	// start := deltas[1000:1100]
	// fmt.Println(start)
	// for i := 0; i < len(deltas); i++ {
	// 	sample := deltas[i : i+100]
	// 	if reflect.DeepEqual(start, sample) {
	// 		fmt.Println(i)
	// 	}
	// }
	fmt.Println(len(game.tower) - 1)
}
func readLine(line string) wind {
	wind := wind{}
	for _, c := range line {
		if c == '<' {
			wind.shape = append(wind.shape, left)
		} else {
			wind.shape = append(wind.shape, right)
		}
	}
	return wind
}

const (
	FLOOR       = 511
	WALLS       = 257
	TOWERWIDDTH = 7

	left direction = iota
	right
)

type towerLevel int
type game struct {
	tower []towerLevel
	wind  wind
	nRock int
}

func (g *game) printTower(n int) {
	for i := len(g.tower) - 1; i >= len(g.tower)-n; i-- {
		l := g.tower[i]
		b := strings.Builder{}
		for j := 8; j >= 0; j-- {
			mask := 1 << j
			r := mask & int(l)
			if r > 0 {
				b.WriteRune('#')
			} else {
				b.WriteRune('.')
			}
		}
		fmt.Println(b.String())
	}
}

func newGame(wind wind) *game {
	return &game{
		tower: []towerLevel{FLOOR},
		wind:  wind,
	}
}

func (g *game) getNextRock() rock {
	g.nRock++
	switch (g.nRock - 1) % 5 {
	case 0:
		return []int{60} //000111100
	case 1:
		return []int{16, 56, 16} //000010000 000111000 000010000
	case 2:
		return []int{56, 8, 8} //000111000 000001000 000001000
	case 3:
		return []int{32, 32, 32, 32} //000100000 000100000 000100000 000100000
	case 4:
		return []int{48, 48} //000110000 000110000
	}
	panic("no other way")
}
func (g *game) getWindDirection() direction {
	return g.wind.get()
}

func (g *game) fallNextRock() {
	fallingRock := g.getNextRock()
	rockLevel := g.getFallingRockStartingLevel()
	for {
		windDirection := g.getWindDirection()
		switch windDirection {
		case left:
			fallinRockNextPosition := fallingRock.moveLeft()
			if g.checkCollision(fallinRockNextPosition, rockLevel) {
				fallingRock = fallinRockNextPosition
			} else {
			}
		case right:
			fallinRockNextPosition := fallingRock.moveRight()
			if g.checkCollision(fallinRockNextPosition, rockLevel) {
				fallingRock = fallinRockNextPosition
			} else {
			}
		}
		if g.checkCollision(fallingRock, rockLevel-1) {
			rockLevel -= 1
			continue
		}
		g.addRockToTower(fallingRock, rockLevel)
		return
	}
}

func (g *game) getFallingRockStartingLevel() int {
	return len(g.tower) + 3
}

func (g *game) checkCollision(rock rock, levelNo int) bool {
	for i, r := range rock {
		level := towerLevel(WALLS)
		if levelNo+i <= len(g.tower)-1 {
			level = g.tower[levelNo+i]
		}
		if int(r)&int(level) != 0 {
			return false
		}
	}
	return true
}

func (g *game) addRockToTower(rock rock, levelNo int) {
	for i, r := range rock {
		level := towerLevel(WALLS)
		if levelNo+i <= len(g.tower)-1 {
			level = g.tower[levelNo+i]
		} else {
			g.tower = append(g.tower, level)
		}
		g.tower[levelNo+i] = towerLevel(int(level) | int(r))
	}
}

type direction int
type wind struct {
	shape []direction
	n     int
}

func (w *wind) get() direction {
	toReturn := w.shape[w.n]
	w.n++
	if w.n >= len(w.shape) {
		w.n = 0
	}
	return toReturn
}

type rock []int

func (r rock) moveLeft() rock {
	nextRock := make(rock, len(r))
	copy(nextRock, r)
	for i := range nextRock {
		nextRock[i] = nextRock[i] << 1
	}
	return nextRock
}

func (r rock) moveRight() rock {
	nextRock := make(rock, len(r))
	copy(nextRock, r)
	for i := range nextRock {
		nextRock[i] = nextRock[i] >> 1
	}
	return nextRock
}
