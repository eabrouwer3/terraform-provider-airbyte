resource "airbyte_workspace" "default" {
  name = "basic_test"
}

data "airbyte_workspace" "by_id" {
  id = airbyte_workspace.default.id
}

data "airbyte_workspace" "by_slug" {
  slug = airbyte_workspace.default.slug
}