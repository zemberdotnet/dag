package dag

import (
	"crypto/md5"
	"errors"
	"fmt"
)

type DAG struct {
	vertices   []*Vertex
	edgesById  map[string]*Edge
	vertexById map[string]*Vertex
}

type Vertex struct {
	id           string
	data         any
	f            func()
	dependents   []string
	dependencies []string
}

type Edge struct {
	head *Vertex
	tail *Vertex
}

func NewDag() *DAG {
	return &DAG{
		vertices:   []*Vertex{},
		vertexById: map[string]*Vertex{},
		edgesById:  map[string]*Edge{},
	}
}

func NewEdge(head *Vertex, tail *Vertex) (*Edge, error) {
	if head == nil || tail == nil {
		return nil, errors.New("vertices must not be nil")
	}

	return &Edge{
		head: head,
		tail: tail,
	}, nil
}

// hash is used to turn an edge into a hash-able value for lookup in a map.
// it depends on the uniqueness of the vertex ids to be valid.
func (e *Edge) hash() string {
	h := md5.New()
	hash := h.Sum([]byte(e.tail.id + e.head.id))
	return string(hash[:])
}

func (d *DAG) AddEdge(e *Edge) error {
	if e == nil {
		return fmt.Errorf("cannot add nil edge")
	}

	if e.head == nil {
		return fmt.Errorf("head vertex must not be nil")
	}

	if e.tail == nil {
		return fmt.Errorf("tail vertex must not be nil")
	}

	_, ok := d.vertexById[e.head.id]
	if !ok {
		return fmt.Errorf("vertex %v does not exist in graph", e.head.id)
	}

	_, ok = d.vertexById[e.tail.id]
	if !ok {
		return fmt.Errorf("vertex %v does not exist in graph", e.tail.id)
	}

	h := e.hash()
	_, ok = d.edgesById[h]
	if ok {
		return fmt.Errorf("edge from %v to %v already exists in graph", e.head.id, e.tail.id)
	}

	d.edgesById[h] = e
	e.head.dependents = append(e.head.dependents, e.tail.id)
	e.tail.dependencies = append(e.tail.dependencies, e.head.id)

	return nil
}

func (d *DAG) AddEdgeByVerticesId(head string, tail string) (*Edge, error) {
	return d.addEdge(head, tail)
}

func (d *DAG) AddEdgeByVertices(head *Vertex, tail *Vertex) (*Edge, error) {
	return d.addEdge(head.id, tail.id)
}

func (d *DAG) addEdge(head string, tail string) (*Edge, error) {
	tailVertex, ok := d.vertexById[tail]
	if !ok {
		return nil, fmt.Errorf("vertex %v does not exist in graph", tail)
	}

	headVertex, ok := d.vertexById[head]
	if !ok {
		return nil, fmt.Errorf("vertex %v does not exist in graph", head)
	}

	e := &Edge{
		head: headVertex,
		tail: tailVertex,
	}

	h := e.hash()
	_, ok = d.edgesById[h]
	if ok {
		return nil, fmt.Errorf("edge from %v to %v already exists in graph", head, tail)
	} else {
		d.edgesById[h] = e
	}

	tailVertex.dependencies = append(tailVertex.dependencies, headVertex.id)
	headVertex.dependents = append(headVertex.dependents, tailVertex.id)

	return e, nil
}

// Vertex-Related Functions

// NewVertex creates a new vertex identified by id (which should be unique)
// and holds arbitrary data.
func NewVertex(id string, data any, f func()) *Vertex {
	return &Vertex{
		id:           id,
		data:         data,
		f:            f,
		dependents:   []string{},
		dependencies: []string{},
	}
}

// AddVertex adds a vertex to the DAG. AddVertex will return an non-nil error
// when adding a vertex whose id already exists in the DAG.
func (d *DAG) AddVertex(v *Vertex) error {
	if v == nil {
		return errors.New("vertex must not be nil")
	}

	_, ok := d.vertexById[v.id]
	if ok {
		return errors.New("vertex with the same id already exists in graph")
	}
	d.vertices = append(d.vertices, v)
	d.vertexById[v.id] = v
	return nil
}

// TopSort performs a topological sort of the vertices in the DAG using Khan's algorithm.
// It returns a slice containing the vertices in topologically sorted order.
func (dag *DAG) TopSort() ([]*Vertex, error) {
	// Initialize the result slice and the vertex indegree counts.
	result := []*Vertex{}
	indegree := make(map[*Vertex]int)

	// Count incoming edges
	for _, vertex := range dag.vertices {
		indegree[vertex] = len(vertex.dependencies)
	}

	// Initialize the queue with the vertices that have an indegree of 0.
	queue := []*Vertex{}
	for vertex, count := range indegree {
		if count == 0 {
			queue = append(queue, vertex)
		}
	}

	// Process the vertices in the queue.
	for len(queue) > 0 {
		vertex := queue[0]
		go queue[0].f()
		// go queue[0]()
		queue = queue[1:]
		result = append(result, vertex)
		for _, dependent := range vertex.dependents {
			indegree[dag.vertexById[dependent]]--
			if indegree[dag.vertexById[dependent]] == 0 {
				queue = append(queue, dag.vertexById[dependent])
			}
		}
	}

	// Check if there are any vertices that were not processed, which means the DAG has a cycle.
	if len(result) < len(dag.vertices) {
		return nil, fmt.Errorf("dag is cyclic and invalid")
	}

	return result, nil
}
