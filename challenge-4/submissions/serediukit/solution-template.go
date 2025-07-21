package main

import (
	"slices"
	"sync"
)

// ConcurrentBFSQueries concurrently processes BFS queries on the provided graph.
// - graph: adjacency list, e.g., graph[u] = []int{v1, v2, ...}
// - queries: a list of starting nodes for BFS.
// - numWorkers: how many goroutines can process BFS queries simultaneously.
//
// Return a map from the query (starting node) to the BFS order as a slice of nodes.
// YOU MUST use concurrency (goroutines + channels) to pass the performance tests.
func ConcurrentBFSQueries(graph map[int][]int, queries []int, numWorkers int) map[int][]int {
	if numWorkers == 0 {
		return map[int][]int{}
	}

	vertexCh := make(chan int, numWorkers)

	res := make(map[int][]int)

	wg := &sync.WaitGroup{}
	mu := &sync.Mutex{}

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			ConcurrentBFS(graph, vertexCh, &res, mu)
			wg.Done()
		}()
	}

	for _, q := range queries {
		vertexCh <- q
	}

	close(vertexCh)

	wg.Wait()

	return res
}

func ConcurrentBFS(graph map[int][]int, vertexCh chan int, res *map[int][]int, mu *sync.Mutex) {
	for v := range vertexCh {
		way := BFS(graph, v)

		mu.Lock()
		(*res)[v] = way
		mu.Unlock()
	}
}

func BFS(graph map[int][]int, startVertex int) []int {
	res := []int{startVertex}

	toVisit := make([]int, 0, len(graph))
	for _, v := range graph[startVertex] {
		if v != startVertex {
			toVisit = append(toVisit, v)
		}
	}

	for len(toVisit) > 0 {
		curVertex := toVisit[0]

		res = append(res, curVertex)

		for _, v := range graph[curVertex] {
			if !slices.Contains(res, v) && !slices.Contains(toVisit, v) {
				toVisit = append(toVisit, v)
			}
		}

		if len(toVisit) > 1 {
			toVisit = toVisit[1:]
		} else {
			break
		}
	}

	return res
}
