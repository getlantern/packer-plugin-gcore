package gcore

import "fmt"

type Artifact struct {
	ImageID   string
	ImageName string
	RegionID  int
	ProjectID int
	StateData map[string]interface{}
}

func (*Artifact) BuilderId() string {
	return BuilderID
}

func (a *Artifact) Files() []string {
	return nil
}

func (a *Artifact) Id() string {
	return a.ImageID
}

func (a *Artifact) String() string {
	return fmt.Sprintf("Gcore image: %s (ID: %s, Region: %d, Project: %d)",
		a.ImageName, a.ImageID, a.RegionID, a.ProjectID)
}

func (a *Artifact) State(name string) interface{} {
	return a.StateData[name]
}

func (a *Artifact) Destroy() error {
	return nil
}
