package main

import (
	"fmt"
	"sync"
)

var mu sync.Mutex

func graphColoringUtilParallel(graph [][]bool, m int, color []int, v int, done chan bool) {
	if v == len(graph) {
		done <- true
		return
	}
	for c := 1; c <= m; c++ {
		mu.Lock()
		if isSafe(graph, color, v, c) {
			color[v] = c
			mu.Unlock()

			nextDone := make(chan bool)
			go graphColoringUtilParallel(graph, m, color, v+1, nextDone)

			if <-nextDone {
				done <- true
				return
			}

			mu.Lock()
			color[v] = 0
		}
		mu.Unlock()
	}

	done <- false
}

func graphColoringParallel(graph [][]bool, m int) ([]int, bool) {
	color := make([]int, len(graph))
	for i := range color {
		color[i] = 0
	}

	done := make(chan bool)
	go graphColoringUtilParallel(graph, m, color, 0, done)

	if <-done {
		fmt.Println("Solution Exists: Following are the assigned colors")
		for i := 0; i < len(graph); i++ {
			fmt.Printf(" %d ", color[i])
		}
		fmt.Println()
		return color, true
	} else {
		fmt.Println("Solution does not exist")
		return color, false
	}
}

// Function to check if it's safe to color the vertex with the given color
func isSafe(graph [][]bool, color []int, v int, c int) bool {
	for i := 0; i < len(graph); i++ {
		if graph[v][i] && c == color[i] {
			return false
		}
	}
	return true
}

// Utility function to solve the graph coloring problem
func graphColoringUtil(graph [][]bool, m int, color []int, v int) bool {
	if v == len(graph) {
		return true
	}

	for c := 1; c <= m; c++ {
		if isSafe(graph, color, v, c) {
			color[v] = c
			if graphColoringUtil(graph, m, color, v+1) {
				return true
			}
			color[v] = 0
		}
	}

	return false
}

// Function to solve the m Coloring problem
func graphColoring(graph [][]bool, m int) ([]int, bool) {
	color := make([]int, len(graph))
	for i := range color {
		color[i] = 0
	}

	if !graphColoringUtil(graph, m, color, 0) {
		fmt.Println("Solution does not exist")
		return color, false
	}

	fmt.Println("Solution Exists: Following are the assigned colors")
	for i := 0; i < len(graph); i++ {
		fmt.Printf(" %d ", color[i])
	}
	fmt.Println()
	return color, true
}

func PrintGraph(graph [][]bool) {
	fmt.Println("Graph Representation:")
	fmt.Print("    ")
	for i := 0; i < len(graph); i++ {
		fmt.Printf("V%d ", i)
	}
	fmt.Println("\n    ----------------")

	for i := 0; i < len(graph); i++ {
		fmt.Printf("V%d | ", i)
		for j := 0; j < len(graph[i]); j++ {
			if graph[i][j] {
				fmt.Print(" * ") // A star (*) represents a connection
			} else {
				fmt.Print(" - ") // A dash (-) represents no connection
			}
		}
		fmt.Println()
	}
}

func PrintColoredGraph(graph [][]bool, color []int) {
	fmt.Println("Colored Graph:")
	for i := 0; i < len(graph); i++ {
		fmt.Printf("Vertex %d is colored with color %d\n", i, color[i])
	}
}

func main() {
	// Representation of graph as an adjacency matrix
	graph := [][]bool{
		{false, true, true, true},
		{true, false, true, false},
		{true, true, false, true},
		{true, false, true, false},
	}
	m := 3 // Number of colors
	PrintGraph(graph)
	//color, _ := graphColoring(graph, m)
	color, success := graphColoringParallel(graph, m)
	if success {
		PrintColoredGraph(graph, color)
	}
}
