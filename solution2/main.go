package main

import (
	"container/heap"
	"fmt"
	"math"
	"math/rand"
	"solution2/loader"
	"solution2/pqueue"
	"solution2/types"
	"sort"

	"github.com/ccssmnn/hego"
)

type State struct {
	TrainAssignment map[string][]string
	TrainPickedUp   map[string][]string
	Route           map[string][]string
	Graph           types.Graph
	Train           map[string]*types.Train
	Package         map[string]*types.Package
}

type Move struct {
	TimeTaken      float64
	Train          string
	StartNode      string
	EndNode        string
	PickedPackage  []string
	DroppedPackage []string
}

func main() {
	train, pkg, graph := loader.Initialize("example.txt")
	t := assignPkgToTrain(graph, train, pkg)
	r := planRoute(graph, t, train, pkg)

	initialState := State{TrainAssignment: t, Route: r, Graph: graph, Train: train, Package: pkg}

	settings := hego.SASettings{}
	settings.MaxIterations = 10000
	settings.Verbose = 10000
	settings.Temperature = 10000
	settings.AnnealingFactor = 0.99

	result, err := hego.SA(initialState, settings)
	if err != nil {
		panic(fmt.Sprintf("Error while running Anneal: %v", err))
	}

	finalEnergy := result.Energy
	s := result.State.(State)

	for trainName := range s.Train {
		route := s.Route[trainName]
		fmt.Println("Train assignment:", s.TrainAssignment[trainName])
		fmt.Printf("Train %s route: %v\n\n", trainName, route)
		//for i := 0; i < len(route)-1; i++ {
		//	move := Move{
		//		TimeTaken:      timeTaken,
		//		Train:          train.Name,
		//		StartNode:      route[i],
		//		EndNode:        route[i+1],
		//		PickedPackage:  train.PickedPackage[i],
		//		DroppedPackage: train.DroppedPackage[i+1],
		//	}
		//	timeTaken += float64(s.Graph[route[i]][route[i+1]])

		//	fmt.Printf("W=%d, T=%s, N1=%s, P1=%v, N2=%s P2=%v\n", int(move.TimeTaken), move.Train, move.StartNode, move.PickedPackage, move.EndNode, move.DroppedPackage)
		//}
	}
	fmt.Printf("Complete simulation in %v ! Minimum time taken to delivery all package: %v", result.Runtime, finalEnergy)
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

func (s State) Neighbor() hego.AnnealingState {
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

			route := planRoute(newState.Graph, newState.TrainAssignment, newState.Train, newState.Package)
			newState.Route = route
		}

	} else {
		if len(newState.TrainAssignment[train1]) > 0 {
			i := rand.Intn(len(newState.TrainAssignment[train1]))
			j := rand.Intn(len(newState.TrainAssignment[train1]))

			newState.TrainAssignment[train1][i], newState.TrainAssignment[train1][j] = newState.TrainAssignment[train1][j], newState.TrainAssignment[train1][i]

			// Reset
			newState.reset()

			route := planRoute(newState.Graph, newState.TrainAssignment, newState.Train, newState.Package)
			newState.Route = route
		}
	}
	return newState
}

func (s State) getRandomTrain() string {
	trainName := make([]string, 0, len(s.Train))
	for n := range s.Train {
		trainName = append(trainName, n)
	}
	tName := trainName[rand.Intn(len(trainName))]
	return tName
}

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

func (s State) reset() {
	for _, each := range s.Package {
		each.Picked = false
	}
	for _, each := range s.Train {
		each.CurrentLocation = each.StartAt
		each.PickedPackage = make(map[int][]string)
		each.DroppedPackage = make(map[int][]string)
	}
}

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

func planRoute(graph types.Graph, assignment map[string][]string, train map[string]*types.Train, pkg map[string]*types.Package) map[string][]string {
	route := make(map[string][]string)
	nodeToPkgMap := make(map[string][]string)

	for _, each := range pkg {
		if _, ok := nodeToPkgMap[each.StartAt]; !ok {
			nodeToPkgMap[each.StartAt] = []string{}
		}
		nodeToPkgMap[each.StartAt] = append(nodeToPkgMap[each.StartAt], each.Name)
	}

	for t, pkgs := range assignment {
		// pickedUp stack to use as stack of delivery job
		pickedUp := make([]string, 0)
		dropped := make([]string, 0)
		for _, name := range pkgs {
			// Skip if package had been picked up by previous route where the train might passed through the node.
			if pkg[name].Picked {
				continue
			}

			_, pickUpPath := shortestDistance(graph, train[t].CurrentLocation, pkg[name].StartAt)
			// Remove deplicate start point where it is previous route's destination
			if len(route[t]) > 0 {
				pickUpPath = pickUpPath[1:]
			}

			// Pickup
			for i := 0; i < len(pickUpPath); i++ {
				train[t].PickedPackage[len(route[t])+i] = pickedUp
				train[t].DroppedPackage[len(route[t])+i] = dropped

				// Check if the path passing thru some other package that assigned to the train, might as well pick up.
				pkgEncounterInThePath := commonStrings(pkgs, nodeToPkgMap[pickUpPath[i]])
				for _, p := range pkgEncounterInThePath {
					pkg[p].Picked = true
					pickedUp = append(pickedUp, p)
					train[t].PickedPackage[len(route[t])+i] = pickedUp
				}
			}

			route[t] = append(route[t], pickUpPath...)
			// Update train current location to picked up package location
			train[t].CurrentLocation = pkg[name].StartAt

			// Drop off the packages
			for len(pickedUp) > 0 {
				p := pickedUp[len(pickedUp)-1]
				pickedUp = pickedUp[:len(pickedUp)-1]

				_, dropOffPath := shortestDistance(graph, train[t].CurrentLocation, pkg[p].Destination)
				// Remove deplicate start point where it is previous route's destination
				if len(route[t]) > 0 {
					dropOffPath = dropOffPath[1:]
				}

				for i := 0; i < len(dropOffPath); i++ {
					train[t].PickedPackage[len(route[t])+i] = pickedUp
					train[t].DroppedPackage[len(route[t])+i] = dropped
					// Check if the path passing thru some other package that assigned to the train, might as well pick up.
					pkgEncounterInThePath := commonStrings(pkgs, nodeToPkgMap[dropOffPath[i]])
					for _, each := range pkgEncounterInThePath {
						if !pkg[each].Picked {
							pkg[each].Picked = true
							pickedUp = append(pickedUp, each)
							train[t].PickedPackage[len(route[t])+i] = pickedUp
						}
					}

					// Check if passing thru some node which is destination of some picked up package
					for j := len(pickedUp) - 1; j >= 0; j-- {
						if pkg[pickedUp[j]].Destination == dropOffPath[i] {
							// Remove from the pickup queue, we don't have to deliver later, as we can drop off now
							train[t].CurrentCapacity += pkg[pickedUp[j]].Weight
							dropped = append(dropped, pickedUp[j])
							pickedUp = append(pickedUp[:j], pickedUp[j+1:]...)
							train[t].PickedPackage[len(route[t])+i] = pickedUp
							train[t].DroppedPackage[len(route[t])+i] = dropped
						}
					}
				}

				// Drop off the main package that this route is intended
				dropped = append(dropped, p)
				train[t].DroppedPackage[len(route[t])+len(dropOffPath)-1] = dropped
				route[t] = append(route[t], dropOffPath...)

				// Update train current location to picked up destination
				train[t].CurrentLocation = pkg[p].Destination
				train[t].CurrentCapacity += pkg[p].Weight
			}
		}
	}
	return route
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
