package component

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Option is a abstract type for making functional options
type Option func(*Component) error

func OptionDependsOn(tfPath, componentNames string) Option {
	return func(c *Component) error {
		if componentNames == "" {
			return nil
		}

		// Validate component names exist
		names := strings.Split(componentNames, ",")
		for _, n := range names {
			metaPath := filepath.Join(tfPath, n, "meta.json")
			comp, err := MakeComponentFromMetaJSON(metaPath)
			if err != nil {
				return errors.New("error adding component dependency: " + err.Error())
			}

			c.DependsOn = append(c.DependsOn, comp)
			// Create vertices of components
			// Make sure there are no cycles
		}
		return nil
	}
}

func OptionName(n string) Option {
	return func(c *Component) error {
		if len(n) == 0 {
			return errors.New("name option can not be empty")
		}

		splitName := strings.Split(n, "/")
		if len(splitName) > 2 {
			return errors.New("can not have more than a single namespace")
		}

		c.NameSpace = splitName[0]
		c.Name = n
		return nil
	}
}

func OptionWorkingDirectory(d string) Option {
	return func(c *Component) error {
		if len(d) == 0 {
			return errors.New("working directory option can not be empty")
		}

		if _, err := os.Stat(d); err != nil {
			return err
		}

		c.WorkingDirectory = d
		return nil
	}
}

func OptionWorkspace(w string) Option {
	return func(c *Component) error {
		for _, i := range acceptableWorkspaces {
			if i == w {
				c.Workspace = w
				return nil
			}
		}
		return fmt.Errorf("workspace invalid: %s is not one of %s", w, acceptableWorkspaces)
	}
}
func OptionAllowedWorkspaces(workspaces []string) Option {
	return func(c *Component) error {
		for _, w := range workspaces {
			for _, i := range acceptableWorkspaces {
				if i == w {
					c.AllowedWorkspaces = append(c.AllowedWorkspaces, w)
					return nil
				}
			}
		}
		return fmt.Errorf("workspace invalid: %s is not one of %s", workspaces, acceptableWorkspaces)
	}
}
