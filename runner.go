package dag

import ("sync")

func runDag(dag *DAG) {
	finished := make(chan string, len(dag.vertices))
	indegree := make(map[*Vertex]int)
	processed := 0
  processedLock := &sync.RWMutex{} 


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
		queue = queue[1:]

		// Start Task
		go func() {
			vertex.f()
			finished <- vertex.id
      processedLock.Lock()
			processed++
      processedLock.Unlock()
		}()

		if len(queue) == 0 {
			select {
			case finishedTask := <-finished:
				finishedVertex := dag.vertexById[finishedTask]
				for _, dependent := range finishedVertex.dependents {
					indegree[dag.vertexById[dependent]]--
					if indegree[dag.vertexById[dependent]] == 0 {
						queue = append(queue, dag.vertexById[dependent])
					}
				}

			default:
        processedLock.RLock()
				if len(queue) == 0 && processed == len(dag.vertices) {
          processedLock.RUnlock()
					return
				}
        processedLock.RUnlock()
			}
		}
	}
}
