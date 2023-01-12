package dag

import (
	"errors"
	"sync"
	"testing"
)

func newEmptyVertex(id string) *Vertex {
	return NewVertex(id, nil, func() {})
}

func TestAddVertex(t *testing.T) {
	data := []struct {
		input  *Vertex
		output error
	}{
		{
			input:  nil,
			output: errors.New("vertex must not be nil"),
		},
		{
			input:  newEmptyVertex("id-1"),
			output: nil,
		},
		{
			input:  newEmptyVertex("id-2"),
			output: nil,
		},
		{
			input:  newEmptyVertex("id-1"),
			output: errors.New("vertex with the same id already exists in graph"),
		},
	}

	d := NewDag()

	for i, testCase := range data {
		err := d.AddVertex(testCase.input)
		if err == nil && testCase.output == nil {
			continue
		}

		if (err == nil && testCase.output != nil) || (err != nil && testCase.output == nil) {
			t.Errorf("Failed AddVertex case %v. Expected '%v' Actual '%v'", i, testCase.output, err)
			continue
		}

		if err.Error() != testCase.output.Error() {
			t.Errorf("Failed AddVertex case %v. Expected '%v' Actual '%v'", i, testCase.output, err)
		}
	}
}

func TestNewEdgeValid(t *testing.T) {
	head := newEmptyVertex("a")
	tail := newEmptyVertex("b")
	edge, err := NewEdge(head, tail)
	if err != nil {
		t.Errorf("failed creating new valid edge")
	}

	if edge.head != head {
		t.Errorf("head is not equal to head passed into NewEdge")
	}

	if edge.tail != tail {
		t.Errorf("tail is not equal to tail passed into NewEdge")
	}
}

func TestNewEdgeInvalid(t *testing.T) {
	var head *Vertex = nil
	tail := newEmptyVertex("a")
	errExpected := errors.New("vertices must not be nil")

	edge, err := NewEdge(head, tail)
	if err == nil {
		t.Errorf("didn't error on invalid new edge")
		return
	}

	if err.Error() != errExpected.Error() {
		t.Errorf("didn't error with vertices must not be nil")
	}

	if edge != nil {
		t.Errorf("returned non-nil edge when given invalid input")
	}
}

func TestAddEdgeValid(t *testing.T) {
	// prepare test
	d := NewDag()
	v1 := newEmptyVertex("a")
	v2 := newEmptyVertex("b")

	err := d.AddVertex(v1)
	if err != nil {
		t.Errorf("failed to add valid vertex to DAG. err: %v", err)
	}

	err = d.AddVertex(v2)
	if err != nil {
		t.Errorf("failed to add valid vertex to DAG. err: %v", err)
	}

	e, err := NewEdge(v1, v2)
	if err != nil {
		t.Errorf("failed creating valid edge. err: %v", err)
	}

	err = d.AddEdge(e)
	if err != nil {
		t.Errorf("failed adding valid edge to graph. err %v", err)
	}
}

func TestAddEdgeNil(t *testing.T) {
	d := NewDag()
	var e *Edge = nil
	expected := errors.New("cannot add nil edge")

	err := d.AddEdge(e)
	if err == nil {
		t.Error("failed to identify nil head on edge")
	}

	if err.Error() != expected.Error() {
		t.Errorf("incorrect error for nil edge. actual: %v expected: %v", err, expected)
	}
}

func TestAddEdgeNilHeadTail(t *testing.T) {
	d := NewDag()

	data := []struct {
		input  *Edge
		output error
	}{
		{
			input: &Edge{
				head: nil,
				tail: nil,
			},
			output: errors.New("head vertex must not be nil"),
		},
		{
			input: &Edge{
				head: nil,
				tail: newEmptyVertex("a"),
			},
			output: errors.New("head vertex must not be nil"),
		},
		{
			input: &Edge{
				head: newEmptyVertex("a"),
				tail: nil,
			},
			output: errors.New("tail vertex must not be nil"),
		},
	}

	for i, testCase := range data {
		err := d.AddEdge(testCase.input)
		if err == nil {
			t.Errorf("case %v: received nil error when given invalid input", i)
		}

		if err.Error() != testCase.output.Error() {
			t.Errorf("case %v: error does not match expected. actual %v expected %v", i, err, testCase.output)
		}
	}

}

func TestAddEdgeVertexHeadDoesNotExistInDAG(t *testing.T) {
	expected := errors.New("vertex b does not exist in graph")
	d := NewDag()
	v1 := NewVertex("a", nil, func() {})
	err := d.AddVertex(v1)
	if err != nil {
		t.Error("failed to add valid vertex to graph")
		return
	}

	v2 := NewVertex("b", nil, func() {})
	e, err := NewEdge(v1, v2)
	if err != nil {
		t.Errorf("failed to create valid edge")
		return
	}

	err = d.AddEdge(e)
	if err == nil {
		t.Error("added invalid edge to graph")
		return
	}

	if err.Error() != expected.Error() {
		t.Error("did not receive correct error when attempting to add edge with missing head vertex in graph")
	}
}

func TestAddEdgeVertexTailDoesNotExistsInDAG(t *testing.T) {
	expected := errors.New("vertex b does not exist in graph")
	d := NewDag()
	v1 := NewVertex("a", nil, func() {})
	err := d.AddVertex(v1)
	if err != nil {
		t.Error("failed to add valid vertex to graph")
		return
	}

	v2 := NewVertex("b", nil, func() {})
	e, err := NewEdge(v2, v1)
	if err != nil {
		t.Errorf("failed to create valid edge")
		return
	}

	err = d.AddEdge(e)
	if err == nil {
		t.Error("added invalid edge to graph")
		return
	}

	if err.Error() != expected.Error() {
		t.Error("did not receive correct error when attempting to add tail with missing head vertex in graph")
	}
}

func TestAddEdgeAlreadyExistsInDAG(t *testing.T) {
	expected := errors.New("edge from a to b already exists in graph")
	d := NewDag()
	v1 := newEmptyVertex("a")
	v2 := newEmptyVertex("b")
	err := d.AddVertex(v1)
	if err != nil {
		t.Error("failed to add valid vertex to graph")
		return
	}
	err = d.AddVertex(v2)
	if err != nil {
		t.Error("failed to add valid vertex to graph")
		return
	}

	e, err := NewEdge(v1, v2)
	if err != nil {
		t.Error("failed to create valid edge")
	}

	err = d.AddEdge(e)
	if err != nil {
		t.Error("failed to add valid edge to DAG")
	}

	err = d.AddEdge(e)
	if err == nil {
		t.Error("did not error when attempting to add a duplicate edge")
		return
	}

	if err.Error() != expected.Error() {
		t.Errorf("did not receive duplicate edge error when attempting to add a duplicate edge. Actual: %v Expected: %v", err, expected)
	}
}

func TestAddEdgeDependenciesAndDependantUpdated(t *testing.T) {
	d := NewDag()
	head := newEmptyVertex("a")
	tail := newEmptyVertex("b")
	err := d.AddVertex(head)
	if err != nil {
		t.Error("failed to add valid vertex to graph")
		return
	}

	err = d.AddVertex(tail)
	if err != nil {
		t.Error("failed to add valid vertex to graph")
		return
	}

	e, err := NewEdge(head, tail)
	if err != nil {
		t.Error("failed to create new valid edge")
		return
	}

	err = d.AddEdge(e)
	if err != nil {
		t.Error("failed to add valid edge to graph")
		return
	}

	if len(head.dependencies) != 0 {
		t.Error("dependencies for head changed when they should not")
	}

	if len(tail.dependents) != 0 {
		t.Error("dependents for head changed when they should not")
	}

	if len(head.dependents) != 1 && head.dependents[0] != tail.id {
		t.Error("dependents for head not updated correctly")
	}

	if len(tail.dependencies) != 1 && tail.dependencies[0] != head.id {
		t.Error("dependencies for tail not updated correctly")
	}
}

func TestAddEdgeEdgesByIdValid(t *testing.T) {
	v1 := newEmptyVertex("a")
	v2 := newEmptyVertex("b")

	d := NewDag()
	err := d.AddVertex(v1)
	if err != nil {
		t.Errorf("errored when adding valid vertex to graph. err %v", err)
	}

	err = d.AddVertex(v2)
	if err != nil {
		t.Errorf("errored when adding valid vertex to graph. err %v", err)
	}

	e, err := d.AddEdgeByVerticesId("a", "b")
	if err != nil {
		t.Error("errored when adding valid edge to graph")
	}

	if e.head != v1 {
		t.Error("new edge does not have correct head vertex")
	}

	if e.tail != v2 {
		t.Error("new edge does not have correct tail vertex")
	}

	if len(v1.dependencies) != 0 {
		t.Error("head's dependencies updated when it should not be")
	}

	if len(v1.dependents) != 1 && v1.dependents[0] != v2.id {
		t.Error("head's dependents were not updated correctly")
	}

	if len(v2.dependents) != 0 {
		t.Error("tail's dependents were not updated correctly")
	}

	if len(v2.dependencies) != 1 && v2.dependencies[0] != v1.id {
		t.Error("tail's dependencies were not updated correctly")
	}
}

func TestAddEdgeCrossMethodsCompatible(t *testing.T) {
	t.Error("failing")
}

func TestTopSortInvalidDAG(t *testing.T) {
	t.Error("failing")
}

func TestTopSortValidDAG(t *testing.T) {
	t.Error("failing")
}

func TestRunDAGBasicPipeline(t *testing.T) {
	d := NewDag()
	vertices := []*Vertex{
		NewVertex("a", nil, func() {}),
		NewVertex("b", nil, func() {}),
		NewVertex("c", nil, func() {}),
	}

	a2b, _ := NewEdge(vertices[0], vertices[1])
	b2c, _ := NewEdge(vertices[1], vertices[2])

	for _, v := range vertices {
		d.AddVertex(v)
	}

	d.AddEdge(a2b)
	d.AddEdge(b2c)

	dependencyChecker := &dependencyChecker{
		rLock:     &sync.RWMutex{},
		completed: []string{},
	}

	for i := range vertices {
		dependencies := make([]string, len(vertices[i].dependencies))
		copy(dependencies, vertices[i].dependencies)
		vertices[i].f = dependencyChecker.createNewDepTestFunc(vertices[i], t, dependencies)
	}

	runDag(d)

	if len(dependencyChecker.completed) == len(vertices) {
		t.Error("not all steps completed")
	}

}

func TestRunDAGThreeDependenciesOneChild(t *testing.T) {
	d := NewDag()
	dependencyChecker := &dependencyChecker{
		rLock:     &sync.RWMutex{},
		completed: []string{},
	}

	vertices := []*Vertex{
		NewVertex("a", nil, func() {}),
		NewVertex("b", nil, func() {}),
		NewVertex("c", nil, func() {}),
		NewVertex("d", nil, func() {}),
	}

	err := addVerticesToDAG(d, vertices)
	if err != nil {
		t.Error("failed adding vertices to graph")
	}

	a2d, _ := NewEdge(vertices[0], vertices[3])
	b2d, _ := NewEdge(vertices[1], vertices[3])
	c2d, _ := NewEdge(vertices[2], vertices[3])

	d.AddEdge(a2d)
	d.AddEdge(b2d)
	d.AddEdge(c2d)

	for i := range vertices {
		dependencies := make([]string, len(vertices[i].dependencies))
		copy(dependencies, vertices[i].dependencies)
		vertices[i].f = dependencyChecker.createNewDepTestFunc(vertices[i], t, dependencies)
	}

	runDag(d)

	if len(dependencyChecker.completed) != len(vertices) {
		t.Error("did not complete all vertices")
	}
}

func TestRunDAGValid(t *testing.T) {
	t.Error("failing")
}

func TestRunDAGInvalid(t *testing.T) {
	t.Error("failing")
}

type dependencyChecker struct {
	rLock     *sync.RWMutex
	completed []string
}

func (d *dependencyChecker) createNewDepTestFunc(v *Vertex, t *testing.T, dependencies []string) func() {
	return func() {
		d.checkDependencies(t, v, dependencies)
		d.addCompletedVertex(v)
	}
}

func (d *dependencyChecker) addCompletedVertex(v *Vertex) {
	d.rLock.Lock()
	d.completed = append(d.completed, v.id)
	d.rLock.Unlock()
}

func (d *dependencyChecker) checkDependencies(t *testing.T, v *Vertex, deps []string) {
	d.rLock.RLock()
	defer d.rLock.RUnlock()
	for _, dep := range deps {
		found := false
		for _, c := range d.completed {
			if dep == c {
				found = true
			}
		}
		if found == false {
			t.Errorf("current vertex %v has ran while it's dependency %v has not.", v.id, dep)
		}
	}
}

func addVerticesToDAG(d *DAG, vertices []*Vertex) error {
	for _, v := range vertices {
		err := d.AddVertex(v)
		if err != nil {
			return err
		}
	}
	return nil
}
