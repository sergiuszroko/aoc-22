package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
)

const (
	X          = 1000
	Y          = 500
	sandStartX = 500
	sandStartY = 0
)

func main() {
	f, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	readings := []sensorReading{}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		readings = append(readings, readLine(line))
	}
	fmt.Println(determineNonOccupied(readings))

}

func readLine(line string) sensorReading {
	r := sensorReading{}
	_, err := fmt.Fscanf(strings.NewReader(line), "Sensor at x=%d, y=%d: closest beacon is at x=%d, y=%d", &r.sensor.x, &r.sensor.y, &r.closestBeacond.x, &r.closestBeacond.y)
	if err != nil {
		panic(err)
	}
	return r
}

type coord struct {
	x, y int
}
type interval struct {
	start int
	end   int
}
type sensorReading struct {
	sensor         coord
	closestBeacond coord
}

func determineNonOccupied(readings []sensorReading) int {
	for i := 0; i < 4000000; i++ {
		c, found := determineNonOccupiedInRow(readings, i)
		if found {
			fmt.Println(c)
			return c.x*4000000 + c.y
		}
	}
	panic("nothing found")
}

func determineNonOccupiedInRow(readings []sensorReading, rowY int) (coord, bool) {
	ranges := []interval{}
	for _, reading := range readings {
		manhattan := getManhattanMetric(reading.sensor, reading.closestBeacond)
		nonOccupiedFromReading, found := getIntervalOnLineInDistance(rowY, reading.sensor, manhattan)
		if !found {
			continue
		}
		ranges = append(ranges, nonOccupiedFromReading)
	}

	x, found := mergeIntervals(ranges)
	return coord{x, rowY}, found
}

func getIntervalOnLineInDistance(rowY int, centre coord, dist int) (interval, bool) {
	yDist := abs(rowY - centre.y)
	if yDist > dist {
		return interval{}, false
	}
	start := centre.x - (dist - yDist)
	if start < 0 {
		start = 0
	}
	end := centre.x + (dist - yDist)
	if end > 4000000 {
		end = 4000000
	}
	return interval{start, end}, true
}

func mergeIntervals(ints []interval) (int, bool) {
	sort.Slice(ints, func(i, j int) bool { return ints[i].start < ints[j].start })
	sum := ints[0]
	if sum.start > 0 {
		return 0, true
	}
	for _, in := range ints[1:] {
		if sum.end+1 < in.start {
			return sum.end + 1, true
		}
		end := in.end
		if sum.end > end {
			end = sum.end
		}
		sum = interval{sum.start, end}
	}
	if sum.end < 4000000 {
		return 4000000, true
	}
	return 0, false
}

func getManhattanMetric(c1, c2 coord) int {
	return abs(c1.x-c2.x) + abs(c1.y-c2.y)
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}
