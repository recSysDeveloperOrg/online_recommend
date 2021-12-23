package cache

import "testing"

type testCase struct {
	Name     string
	InitKeys []int
	MoreOps  []*operation
}

type operation struct {
	OpType uint8
	Key    int
}

const (
	opTypeGetTop     = 0
	opTypePopTop     = 1
	opTypeReplaceTop = 2
)

type intHeapComparator int

func (a intHeapComparator) Compare(b interface{}) int {
	if bCom, ok := b.(intHeapComparator); ok {
		return int(a) - int(bCom)
	}
	panic("???")
}

func Test_Heap(t *testing.T) {
	testCases := []*testCase{
		{
			Name:     "simpleCase",
			InitKeys: []int{5, 4, 3, 2, 0, 1},
			MoreOps: []*operation{
				{
					OpType: opTypeGetTop,
					Key:    0,
				},
				{
					OpType: opTypePopTop,
				},
				{
					OpType: opTypeGetTop,
					Key:    1,
				},
				{
					OpType: opTypePopTop,
				},
				{
					OpType: opTypeGetTop,
					Key:    2,
				},
				{
					OpType: opTypePopTop,
				},
				{
					OpType: opTypeGetTop,
					Key:    3,
				},
				{
					OpType: opTypePopTop,
				},
				{
					OpType: opTypeGetTop,
					Key:    4,
				},
				{
					OpType: opTypePopTop,
				},
				{
					OpType: opTypeReplaceTop,
					Key:    3,
				},
				{
					OpType: opTypeGetTop,
					Key:    3,
				},
			},
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			heapNodes := make([]*HeapNode, 0, len(testCase.InitKeys))
			for _, v := range testCase.InitKeys {
				heapNodes = append(heapNodes, &HeapNode{
					Key: intHeapComparator(v),
					Value: struct {
					}{},
				})
			}
			heap := NewHeap(heapNodes)

			for _, moreOp := range testCase.MoreOps {
				switch moreOp.OpType {
				case opTypeGetTop:
					topKey := heap.TopKey()
					if topKey.Compare(intHeapComparator(moreOp.Key)) != 0 {
						t.Fatalf("topKey:%d,expect:%d\n", topKey, moreOp.Key)
					}
				case opTypePopTop:
					heap.PopNode()
				case opTypeReplaceTop:
					heap.ReplaceTop(intHeapComparator(moreOp.Key), struct {
					}{})
				}
			}
		})
	}
}
