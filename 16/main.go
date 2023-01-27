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
	valves := map[string]valve{}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		valve := readLine(line)
		valves[valve.name] = valve
	}
	intereseting := findInterestings(valves)
	mapTunnels := findPaths(valves, intereseting)
	g := game{
		costs:    mapTunnels,
		vals:     getVals(valves),
		stateMap: getStateMap(valves),
	}
	ans := g.g("AA", 30, 0, answer{}, 0)
	max1 := 0
	for _, v := range ans {
		if max1 < v {
			max1 = v
		}
	}
	fmt.Println(max1)
	ans = g.g("AA", 26, 0, answer{}, 0)
	max2 := 0
	for k1, v1 := range ans {
		for k2, v2 := range ans {
			if !k1.isSeparate(k2) {
				continue
			}
			cand := v1 + v2
			if cand > max2 {
				max2 = cand
			}
		}
	}
	fmt.Println(max2)

}
func readLine(line string) valve {
	var name string
	var pressure int
	parts := strings.Split(line, ";")
	_, err := fmt.Fscanf(strings.NewReader(parts[0]), "Valve %s has flow rate=%d", &name, &pressure)
	if err != nil {
		panic(err)
	}
	tunnels := strings.Split(parts[1], ", ")
	tunnels[0] = tunnels[0][len(tunnels[0])-2:]
	return valve{
		name:     name,
		pressure: pressure,
		tunnels:  tunnels,
	}
}

func getVals(m map[string]valve) map[string]int {
	vals := make(map[string]int)
	for k, v := range m {
		vals[k] = v.pressure
	}
	return vals
}
func getStateMap(m map[string]valve) map[string]int {
	vals := make(map[string]int)
	i := 0
	for k := range m {
		vals[k] = i
		i++
	}
	return vals
}

type valve struct {
	name     string
	pressure int
	tunnels  []string
}

type path struct {
	start string
	end   string
	cost  int
}

func findInterestings(valves map[string]valve) map[string]struct{} {
	m := make(map[string]struct{}, 0)
	m["AA"] = struct{}{}
	for k, v := range valves {
		if v.pressure == 0 {
			continue
		}
		m[k] = struct{}{}
	}
	return m
}

func findPaths(valves map[string]valve, interesting map[string]struct{}) map[string]map[string]int {
	result := make(map[string]map[string]int)
	for t := range interesting {
		paths := searchPath(valves, t, 0, t, map[string]struct{}{}, interesting)
		result[t] = make(map[string]int)
		for _, path := range paths {
			best, ok := result[path.start][path.end]
			if !ok || best > path.cost {
				result[path.start][path.end] = path.cost
			}
		}
	}
	return result
}

func searchPath(valves map[string]valve, start string, cost int, current string, visited map[string]struct{}, interesting map[string]struct{}) []path {
	res := []path{}
	newVisited := newVisited(visited, current)
	currentValve := valves[current]

	for _, t := range currentValve.tunnels {
		if _, ok := visited[t]; ok {
			continue
		}
		if _, ok := interesting[t]; ok {
			res = append(res, path{start: start, end: t, cost: cost + 1})
		}
		res = append(res, searchPath(valves, start, cost+1, t, newVisited, interesting)...)
	}
	return res

}

func newVisited(old map[string]struct{}, visited string) map[string]struct{} {
	m := make(map[string]struct{}, len(old))
	for k, v := range old {
		m[k] = v
	}
	m[visited] = struct{}{}
	return m
}

type visited int

func (v visited) isSeparate(v2 visited) bool {
	return (v & v2) == 0
}

func (v visited) isVisited(s string, stateMap map[string]int) bool {
	return 1<<stateMap[s]&v != 0
}

func (v visited) add(s string, stateMap map[string]int) visited {
	return 1<<stateMap[s] | v
}

type answer map[visited]int

func (a answer) get(s visited) (int, bool) {
	v, ok := a[s]
	return v, ok
}

func (a answer) put(s visited, v int) {
	a[s] = v
}

type game struct {
	costs    map[string]map[string]int
	vals     map[string]int
	stateMap map[string]int
}

func (g *game) g(goTo string, budget int, vis visited, answers answer, acc int) answer {
	ans, ok := answers.get(vis)
	if !ok {
		ans = 0
	}
	if ans < acc {
		ans = acc
	}
	answers.put(vis, ans)
	for next, cost := range g.costs[goTo] {
		newBudget := budget - cost - 1
		if newBudget < 0 {
			continue
		}
		if vis.isVisited(next, g.stateMap) {
			continue
		}
		answers = g.g(next, newBudget, vis.add(next, g.stateMap), answers, newBudget*g.vals[next]+acc)
	}
	return answers
}
