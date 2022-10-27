resource "airbyte_workspace" "test" {
  name = "test_workspace"
}

resource "airbyte_source_definition" "test" {
  name              = "test_source_definition"
  docker_repository = "eabrouwer3/airbyte-test-data-source"
  docker_image_tag  = "0.0.1"
  documentation_url = "https://github.com/eabrouwer3/airbyte-test-data-source"
}

resource "airbyte_source" "test" {
  definition_id            = airbyte_source_definition.test.id
  workspace_id             = airbyte_workspace.test.id
  name                     = "test_source"
  connection_configuration = jsonencode({})
}

# Get schema catalog for the source above
data "airbyte_source_schema_catalog" "test" {
  source_id = airbyte_source.test.id
}