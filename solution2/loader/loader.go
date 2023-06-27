package loader

import (
	"bufio"
	"fmt"
	"os"
	"solution2/types"
	"strconv"
	"strings"
)

func Initialize(path string) (map[string]*types.Train, map[string]*types.Package, types.Graph) {
	// Initialize variables
	var numStations, numEdges, numDeliveries, numTrains int
	var err error
	// Open file
	file, err := os.Open(path)
	if err != nil {
		panic(fmt.Sprintln("Error opening file:", err))
	}
	scanner := bufio.NewScanner(file)
	graph := make(types.Graph)

	// Read number of stations
	scanner.Scan()
	numStations, err = strconv.Atoi(scanner.Text())
	if err != nil {
		panic(fmt.Sprintln("Error reading number of stations:", err))
	}
	// Read station names
	for i := 0; i < numStations; i++ {
		scanner.Scan()
		graph[scanner.Text()] = make(map[string]int)
	}

	// skip next line
	scanner.Scan()

	// Read number of edges
	scanner.Scan()
	numEdges, err = strconv.Atoi(scanner.Text())
	if err != nil {
		panic(fmt.Sprintln("Error reading number of edges:", err))
	}
	// Read edges
	for i := 0; i < numEdges; i++ {
		scanner.Scan()
		edgeInfo := strings.Split(scanner.Text(), ",")

		weight, err := strconv.Atoi(edgeInfo[3])
		if err != nil {
			panic(fmt.Sprintln("Error reading edge weight:", err))
		}

		graph[edgeInfo[1]][edgeInfo[2]] = weight
		graph[edgeInfo[1]][edgeInfo[1]] = 0
		graph[edgeInfo[2]][edgeInfo[1]] = weight
	}

	// skip next line
	scanner.Scan()

	// Read number of deliveries
	scanner.Scan()
	numDeliveries, err = strconv.Atoi(scanner.Text())
	if err != nil {
		panic(fmt.Sprintln("Error reading number of deliveries:", err))
	}

	// Read deliveries
	pkg := make(map[string]*types.Package)
	for i := 0; i < numDeliveries; i++ {
		scanner.Scan()
		deliveryInfo := strings.Split(scanner.Text(), ",")
		weight, err := strconv.Atoi(deliveryInfo[1])
		if err != nil {
			panic(fmt.Sprintln("Error reading package weight:", err))
		}
		pkg[deliveryInfo[0]] = &types.Package{Name: deliveryInfo[0], Weight: weight, StartAt: deliveryInfo[2], Destination: deliveryInfo[3]}
	}

	// skip next line
	scanner.Scan()

	// Read number of trains
	scanner.Scan()
	numTrains, err = strconv.Atoi(scanner.Text())
	if err != nil {
		panic(fmt.Sprintln("Error reading number of trains:", err))
	}

	// Read trains
	train := make(map[string]*types.Train)
	for i := 0; i < numTrains; i++ {
		scanner.Scan()
		trainInfo := strings.Split(scanner.Text(), ",")
		capacity, err := strconv.Atoi(trainInfo[1])
		if err != nil {
			panic(fmt.Sprintln("Error reading train capacity:", err))
		}
		train[trainInfo[0]] = &types.Train{Capacity: capacity, StartAt: trainInfo[2], CurrentLocation: trainInfo[2], Name: trainInfo[0], CurrentCapacity: capacity, PickedPackage: make(map[int][]string), DroppedPackage: make(map[int][]string)}
	}
	return train, pkg, graph
}
