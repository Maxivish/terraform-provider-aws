---
subcategory: "Connect"
layout: "aws"
page_title: "AWS: aws_connect_instance_storage_config"
description: |-
  Terraform resource for managing an AWS Connect InstanceStorageConfig.
---

# Resource: aws_connect_instance_storage_config

Terraform resource for managing an AWS Connect InstanceStorageConfig.

## Example Usage

### Basic Usage

```terraform
resource "aws_connect_instance_storage_config" "example" {
}
```

## Argument Reference

The following arguments are required:

* `example_arg` - (Required) Concise argument description.

The following arguments are optional:

* `optional_arg` - (Optional) Concise argument description.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `arn` - ARN of the InstanceStorageConfig.
* `example_attribute` - Concise description.

## Timeouts

`aws_connect_instance_storage_config` provides the following [Timeouts](https://www.terraform.io/docs/configuration/blocks/resources/syntax.html#operation-timeouts) configuration options:

* `create` - (Optional, Default: `60m`)
* `update` - (Optional, Default: `180m`)
* `delete` - (Optional, Default: `90m`)

## Import

Connect InstanceStorageConfig can be imported using the `example_id_arg`, e.g.,

```
$ terraform import aws_connect_instance_storage_config.example rft-8012925589
```
