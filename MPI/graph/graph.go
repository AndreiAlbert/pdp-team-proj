package graph

type Graph struct {
	cntNodes      int
	adjacencyList []map[int]struct{}
}

func NewGraph(cntNodes int) *Graph {
	g := &Graph{
		cntNodes:      cntNodes,
		adjacencyList: make([]map[int]struct{}, cntNodes),
	}
	for i := range g.adjacencyList {
		g.adjacencyList[i] = make(map[int]struct{})
	}
	return g
}

func (g *Graph) AddEdge(fromVertex, toVertex int) {
	g.adjacencyList[fromVertex][toVertex] = struct{}{}
	g.adjacencyList[toVertex][fromVertex] = struct{}{}
}

func (g *Graph) IsEdge(fromVertex, toVertex int) bool {
	_, exists := g.adjacencyList[fromVertex][toVertex]
	return exists
}

func (g *Graph) GetCntNodes() int {
	return g.cntNodes
}
