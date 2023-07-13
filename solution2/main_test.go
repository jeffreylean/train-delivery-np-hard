package main

import (
	"solution2/anneal"
	"solution2/loader"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test1(t *testing.T) {
	train, pkg, graph := loader.Initialize("test/test1.txt")
	asgn := assignPkgToTrain(graph, train, pkg)
	r, m := planRoute(graph, asgn, train, pkg)

	initialState := State{TrainAssignment: asgn, Route: r, Move: m, Graph: graph, Train: train, Package: pkg}

	s := anneal.Init(initialState, anneal.Config{Iteration: 10000, Temperature: 25000, AneallingFactor: 0.99})
	assert.Equal(t, 70, int(s.Energy()))
}

func Test2(t *testing.T) {
	train, pkg, graph := loader.Initialize("test/test2.txt")
	asgn := assignPkgToTrain(graph, train, pkg)
	r, m := planRoute(graph, asgn, train, pkg)

	initialState := State{TrainAssignment: asgn, Route: r, Move: m, Graph: graph, Train: train, Package: pkg}

	s := anneal.Init(initialState, anneal.Config{Iteration: 10000, Temperature: 25000, AneallingFactor: 0.99})
	assert.Equal(t, 40, int(s.Energy()))
}

func Test3(t *testing.T) {
	train, pkg, graph := loader.Initialize("test/test3.txt")
	asgn := assignPkgToTrain(graph, train, pkg)
	r, m := planRoute(graph, asgn, train, pkg)

	initialState := State{TrainAssignment: asgn, Route: r, Move: m, Graph: graph, Train: train, Package: pkg}

	s := anneal.Init(initialState, anneal.Config{Iteration: 10000, Temperature: 25000, AneallingFactor: 0.99})
	assert.Equal(t, 26, int(s.Energy()))
}

func Test4(t *testing.T) {
	train, pkg, graph := loader.Initialize("test/test4.txt")
	asgn := assignPkgToTrain(graph, train, pkg)
	r, m := planRoute(graph, asgn, train, pkg)

	initialState := State{TrainAssignment: asgn, Route: r, Move: m, Graph: graph, Train: train, Package: pkg}

	s := anneal.Init(initialState, anneal.Config{Iteration: 10000, Temperature: 25000, AneallingFactor: 0.99})
	assert.Equal(t, 25, int(s.Energy()))
}
