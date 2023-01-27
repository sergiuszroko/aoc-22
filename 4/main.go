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
	var s int
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		txt := scanner.Text()
		var a1, b1, a2, b2 int
		fmt.Println(txt)
		_, err := fmt.Fscanf(strings.NewReader(txt), "%d-%d,%d-%d\n", &a1, &b1, &a2, &b2)
		if err != nil {
			log.Fatal(err)
		}
		one := interval{a1, b1}
		two := interval{a2, b2}
		if contain(one, two) {
			fmt.Printf("%v %v\n", one, two)
			s += 1
		}
	}
	fmt.Printf("n of overlapping is %d\n", s)
}

type interval struct {
	a, b int
}

func contain(one, two interval) bool {
	if one.a <= two.a && two.b <= one.b {
		return true
	}
	if two.a <= one.a && one.b <= two.b {
		return true
	}
	if two.a <= one.b && one.b <= two.b {
		return true
	}
	if two.a <= one.a && one.a <= two.b {
		return true
	}
	return false
}
