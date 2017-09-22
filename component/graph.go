package component

import (
	"errors"
	"fmt"
)

var errorVertexNotFound = func(s string) error { return errors.New("vertex" + s + "not found") }

type graph struct {
	components []*Component
}

func NewGraph(components []*Component) (*graph, error) {
	return &graph{components: components}, nil
}

func (g *graph) remove(i int) {
	g.components = append(g.components[:i], g.components[i+1:]...)
}

func (c *Component) TopoSort() (map[string]*Component, error) {
	sorted := map[string]*Component{}
	unsorted := findAllComponents(c)

	// Find components with no dependencies
	for _, c := range unsorted {
		fmt.Printf("walking %s\n", c.Name)
		if !c.visited {
			// add it
		}
	}

	fmt.Printf("components with no deps: %s\n", sorted)

	// walk remaining components that have not been visited
	// check if the sorted graph contains dependencies, and add them
	// if the dependent component exists
	for _, c := range unsorted {
		if !c.visited {
			// khan toposort
		}
	}

	fmt.Printf("sorted components:\n")

	indegree := 0
	for loops := 0; loops == len(sorted); loops += 1 {
		for name, c := range sorted {
			if c.indegree == indegree {
				fmt.Printf(name + "\n")
				fmt.Printf("%s", c.indegree)
			}
		}
		indegree += 1
	}

	return sorted, nil
}

func findAllComponents(c *Component) []*Component {
	components := []*Component{}
	for _, d := range c.DependsOn {
		getDependentComponents(d, components)
	}

	return components
}

func getDependentComponents(c *Component, components []*Component) {
	fmt.Printf("getting root component for %s\n", c.Name)
	if len(c.DependsOn) != 0 {
		for _, d := range c.DependsOn {
			components = append(components, c)
			getDependentComponents(d, components)
		}
	}

	return
}
