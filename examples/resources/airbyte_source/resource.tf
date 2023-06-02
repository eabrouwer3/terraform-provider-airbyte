resource "airbyte_workspace" "test" {
  name = "test_workspace"
}

# Super basic custom source
resource "airbyte_source_definition" "custom" {
  name              = "custom_source_definition"
  docker_repository = "eabrouwer3/airbyte-test-data-source"
  docker_image_tag  = "0.0.1"
  documentation_url = "https://github.com/eabrouwer3/airbyte-test-data-source"
}

resource "airbyte_source" "custom" {
  definition_id = airbyte_source_definition.custom.id
  workspace_id  = airbyte_workspace.test.id
  name          = "custom_source"
  # The source definition above takes no parameters
  # Note that no terraform validation happens for this - just errors directly from the API
  connection_configuration = jsonencode({})
}

# More complex with existing source definition
# Uses E2E Testing Source (https://docs.airbyte.com/integrations/sources/e2e-test/)
resource "airbyte_source" "existing" {
  # Find the definition_id for an existing source here https://github.com/airbytehq/airbyte/tree/master/airbyte-integrations/connectors
  # Look for metadata.yaml file, e.g. https://github.com/airbytehq/airbyte/blob/master/airbyte-integrations/connectors/source-e2e-test/metadata.yaml
  definition_id = "d53f9084-fa6b-4a5a-976c-5b8392f4ad8a"
  workspace_id  = airbyte_workspace.test.id
  name          = "e2e_source"
  # Find the spec either in the docs for the connector
  # Or, find it here: https://github.com/airbytehq/airbyte/tree/master/airbyte-integrations/connectors
  # Look for src/main/resources/spec.json file, e.g. https://github.com/airbytehq/airbyte/blob/master/airbyte-integrations/connectors/source-e2e-test/src/main/resources/spec.json
  connection_configuration = jsonencode({
    type = "CONTINUOUS_FEED"
    mock_catalog = {
      type               = "SINGLE_STREAM"
      stream_name        = "data_stream"
      stream_schema      = "{ \"type\": \"object\", \"properties\": { \"column1\": { \"type\": \"string\" } } }"
      stream_duplication = 1
    }
  })
}
