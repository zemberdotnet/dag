package dag

import (
	"errors"
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

func TestAddEdgeEdgesByIdUpdate(t *testing.T) {
	t.Errorf("failing")
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
  runOrder := make(chan string, 3)
  expectedRunOrder := []string{"a", "b", "c"}
  verticies := []*Vertex{
    NewVertex("a", nil, func () {
      runOrder <- "a"
    }),
    NewVertex("b",nil, func () {
      runOrder <- "b"
    }),
    NewVertex("c", nil, func () {
      runOrder <- "c"
    }),
  }

  a2b, _ :=  NewEdge(verticies[0], verticies[1])
  b2c, _ := NewEdge(verticies[1], verticies[2])

  d := NewDag()
  for _, v := range verticies {
    d.AddVertex(v)
  }

  d.AddEdge(a2b)
  d.AddEdge(b2c)

  runDag(d)

  for i := 0; i < len(runOrder); i++ {
    v := <- runOrder
    if v != expectedRunOrder[i] {
      t.Error("dag did not execute in correct order")
    }
  }

}

func TestRunDAGThreeDependenciesOneChild(t *testing.T) {
  runOrder := make(chan string, 4)

  verticies := []*Vertex{
    NewVertex("a", nil, func() {
      runOrder <- "a"
    }),
 NewVertex("b", nil, func() {
      runOrder <- "a"
    }),
 NewVertex("c", nil, func() {
      runOrder <- "a"
    }),
 NewVertex("d", nil, func() {
      runOrder <- "a"
    }),
  }  

  a2d, _ := NewEdge(verticies[0], verticies[3])
  b2d, _ := NewEdge(verticies[1], verticies[3])
  c2d, _ := NewEdge(verticies[2], verticies[3])

  d := NewDag()

  d.AddEdge(a2d)
  d.AddEdge(b2d)
  d.AddEdge(c2d)


  runDag(d)
  
  for i := 0; i < len(runOrder); i++ {
    v := <- runOrder
    if i == 0 || i == 1 || i == 2 {
      if v != verticies[0].id || v != verticies[1].id || v != verticies[2].id {
        t.Error("dag did not execute in correct order")
      }
    }

    if i == 3 {
      if v !=  verticies[3].id {
        t.Error("dag did not execute in correct order")
      }
  }
  
  if len(runOrder) != 0 {
    t.Error("an unexpected number of steps occured")
  }
}
}

func TestRunDAGValid(t *testing.T) {
	t.Error("failing")
}

func TestRunDAGInvalid(t *testing.T) {
	t.Error("failing")
}
