package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

const (
	rock choice = iota
	paper
	scissors

	win outcome = iota - 4
	draw
	lose
)

type choice int
type outcome int

func main() {
	f, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	var sum int
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		text := scanner.Text()
		they := text[0:1]
		us := text[2:3]
		s := score(they, us)
		sum += s
	}
	fmt.Printf("final score %d\n", sum)
}

func score(t, u string) int {
	they := mapChoice(t)
	out := mapOutcome(u)
	us := getChoice(they, out)
	p1 := int(us) + 1
	p2 := (int(out) - 1) * (-3)
	return p1 + p2
}

func getChoice(they choice, out outcome) choice {
	return choice((int(they) - int(out) + 3) % 3)
}

func mapOutcome(s string) outcome {
	switch s {
	case "X":
		return lose
	case "Y":
		return draw
	case "Z":
		return win
	}
	panic("no input")
}
func mapChoice(s string) choice {
	switch s {
	case "A":
		return rock
	case "B":
		return paper
	case "C":
		return scissors
	}
	panic("no input")
}
