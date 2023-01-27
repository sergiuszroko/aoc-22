package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
)

const input = "input.txt"

func main() {
	f, err := os.Open(input)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	var maxs [3]int
	var current int
	for scanner.Scan() {
		text := scanner.Text()
		if text == "" {
			addToMaxs(&maxs, current)
			fmt.Println(current)
			current = 0
			continue
		}
		n, err := strconv.Atoi(text)
		if err != nil {
			log.Fatal(err)
		}
		current += n
	}
	addToMaxs(&maxs, current)
	fmt.Printf("calories in top 3 backpacks is %d\n", sum(maxs))

}

func addToMaxs(maxs *[3]int, current int) {
	for i, el := range maxs {
		if el < current {
			maxs[i] = current
			current = el
		}
	}
}

func sum(maxs [3]int) int {
	var sum int
	for _, el := range maxs {
		sum += el
	}
	return sum
}
