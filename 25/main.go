package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	txt, err := os.ReadFile("input.txt")
	if err != nil {
		panic(err)
	}
	sns := readNumbers(string(txt))
	ns := make([]int, len(sns))
	for i, n := range sns {
		ns[i] = fromSnafuToInt(n)
	}
	s := 0
	for _, n := range ns {
		s += n
	}
	fmt.Printf("result %s\n", fromIntToSnafu(s))
}
func readNumbers(txt string) []snafu {
	lines := strings.Split(txt, "\n")
	sns := make([]snafu, len(lines))
	for i, line := range lines {
		ss := strings.Split(line, "")
		sns[i] = ss
	}
	return sns
}

type snafu []string

func fromSnafuToInt(sn snafu) int {
	s := 0
	r := 1
	for i := len(sn) - 1; i >= 0; i-- {
		c := changeStringTo10(sn[i])
		s = s + r*c
		r *= 5
	}
	return s
}

func fromIntToSnafu(n int) snafu {
	int5 := from10To5(n)
	sn := snafu{}
	last := false
	for i := 0; i < len(int5); i++ {
		c, carry := change5ToString(int5[i])
		sn = append(sn, c)
		if carry && i < len(int5)-1 {
			int5[i+1]++
		}
		if carry && i >= len(int5)-1 {
			last = true
		}
	}
	if last {
		sn = append(sn, "1")
	}
	for i, j := 0, len(sn)-1; i < j; i, j = i+1, j-1 {
		sn[i], sn[j] = sn[j], sn[i]
	}

	return sn
}

func from10To5(n int) []int {
	stack := []int{}
	for {
		if n <= 0 {
			break
		}
		rem := n % 5
		stack = append(stack, rem)
		n = n / 5
	}
	return stack
}

func change5ToString(n int) (string, bool) {
	switch n {
	case 5:
		return "0", true
	case 4:
		return "-", true
	case 3:
		return "=", true
	case 2:
		return "2", false
	case 1:
		return "1", false
	case 0:
		return "0", false
	}
	panic(1)
}
func changeStringTo10(s string) int {
	switch s {
	case "0":
		return 0
	case "1":
		return 1
	case "2":
		return 2
	case "-":
		return -1
	case "=":
		return -2
	}
	panic(1)
}
