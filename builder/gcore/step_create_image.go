package gcore

import (
	"context"
	"fmt"

	"github.com/hashicorp/packer-plugin-sdk/multistep"
	"github.com/hashicorp/packer-plugin-sdk/packer"

	gcoresdk "github.com/G-Core/gcore-go"
	"github.com/G-Core/gcore-go/cloud"
)

// StepCreateImage creates a Gcore image from the instance's boot volume.
type StepCreateImage struct{}

func (s *StepCreateImage) Run(ctx context.Context, state multistep.StateBag) multistep.StepAction {
	config := state.Get("config").(*Config)
	client := state.Get("gcore_client").(*gcoresdk.Client)
	volumeID := state.Get("boot_volume_id").(string)
	ui := state.Get("ui").(packer.Ui)

	ui.Say(fmt.Sprintf("Creating image '%s' from boot volume %s...", config.ImageName, volumeID))

	params := cloud.InstanceImageNewFromVolumeParams{
		VolumeID: volumeID,
		Name:     config.ImageName,
	}

	image, err := client.Cloud.Instances.Images.NewFromVolumeAndPoll(ctx, params)
	if err != nil {
		err = fmt.Errorf("error creating image from volume: %w", err)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	ui.Say(fmt.Sprintf("Image created: %s (ID: %s)", image.Name, image.ID))
	state.Put("image_id", image.ID)

	return multistep.ActionContinue
}

func (s *StepCreateImage) Cleanup(state multistep.StateBag) {}
