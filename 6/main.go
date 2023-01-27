package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

const chunkSize = 14

func main() {
	f, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	bs, err := ioutil.ReadAll(f)
	if err != nil {
		log.Fatal(err)
	}
	var end int
	for i := chunkSize; i < len(bs); i++ {
		chunk := bs[i-chunkSize : i]
		if isAllElementDifferent(chunk) {
			end = i
			break
		}
	}
	fmt.Println(end)
}

func isAllElementDifferent(chunk []byte) bool {
	s := make(map[byte]struct{}, chunkSize)
	for _, el := range chunk {
		_, ok := s[el]
		if ok {
			return false
		}
		s[el] = struct{}{}
	}
	return true
}
