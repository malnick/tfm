package component

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	logging "github.com/op/go-logging"
)

var (
	log = logging.MustGetLogger("tfm.component")

	acceptableWorkspaces = []string{
		WorkspaceProduction,
		WorkspaceDevelopment,
	}
)

type Component struct {
	Name              string            `json:"name"`
	NameSpace         string            `json:"namespace"`
	IsApplied         bool              `json:"is_applied"`
	AllowedWorkspaces []string          `json:"allowed_workspaces"`
	DependsOn         []*Component      `json:"depends_on"`
	Status            map[string]string `json:"status"`

	Workspace        string `json:"-"`
	WorkingDirectory string `json:"-"`
	ComponentPath    string `json:"-"`
	RemoteState      struct {
		S3Bucket   string
		AWSProfile string
		AWSRegion  string
	} `json:"-"`

	visited  bool
	indegree int
}

// New returns a new instance of a Component with functional options set
func New(options ...Option) (*Component, error) {
	c := &Component{
		AllowedWorkspaces: []string{},
		DependsOn:         []*Component{},
		WorkingDirectory:  "terraform",
		Workspace:         WorkspaceDevelopment,
		Status:            map[string]string{WorkspaceDevelopment: ""},
	}

	for _, o := range options {
		if err := o(c); err != nil {
			return c, err
		}
	}

	c.ComponentPath = strings.Join([]string{c.WorkingDirectory, c.Name}, "/")

	log.Debugf("checking that %s exists...", c.ComponentPath)
	if _, err := os.Stat(c.ComponentPath); err != nil {
		log.Warningf("%s does not exist, creating...", c.ComponentPath)
		if err := os.MkdirAll(c.ComponentPath, 0755); err != nil {
			return c, err
		}
	}

	return c, nil
}

func (c *Component) WriteMetaJSON() error {
	data, err := json.MarshalIndent(&c, "", "    ")
	if err != nil {
		return err
	}

	metaPath := strings.Join([]string{c.ComponentPath, "meta.json"}, "/")

	log.Infof("writiting meta.json for %s at %s", c.Name, metaPath)

	return ioutil.WriteFile(metaPath, data, 0644)
}

func (c *Component) ReadMetaJSON() error {
	metaPath := strings.Join([]string{c.ComponentPath, "meta.json"}, "/")
	data, err := ioutil.ReadFile(metaPath)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, &c)
}

func (c *Component) StatComponentFilesOrCreate() error {
	type terraformDefaults struct {
		path     string
		template string
	}

	defaults := map[string]terraformDefaults{
		"main": terraformDefaults{
			path:     filepath.Join(c.ComponentPath, "main.tf"),
			template: mainTemplate,
		},

		"vars": terraformDefaults{
			path:     filepath.Join(c.ComponentPath, "vars.tf"),
			template: varsTemplate,
		},
	}

	for _, td := range defaults {
		if _, err := os.Stat(td.path); err != nil {
			log.Warningf("%s does not exist, writing default...", td.path)
			if !os.IsNotExist(err) {
				return err
			}

			tmpl, err := template.New("").Parse(td.template)
			if err != nil {
				return err
			}

			var rendered bytes.Buffer
			if err := tmpl.Execute(&rendered, c); err != nil {
				return err
			}

			log.Infof("rendered main.tf:\n%s", rendered.String())

			if err := ioutil.WriteFile(td.path, rendered.Bytes(), 0644); err != nil {
				return err
			}
		}
	}

	return nil
}

func Create(c *Component) error {
	if err := c.StatComponentFilesOrCreate(); err != nil {
		return err
	}

	return c.WriteMetaJSON()
}

func (c *Component) Execute(action string) error {
	if err := SelectWorkspace(c); err != nil {
		return err
	}

	var planBuffer bytes.Buffer
	log.Infof("running %s action for %s in %s workspace", action, c.Name, c.Workspace)

	cmd := exec.Command("terraform", action)
	cmd.Dir = c.ComponentPath
	cmd.Stdout = &planBuffer
	cmd.Stderr = &planBuffer
	cmd.Env = []string{
		fmt.Sprintf("AWS_PROFILE=%s", c.Workspace),
		"PATH=/bin",
	}

	err := cmd.Run()
	c.Status[c.Workspace] = planBuffer.String()

	if err != nil {
		return errors.New(planBuffer.String())
	}

	if action == "apply" {
		c.IsApplied = true
	}
	return nil
}

func Run(action string, c *Component) error {
	if err := c.Execute(action); err != nil {
		log.Error(err)
	}
	return c.WriteMetaJSON()
}

func MakeComponentFromMetaJSON(path string) (*Component, error) {
	var c = &Component{}
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(data, &c); err != nil {
		return nil, err
	}

	return c, nil
}
