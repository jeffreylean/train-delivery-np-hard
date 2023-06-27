package types

type Train struct {
	Capacity        int
	CurrentLocation string
	Name            string
	CurrentCapacity int
	PickedPackage   map[int][]string
	DroppedPackage  map[int][]string
	StartAt         string
}

type Graph map[string]map[string]int

type Package struct {
	Weight      int
	StartAt     string
	Destination string
	Name        string
	Picked      bool
}
