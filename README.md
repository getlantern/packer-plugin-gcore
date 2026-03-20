# packer-plugin-gcore

Packer plugin for building images on [Gcore Cloud](https://gcore.com/cloud).

## Usage

Add the plugin to your Packer HCL config:

```hcl
packer {
  required_plugins {
    gcore = {
      version = ">= 0.1.0"
      source  = "github.com/getlantern/gcore"
    }
  }
}
```

Define a source:

```hcl
variable "gcore_api_key" {
  type      = string
  sensitive = true
  default   = env("GCORE_API_KEY")
}

source "gcore" "lantern-box" {
  api_key         = var.gcore_api_key
  project_id      = 12345
  region_id       = 8    # Luxembourg

  source_image_id = "abc-123"  # Ubuntu 24.04 base image ID
  flavor_id       = "g1-standard-1-2"
  volume_size     = 20

  image_name      = "lantern-box-${var.lantern_box_version}"
  ssh_username    = "ubuntu"
}
```

Use in a build:

```hcl
build {
  sources = ["source.gcore.lantern-box"]

  provisioner "file" { ... }
  provisioner "shell" { ... }
}
```

## Configuration Reference

### Required

| Parameter | Description |
|-----------|-------------|
| `api_key` | Gcore API key |
| `project_id` | Gcore project ID |
| `region_id` | Gcore region ID |
| `source_image_id` | Base image ID to boot from |
| `flavor_id` | Instance flavor (size) |
| `image_name` | Name for the output image |

### Optional

| Parameter | Default | Description |
|-----------|---------|-------------|
| `instance_name` | `packer-<build>` | Name for the temporary build instance |
| `volume_size` | `20` | Boot volume size in GB |
| `volume_type` | `ssd_hiiops` | Volume type |
| `network_id` | (external) | Network to attach; defaults to external |
| `keypair_name` | | SSH keypair name (if pre-created) |
| `user_data` | | Cloud-init user data |
| `ssh_username` | `ubuntu` | SSH user for provisioning |
| `image_tags` | | Tags for the output image |

## How It Works

1. Creates a VM instance from the base image
2. Waits for a public IPv4 address
3. Connects via SSH and runs Packer provisioners
4. Stops the instance
5. Creates an image from the boot volume via `InstanceImage.NewFromVolume()`
6. Destroys the temporary instance

## Building

```bash
go build -o packer-plugin-gcore
```

## License

MPL-2.0
