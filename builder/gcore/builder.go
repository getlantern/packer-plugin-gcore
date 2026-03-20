package gcore

import (
	"context"
	"fmt"

	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/hashicorp/packer-plugin-sdk/communicator"
	"github.com/hashicorp/packer-plugin-sdk/multistep"
	"github.com/hashicorp/packer-plugin-sdk/multistep/commonsteps"
	"github.com/hashicorp/packer-plugin-sdk/packer"
)

const BuilderID = "getlantern.gcore"

type Builder struct {
	config Config
	runner multistep.Runner
}

func (b *Builder) ConfigSpec() hcldec.ObjectSpec {
	// TODO: Generate proper HCL2 spec using packer-sdc generate.
	// For now, return nil which makes Packer fall back to JSON config.
	return nil
}

func (b *Builder) Prepare(raws ...interface{}) ([]string, []string, error) {
	warnings, err := b.config.Prepare(raws...)
	if err != nil {
		return nil, warnings, err
	}
	return nil, warnings, nil
}

func (b *Builder) Run(ctx context.Context, ui packer.Ui, hook packer.Hook) (packer.Artifact, error) {
	state := new(multistep.BasicStateBag)
	state.Put("config", &b.config)
	state.Put("hook", hook)
	state.Put("ui", ui)

	steps := []multistep.Step{
		&StepCreateInstance{},
		&StepWaitForIP{},
		&communicator.StepConnect{
			Config:    &b.config.Comm,
			Host:      commHost,
			SSHConfig: b.config.Comm.SSHConfigFunc(),
		},
		&commonsteps.StepProvision{},
		&StepStopInstance{},
		&StepCreateImage{},
	}

	b.runner = commonsteps.NewRunner(steps, b.config.PackerConfig, ui)
	b.runner.Run(ctx, state)

	if rawErr, ok := state.GetOk("error"); ok {
		return nil, rawErr.(error)
	}

	imageID, ok := state.GetOk("image_id")
	if !ok {
		return nil, fmt.Errorf("image_id not found in state — image creation may have failed")
	}

	artifact := &Artifact{
		ImageID:   imageID.(string),
		ImageName: b.config.ImageName,
		RegionID:  b.config.RegionID,
		ProjectID: b.config.ProjectID,
		StateData: map[string]interface{}{
			"generated_data": state.Get("generated_data"),
		},
	}

	return artifact, nil
}

// commHost extracts the instance IP from the state for the SSH communicator.
func commHost(state multistep.StateBag) (string, error) {
	ip, ok := state.GetOk("instance_ip")
	if !ok {
		return "", fmt.Errorf("instance_ip not found in state")
	}
	return ip.(string), nil
}
