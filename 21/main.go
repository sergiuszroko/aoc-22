package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	f, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	g := game{
		monkeysWaiting:    make(map[string]*monkey),
		monkeyYelling:     map[string]int{},
		solved:            map[string]int{},
		waiting:           map[string]*monkey{},
		copyMokeysWaiting: map[string]*monkey{},
	}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, ": ")
		name := parts[0]
		partsParts := strings.Split(parts[1], " ")
		if len(partsParts) == 1 {
			n, err := strconv.Atoi(partsParts[0])
			if err != nil {
				panic(err)
			}
			g.addYelling(name, n)
			continue
		}
		left := partsParts[0]
		var op operation
		switch partsParts[1] {
		case "+":
			op = add
		case "-":
			op = substraction
		case "*":
			op = multiplication
		case "/":
			op = division
		}
		right := partsParts[2]
		g.addWaiting(left, right, op, name)
	}
	g.solve()
	fmt.Println(g.solved["root"])
	fmt.Println(g.solved["bjft"])
	fmt.Println(g.solved["lntp"])
	g.markHumn("humn")
	fmt.Println(g.solve2())
}

const (
	add operation = iota
	substraction
	multiplication
	division
)

type monkey struct {
	left      number
	operation operation
	right     number
	name      string
	humn      bool
}

func (m *monkey) add(s string, n int) bool {
	if m.left.name == s {
		if m.left.known {
			return false
		}
		m.left.known = true
		m.left.n = n
		return true
	}
	if m.right.name == s {
		if m.right.known {
			return false
		}
		m.right.known = true
		m.right.n = n
		return true
	}
	return false
}
func (m *monkey) solve() (int, bool) {
	if !m.left.known || !m.right.known {
		return 0, false
	}
	switch m.operation {
	case add:
		return m.left.n + m.right.n, true
	case substraction:
		return m.left.n - m.right.n, true
	case multiplication:
		return m.left.n * m.right.n, true
	case division:
		return m.left.n / m.right.n, true
	}
	panic(1)
}

type number struct {
	known bool
	name  string
	n     int
}
type operation int

type game struct {
	monkeysWaiting    map[string]*monkey
	monkeyYelling     map[string]int
	solved            map[string]int
	waiting           map[string]*monkey
	copyMokeysWaiting map[string]*monkey
}

func (g *game) addYelling(s string, n int) {
	g.monkeyYelling[s] = n
}
func (g *game) addWaiting(left, right string, op operation, name string) {
	m := &monkey{
		left: number{
			known: false,
			name:  left,
		},
		operation: op,
		right: number{
			known: false,
			name:  right,
		},
		name: name,
	}
	g.monkeysWaiting[left] = m
	g.monkeysWaiting[right] = m
	g.copyMokeysWaiting[left] = m
	g.copyMokeysWaiting[right] = m
	g.waiting[name] = m
}

func (g *game) solveMonkey(s string, n int) {
	if s == "root" {
		return
	}
	m := g.monkeysWaiting[s]
	ok := m.add(s, n)
	if !ok {
		panic(ok)
	}
	nSolved, ok := m.solve()
	if ok {
		g.solved[m.name] = nSolved
		g.solveMonkey(m.name, nSolved)
	}
	delete(g.monkeysWaiting, s)
}

func (g *game) solve() {
	for k, v := range g.monkeyYelling {
		g.solveMonkey(k, v)
	}
}

func (g *game) markHumn(s string) {
	m := g.copyMokeysWaiting[s]
	m.humn = true
	g.copyMokeysWaiting[s] = m
	if m.name == "root" {
		return
	}
	g.markHumn(m.name)
}

func (g *game) solve2() int {
	r := g.waiting["root"]
	var humn *monkey
	var other *monkey
	left := g.waiting[r.left.name]
	right := g.waiting[r.right.name]
	if left.humn {
		humn = left
		other = right
	} else {
		humn = right
		other = left
	}
	otherN := g.solved[other.name]
	return g.returnHuman(humn, otherN)
}

func (g *game) returnHuman(m *monkey, n int) int {
	if m.name == "humn" {
		return n
	}
	left, leftOk := g.waiting[m.left.name]
	if !leftOk {
		left = &monkey{
			humn: m.left.name == "humn",
			name: m.left.name,
		}
	}
	right, rightOk := g.waiting[m.right.name]
	if !rightOk {
		right = &monkey{
			humn: m.right.name == "humn",
			name: m.right.name,
		}
	}
	var humn, other *monkey
	if left.humn {
		humn = left
		other = right
	} else {
		humn = right
		other = left
	}
	otherN, ok := g.solved[other.name]
	if !ok {
		otherN = g.monkeyYelling[other.name]
	}
	switch m.operation {
	case add:
		return g.returnHuman(humn, n-otherN)
	case multiplication:
		return g.returnHuman(humn, n/otherN)
	case substraction:
		if left.humn {
			return g.returnHuman(humn, n+otherN) // humn - other = n
		}
		return g.returnHuman(humn, otherN-n) // other - humn = n =>
	case division:
		if left.humn {
			return g.returnHuman(humn, n*otherN)
		}
		return g.returnHuman(humn, otherN/n)
	}
	panic((1))
}
