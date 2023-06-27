package pqueue

type Node struct {
	Name     string
	Duration int
	index    int
}

// Min heap
type DistancePQ []*Node

func (pq DistancePQ) Len() int {
	return len(pq)
}

func (pq DistancePQ) Less(i, j int) bool {
	return pq[i].Duration < pq[j].Duration
}

func (pq DistancePQ) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *DistancePQ) Pop() any {
	old := *pq
	item := old[len(old)-1]
	old[len(old)-1] = nil
	item.index = -1
	*pq = old[:len(old)-1]

	return item
}

func (pq *DistancePQ) Push(node any) {
	i := node.(*Node)
	i.index = len(*pq)
	*pq = append(*pq, i)
}

// Max heap for assignment
type Assignment struct {
	Train               string
	Pkg                 string
	Path                []string
	Distance            int
	Weight              int
	WeightDistanceRatio float32
	Action              int
	index               int
}
type AssignmentPQ []*Assignment

func (pq AssignmentPQ) Len() int {
	return len(pq)
}

func (pq AssignmentPQ) Less(i, j int) bool {
	return pq[i].WeightDistanceRatio > pq[j].WeightDistanceRatio
}

func (pq AssignmentPQ) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *AssignmentPQ) Pop() any {
	old := *pq
	item := old[len(old)-1]
	old[len(old)-1] = nil
	item.index = -1
	*pq = old[:len(old)-1]

	return item
}

func (pq *AssignmentPQ) Push(node any) {
	i := node.(*Assignment)
	i.index = len(*pq)
	*pq = append(*pq, i)
}
