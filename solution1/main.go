package main

import (
	"container/heap"
	"fmt"
	"math"
	"solution1/pkg/loader"
	"solution1/pkg/types"
	"solution1/pqueue"
)

const (
	Pickup = iota
	DeliverToDestination
	MoveToPackage
)

type Move struct {
	TimeTaken      int
	Train          string
	StartNode      string
	EndNode        string
	PickedPackage  []string
	DroppedPackage []string
}

type Movement struct {
	Move      []Move
	TimeTaken int
}

func (m Movement) Print() {
	for _, each := range m.Move {
		fmt.Printf("W=%d, T=%s, N1=%s, P1=%v, N2=%s P2=%v\n", each.TimeTaken, each.Train, each.StartNode, each.PickedPackage, each.EndNode, each.DroppedPackage)
	}
	fmt.Printf("// Takes %d mintues total.", m.TimeTaken)
}

func main() {
	train, pkg, g := loader.Initialize("example.txt")
	movement := Movement{Move: make([]Move, 0), TimeTaken: 0}
	// Queue for the assignment. The assignment are store by chunk, each chunk contain assignment of multiple trains at a point of time.
	// All assignment are being execute sequentially.
	queue := make([][]pqueue.Assignment, 0)

	// Assign first package to each train
	assignment := assignPackage(pkg, train, g)
	a := make([]pqueue.Assignment, len(assignment))
	i := 0
	for _, each := range assignment {
		a[i] = each
	}
	queue = append(queue, a)

	for len(queue) > 0 {
		asgnt := queue[0]
		queue = queue[1:]
		for _, each := range asgnt {
			switch each.Action {
			// Pickup the package
			case Pickup:
				t := train[each.Train]
				p := pkg[each.Pkg]

				mv := pickupPkg(t, p, g, each.Path, &movement.TimeTaken)
				movement.Move = append(movement.Move, mv...)
			// Deliver the package
			case DeliverToDestination:
				t := train[each.Train]
				p := pkg[each.Pkg]
				mv := dropOffPackage(t, p, g, each.Path, &movement.TimeTaken)
				movement.Move = append(movement.Move, mv...)
			}
		}
		asn := make([]pqueue.Assignment, 0)
		for _, as := range deliveryOrPickUp(pkg, train, g) {
			if as.Action != -1 {
				asn = append(asn, as)
			}
		}
		if len(asn) > 0 {
			queue = append(queue, asn)
		}
	}
	movement.Print()
}

func shortest(graph types.Graph, start, end string) (int, []string) {
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

// Assign closest package to train
func assignPackage(pkg map[string]*types.Package, train map[string]*types.Train, graph types.Graph) map[string]pqueue.Assignment {
	assignment := make(map[string]pqueue.Assignment)
	pq := make(pqueue.AssignmentPQ, 0)
	heap.Init(&pq)
	pkgAssigned := make(map[string]bool)
	trainAssigned := make(map[string]bool)

	for _, t := range train {
		for _, p := range pkg {
			if p.Picked {
				continue
			}
			dist, path := shortest(graph, t.CurrentLocation, p.StartAt)
			if t.CurrentCapacity >= p.Weight {
				// If there's only one train, we don't need to worry about optimal assignment on weight and distance for difference train
				// If there's only one train, just go with the closest package at the time.
				// The default will be just using inverse of distance as the priority, the furthers the lower priority.
				// Use inverse here because the priority queue is a max heap.
				priority := float32(1 / (dist + 1))
				if len(train) > 1 {
					// If there's more than one train, we cannot just prioritize distance, there might be some large package in the far where only particular
					// train able to carry. If prioritze shorter package, the train will get occupied with some lighter.
					// So here a ratio is being used for heuristic method to handle this, where heavier and further package will be prioritize.
					priority = float32(p.Weight) / float32((1 + dist))
				}
				heap.Push(&pq, &pqueue.Assignment{
					Train:               t.Name,
					Pkg:                 p.Name,
					Distance:            dist,
					Weight:              p.Weight,
					WeightDistanceRatio: priority,
					Path:                path,
					Action:              Pickup,
				})
			}
		}
	}

	// We had done some local calculation of best package to assign to for each train.
	// But we want the best package to being assign to a train in global context, so max heap priority queue is used, to only assign package to the train with
	// highest ratio.
	for pq.Len() > 0 {
		popped := heap.Pop(&pq).(*pqueue.Assignment)
		if !pkgAssigned[popped.Pkg] && !trainAssigned[popped.Train] {
			pkgAssigned[popped.Pkg] = true
			trainAssigned[popped.Train] = true
			assignment[popped.Train] = *popped
		}
	}
	return assignment
}

func pickupPkg(train *types.Train, pkg *types.Package, graph types.Graph, path []string, timeTaken *int) []Move {
	movement := make([]Move, 0)
	// if package and train in the same location
	if len(path) == 1 {
		move := Move{
			TimeTaken:      *timeTaken,
			Train:          train.Name,
			StartNode:      path[0],
			EndNode:        path[0],
			PickedPackage:  []string{pkg.Name},
			DroppedPackage: make([]string, 0),
		}
		train.CurrentCapacity -= pkg.Weight
		train.PickedPackage = append(train.PickedPackage, pkg.Name)
		pkg.Picked = true
		movement = append(movement, move)
	} else {
		for i := 0; i < len(path)-1; i++ {
			move := Move{
				TimeTaken:      *timeTaken,
				Train:          train.Name,
				StartNode:      path[i],
				EndNode:        path[i+1],
				PickedPackage:  train.PickedPackage,
				DroppedPackage: train.DroppedPackage,
			}
			train.CurrentLocation = path[i+1]
			*timeTaken = *timeTaken + graph[path[i]][path[i+1]]
			if path[i+1] == pkg.StartAt {
				//Pickup
				train.CurrentCapacity -= pkg.Weight
				train.PickedPackage = append(train.PickedPackage, pkg.Name)
				pkg.Picked = true
			}
			movement = append(movement, move)
		}
	}
	return movement
}

func dropOffPackage(train *types.Train, pkg *types.Package, graph types.Graph, path []string, timeTaken *int) []Move {
	movement := make([]Move, 0)
	for i := 0; i < len(path)-1; i++ {
		move := Move{
			TimeTaken:      *timeTaken,
			Train:          train.Name,
			StartNode:      path[i],
			EndNode:        path[i+1],
			PickedPackage:  train.PickedPackage,
			DroppedPackage: train.DroppedPackage,
		}
		train.CurrentLocation = path[i+1]
		*timeTaken = *timeTaken + graph[path[i]][path[i+1]]
		if path[i+1] == pkg.Destination {
			//Drop Off package
			train.CurrentCapacity += pkg.Weight
			// Remove from pickedup slice
			for i, each := range train.PickedPackage {
				if each == pkg.Name {
					train.PickedPackage = append(train.PickedPackage[:i], train.PickedPackage[i+1:]...)
				}
			}
			train.DroppedPackage = append(train.DroppedPackage, pkg.Name)
			move.PickedPackage = train.PickedPackage
			move.DroppedPackage = train.DroppedPackage
		}
		movement = append(movement, move)
	}
	return movement
}

// Function to decide whether the next assignment should be delivering picked up package or conitnue pick up next package.
func deliveryOrPickUp(pkg map[string]*types.Package, train map[string]*types.Train, graph types.Graph) map[string]pqueue.Assignment {
	a := make(map[string]pqueue.Assignment)
	// Get assignment of next package
	asgn := assignPackage(pkg, train, graph)
	// Check if the next assignment for each train is optimal choice or not
	// Compare if the next assignment or deliver the picked up package is use  lesser time
	for _, t := range train {
		// Initial value
		minDist := math.MaxInt32
		minPkg := ""
		minAction := -1
		minPath := []string{}

		// If there's next pickup assignment for the train
		if a, ok := asgn[t.Name]; ok {
			minDist = a.Distance
			minAction = Pickup
			minPath = a.Path
			minPkg = a.Pkg
		}
		// If the train has pickedup package, and if the destination of the pickedup package is shorter than the next pickup assignment, overwrite the Assignment
		// to deliver.
		if len(t.PickedPackage) > 0 {
			for _, each := range t.PickedPackage {
				deliverDist, deliverPath := shortest(graph, t.CurrentLocation, pkg[each].Destination)
				if deliverDist < minDist {
					minAction = DeliverToDestination
					minDist = deliverDist
					minPath = deliverPath
					minPkg = pkg[each].Name
				}
			}
		}
		a[t.Name] = pqueue.Assignment{Train: t.Name, Pkg: minPkg, Path: minPath, Distance: minDist, Action: minAction}
	}
	return a
}
