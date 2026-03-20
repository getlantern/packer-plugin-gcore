package gcore

import (
	"context"
	"fmt"

	"github.com/hashicorp/packer-plugin-sdk/multistep"
	"github.com/hashicorp/packer-plugin-sdk/packer"

	gcoresdk "github.com/G-Core/gcore-go"
	"github.com/G-Core/gcore-go/cloud"
)

// StepStopInstance stops the instance before creating an image from its boot volume.
type StepStopInstance struct{}

func (s *StepStopInstance) Run(ctx context.Context, state multistep.StateBag) multistep.StepAction {
	client := state.Get("gcore_client").(*gcoresdk.Client)
	instanceID := state.Get("instance_id").(string)
	ui := state.Get("ui").(packer.Ui)

	ui.Say("Stopping instance before image creation...")

	stopParams := cloud.InstanceActionParams{
		OfBasicActionInstanceSerializer: &cloud.InstanceActionParamsBodyBasicActionInstanceSerializer{
			Action: "stop",
		},
	}

	if _, err := client.Cloud.Instances.ActionAndPoll(ctx, instanceID, stopParams); err != nil {
		err = fmt.Errorf("error stopping instance: %w", err)
		state.Put("error", err)
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	ui.Say("Instance stopped")
	return multistep.ActionContinue
}

func (s *StepStopInstance) Cleanup(state multistep.StateBag) {}
