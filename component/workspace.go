package component

import (
	"bytes"
	"errors"
	"os/exec"
)

func SelectWorkspace(c *Component) error {
	var wsBuffer bytes.Buffer

	log.Infof("checking out %s workspace...", c.Workspace)

	wsCmd := exec.Command("terraform", "workspace", "select", c.Workspace)
	wsCmd.Dir = c.ComponentPath
	wsCmd.Stdout = &wsBuffer
	wsCmd.Stderr = &wsBuffer

	wsErr := wsCmd.Run()
	c.Status[c.Workspace] = wsBuffer.String()

	if wsErr != nil {
		return errors.New(c.Status[c.Workspace])
	}

	return wsErr
}
