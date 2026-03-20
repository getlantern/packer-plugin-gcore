package version

import "github.com/hashicorp/packer-plugin-sdk/version"

var (
	Version           = "0.1.0"
	VersionPrerelease = ""
	PluginVersion     = version.InitializePluginVersion(Version, VersionPrerelease)
)
