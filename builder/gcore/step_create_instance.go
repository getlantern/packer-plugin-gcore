package gcore

import (
	"context"
	"fmt"

	"github.com/hashicorp/packer-plugin-sdk/multistep"
	"github.com/hashicorp/packer-plugin-sdk/packer"

	gcoresdk "github.com/G-Core/gcore-go"
	"github.com/G-Core/gcore-go/cloud"
	"github.com/G-Core/gcore-go/option"
	"github.com/G-Core/gcore-go/shared/constant"
)

// StepCreateInstance creates a Gcore cloud VM from a base image.
type StepCreateInstance struct {
	client   *gcoresdk.Client
	instance *cloud.Instance
}

func (s *StepCreateInstance) Run(ctx context.Context, state multistep.StateBag) multistep.StepAction {
	config := state.Get("config").(*Config)
	ui := state.Get("ui").(packer.Ui)

	ui.Say("Creating Gcore cloud instance...")

	client := gcoresdk.NewClient(
		option.WithAPIKey(config.APIKey),
		option.WithCloudProjectID(int64(config.ProjectID)),
		option.WithCloudRegionID(int64(config.RegionID)),
	)
	s.client = &client

	interfaces := []cloud.InstanceNewParamsInterfaceUnion{
		{OfExternal: &cloud.InstanceNewParamsInterfaceExternal{
			Type: constant.ValueOf[constant.External](),
		}},
	}

	volumes := []cloud.InstanceNewParamsVolumeUnion{
		{OfImage: &cloud.InstanceNewParamsVolumeImage{
			Source:    constant.ValueOf[constant.Image](),
			ImageID:   config.ImageID,
			Size:      gcoresdk.Int(int64(config.VolumeSize)),
			TypeName:  config.VolumeType,
			BootIndex: gcoresdk.Int(0),
		}},
	}

	params := cloud.InstanceNewParams{
		Name:       gcoresdk.String(config.InstanceName),
		Flavor:     config.FlavorID,
		Interfaces: interfaces,
		Volumes:    volumes,
	}

	result, err := client.Cloud.Instances.NewAndPoll(ctx, params)
	if err != nil {
		err = fmt.Errorf("error creating instance: %w", err)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	s.instance = result

	ui.Say(fmt.Sprintf("Instance created: %s (ID: %s)", result.Name, result.ID))
	state.Put("instance_id", result.ID)
	state.Put("gcore_client", s.client)

	// Find the boot volume ID for later image creation
	if len(result.Volumes) > 0 {
		state.Put("boot_volume_id", result.Volumes[0].ID)
		ui.Say(fmt.Sprintf("Boot volume: %s", result.Volumes[0].ID))
	}

	return multistep.ActionContinue
}

func (s *StepCreateInstance) Cleanup(state multistep.StateBag) {
	ui := state.Get("ui").(packer.Ui)

	if s.instance == nil {
		return
	}

	ui.Say("Destroying Gcore instance...")

	// Delete instance and associated floating IPs, but NOT volumes
	// (volumes may be needed if image creation hasn't happened yet)
	params := cloud.InstanceDeleteParams{
		DeleteFloatings: gcoresdk.Bool(true),
	}

	if err := s.client.Cloud.Instances.DeleteAndPoll(context.Background(), s.instance.ID, params); err != nil {
		ui.Error(fmt.Sprintf("Error destroying instance: %s", err))
	} else {
		ui.Say("Instance destroyed")
	}
}
