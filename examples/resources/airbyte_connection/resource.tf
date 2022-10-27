resource "airbyte_workspace" "test" {
  name = "test_workspace"
}

# Most basic possible example - custom sources/destinations with zero configuration
resource "airbyte_source_definition" "custom" {
  name              = "test_source_definition"
  docker_repository = "eabrouwer3/airbyte-test-data-source"
  docker_image_tag  = "0.0.1"
  documentation_url = "https://github.com/eabrouwer3/airbyte-test-data-source"
}

resource "airbyte_source" "custom" {
  definition_id            = airbyte_source_definition.custom.id
  workspace_id             = airbyte_workspace.test.id
  name                     = "test_source"
  connection_configuration = jsonencode({})
}

resource "airbyte_destination_definition" "custom" {
  name              = "test_destination_definition"
  docker_repository = "eabrouwer3/airbyte-test-data-destination"
  docker_image_tag  = "0.0.1"
  documentation_url = "https://github.com/eabrouwer3/airbyte-test-data-destination"
}

resource "airbyte_destination" "custom" {
  definition_id            = airbyte_destination_definition.custom.id
  workspace_id             = airbyte_workspace.test.id
  name                     = "test_destination"
  connection_configuration = jsonencode({})
}

data "airbyte_source_schema_catalog" "custom" {
  source_id = airbyte_source.custom.id
}

resource "airbyte_connection" "custom" {
  source_id      = airbyte_source.custom.id
  destination_id = airbyte_destination.custom.id
  status         = "active"
  sync_catalog   = data.airbyte_source_schema_catalog.custom.sync_catalog
}

# More complex E2E Testing setup with some custom configuration
resource "airbyte_source" "e2e" {
  # Find the definition_id for an existing source here: https://github.com/airbytehq/airbyte/blob/master/airbyte-config/init/src/main/resources/seed/source_definitions.yaml
  definition_id = "d53f9084-fa6b-4a5a-976c-5b8392f4ad8a"
  workspace_id  = airbyte_workspace.test.id
  name          = "e2e_source"
  # Find the spec either in the docs for the connector
  # Or, find it here: https://github.com/airbytehq/airbyte/blob/master/airbyte-config/init/src/main/resources/seed/source_specs.yaml
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

resource "airbyte_destination" "e2e" {
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

data "airbyte_source_schema_catalog" "e2e" {
  source_id = airbyte_source.e2e.id
}

resource "airbyte_connection" "e2e" {
  name           = "E2E Testing Connection"
  source_id      = airbyte_source.e2e.id
  destination_id = airbyte_destination.e2e.id
  status         = "inactive"
  sync_catalog = [{
    source_schema = data.airbyte_source_schema_catalog.custom.sync_catalog.0.source_schema
    # Config some of the destination settings
    destination_config = merge(
      data.airbyte_source_schema_catalog.custom.sync_catalog.0.destination_config,
      {
        alias_name            = "data_stream_destination_alias"
        destination_sync_mode = "overwrite"
        sync_mode             = "full_refresh"
      }
    )
  }]
  # Set up a time schedule
  schedule_type = "basic"
  basic_schedule = {
    time_unit = "hours"
    units     = 24
  }
  # Add some special resource requirements
  resource_requirements = {
    cpu_request    = "0.5"
    cpu_limit      = "0.5"
    memory_request = "500Mi"
    memory_limit   = "500Mi"
  }
}