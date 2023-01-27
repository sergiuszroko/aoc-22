package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

const (
	ore robot = iota
	clay
	obsidian
	geode
)

type robot int

func (r robot) getNext() robot {
	return (r + 1) % 3
}

func main() {
	f, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	s := 1
	for scanner.Scan() {
		line := scanner.Text()
		_, sim := readLine(line)
		max := sim.Sim(ore)
		g := sim.Sim(clay)
		if g > max {
			max = g
		}
		g = sim.Sim(obsidian)
		if g > max {
			max = g
		}
		fmt.Println(max)
		s *= max
	}
	fmt.Println(s)
}
func readLine(l string) (int, simulation) {
	sim := simulation{maxTime: 32}
	var id int
	_, err := fmt.Fscanf(
		strings.NewReader(l),
		"Blueprint %d: Each ore robot costs %d ore. Each clay robot costs %d ore. Each obsidian robot costs %d ore and %d clay. Each geode robot costs %d ore and %d obsidian.",
		&id,
		&sim.buyOre.oreCost,
		&sim.buyClay.oreCost,
		&sim.buyObsidian.oreCost, &sim.buyObsidian.clayCost,
		&sim.buyGeode.oreCost, &sim.buyGeode.obsidianCost,
	)
	if err != nil {
		panic(err)
	}
	maxOreCost := sim.buyOre.oreCost
	if maxOreCost < sim.buyClay.oreCost {
		maxOreCost = sim.buyClay.oreCost
	}
	if maxOreCost < sim.buyObsidian.oreCost {
		maxOreCost = sim.buyObsidian.oreCost
	}
	if maxOreCost < sim.buyGeode.oreCost {
		maxOreCost = sim.buyGeode.oreCost
	}
	sim.buyOre.oreCostOther = maxOreCost
	sim.buyClay.clayCostOther = sim.buyObsidian.clayCost
	sim.buyObsidian.obsidianCostOther = sim.buyGeode.obsidianCost
	return id, sim
}

type bank struct {
	ore      int
	clay     int
	obsidian int
}
type team struct {
	ore      int
	clay     int
	obsidian int
	geode    int
}

type simulation struct {
	buyOre      buyOre
	buyClay     buyClay
	buyObsidian buyObsidian
	buyGeode    buyGeode
	buyNext     buy
	maxTime     int
}

func (s simulation) Sim(first robot) int {
	s.buyNext = buy{
		buyOre:        s.buyOre,
		buyClay:       s.buyClay,
		buyObsidian:   s.buyObsidian,
		whatToBuyNext: first,
	}
	return s.simulate(0, bank{}, team{ore: 1}, 0, first)
}

func (s simulation) simulate(time int, b bank, t team, geode int, first robot) int {
	s.buyNext.whatToBuyNext = first
	maxGeode := 0
	for {
		time++

		if time >= s.maxTime+1 {
			if maxGeode > geode {
				return maxGeode
			}
			return geode
		}

		b.ore += t.ore
		b.clay += t.clay
		b.obsidian += t.obsidian
		geode += t.geode

		if s.buyGeode.shouldBuy(b, t) {
			b, t = s.buyGeode.buy(b, t)
		} else if s.buyNext.canBay(b, t, s.maxTime-time) {
			b, t = s.buyNext.buy(b, t)
			maxGeodeOre := s.simulate(time, b, t, geode, first.getNext())
			if maxGeode < maxGeodeOre {
				maxGeode = maxGeodeOre
			}
			maxGeodeOre = s.simulate(time, b, t, geode, first.getNext().getNext())
			if maxGeode < maxGeodeOre {
				maxGeode = maxGeodeOre
			}
		}
	}
}

type buy struct {
	whatToBuyNext robot
	buyOre        buyOre
	buyClay       buyClay
	buyObsidian   buyObsidian
}

func (bu *buy) canBay(b bank, t team, timeLeft int) bool {
	switch bu.whatToBuyNext {
	case ore:
		return bu.buyOre.shouldBuy(b, t, timeLeft)
	case clay:
		return bu.buyClay.shouldBuy(b, t, timeLeft)
	case obsidian:
		return bu.buyObsidian.shouldBuy(b, t, timeLeft)
	}
	panic(1)
}

func (bu *buy) buy(b bank, t team) (bank, team) {
	switch bu.whatToBuyNext {
	case ore:
		return bu.buyOre.buy(b, t)
	case clay:
		return bu.buyClay.buy(b, t)
	case obsidian:
		return bu.buyObsidian.buy(b, t)
	}
	panic(1)
}

type buyOre struct {
	oreCost      int
	oreCostOther int
}

func (bo buyOre) shouldBuy(b bank, t team, timeLeft int) bool {
	if b.ore-t.ore < bo.oreCost {
		return false
	}
	if timeLeft*bo.oreCostOther <= b.ore+t.ore {
		return false
	}
	if bo.oreCostOther <= t.ore {
		return false
	}
	return true
}
func (bo buyOre) buy(b bank, t team) (bank, team) {
	b.ore -= bo.oreCost
	t.ore++
	return b, t
}

type buyClay struct {
	oreCost       int
	clayCostOther int
}

func (bo buyClay) shouldBuy(b bank, t team, timeLeft int) bool {
	if b.ore-t.ore < bo.oreCost {
		return false
	}
	if timeLeft*bo.clayCostOther <= b.clay+t.clay {
		return false
	}
	if bo.clayCostOther <= t.clay {
		return false
	}
	return true
}
func (bo buyClay) buy(b bank, t team) (bank, team) {
	b.ore -= bo.oreCost
	t.clay++
	return b, t
}

type buyObsidian struct {
	oreCost           int
	clayCost          int
	obsidianCostOther int
}

func (bo buyObsidian) shouldBuy(b bank, t team, timeleft int) bool {
	if b.ore-t.ore < bo.oreCost {
		return false
	}
	if b.clay-t.clay < bo.clayCost {
		return false
	}
	if timeleft*bo.obsidianCostOther <= b.obsidian+t.obsidian {
		return false
	}
	if bo.obsidianCostOther <= t.obsidian {
		return false
	}
	return true
}
func (bo buyObsidian) buy(b bank, t team) (bank, team) {
	b.clay -= bo.clayCost
	b.ore -= bo.oreCost
	t.obsidian++
	return b, t
}

type buyGeode struct {
	oreCost      int
	obsidianCost int
}

func (bo buyGeode) shouldBuy(b bank, t team) bool {
	if b.ore-t.ore < bo.oreCost {
		return false
	}
	if b.obsidian-t.obsidian < bo.obsidianCost {
		return false
	}
	return true
}
func (bo buyGeode) buy(b bank, t team) (bank, team) {
	b.obsidian -= bo.obsidianCost
	b.ore -= bo.oreCost
	t.geode++
	return b, t
}
