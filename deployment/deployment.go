package deployment

import (
	"fmt"

	"github.com/malnick/genesis/component"
)

// Deployment abstracts the business logic of a given set of components
type Deployment struct {
	Components []component.Component
	Workspaces []string
}

// New returns a new instance of a Deployment with functional options set
func New(options ...Option) (*Deployment, error) {
	d := &Deployment{
		Components: []component.Component{},
		Workspaces: []string{},
	}

	for _, o := range options {
		if err := o(d); err != nil {
			return d, err
		}
	}

	return d, nil
}

// Option is a abstract type for making functional options
type Option func(*Deployment) error

const (
	workspaceProduction  = "production"
	workspaceDevelopment = "development"
)

var acceptableWorkspaces = []string{
	workspaceProduction,
	workspaceDevelopment,
}

func OptionAddComponent(c component.Component) Option {
	return func(d *Deployment) error {
		d.Components = append(d.Components, c)
		return nil
	}
}

func OptionWorkspace(w string) Option {
	return func(d *Deployment) error {
		for _, i := range acceptableWorkspaces {
			if i == w {
				d.Workspaces = append(d.Workspaces, w)
				return nil
			}
		}
		return fmt.Errorf("workspace invalid: %s is not one of %s", w, acceptableWorkspaces)
	}
}
