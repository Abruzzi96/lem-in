package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

type AntV2 struct {
	id       int
	location string
	path     []string
	step     int
}

type PathV2 struct {
	nodes []string
}

type GraphV2 struct {
	connections map[string][]string
	paths       []PathV2
	ants        []AntV2
	start       string
	end         string
	traffic     map[string]bool
}

func createGraphV2() *GraphV2 {
	return &GraphV2{
		connections: make(map[string][]string),
		traffic:     make(map[string]bool),
	}
}

func (g *GraphV2) addEdgeV2(from, to string) {
	g.connections[from] = append(g.connections[from], to)
	g.connections[to] = append(g.connections[to], from)
}

func (g *GraphV2) parseFileV2(filename string) (int, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return 0, err
	}

	lines := strings.Split(string(data), "\n")
	antCount, err := strconv.Atoi(lines[0])
	if err != nil || antCount <= 0 {
		return 0, fmt.Errorf("invalid number of ants")
	}

	for i := 1; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])
		if line == "##start" {
			g.start = strings.Fields(lines[i+1])[0]
		} else if line == "##end" {
			g.end = strings.Fields(lines[i+1])[0]
		} else if strings.Contains(line, "-") {
			nodes := strings.Split(line, "-")
			if len(nodes) == 2 {
				g.addEdgeV2(strings.TrimSpace(nodes[0]), strings.TrimSpace(nodes[1]))
			}
		}
	}

	if g.start == "" || g.end == "" {
		return 0, fmt.Errorf("start or end node missing")
	}

	return antCount, nil
}

func (g *GraphV2) findPathsV2(current, destination string, visited map[string]bool, path []string) {
	visited[current] = true
	path = append(path, current)

	if current == destination {
		g.paths = append(g.paths, PathV2{nodes: append([]string(nil), path...)})
	} else {
		for _, neighbor := range g.connections[current] {
			if !visited[neighbor] {
				g.findPathsV2(neighbor, destination, visited, path)
			}
		}
	}

	visited[current] = false
}

func (g *GraphV2) distinctPathsV2() []PathV2 {
	var uniquePaths []PathV2

	for _, path := range g.paths {
		isValid := true
		for _, uPath := range uniquePaths {
			if hasOverlapV2(uPath.nodes, path.nodes) {
				isValid = false
				break
			}
		}
		if isValid {
			uniquePaths = append(uniquePaths, path)
		}
	}

	return uniquePaths
}
// cakismayan ve temiz bir yol bulmaya yonelik
// filter-validation sistemi

func hasOverlapV2(path1, path2 []string) bool {
	nodeSet := make(map[string]struct{})
	for _, node := range path1[1 : len(path1)-1] {
		nodeSet[node] = struct{}{}
	}
	for _, node := range path2[1 : len(path2)-1] {
		if _, exists := nodeSet[node]; exists {
			return true
		}
	}
	return false
}
// bu func dahil
func (g *GraphV2) initializeAntsV2(count int) {
	for i := 1; i <= count; i++ {
		g.ants = append(g.ants, AntV2{id: i, location: g.start})
	}
}


func (g *GraphV2) printTunnels() {
	fmt.Println("Tunnels:")
	seen := make(map[string]bool)
	for node, connections := range g.connections {
		for _, neighbor := range connections {
			if node < neighbor {
				tunnel := node + "-" + neighbor
				if !seen[tunnel] {
					fmt.Println(tunnel)
					seen[tunnel] = true
				}
			} else {
				tunnel := neighbor + "-" + node
				if !seen[tunnel] {
					fmt.Println(tunnel)
					seen[tunnel] = true
				}
			}
		}
	}
}

