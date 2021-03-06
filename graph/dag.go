package graph

import (
	"fmt"
)

type topoOrder struct {
	Graph Graph
	Order []Node  // Order in topographical sort
	Cycle []Node  // If not DAG, which nodes participate in a cycle
	black nodeSet // Fully processed nodes
	gray  nodeSet // Nodes currently being processed
}

// TODO(ddn): For performance, consider avoiding recursion

// Visit leaves, then nodes next to leaves, etc.
func (a *topoOrder) visit(n Node) {
	if a.gray[n] {
		if !a.black[n] {
			a.black[n] = true
			a.Cycle = append(a.Cycle, n)
		}
		return
	}
	if a.black[n] {
		return
	}
	a.gray[n] = true
	for i, num := 0, a.Graph.NumOuts(n); i < num; i += 1 {
		a.visit(a.Graph.Out(n, i))
	}
	delete(a.gray, n)
	if !a.black[n] {
		a.black[n] = true
		a.Order = append(a.Order, n)
	}
}

func (a *topoOrder) cycleError() error {
	if len(a.Cycle) == 0 {
		return nil
	} else {
		return fmt.Errorf("cycle containing %q", a.Cycle[0])
	}
}

// Run topographic sort
func topoSort(opt TopoSortOpt) *topoOrder {
	g := opt.Graph
	if opt.NodeOrder != nil {
		g = makeSortedGraph(opt.Graph, opt.NodeOrder)
	}
	to := &topoOrder{
		Graph: g,
		black: make(nodeSet),
		gray:  make(nodeSet),
	}
	for i, num := 0, g.NumNodes(); i < num; i += 1 {
		n := g.Node(i)
		if _, seen := to.black[n]; !seen {
			to.visit(n)
		}
	}
	return to
}

// Return nil if graph is acyclic. If graph contains a cycle, return error.
func IsDag(g Graph) error {
	return topoSort(TopoSortOpt{Graph: g}).cycleError()
}

type TopoSortOpt struct {
	Graph     Graph
	NodeOrder LessThan // Optional argument to ensure deterministic output
}

// Return topological sort of graph. If edge (a, b) is in g, then b < a in the
// resulting order.  Returns an error if graph contains a cycle.
func TopoSort(opt TopoSortOpt) ([]Node, error) {
	to := topoSort(opt)
	if err := to.cycleError(); err != nil {
		return nil, err
	} else {
		return to.Order, nil
	}
}

// Compute transitive reduction of a graph. Relatively expensive operation:
// O(nm).
func TransitiveReduction(graph Graph) (Graph, error) {
	// Find some node in g not in ns
	someNode := func(g Graph, ns []Node) Node {
		nodes := make(map[Node]bool)
		for _, n := range ns {
			nodes[n] = true
		}
		for i, inum := 0, g.NumNodes(); i < inum; i += 1 {
			n := g.Node(i)
			if !nodes[n] {
				return n
			}
		}
		return nil
	}

	ret := &qgraph{
		Outs: make(map[Node][]Node),
	}

	dag := Schedule(graph)
	for len(dag.Roots) > 0 {
		// In DAG, solving shortest path with -w() is the solution to the
		// longest path problem
		dist := ShortestPath(ShortestPathOpt{
			Graph:   graph,
			Sources: dag.Roots,
			Weight: func(x, y Node) int {
				return -1
			},
		})

		for _, src := range dag.Roots {
			ret.Nodes = append(ret.Nodes, src)
			for i, inum := 0, graph.NumOuts(src); i < inum; i += 1 {
				dst := graph.Out(src, i)
				if dist[dst] == -1 {
					ret.Outs[src] = append(ret.Outs[src], dst)
				}
			}
		}

		var next []Node
		for _, src := range dag.Roots {
			next = append(next, dag.Visit(src)...)
		}
		dag.Roots = next
	}

	if len(ret.Nodes) < graph.NumNodes() {
		// TODO(ddn): transitive reductions exist for cyclic graphs but we just
		// can't use SSSP to find them
		return nil, fmt.Errorf("not yet implemented: cycle containing %q", someNode(graph, ret.Nodes))
	}

	return ret, nil
}
