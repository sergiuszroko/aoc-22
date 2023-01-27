package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
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
		first := scanner.Text()
		scanner.Scan()
		second := scanner.Text()
		scanner.Scan()
		third := scanner.Text()
		f := compare(collapse((first)), collapse(second), collapse(third))
		p := getPriority(f)
		s += p
	}
	fmt.Printf("sum o priorities is %d\n", s)
}

func collapse(s string) map[rune]struct{} {
	set := make(map[rune]struct{})
	for _, e := range s {
		set[e] = struct{}{}
	}
	return set
}

func compare(set1, set2, set3 map[rune]struct{}) rune {
	var found rune
	for k := range set1 {
		_, ok := set2[k]
		if ok {
			_, ok := set3[k]
			if ok {
				found = k
			}
		}
	}
	return found
}

func getPriority(r rune) int {
	if r >= 96 {
		return int(r - 96)
	}
	return int(r - 38)
}
