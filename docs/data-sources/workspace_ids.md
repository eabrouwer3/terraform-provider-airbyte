---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "airbyte_workspace_ids Data Source - terraform-provider-airbyte"
subcategory: ""
description: |-
  Get all Airbyte Workspace ids (first will always be the default one created by Airbyte on launch)
---

# airbyte_workspace_ids (Data Source)

Get all Airbyte Workspace ids (first will always be the default one created by Airbyte on launch)

## Example Usage

```terraform
data "airbyte_workspace_ids" "all" {}

resource "airbyte_source" "test" {
  # First workspace returned will be the default created by Airbyte when bootstrapped
  workspace_id = data.airbyte_workspace_ids.all.ids.0
  # ...
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Read-Only

- `id` (String) The ID of this resource.
- `ids` (List of String) Workspace Id List


