package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type cpu struct {
	register int
	cycle    int
}

type game struct {
	cpu     cpu
	printer *printer
}

type printer struct {
	screen strings.Builder
}

func (p *printer) print(register, cycle int) {
	currentPosition := (cycle - 1) % 40
	if register-1 <= currentPosition && currentPosition <= register+1 {
		p.screen.WriteString("#")
	} else {
		p.screen.WriteString(".")
	}
	if currentPosition == 39 {
		p.screen.WriteString("\n")
	}
}
func (p *printer) getScreen() string {
	return p.screen.String()
}

func newGame() *game {
	return &game{
		cpu: cpu{
			register: 1,
			cycle:    1,
		},
		printer: &printer{},
	}
}

func (g *game) doNoop() {
	g.hookObservation()
	g.cpu.cycle++
}

func (g *game) doAdd(v int) {
	g.hookObservation()
	g.cpu.cycle++
	g.hookObservation()
	g.cpu.cycle++
	g.cpu.register += v
}

func (g *game) hookObservation() {
	g.printer.print(g.cpu.register, g.cpu.cycle)
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
	g := newGame()
	for scanner.Scan() {
		txt := scanner.Text()
		if strings.HasPrefix(txt, "noop") {
			g.doNoop()
		}
		if strings.HasPrefix(txt, "addx") {
			parts := strings.Split(txt, " ")
			v, err := strconv.Atoi(parts[1])
			if err != nil {
				log.Fatal(err)
			}
			g.doAdd(v)
		}

	}
	fmt.Println((g.printer.getScreen()))
}
