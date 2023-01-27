package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	bs, err := os.ReadFile("input.txt")
	if err != nil {
		panic(err)
	}
	g := readGame(string(bs))
	finalPosition := g.solve()
	rs := drawGame(g)
	drawCrumbs(rs, g)
	draw(rs)
	fmt.Println(1000*finalPosition.c.y + 4*finalPosition.c.x + int(finalPosition.d))
}
func drawGame(g *game) [][]rune {
	rs := make([][]rune, len(g.m.m))
	for i, l := range g.m.m {
		rs[i] = make([]rune, len(l))
		for j, c := range l {
			switch c {
			case empty:
				rs[i][j] = 'o'
			case open:
				rs[i][j] = '.'
			case wall:
				rs[i][j] = '#'
			}
		}
	}
	return rs
}
func drawCrumbs(rs [][]rune, g *game) {
	for _, c := range g.crumbs {
		rs[c.c.y][c.c.x] = c.r
	}
}
func draw(rs [][]rune) {
	b := strings.Builder{}
	for _, rs := range rs {
		for _, r := range rs {
			b.WriteRune(r)
		}
		b.WriteRune('\n')
	}
	fmt.Println(b.String())
}

func readGame(txt string) *game {
	lines := strings.Split(txt, "\n")
	p := position{}
	m := [][]tileType{}
	// fill first with empty
	lineLen := len(lines[0]) + 2
	first := make([]tileType, lineLen)
	for i := range first {
		first[i] = empty
	}
	m = append(m, first)
	// fill middle
	for _, s := range lines {
		if s == "" {
			break
		}
		mapLine := make([]tileType, lineLen)
		for i := range mapLine {
			mapLine[i] = empty
		}
		mapLine[0] = empty
		for i, c := range s {
			t := empty
			switch c {
			case ' ':
				t = empty
			case '.':
				t = open
			case '#':
				t = wall
			default:
				panic("bad tiletype")
			}
			mapLine[i+1] = t
		}
		mapLine[len(mapLine)-1] = empty
		m = append(m, mapLine)
	}
	// fill last
	last := make([]tileType, lineLen)
	for i := range last {
		last[i] = empty
	}
	m = append(m, last)

	for i, c := range m[1] {
		if c == open {
			p.c.x = i
			p.c.y = 1
			break
		}
	}

	orders := []string{}
	splitR := strings.Split(lines[len(lines)-1], "R")
	for _, s := range splitR {
		splitL := strings.Split(s, "L")
		for _, l := range splitL {
			orders = append(orders, l)
			orders = append(orders, "L")
		}
		orders[len(orders)-1] = "R"
	}
	orders = orders[:len(orders)-1]

	return &game{
		p:      p,
		m:      spaceMap{m: m, e: getEdges()},
		orders: orders,
		crumbs: make([]crumb, 0),
	}
}

type direction int
type tileType int

type coord struct {
	x int
	y int
}

const (
	east direction = iota
	south
	west
	north

	empty tileType = iota
	open
	wall
)

type game struct {
	p      position
	m      spaceMap
	orders []string
	crumbs []crumb
}
type crumb struct {
	r rune
	c coord
}

func (g *game) solve() position {
	for i, o := range g.orders[:len(g.orders)-1] {
		if i%2 == 1 {
			if o != "R" && o != "L" {
				panic("bad rotation")
			}
			g.p = g.p.rotate(o == "R")
		} else {
			n, err := strconv.Atoi(o)
			if err != nil {
				panic(err)
			}
			g.p = g.p.walk(n, g.m, &g.crumbs)
		}
	}
	n, err := strconv.Atoi(g.orders[len(g.orders)-1])
	if err != nil {
		panic(err)
	}
	return g.p.walk(n, g.m, &g.crumbs)
}

type position struct {
	c coord
	d direction
}

func (p position) walk(n int, m spaceMap, crumbs *[]crumb) position {
	startPos := p
	for i := 0; i < n; i++ {
		nextPos, ok := startPos.walkOne(m, crumbs)
		if !ok {
			break
		}
		startPos = nextPos
	}
	return startPos
}

func (p position) walkOne(m spaceMap, crumbs *[]crumb) (position, bool) {
	nextCoord := p.c
	direction := p.d
	var r rune
	switch direction {
	case east:
		r = '>'
		nextCoord.x++
	case south:
		r = 'v'
		nextCoord.y++
	case west:
		r = '<'
		nextCoord.x--
	case north:
		r = '^'
		nextCoord.y--
	default:
		panic(direction)
	}
	crumb := crumb{
		c: nextCoord,
		r: r,
	}
	nextTile := m.m[nextCoord.y][nextCoord.x]
	switch nextTile {
	case empty:
	case open:
		*crumbs = append(*crumbs, crumb)
		return position{c: nextCoord, d: p.d}, true
	case wall:
		return p, false
	}
	otherSidePosition := m.findOtherSide(p)
	nextTile = m.m[otherSidePosition.c.y][otherSidePosition.c.x]
	crumb.c = otherSidePosition.c
	switch nextTile {
	case open:
		*crumbs = append(*crumbs, crumb)
		return otherSidePosition, true
	case wall:
		return p, false
	}
	panic("error")
}

func (p position) rotate(right bool) position {
	delta := direction(3)
	if right {
		delta = direction(1)
	}
	return position{c: p.c, d: (p.d + delta) % 4}
}

type spaceMap struct {
	m [][]tileType
	e []pair
}

func (m spaceMap) findOtherSide(p position) position {
	for _, pa := range m.e {
		if pa.isOn(p) {
			return pa.getNewPosition(p)
		}
	}
	panic("nothing found")
}

type pair struct {
	f edge
	s edge
}

func (pa pair) isOn(p position) bool {
	if pa.f.isOn(p) {
		return true
	}
	if pa.s.isOn(p) {
		return true
	}
	return false
}

func (pa pair) getNewPosition(p position) position {
	if pa.f.isOn(p) {
		d := pa.f.toParameter(p.c)
		newCord := pa.s.fromParameter(d)
		return position{c: newCord, d: pa.s.outDirection}
	}
	if pa.s.isOn(p) {
		d := pa.s.toParameter(p.c)
		newCord := pa.f.fromParameter(d)
		return position{c: newCord, d: pa.f.outDirection}
	}
	panic("not on this edge")
}

type edge struct {
	start        coord
	end          coord
	outDirection direction
}

func (e edge) isOn(p position) bool {
	if (p.d+2)%4 != e.outDirection {
		return false
	}
	c := p.c
	if e.start.x == e.end.x {
		if c.x != e.start.x {
			return false
		}
		if e.start.y > e.end.y {
			return e.start.y >= c.y && c.y >= e.end.y
		}
		return e.start.y <= c.y && c.y <= e.end.y
	}
	if e.start.y == e.end.y {
		if c.y != e.start.y {
			return false
		}
		if e.start.x > e.end.x {
			return e.start.x >= c.x && c.x >= e.end.x
		}
		return e.start.x <= c.x && c.x <= e.end.x
	}
	panic("diagonal")
}
func (e edge) toParameter(c coord) int {
	if e.start.x == e.end.x {
		t := c.y - e.start.y
		if t < 0 {
			return -t
		}
		return t
	}
	if e.start.y == e.end.y {
		t := c.x - e.start.x
		if t < 0 {
			return -t
		}
		return t
	}
	panic("diagonal")
}
func (e edge) fromParameter(i int) coord {
	if e.start.x == e.end.x {
		if e.start.y > e.end.y {
			i = -i
		}
		return coord{x: e.start.x, y: e.start.y + i}
	}
	if e.start.y == e.end.y {
		if e.start.x > e.end.x {
			i = -i
		}
		return coord{x: e.start.x + i, y: e.start.y}
	}
	panic("diagonal")
}

func getEdges() []pair {
	return []pair{
		{
			f: edge{start: coord{x: 150, y: 50}, end: coord{x: 101, y: 50}, outDirection: north},
			s: edge{start: coord{x: 100, y: 100}, end: coord{x: 100, y: 51}, outDirection: west},
		},
		{
			f: edge{start: coord{x: 100, y: 150}, end: coord{x: 51, y: 150}, outDirection: north},
			s: edge{start: coord{x: 50, y: 200}, end: coord{x: 50, y: 151}, outDirection: west},
		},
		{
			f: edge{start: coord{x: 51, y: 51}, end: coord{x: 51, y: 100}, outDirection: east},
			s: edge{start: coord{x: 1, y: 101}, end: coord{x: 50, y: 101}, outDirection: south},
		},
		{
			f: edge{start: coord{x: 150, y: 50}, end: coord{x: 150, y: 1}, outDirection: west},
			s: edge{start: coord{x: 100, y: 101}, end: coord{x: 100, y: 150}, outDirection: west},
		},
		{
			f: edge{start: coord{x: 1, y: 150}, end: coord{x: 1, y: 101}, outDirection: east},
			s: edge{start: coord{x: 51, y: 1}, end: coord{x: 51, y: 50}, outDirection: east},
		},

		{
			f: edge{start: coord{x: 51, y: 1}, end: coord{x: 100, y: 1}, outDirection: south},
			s: edge{start: coord{x: 1, y: 151}, end: coord{x: 1, y: 200}, outDirection: east},
		},
		{
			f: edge{start: coord{x: 101, y: 1}, end: coord{x: 150, y: 1}, outDirection: south},
			s: edge{start: coord{x: 1, y: 200}, end: coord{x: 50, y: 200}, outDirection: north},
		},
	}
}
