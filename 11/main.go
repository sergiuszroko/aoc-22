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
	scanner := bufio.NewScanner(f)
	if err != nil {
		log.Fatal(err)
	}
	g := game{}
	for scanner.Scan() {
		m := readMonkey(scanner)
		g.addMonkey(m)
		scanner.Scan()
	}
	g.initialize()
	g.g()
	first, second := g.getMostActive()
	fmt.Println(first)
	fmt.Println(second)
	fmt.Println(first * second)
}

func readMonkey(scanner *bufio.Scanner) monkey {
	m := monkey{}
	scanner.Scan()
	itemsLine := scanner.Text()
	itemsLine = strings.TrimPrefix(itemsLine, "  Starting items: ")
	listItemsString := strings.Split(itemsLine, ", ")
	listItems := make([]int, len(listItemsString))
	for i, it := range listItemsString {
		tt, err := strconv.Atoi(it)
		if err != nil {
			log.Fatal(err)
		}
		listItems[i] = tt
	}
	m.items = listItems

	scanner.Scan()
	operationLine := scanner.Text()
	if strings.Contains(operationLine, "*") {
		nString := operationLine[len("  Operation: new = old * "):]
		if nString == "old" {
			m.inspect = newPowFunc()
		} else {
			n, err := strconv.Atoi(nString)
			if err != nil {
				log.Fatal(err)
			}
			m.inspect = newMultFunc(n)
		}
	} else {
		nString := operationLine[len("  Operation: new = old + "):]
		n, err := strconv.Atoi(nString)
		if err != nil {
			log.Fatal(err)
		}
		m.inspect = newAddFunc(n)
	}

	t := testFunc{}
	scanner.Scan()
	testLine := scanner.Text()
	nString := testLine[len("  Test: divisible by "):]
	n, err := strconv.Atoi(nString)
	if err != nil {
		log.Fatal(err)
	}
	t.divisor = n

	scanner.Scan()
	trueLine := scanner.Text()
	nString = trueLine[len("    If true: throw to monkey "):]
	n, err = strconv.Atoi(nString)
	if err != nil {
		log.Fatal(err)
	}
	t.ifTrue = n

	scanner.Scan()
	falseLine := scanner.Text()
	nString = falseLine[len("    If false: throw to monkey "):]
	n, err = strconv.Atoi(nString)
	if err != nil {
		log.Fatal(err)
	}
	t.ifFalse = n

	m.test = t
	return m
}

type game struct {
	monkeys  []monkey
	throws   []int
	truncate int
}

func (g *game) initialize() {
	s := 1
	for _, m := range g.monkeys {
		s *= m.test.divisor
	}
	g.truncate = s
}

func (g *game) getMostActive() (int, int) {
	var firstn, secondn int
	for _, n := range g.throws {
		if n < secondn {
			continue
		}
		if n < firstn {
			secondn = n
			continue
		}
		secondn = firstn
		firstn = n
	}
	return firstn, secondn
}

func (g *game) addMonkey(m monkey) {
	g.monkeys = append(g.monkeys, m)
	g.throws = append(g.throws, 0)
}

func (g *game) g() {
	for i := 0; i < 10000; i++ {
		for j := range g.monkeys {
			for {
				fmt.Printf("Monkey %d\n", j)
				if !g.monkeys[j].canThrow() {
					break
				}
				g.throws[j]++
				monkeyN, item := g.monkeys[j].throw()
				fmt.Printf("throws %d to %d\n", item, monkeyN)
				item = item % g.truncate
				g.monkeys[monkeyN].catch(item)
			}
		}
	}
}

type monkey struct {
	items   []int
	inspect inspectFunc
	test    testFunc
}

func (m *monkey) catch(item int) {
	m.items = append(m.items, item)
}

func (m *monkey) canThrow() bool {
	return len(m.items) > 0
}
func (m *monkey) throw() (int, int) {
	it := m.items[0]
	fmt.Printf("got item %d\n", it)
	m.items = m.items[1:len(m.items)]
	newItem := m.inspect(it)
	return m.test.getNewMonkey(newItem), newItem
}

type inspectFunc func(int) int

func newAddFunc(n int) inspectFunc {
	return func(i int) int {
		return i + n
	}
}
func newMultFunc(n int) inspectFunc {
	return func(i int) int {
		return i * n
	}
}
func newPowFunc() inspectFunc {
	return func(i int) int {
		return i * i
	}
}

type testFunc struct {
	divisor int
	ifTrue  int
	ifFalse int
}

func (t *testFunc) getNewMonkey(n int) int {
	div := n % t.divisor
	if div == 0 {
		return t.ifTrue
	}
	return t.ifFalse
}
