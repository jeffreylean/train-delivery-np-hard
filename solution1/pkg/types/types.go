package types

type Train struct {
	Capacity        int
	CurrentLocation string
	Name            string
	CurrentCapacity int
	PickedPackage   []string
	DroppedPackage  []string
}

type Graph map[string]map[string]int

type Package struct {
	Weight      int
	StartAt     string
	Destination string
	Name        string
	Picked      bool
}
