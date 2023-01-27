package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

const (
	integer typ = iota
	list

	left ordering = iota
	noorder
	right
)

type typ int
type ordering int

func main() {
	f, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	var packets []packetWithEnvelope
	p1 := readPacket("[[2]]")
	p2 := readPacket("[[6]]")
	packets = append(packets, packetWithEnvelope{packet: p1, id: 0}, packetWithEnvelope{packet: p2, id: 1})
	scanner := bufio.NewScanner(f)
	i := 2
	for scanner.Scan() {
		line1 := scanner.Text()
		p1 := readPacket(line1)
		scanner.Scan()
		line2 := scanner.Text()
		p2 := readPacket(line2)
		packets = append(packets, packetWithEnvelope{packet: p1, id: i}, packetWithEnvelope{packet: p2, id: i + 1})
		scanner.Scan()
		i++
		i++
	}
	sort.SliceStable(packets, func(i, j int) bool {
		p1 := packets[i]
		p2 := packets[j]
		res := compare(p1.packet, p2.packet)
		return res == left
	})
	id2 := 0
	id6 := 0
	for i, p := range packets {
		if p.id == 0 {
			id2 = i + 1
		}
		if p.id == 1 {
			id6 = i + 1
		}
	}
	fmt.Println(id2 * id6)

}
func readPacket(line string) value {
	trimmed := strings.TrimPrefix(line, "[")
	v, _ := readList(trimmed)
	return v
}

func readList(line string) (value, int) {
	i := 0
	var vals []value
	var b strings.Builder
	for {
		c := line[i]
		if c == '[' {
			li, n := readList(line[i+1:])
			vals = append(vals, li)
			i += n + 1
			continue
		}
		if c == ']' {
			s := b.String()
			b.Reset()
			if s == "" {
				i++
				break
			}
			n, err := strconv.Atoi(s)
			if err != nil {
				panic(err)
			}
			vals = append(vals, newInt(n))
			i++
			break
		}
		if c == ',' {
			s := b.String()
			b.Reset()
			if s == "" {
				i++
				continue
			}
			n, err := strconv.Atoi(s)
			if err != nil {
				panic(err)
			}
			vals = append(vals, newInt(n))
			i++
			continue
		}
		b.WriteByte(c)
		i++
	}
	return newList(vals...), i
}

type packetWithEnvelope struct {
	packet value
	id     int
}
type value struct {
	typ     typ
	integer int
	list    []value
}

func (v value) isInteger() bool {
	return v.typ == integer
}

func (v value) isList() bool {
	return v.typ == list
}

func newInt(n int) value {
	return value{
		typ:     integer,
		integer: n,
	}
}

func newList(values ...value) value {
	return value{
		typ:  list,
		list: values,
	}
}

func compare(v1, v2 value) ordering {
	if v1.isInteger() {
		if v2.isList() {
			return compare(newList(v1), v2)
		}
		return compareInts(v1.integer, v2.integer)
	}
	if v2.isInteger() {
		return compare(v1, newList(v2))
	}
	return compareList(v1.list, v2.list)
}

func compareInts(v1, v2 int) ordering {
	if v1 < v2 {
		return left
	}
	if v1 == v2 {
		return noorder
	}
	return right
}
func compareList(v1, v2 []value) ordering {
	lenV1 := len(v1)
	lenV2 := len(v2)
	i := 0
	for {
		if lenV1 < i+1 {
			if lenV1 == lenV2 {
				return noorder
			}
			return left
		}
		if lenV2 < i+1 {
			return right
		}
		switch compare(v1[i], v2[i]) {
		case left:
			return left
		case right:
			return right
		case noorder:
			i++
			continue
		}
	}
}
