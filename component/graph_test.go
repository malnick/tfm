package component

import "testing"

func TestSortComponents(t *testing.T) {
	components := []*Component{
		&Component{
			Name: "datascience/dcos",
			DependsOn: []*Component{
				&Component{
					Name: "managment/puppet",
				},
				&Component{
					Name: "datascience/cassandra",
					DependsOn: []*Component{
						&Component{
							Name: "managment/puppet",
						},
					},
				},
				&Component{
					Name: "datascience/kafka",
					DependsOn: []*Component{
						&Component{
							Name: "managment/puppet",
						},
					},
				},
			},
		},
	}

	// expected toposort to be puppet -> kafak -||-> cassandra -> dcos
	graph, err := NewGraph(components)
	if err != nil {
		t.Error(err)
	}

	sorted, err := graph.TopoSort()
	if err != nil {
		t.Error(err)
	}

	if len(sorted) != 4 {
		t.Error("expected 4 components, got", len(sorted))
	}

	if sorted["management/puppet"].indegree != 0 {
		t.Error("expected management/puppet, got", sorted["management/puppet"].indegree)
	}

	if sorted["datascience/dcos"].indegree != 3 {
		t.Error("expected datascience/dcos, got", sorted["datascience/dcos"].indegree)
	}
}
