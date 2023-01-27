package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
)

const (
	KEY      = 811589153
	MIXTIMES = 10
)

func main() {
	f, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	ns := []int{}
	for scanner.Scan() {
		line := scanner.Text()
		n, err := strconv.Atoi(line)
		if err != nil {
			panic(err)
		}
		ns = append(ns, n*KEY)
	}
	n := len(ns)
	el := newCycliList(ns)
	first := el
	var zero *listElement
	for i := 0; i < MIXTIMES; i++ {
		el = first
		for {
			if el.n == 0 {
				zero = el
				el = el.nextOriginal()
				continue
			}
			eln := el.n
			el.cut()
			target := el.move(eln % (n - 1))
			if eln < 0 {
				target.insertPrevious(el)
			} else {
				el.cut()
				target.insertNext(el)
			}
			if el.nextOriginal() == nil {
				break
			}
			el = el.nextOriginal()
		}
	}
	n1k := zero.move(1000 % n)
	n2k := zero.move(2000 % n)
	n3k := zero.move(3000 % n)
	fmt.Println(n1k.n, n2k.n, n3k.n)
	fmt.Println(n1k.n + n2k.n + n3k.n)
}

func newCycliList(ns []int) *listElement {
	if len(ns) == 0 {
		panic("list len = 0")
	}
	first := &listElement{
		n: ns[0],
	}
	last := first
	for _, n := range ns[1:] {
		cur := &listElement{
			n: n,
		}
		last.next = cur
		last.originalNext = cur
		cur.previous = last
		last = cur
	}
	first.previous = last
	last.next = first
	return first
}

type listElement struct {
	n            int
	previous     *listElement
	next         *listElement
	originalNext *listElement
}

func (le *listElement) nextOriginal() *listElement {
	return le.originalNext
}

func (le *listElement) move(n int) *listElement {
	if n == 0 {
		return le
	}
	if n < 0 {
		return le.previous.move(n + 1)
	}
	return le.next.move(n - 1)
}

func (le *listElement) cut() {
	le.previous.next = le.next
	le.next.previous = le.previous
}

func (left *listElement) insertNext(middle *listElement) {
	right := left.next
	left.next = middle
	middle.next = right
	right.previous = middle
	middle.previous = left
}

func (right *listElement) insertPrevious(middle *listElement) {
	left := right.previous
	left.next = middle
	middle.next = right
	right.previous = middle
	middle.previous = left
}
