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
	file typ = iota
	dir
)

type typ int

type parser struct {
	s *state
}

func newParser() *parser {
	return &parser{
		s: newState(),
	}

}
func (p *parser) parse(line string) {
	if line == "$ ls" {
		return
	}
	if line == "$ cd /" {
		p.s.navigateRoot()
		return
	}
	if line == "$ cd .." {
		p.s.navigateUp()
		return
	}
	if strings.HasPrefix(line, "$ cd") {
		p.s.navigateDown(line[5:])
		return
	}
	if strings.HasPrefix(line, "dir") {
		p.s.addFiles(makeDir(line[4:]))
		return
	}
	parts := strings.Split(line, " ")
	size, err := strconv.Atoi(parts[0])
	if err != nil {
		panic(err)
	}
	p.s.addFiles(makeFile(parts[1], size))
}

type state struct {
	current *fileDir
	path    string
	history []*fileDir
}

func newState() *state {
	root := makeDir("/")
	return &state{
		path:    "/",
		current: root,
		history: []*fileDir{root},
	}
}

func (s *state) navigateRoot() {
	s.path = "/"
	s.current = s.history[0]
	s.history = []*fileDir{}
}
func (s *state) navigateUp() {
	s.path = "/"
	s.current = s.history[len(s.history)-1]
	s.history = s.history[:len(s.history)-1]
}
func (s *state) navigateDown(name string) {
	s.path += ("/" + name)
	for _, f := range s.current.children {
		if f.name != name {
			continue
		}
		if f.typ == file {
			panic("trying to go down on file")
		}
		s.history = append(s.history, s.current)
		s.current = f
		return
	}
	panic("could not find this file")
}

func (s *state) addFiles(fs ...*fileDir) {
	s.current.addFiles(fs)
}

type fileDir struct {
	name     string
	size     int
	children []*fileDir
	typ      typ
}

func (fd *fileDir) getTree() string {
	b := strings.Builder{}
	switch fd.typ {
	case file:
		b.WriteString(fd.name)
		b.WriteString(" (file) ")
		b.WriteString(fmt.Sprint(fd.size))
	case dir:
		b.WriteString(fd.name)
		b.WriteString(" (dir) ")
		b.WriteString(fmt.Sprint(fd.getSize()))
		b.WriteString("\n")
		for _, c := range fd.children {
			s := c.getTree()
			lines := strings.Split(s, "\n")
			for _, l := range lines {
				b.WriteString("    ")
				b.WriteString(l)
				b.WriteString("\n")
			}
		}
	}
	return b.String()
}

func (fd *fileDir) getSize() int {
	switch fd.typ {
	case file:
		return fd.size
	case dir:
		sum := 0
		for _, c := range fd.children {
			sum += c.getSize()
		}
		return sum
	}
	panic("unknown file typ")
}

func getSize100000(fd *fileDir, resp *[]*fileDir) {
	if fd.typ == file {
		return
	}
	if fd.getSize() <= 100000 {
		*resp = append(*resp, fd)
	}
	for _, c := range fd.children {
		getSize100000(c, resp)
	}
}
func searchToDelete(fd *fileDir, currentLeast *fileDir, leastBound int) *fileDir {
	if fd.typ == file {
		return currentLeast
	}
	if fd.getSize() >= leastBound {
		if fd.getSize() <= currentLeast.getSize() {
			currentLeast = fd
		}
	}
	for _, c := range fd.children {
		currentLeast = searchToDelete(c, currentLeast, leastBound)
	}
	return currentLeast
}

func (fd *fileDir) addFiles(files []*fileDir) {
	fd.children = append(fd.children, files...)
}

func makeFile(name string, size int) *fileDir {
	return &fileDir{
		size: size,
		name: name,
		typ:  file,
	}
}

func makeDir(name string, files ...*fileDir) *fileDir {
	return &fileDir{
		name:     name,
		children: files,
		typ:      dir,
	}
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
	p := newParser()
	for scanner.Scan() {
		p.parse(scanner.Text())
	}
	p.s.navigateRoot()
	totalSize := 70000000
	usedSize := p.s.current.getSize()
	fmt.Printf("used size is %d\n", usedSize)
	unusedSize := totalSize - usedSize
	fmt.Printf("unused size is %d\n", unusedSize)
	requiredSize := 30000000
	toFreeSize := requiredSize - unusedSize
	fmt.Printf("size to free is %d\n", toFreeSize)
	least := searchToDelete(p.s.current, p.s.current, toFreeSize)
	fmt.Println(least.getSize())
}
