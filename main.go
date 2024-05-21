package main

import (
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
    "time"
)

func (g *GraphV2) moveAntsV2(paths []PathV2) {
	for {
		moved := false
		var moves []string

		for i := range g.ants {
			ant := &g.ants[i]
			if ant.location == g.end {
				continue
			}

			if ant.location == g.start {
				for _, path := range paths {
					nextPos := path.nodes[1]
					if !g.traffic[nextPos] {
						ant.path = path.nodes
						ant.step = 1
						ant.location = nextPos
						g.traffic[nextPos] = true
						moved = true
						moves = append(moves, fmt.Sprintf("L%d-%s", ant.id, ant.location))
						break
					}
				}
			} else if ant.step < len(ant.path)-1 {
				nextPos := ant.path[ant.step+1]
				if !g.traffic[nextPos] {
					ant.step++
					ant.location = nextPos
					if ant.location == g.end {
						moves = append(moves, fmt.Sprintf("L%d-end", ant.id))
					} else {
						g.traffic[nextPos] = true
						moves = append(moves, fmt.Sprintf("L%d-%s", ant.id, ant.location))
					}
					moved = true
				} else {
					moves = append(moves, fmt.Sprintf("L%d-%s", ant.id, ant.location))
				}
			} else {
				moves = append(moves, fmt.Sprintf("L%d-%s", ant.id, ant.location))
			}
		}

		if len(moves) > 0 {
			fmt.Println(strings.Join(moves, " "))
		}

		g.traffic = make(map[string]bool)

		if !moved {
			break
		}
	}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run main.go <input_file>")
		os.Exit(1)
	}

	filename := os.Args[1]
    startTime := time.Now()
	graph := createGraphV2()
	antCount, err := graph.parseFileV2(filename)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	visited := make(map[string]bool)
	graph.findPathsV2(graph.start, graph.end, visited, []string{})

	sort.Slice(graph.paths, func(i, j int) bool {
		return len(graph.paths[i].nodes) < len(graph.paths[j].nodes)
	})

	validPaths := graph.distinctPathsV2()
	graph.initializeAntsV2(antCount)
	graph.moveAntsV2(validPaths)
    graph.printTunnels()
    endTime := time.Now() // End measuring time
	elapsedTime := endTime.Sub(startTime).Seconds()
    fmt.Printf("Execution time: %.6f seconds\n", elapsedTime)
}
