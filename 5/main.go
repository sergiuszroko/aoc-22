package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	data := [][]rune{
		{'R', 'S', 'L', 'F', 'Q'},
		{'N', 'Z', 'Q', 'G', 'P', 'T'},
		{'S', 'M', 'Q', 'B'},
		{'T', 'G', 'Z', 'J', 'H', 'C', 'B', 'Q'},
		{'P', 'H', 'M', 'B', 'N', 'F', 'S'},
		{'P', 'C', 'Q', 'N', 'S', 'L', 'V', 'G'},
		{'W', 'C', 'F'},
		{'Q', 'H', 'G', 'Z', 'W', 'V', 'P', 'M'},
		{'G', 'Z', 'D', 'L', 'C', 'N', 'R'},
	}
	f, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		txt := scanner.Text()
		var n, from, to int
		_, err := fmt.Fscanf(strings.NewReader(txt), "move %d from %d to %d\n", &n, &from, &to)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(n, from, to)
		fmt.Println(data)
		data[to-1] = append(data[to-1], data[from-1][len(data[from-1])-n:len(data[from-1])]...)
		fmt.Println("b")
		data[from-1] = data[from-1][:len(data[from-1])-n]
	}
	for _, el := range data {
		fmt.Print(string(el[len(el)-1]))
	}
	fmt.Println()
}
