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
