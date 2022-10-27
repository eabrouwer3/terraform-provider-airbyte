resource "airbyte_workspace" "test" {
  name = "test_workspace"
}

# Super basic custom destination
resource "airbyte_destination_definition" "test" {
  name              = "test_destination_definition"
  docker_repository = "eabrouwer3/airbyte-test-data-destination"
  docker_image_tag  = "0.0.1"
  documentation_url = "https://github.com/eabrouwer3/airbyte-test-data-destination"
}

resource "airbyte_destination" "test" {
  definition_id = airbyte_destination_definition.test.id
  workspace_id  = airbyte_workspace.test.id
  name          = "test_destination"
  # The destination definition above takes no parameters
  # Note that no terraform validation happens for this - just errors directly from the API
  connection_configuration = jsonencode({})
}

# More complex with existing destination definition
# Uses E2E Testing Destination (https://docs.airbyte.com/integrations/destinations/e2e-test/)
resource "airbyte_destination" "existing" {
  # Find the definition_id for an existing source here: https://github.com/airbytehq/airbyte/blob/master/airbyte-config/init/src/main/resources/seed/destination_definitions.yaml
  definition_id = "2eb65e87-983a-4fd7-b3e3-9d9dc6eb8537"
  workspace_id  = airbyte_workspace.test.id
  name          = "e2e_destination"
  # Find the spec either in the docs for the connector
  # Or, find it here: https://github.com/airbytehq/airbyte/blob/master/airbyte-config/init/src/main/resources/seed/destination_specs.yaml
  connection_configuration = jsonencode({
    type = "LOGGING"
    logging_config = {
      logging_type    = "FirstN"
      max_entry_count = 100
    }
  })
}