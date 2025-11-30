package util

import (
	"fmt"
)

type PriorityQueue[T any] struct {
	data []T
	less func(a, b T) bool
}

func NewPriorityQueue[T any](data []T, less func(a, b T) bool) PriorityQueue[T] {
	pq := PriorityQueue[T]{
		data: data,
		less: less,
	}
	pq.heapify()

	return pq
}

// Push inserts x into the queue.
func (q *PriorityQueue[T]) Push(x T) {
	q.data = append(q.data, x)
	q.siftUp(len(q.data) - 1)
}

// Pop removes and returns the max element.
func (q *PriorityQueue[T]) Pop() (T, bool) {
	if len(q.data) == 0 {
		var zero T
		return zero, false
	}

	max := q.data[0]
	last := q.data[len(q.data)-1]
	q.data = q.data[:len(q.data)-1]

	if len(q.data) > 0 {
		q.data[0] = last
		q.siftDown(0)
	}

	return max, true
}

func (q *PriorityQueue[T]) Print() {
	if len(q.data) == 0 {
		return
	}

	lines := [][]string{
		{fmt.Sprint(q.data[0])},
	}

	current := 0
	level := 1
	width := 1
	for {
		left := 2*current + 1
		if left < len(q.data) {
			var line []string
			current = left
			for i := left; i < left+pow(2, level) && i < len(q.data); i++ {
				line = append(line, fmt.Sprint(q.data[i]))
			}

			lines = append(lines, line)
			width = pow(2, level)
			fmt.Println(level, width)
			level++
		} else {
			break
		}
	}

	fmt.Println(width)
	for _, line := range lines {
		fmt.Println(line)
	}
}

func pow(base, exponent int) int {
	total := 1
	for i := 0; i < exponent; i++ {
		total *= base
	}

	return total
}

// sift down from the last parent to setup initial heap order
func (q *PriorityQueue[T]) heapify() {
	for i := len(q.data)/2 - 1; i >= 0; i-- {
		q.siftDown(i)
	}
}

// siftDown moves the current item down until heap order is restored.
func (q *PriorityQueue[T]) siftDown(current int) {
	for {
		left := 2*current + 1
		right := left + 1
		smallest := current

		if left < len(q.data) && q.less(q.data[left], q.data[smallest]) {
			smallest = left
		}
		if right < len(q.data) && q.less(q.data[right], q.data[smallest]) {
			smallest = right
		}
		if smallest == current {
			break
		}

		q.data[current], q.data[smallest] = q.data[smallest], q.data[current]
		current = smallest
	}
}

// siftUp moves the current item up until heap order is restored.
func (q *PriorityQueue[T]) siftUp(current int) {
	for current > 0 {
		parent := (current - 1) / 2
		if !q.less(q.data[current], q.data[parent]) {
			break
		}
		q.data[current], q.data[parent] = q.data[parent], q.data[current]
		current = parent
	}
}
