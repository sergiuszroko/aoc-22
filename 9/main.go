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
	up direction = iota
	down
	left
	right
)

type direction int

type cords struct {
	x, y int
}

type instruction struct {
	direction direction
	value     int
}

func moveHead(direction direction, head cords) cords {
	switch direction {
	case up:
		head.y++
	case down:
		head.y--
	case left:
		head.x--
	case right:
		head.x++
	}
	return head
}

func updatePosition(head, tail cords) cords {
	for {
		oldTail := tail
		tail = changePosition(head, tail)
		if tail == oldTail {
			break
		}
	}
	return tail
}
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func changePosition(head, tail cords) cords {
	diffX := head.x - tail.x
	diffY := head.y - tail.y
	switch {
	case diffX == 0 && diffY > 1:
		tail.y++
	case diffX == 0 && diffY < -1:
		tail.y--
	case diffX > 1 && diffY == 0:
		tail.x++
	case diffX < -1 && diffY == 0:
		tail.x--
	case abs(diffX) == 1 && abs(diffY) == 1:
	case diffX > 0 && diffY > 0:
		tail.x++
		tail.y++
	case diffX > 0 && diffY < -0:
		tail.x++
		tail.y--
	case diffX < -0 && diffY > 0:
		tail.x--
		tail.y++
	case diffX < -0 && diffY < -0:
		tail.x--
		tail.y--
	}

	return tail
}

func countVisited(poss []cords) int {
	s := make(map[cords]struct{})
	for _, el := range poss {
		s[el] = struct{}{}
	}
	return len(s)
}
func readDir(d string) direction {
	switch d {
	case "U":
		return up
	case "D":
		return down
	case "L":
		return left
	case "R":
		return right
	}
	panic("cant read direction")
}

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
	tailPoss := make([]cords, 0)
	currentKnots := [10]cords{}
	for scanner.Scan() {
		txt := scanner.Text()
		parts := strings.Split(txt, " ")
		dir := readDir(parts[0])
		value, err := strconv.Atoi(parts[1])
		if err != nil {
			log.Fatal(err)
		}
		for i := 0; i < value; i++ {
			head := currentKnots[0]
			head = moveHead(dir, head)
			currentKnots[0] = head
			for j := 0; j < 9; j++ {
				head := currentKnots[j]
				tail := currentKnots[j+1]
				currentTail := updatePosition(head, tail)
				currentKnots[j+1] = currentTail
			}
			tailPoss = append(tailPoss, currentKnots[9])
		}
	}
	fmt.Println(countVisited(tailPoss))
}
