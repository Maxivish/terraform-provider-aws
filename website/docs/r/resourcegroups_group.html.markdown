---
subcategory: "Resource Groups"
layout: "aws"
page_title: "AWS: aws_resourcegroups_group"
description: |-
  Provides a Resource Group.
---

# Resource: aws_resourcegroups_group

Provides a Resource Group.

## Example Usage

```terraform
resource "aws_resourcegroups_group" "test" {
  name = "test-group"

  resource_query {
    query = <<JSON
{
  "ResourceTypeFilters": [
    "AWS::EC2::Instance"
  ],
  "TagFilters": [
    {
      "Key": "Stage",
      "Values": ["Test"]
    }
  ]
}
JSON
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The resource group's name. A resource group name can have a maximum of 127 characters, including letters, numbers, hyphens, dots, and underscores. The name cannot start with `AWS` or `aws`.
* `description` - (Optional) A description of the resource group.
* `resource_query` - (Required) A `resource_query` block. Resource queries are documented below.
* `tags` - (Optional) Key-value map of resource tags. If configured with a provider [`default_tags` configuration block](https://registry.terraform.io/providers/hashicorp/aws/latest/docs#default_tags-configuration-block) present, tags with matching keys will overwrite those defined at the provider-level.

An `resource_query` block supports the following arguments:

* `query` - (Required) The resource query as a JSON string.
* `type` - (Required) The type of the resource query. Defaults to `TAG_FILTERS_1_0`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `arn` - The ARN assigned by AWS for this resource group.
* `tags_all` - A map of tags assigned to the resource, including those inherited from the provider [`default_tags` configuration block](https://registry.terraform.io/providers/hashicorp/aws/latest/docs#default_tags-configuration-block).

## Import

Resource groups can be imported using the `name`, e.g.,

```
$ terraform import aws_resourcegroups_group.foo resource-group-name
```
