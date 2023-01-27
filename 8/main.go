package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
)

func main1() {
	f, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	if err != nil {
		log.Fatal(err)
	}
	matrix := make([][]int, 0)
	for scanner.Scan() {
		txt := scanner.Text()
		line := make([]int, len(txt))
		for i, d := range txt {
			n, err := strconv.Atoi(string(d))
			if err != nil {
				log.Fatal(err)
			}
			line[i] = n
		}
		matrix = append(matrix, line)
	}
	nLines := len(matrix)
	nColumns := len(matrix[0])
	matrixVisible := make([][]bool, nLines)
	for i := range matrixVisible {
		matrixVisible[i] = make([]bool, nColumns)
	}
	for i, l := range matrix {
		h := l[0]
		matrixVisible[i][0] = true
		for j := 0; j < len(l); j++ {
			t := l[j]
			if t > h {
				h = t
				matrixVisible[i][j] = true
			}
		}
		h = l[len(l)-1]
		matrixVisible[i][len(l)-1] = true
		for j := len(l) - 1; j >= 0; j-- {
			t := l[j]
			if t > h {
				h = t
				matrixVisible[i][j] = true
			}
		}
	}
	for j := 0; j < nColumns; j++ {
		h := matrix[0][j]
		matrixVisible[0][j] = true
		for i, l := range matrix {
			t := l[j]
			if t > h {
				h = t
				matrixVisible[i][j] = true
			}
		}
		h = matrix[nLines-1][j]
		matrixVisible[nLines-1][j] = true
		for i := nLines - 1; i >= 0; i-- {
			t := matrix[i][j]
			if t > h {
				h = t
				matrixVisible[i][j] = true
			}
		}
	}
	for _, l := range matrixVisible {
		for _, t := range l {
			if t {
				fmt.Print("t ")
			} else {
				fmt.Print("f ")
			}
		}
		fmt.Println()
	}
	s := 0
	for _, l := range matrixVisible {
		for _, t := range l {
			if t {
				s++
			}
		}
	}
	fmt.Println(s)
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
	matrix := make([][]int, 0)
	for scanner.Scan() {
		txt := scanner.Text()
		line := make([]int, len(txt))
		for i, d := range txt {
			n, err := strconv.Atoi(string(d))
			if err != nil {
				log.Fatal(err)
			}
			line[i] = n
		}
		matrix = append(matrix, line)
	}
	nLines := len(matrix)
	nColumns := len(matrix[0])
	matrixScore := make([][]int, nLines)
	for i := range matrixScore {
		matrixScore[i] = make([]int, nColumns)
	}
	for i, line := range matrix {
		for j, t := range line {
			d := 1
			for {
				if i+d >= nColumns {
					d--
					break
				}
				if matrix[i+d][j] >= t {
					break
				}
				d++
			}
			u := 1
			for {
				if i-u < 0 {
					u--
					break
				}
				if matrix[i-u][j] >= t {
					break
				}
				u++
			}
			r := 1
			for {
				if j+r >= nLines {
					r--
					break
				}
				if matrix[i][j+r] >= t {
					break
				}
				r++
			}
			l := 1
			for {
				if j-l < 0 {
					l--
					break
				}
				if matrix[i][j-l] >= t {
					break
				}
				l++
			}
			matrixScore[i][j] = (u) * (d) * (l) * (r)
		}
	}
	s := 0
	for _, l := range matrixScore {
		for _, t := range l {
			if t > s {
				s = t
			}
		}
	}
	fmt.Println(s)
}
