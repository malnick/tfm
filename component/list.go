package component

import (
	"fmt"
	"io/ioutil"
)

func (c *Component) List(details bool) error {
	metaPath := fmt.Sprintf("%s/meta.json", c.ComponentPath)
	f, err := ioutil.ReadFile(metaPath)
	if err != nil {
		return err
	}

	fmt.Printf("%s", f)

	return nil
}
