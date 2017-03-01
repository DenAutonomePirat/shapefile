package main

import "container/heap"

// astar is an A* pathfinding implementation.

// node is a wrapper to store A* data for a Point node.
type node struct {
	point  Point
	cost   float64
	rank   float64
	parent *node
	open   bool
	closed bool
	index  int
}

// nodeMap is a collection of nodes keyed by Point nodes for quick reference.
type nodeMap map[Point]*node

// get gets the Point object wrapped in a node, instantiating if required.
func (nm nodeMap) get(p Point) *node {
	n, ok := nm[p]
	if !ok {
		n = &node{
			point: p,
		}
		nm[p] = n
	}
	return n
}

// Path calculates a short path and the distance between the two Point nodes.
//
// If no path is found, found will be false.
func Path(from, to Point) (path []Point, distance float64, found bool) {
	nm := nodeMap{}
	nq := &priorityQueue{}
	heap.Init(nq)
	fromNode := nm.get(from)
	fromNode.open = true
	heap.Push(nq, fromNode)
	for {
		if nq.Len() == 0 {
			// There's no path, return found false.
			return
		}
		current := heap.Pop(nq).(*node)
		current.open = false
		current.closed = true

		if current == nm.get(to) {
			// Found a path to the goal.
			p := []Point{}
			curr := current
			for curr != nil {
				p = append(p, curr.point)
				curr = curr.parent
			}
			return p, current.cost, true
		}

		for _, neighbor := range current.point.PathNeighbors() {
			cost := current.cost + current.point.PathNeighborCost(neighbor)
			neighborNode := nm.get(neighbor)
			if cost < neighborNode.cost {
				if neighborNode.open {
					heap.Remove(nq, neighborNode.index)
				}
				neighborNode.open = false
				neighborNode.closed = false
			}
			if !neighborNode.open && !neighborNode.closed {
				neighborNode.cost = cost
				neighborNode.open = true
				neighborNode.rank = cost + neighbor.PathEstimatedCost(to)
				neighborNode.parent = current
				heap.Push(nq, neighborNode)
			}
		}
	}
}
