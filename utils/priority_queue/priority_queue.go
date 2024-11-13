package priority_queue

import (
	"cmp"
	"container/heap"
	"fmt"
	"github.com/BooleanCat/go-functional/v2/it"
	"slices"
)

// From: https://pkg.go.dev/container/heap@go1.23.2#pkg-overview, but genericized
type Item[P cmp.Ordered, V any] struct {
	value    V // The value of the item; arbitrary.
	priority P // The priority of the item in the queue.
	// The index is needed by update and is maintained by the heap.Interface methods.
	index int // The index of the item in the heap.
}

type PriorityQueueBase[P cmp.Ordered, V any] []*Item[P, V]

type PriorityFunc[P cmp.Ordered, V any] func(value V) P

type PriorityQueue[P cmp.Ordered, V any] struct {
	base       PriorityQueueBase[P, V]
	priorityFn PriorityFunc[P, V]
}

func NewQueue[P cmp.Ordered, V any](priorityFn PriorityFunc[P, V], items ...V) PriorityQueue[P, V] {
	ret := PriorityQueue[P, V]{base: make([]*Item[P, V], len(items)), priorityFn: priorityFn}
	i := 0
	for _, v := range items {
		ret.base[i] = &Item[P, V]{
			value:    v,
			priority: priorityFn(v),
			index:    i,
		}
		i++
	}
	heap.Init(&ret)
	return ret
}

func (pq PriorityQueue[P, V]) String() string {
	return fmt.Sprint(slices.Collect(it.Map(slices.Values(pq.base), func(item *Item[P, V]) string {
		return fmt.Sprintf("%d: %s", item.priority, item.value)
	})))
}

func (pq PriorityQueue[P, V]) Len() int { return len(pq.base) }

func (pq PriorityQueue[P, V]) Less(i, j int) bool {
	// We want Pop to give us the highest, not lowest, priority so we use greater than here.
	return pq.base[i].priority < pq.base[j].priority
}

func (pq PriorityQueue[P, V]) Swap(i, j int) {
	pq.base[i], pq.base[j] = pq.base[j], pq.base[i]
	pq.base[i].index = i
	pq.base[j].index = j
}

// Don't call directly use Append instead
func (pq *PriorityQueue[P, V]) Push(x any) {
	n := len((*pq).base)
	item := x.(*Item[P, V])
	item.index = n
	(*pq).base = append((*pq).base, item)
}

func (pq *PriorityQueue[P, V]) Append(v V) {
	heap.Push(pq, &Item[P, V]{
		value:    v,
		priority: pq.priorityFn(v),
	})
}

// Don't call directly use Poll instead
func (pq *PriorityQueue[P, V]) Pop() any {
	old := (*pq).base
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // don't stop the GC from reclaiming the item eventually
	item.index = -1 // for safety
	(*pq).base = old[0 : n-1]
	return item
}

func (pq *PriorityQueue[P, V]) Poll() V {
	return heap.Pop(pq).(*Item[P, V]).value
}

func (pq PriorityQueue[P, V]) Peek() V {
	item := pq.base[0]
	return item.value
}

// update modifies the priority and value of an Item in the queue.
func (pq *PriorityQueue[P, V]) update(item *Item[P, V], value V) {
	item.value = value
	item.priority = pq.priorityFn(value)
	heap.Fix(pq, item.index)
}
