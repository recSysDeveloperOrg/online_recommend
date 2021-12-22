package cache

type HeapComparator interface {
	Compare(another interface{}) int
}

type HeapNode struct {
	Key   HeapComparator
	Value interface{}
}

type Heap struct {
	nodes    []*HeapNode
	nCurrent int
}

func NewHeap(heapNodes []*HeapNode) *Heap {
	heap := &Heap{
		nodes:    heapNodes,
		nCurrent: len(heapNodes),
	}
	heap.buildHeap()

	return heap
}

func left(i int) int {
	return i<<1 + 1
}

func right(i int) int {
	return i<<1 + 2
}

func (h *Heap) pushDown(index int) {
	oldNode := h.nodes[index]
	for index < (h.nCurrent >> 1) {
		childIndex, childNode := left(index), h.nodes[left(index)]
		if right(index) < h.nCurrent && h.nodes[right(index)].Key.Compare(h.nodes[left(index)].Key) < 0 {
			childNode = h.nodes[right(index)]
			childIndex = right(index)
		}
		if oldNode.Key.Compare(childNode.Key) <= 0 {
			break
		}
		h.nodes[index] = childNode
		index = childIndex
	}
	h.nodes[index] = oldNode
}

func (h *Heap) buildHeap() {
	for i := (h.nCurrent >> 1) - 1; i >= 0; i-- {
		h.pushDown(i)
	}
}

func (h *Heap) ReplaceTop(k HeapComparator, v interface{}) interface{} {
	oldTop := h.nodes[0]
	h.nodes[0] = &HeapNode{
		Key:   k,
		Value: v,
	}
	h.pushDown(0)

	return oldTop.Value
}

func (h *Heap) TopValue() interface{} {
	return h.nodes[0].Value
}

func (h *Heap) TopKey() HeapComparator {
	return h.nodes[0].Key
}

func (h *Heap) PopNode() *HeapNode {
	top := h.nodes[0]
	h.nodes[0] = h.nodes[h.nCurrent-1]
	h.nCurrent--
	h.pushDown(0)

	return top
}

func (h *Heap) PopValue() interface{} {
	return h.PopNode().Value
}

func (h *Heap) PopKey() HeapComparator {
	return h.PopNode().Key
}

func (h *Heap) PopValues() []interface{} {
	res, i := make([]interface{}, len(h.nodes)), 0
	for h.nCurrent > 0 {
		res[i] = h.PopValue()
		i++
	}

	return res
}
