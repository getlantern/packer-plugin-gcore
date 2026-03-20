package gcore

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/packer-plugin-sdk/multistep"
	"github.com/hashicorp/packer-plugin-sdk/packer"

	gcoresdk "github.com/G-Core/gcore-go"
	"github.com/G-Core/gcore-go/cloud"
)

// StepWaitForIP polls the instance until it has a public IPv4 address.
type StepWaitForIP struct{}

func (s *StepWaitForIP) Run(ctx context.Context, state multistep.StateBag) multistep.StepAction {
	config := state.Get("config").(*Config)
	client := state.Get("gcore_client").(*gcoresdk.Client)
	instanceID := state.Get("instance_id").(string)
	ui := state.Get("ui").(packer.Ui)

	ui.Say("Waiting for instance to get a public IP...")

	deadline := time.After(5 * time.Minute)
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			state.Put("error", ctx.Err())
			return multistep.ActionHalt
		case <-deadline:
			err := fmt.Errorf("timed out waiting for public IP on instance %s", instanceID)
			state.Put("error", err)
			ui.Error(err.Error())
			return multistep.ActionHalt
		case <-ticker.C:
			instance, err := client.Cloud.Instances.Get(ctx, instanceID, cloud.InstanceGetParams{})
			if err != nil {
				ui.Say(fmt.Sprintf("Waiting... (error checking: %s)", err))
				continue
			}

			// Look for a public IPv4 address in the instance addresses
			for _, addrList := range instance.Addresses {
				for _, a := range addrList {
					if a.Type == "fixed" && a.Addr != "" {
						ip := a.Addr
						ui.Say(fmt.Sprintf("Instance IP: %s", ip))
						state.Put("instance_ip", ip)
						config.Comm.SSHHost = ip
						return multistep.ActionContinue
					}
				}
			}
		}
	}
}

func (s *StepWaitForIP) Cleanup(state multistep.StateBag) {}
