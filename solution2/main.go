package main

import (
	"container/heap"
	"fmt"
	"math"
	"math/rand"
	"solution2/anneal"
	"solution2/loader"
	"solution2/pqueue"
	"solution2/types"
	"sort"
)

type State struct {
	TrainAssignment map[string][]string
	TrainPickedUp   map[string][]string
	Move            []Move
	Route           map[string][]string
	Graph           types.Graph
	Train           map[string]*types.Train
	Package         map[string]*types.Package
}

type Move struct {
	TimeTaken      int
	Train          string
	StartNode      string
	EndNode        string
	PickedPackage  []string
	DroppedPackage []string
}

func main() {
	train, pkg, graph := loader.Initialize("example.txt")
	t := assignPkgToTrain(graph, train, pkg)
	r, m := planRoute(graph, t, train, pkg)

	initialState := State{TrainAssignment: t, Route: r, Move: m, Graph: graph, Train: train, Package: pkg}

	s := anneal.Init(initialState, anneal.Config{Iteration: 10000, Temperature: 25000, AneallingFactor: 0.99})
	s.PrintMovement()
}

func (s State) Energy() float64 {
	var timeTaken float64

	for trainName := range s.Train {
		route := s.Route[trainName]
		for i := 0; i < len(route)-1; i++ {
			timeTaken += float64(s.Graph[route[i]][route[i+1]])
		}
	}
	return timeTaken
}

func (s State) PrintMovement() {
	for _, each := range s.Move {
		fmt.Printf("W=%d, T=%s, N1=%s, P1=%v, N2=%s P2=%v\n", each.TimeTaken, each.Train, each.StartNode, each.PickedPackage, each.EndNode, each.DroppedPackage)
	}
	fmt.Printf("// Takes %d mintues total.", int(s.Energy()))
}

func (s State) Neighbor() anneal.State {
	newState := s
	// Generate 2 random train
	train1 := s.getRandomTrain()
	train2 := s.getRandomTrain()
	if rand.Float64() > 0.5 && train1 != train2 && len(newState.TrainAssignment[train1]) > 0 {
		// Swap train1's package assignment to train 2
		i := rand.Intn(len(newState.TrainAssignment[train1]))
		pkgToReassign := newState.TrainAssignment[train1][i]
		if newState.Train[train2].CurrentCapacity >= newState.Package[pkgToReassign].Weight {
			// Remove from the train1
			newState.TrainAssignment[train1] = append(newState.TrainAssignment[train1][:i], newState.TrainAssignment[train1][i+1:]...)
			newState.Train[train1].CurrentCapacity += newState.Package[pkgToReassign].Weight
			// Assign to train2
			newState.TrainAssignment[train2] = append(newState.TrainAssignment[train2], pkgToReassign)
			newState.Train[train2].CurrentCapacity -= newState.Package[pkgToReassign].Weight

			// reset the package to not picked up
			newState.reset()

			route, move := planRoute(newState.Graph, newState.TrainAssignment, newState.Train, newState.Package)
			newState.Route = route
			newState.Move = move
		}

	} else {
		// Swap package within a train and regenerate new route for the train
		if len(newState.TrainAssignment[train1]) > 0 {
			i := rand.Intn(len(newState.TrainAssignment[train1]))
			j := rand.Intn(len(newState.TrainAssignment[train1]))

			newState.TrainAssignment[train1][i], newState.TrainAssignment[train1][j] = newState.TrainAssignment[train1][j], newState.TrainAssignment[train1][i]

			// Reset
			newState.reset()

			route, move := planRoute(newState.Graph, newState.TrainAssignment, newState.Train, newState.Package)
			newState.Route = route
			newState.Move = move
		}
	}
	return newState
}

// Get random train
func (s State) getRandomTrain() string {
	trainName := make([]string, 0, len(s.Train))
	for n := range s.Train {
		trainName = append(trainName, n)
	}
	tName := trainName[rand.Intn(len(trainName))]
	return tName
}

// Reset some state, typically every iteration
func (s State) reset() {
	for _, each := range s.Package {
		each.Picked = false
	}
	for _, each := range s.Train {
		each.CurrentLocation = each.StartAt
		each.PickedPackage = make([]string, 0)
		each.DroppedPackage = make([]string, 0)
	}
}

// Djikstra shortest distance
func shortestDistance(graph types.Graph, start, end string) (int, []string) {
	pq := make(pqueue.DistancePQ, len(graph))
	duration := make(map[string]int)
	prev := make(map[string]string)

	i := 0
	for v := range graph {
		duration[v] = math.MaxInt32
		if v == start {
			duration[v] = 0
		}
		pq[i] = &pqueue.Node{Name: v, Duration: duration[v]}
		i++
	}

	heap.Init(&pq)

	for pq.Len() > 0 {
		node := heap.Pop(&pq).(*pqueue.Node)
		for dest, dur := range graph[node.Name] {
			alt := duration[node.Name] + dur
			if alt < duration[dest] {
				duration[dest] = alt
				prev[dest] = node.Name
				for i := 0; i < pq.Len(); i++ {
					if pq[i].Name == dest {
						pq[i].Duration = alt
						heap.Fix(&pq, i)
					}
				}
			}
		}
	}

	path := make([]string, 0)
	temp := end
	for temp != "" {
		path = append([]string{temp}, path...)
		temp = prev[temp]
	}
	return duration[end], path
}

// Randomly travel the nodes using DFS
func randomGraphTravel(graph types.Graph, start, end string) []string {
	visited := make(map[string]bool)
	stack := [][]string{{start}}

	for len(stack) > 0 {
		path := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		curr := path[len(path)-1]
		if curr == end {
			return path
		}

		if !visited[curr] {
			visited[curr] = true

			neighbor := make([]string, 0)
			for each := range graph[curr] {
				neighbor = append(neighbor, each)
			}
			// Shuffle the neighbour sequence for the randomness
			if len(neighbor) > 0 {
				rand.Shuffle(len(neighbor), func(i, j int) {
					neighbor[i], neighbor[j] = neighbor[j], neighbor[i]
				})

				for _, n := range neighbor {
					if !visited[n] {
						newPath := make([]string, len(path))
						copy(newPath, path)
						newPath = append(newPath, n)
						stack = append(stack, newPath)
					}
				}
			}
		}
	}
	// No path
	return nil
}

// Get random neighbor node
func getRandomNeighborNode(graph types.Graph, currNode string) string {
	neighbor := make([]string, 0, len(graph[currNode]))
	for n := range graph[currNode] {
		neighbor = append(neighbor, n)
	}
	nName := neighbor[rand.Intn(len(neighbor))]
	return nName
}

// Package assignment
func assignPkgToTrain(graph types.Graph, train map[string]*types.Train, pkg map[string]*types.Package) map[string][]string {
	trainKey := make([]string, 0, len(train))
	// Generate key
	for _, each := range train {
		trainKey = append(trainKey, each.Name)
	}

	trainAssgn := make(map[string][]string)
	for _, each := range train {
		trainAssgn[each.Name] = []string{}
	}
	// sort the package by it's weight, make sure the heavier package is prioritize
	sortedPkg := make([]*types.Package, 0, len(pkg))
	for _, each := range pkg {
		sortedPkg = append(sortedPkg, each)
	}
	sort.Slice(sortedPkg, func(x, y int) bool {
		return sortedPkg[x].Weight > sortedPkg[y].Weight
	})

	for _, p := range sortedPkg {
		for {
			randTrainKey := trainKey[rand.Intn(len(train))]
			// Check the capacity
			t := train[randTrainKey]
			if t.CurrentCapacity >= p.Weight {
				trainAssgn[t.Name] = append(trainAssgn[t.Name], p.Name)
				t.CurrentCapacity -= p.Weight
				break
			}
		}
	}
	return trainAssgn
}

// Create route for train to deliver assigned package
func planRoute(graph types.Graph, assignment map[string][]string, train map[string]*types.Train, pkg map[string]*types.Package) (map[string][]string, []Move) {
	route := make(map[string][]string)
	nodeToPkgMap := make(map[string][]string)
	move := make([]Move, 0)
	timeTaken := 0

	for _, each := range pkg {
		if _, ok := nodeToPkgMap[each.StartAt]; !ok {
			nodeToPkgMap[each.StartAt] = []string{}
		}
		nodeToPkgMap[each.StartAt] = append(nodeToPkgMap[each.StartAt], each.Name)
	}

	for t, pkgs := range assignment {
		// pickedUp stack to use as stack of delivery job
		train[t].PickedPackage = make([]string, 0)
		train[t].DroppedPackage = make([]string, 0)

		// The pkg loop here basically generate route for picking up a pkg and drop the package one at a time
		for _, name := range pkgs {
			// Skip if package had been picked up by previous route where the train might passed through the node.
			if pkg[name].Picked {
				continue
			}

			pickUpPath := randomGraphTravel(graph, train[t].CurrentLocation, pkg[name].StartAt)

			// Pickup
			for i := 0; i < len(pickUpPath)-1; i++ {
				m := Move{
					TimeTaken:      timeTaken,
					Train:          train[t].Name,
					StartNode:      pickUpPath[i],
					EndNode:        pickUpPath[i+1],
					PickedPackage:  train[t].PickedPackage,
					DroppedPackage: train[t].DroppedPackage,
				}
				move = append(move, m)
				timeTaken += graph[pickUpPath[i]][pickUpPath[i+1]]

				// Check if the path passing thru some other package that assigned to the train, might as well pick up.
				pkgEncounterInThePath := commonStrings(pkgs, nodeToPkgMap[pickUpPath[i+1]])
				for _, each := range pkgEncounterInThePath {
					if !pkg[each].Picked {
						pkg[each].Picked = true
						train[t].PickedPackage = append(train[t].PickedPackage, each)
					}
				}
			}

			route[t] = append(route[t], pickUpPath...)
			// Update train current location to picked up package location
			train[t].CurrentLocation = pkg[name].StartAt

			// Drop off the packages
			for len(train[t].PickedPackage) > 0 {
				p := train[t].PickedPackage[len(train[t].PickedPackage)-1]
				dropOffPath := randomGraphTravel(graph, train[t].CurrentLocation, pkg[p].Destination)

				for i := 0; i < len(dropOffPath)-1; i++ {
					m := Move{
						TimeTaken:      timeTaken,
						Train:          train[t].Name,
						StartNode:      dropOffPath[i],
						EndNode:        dropOffPath[i+1],
						PickedPackage:  train[t].PickedPackage,
						DroppedPackage: train[t].DroppedPackage,
					}
					move = append(move, m)
					timeTaken += graph[dropOffPath[i]][dropOffPath[i+1]]

					// Check if the path passing thru some other package that assigned to the train, might as well pick up.
					pkgEncounterInThePath := commonStrings(pkgs, nodeToPkgMap[dropOffPath[i+1]])
					for _, each := range pkgEncounterInThePath {
						if !pkg[each].Picked {
							pkg[each].Picked = true
							train[t].PickedPackage = append(train[t].PickedPackage, each)
						}
					}

					// Check if passing thru some node which is destination of some picked up package
					for j := len(train[t].PickedPackage) - 1; j >= 0; j-- {
						if pkg[train[t].PickedPackage[j]].Destination == dropOffPath[i+1] {
							// Remove from the pickup queue, we don't have to deliver later, as we can drop off now
							train[t].CurrentCapacity += pkg[train[t].PickedPackage[j]].Weight
							train[t].DroppedPackage = append(train[t].DroppedPackage, train[t].PickedPackage[j])
							train[t].PickedPackage = append(train[t].PickedPackage[:j], train[t].PickedPackage[j+1:]...)
						}
					}
				}

				// Update train current location to picked up destination
				route[t] = append(route[t], dropOffPath...)
				train[t].CurrentLocation = pkg[p].Destination
				train[t].CurrentCapacity += pkg[p].Weight
			}
		}
	}
	return route, move
}

func commonStrings(arr1, arr2 []string) []string {
	stringMap := make(map[string]bool)
	for _, s := range arr1 {
		stringMap[s] = true
	}

	commonStrings := make([]string, 0)
	for _, s := range arr2 {
		if stringMap[s] {
			commonStrings = append(commonStrings, s)
		}
	}
	return commonStrings
}
