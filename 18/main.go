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
	d := newDroplet()
	for scanner.Scan() {
		line := scanner.Text()
		c := readLine(line)
		d.add(c)
	}
	fmt.Println(d.count())
	clz := clusterize{
		max:     30,
		min:     -5,
		d:       d,
		visited: map[cord]struct{}{},
	}
	clz.check()
	// clz.printClusters()
	fill := filler{
		clz: clz,
		d:   d,
	}
	fill.fillAllExceptBiggest()
	fmt.Println(d.count())
}
func readLine(l string) cord {
	c := cord{}
	_, err := fmt.Fscanf(strings.NewReader(l), "%d,%d,%d", &c.x, &c.y, &c.z)
	if err != nil {
		panic(err)
	}
	return c
}

type cord struct {
	x, y, z int
}

func (c *cord) neighbours() []cord {
	return []cord{
		{c.x - 1, c.y, c.z},
		{c.x + 1, c.y, c.z},
		{c.x, c.y - 1, c.z},
		{c.x, c.y + 1, c.z},
		{c.x, c.y, c.z - 1},
		{c.x, c.y, c.z + 1},
	}
}

type droplet struct {
	xyzs map[cord]int
}

func newDroplet() *droplet {
	return &droplet{
		xyzs: make(map[cord]int),
	}
}

func (d *droplet) add(c cord) {
	s := 6
	for _, n := range c.neighbours() {
		nc, ok := d.xyzs[n]
		if ok {
			if nc <= 0 {
				panic("no possible")
			}
			d.xyzs[n]--
			s--
		}
	}
	d.xyzs[c] = s
}

func (d *droplet) count() int {
	s := 0
	for _, c := range d.xyzs {
		s += c
	}
	return s
}

type clusterize struct {
	clusters []cluster
	max, min int
	d        *droplet
	visited  map[cord]struct{}
}

func (clz *clusterize) printClusters() {
	for _, cl := range clz.clusters {
		hint := strings.Builder{}
		i := 0
		for k := range cl.cords {
			if i > 10 {
				break
			}
			hint.WriteString(fmt.Sprintf("%v,", k))
		}
		fmt.Printf("cluster size %d: %v\n", len(cl.cords), hint.String())
	}
}

func (clz *clusterize) check() {
	for i := 0; i < clz.max; i++ {
		for j := 0; j < clz.max; j++ {
			for k := 0; k < clz.max; k++ {
				c := cord{i, j, k}
				if cl, ok := clz.checkCord(c); ok {
					clz.clusters = append(clz.clusters, newCluster(cl))
				}
			}
		}
	}
}

func (clz *clusterize) checkCord(c cord) ([]cord, bool) {
	if ok := clz.checkBoundary(c); !ok {
		return nil, false
	}
	if _, ok := clz.visited[c]; ok {
		return nil, false
	}
	clz.visited[c] = struct{}{}
	if _, ok := clz.d.xyzs[c]; ok {
		return nil, false
	}
	for _, cl := range clz.clusters {
		if _, ok := cl.cords[c]; ok {
			return nil, false
		}
	}
	r := []cord{c}
	for _, n := range c.neighbours() {
		if more, ok := clz.checkCord(n); ok {
			r = append(r, more...)
		}
	}
	return r, true
}

func (clz *clusterize) checkBoundary(c cord) bool {
	if c.x < clz.min {
		return false
	}
	if c.y < clz.min {
		return false
	}
	if c.z < clz.min {
		return false
	}
	if c.x > clz.max {
		return false
	}
	if c.y > clz.max {
		return false
	}
	if c.z > clz.max {
		return false
	}
	return true
}

type cluster struct {
	cords map[cord]struct{}
}

func newCluster(cs []cord) cluster {
	cl := cluster{
		cords: make(map[cord]struct{}),
	}
	for _, c := range cs {
		cl.cords[c] = struct{}{}
	}
	return cl
}

type filler struct {
	clz clusterize
	d   *droplet
}

func (f *filler) fillAllExceptBiggest() {
	indexOfBiggest := f.findBiggest()
	for i, cl := range f.clz.clusters {
		if i == indexOfBiggest {
			continue
		}
		for k := range cl.cords {
			f.d.add(k)
		}
	}
}
func (f *filler) findBiggest() int {
	var idx, size int
	for i, cl := range f.clz.clusters {
		if size < len(cl.cords) {
			idx = i
			size = len(cl.cords)
		}
	}
	return idx
}
