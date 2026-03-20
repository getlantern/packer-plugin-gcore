package gcore

import (
	"fmt"

	"github.com/hashicorp/packer-plugin-sdk/common"
	"github.com/hashicorp/packer-plugin-sdk/communicator"
	"github.com/hashicorp/packer-plugin-sdk/packer"
	"github.com/hashicorp/packer-plugin-sdk/template/config"
	"github.com/hashicorp/packer-plugin-sdk/template/interpolate"
)

type Config struct {
	common.PackerConfig    `mapstructure:",squash"`
	Comm communicator.Config `mapstructure:",squash"`

	// Gcore API credentials
	APIKey    string `mapstructure:"api_key" required:"true"`
	ProjectID int    `mapstructure:"project_id" required:"true"`
	RegionID  int    `mapstructure:"region_id" required:"true"`

	// Instance configuration
	FlavorID     string `mapstructure:"flavor_id" required:"true"`
	ImageID      string `mapstructure:"source_image_id" required:"true"`
	InstanceName string `mapstructure:"instance_name"`
	NetworkID    string `mapstructure:"network_id"`
	KeypairName  string `mapstructure:"keypair_name"`
	UserData     string `mapstructure:"user_data"`

	// Volume configuration
	VolumeSize int    `mapstructure:"volume_size"`
	VolumeType string `mapstructure:"volume_type"`

	// Image output configuration
	ImageName string            `mapstructure:"image_name" required:"true"`
	ImageTags map[string]string `mapstructure:"image_tags"`

	ctx interpolate.Context
}

func (c *Config) Prepare(raws ...interface{}) ([]string, error) {
	err := config.Decode(c, &config.DecodeOpts{
		PluginType:         "gcore",
		Interpolate:        true,
		InterpolateContext: &c.ctx,
		InterpolateFilter: &interpolate.RenderFilter{
			Exclude: []string{},
		},
	}, raws...)
	if err != nil {
		return nil, err
	}

	var errs *packer.MultiError

	if c.APIKey == "" {
		errs = packer.MultiErrorAppend(errs, fmt.Errorf("api_key is required"))
	}
	if c.ProjectID == 0 {
		errs = packer.MultiErrorAppend(errs, fmt.Errorf("project_id is required"))
	}
	if c.RegionID == 0 {
		errs = packer.MultiErrorAppend(errs, fmt.Errorf("region_id is required"))
	}
	if c.FlavorID == "" {
		errs = packer.MultiErrorAppend(errs, fmt.Errorf("flavor_id is required"))
	}
	if c.ImageID == "" {
		errs = packer.MultiErrorAppend(errs, fmt.Errorf("source_image_id is required"))
	}
	if c.ImageName == "" {
		errs = packer.MultiErrorAppend(errs, fmt.Errorf("image_name is required"))
	}

	if c.InstanceName == "" {
		c.InstanceName = fmt.Sprintf("packer-%s", c.PackerBuildName)
	}
	if c.VolumeSize == 0 {
		c.VolumeSize = 20
	}
	if c.VolumeType == "" {
		c.VolumeType = "ssd_hiiops"
	}

	if c.Comm.SSHUsername == "" {
		c.Comm.SSHUsername = "ubuntu"
	}

	if errs != nil && len(errs.Errors) > 0 {
		return nil, errs
	}

	return nil, nil
}
