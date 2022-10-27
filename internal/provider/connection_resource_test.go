package provider

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceConnection(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceConnection,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("airbyte_connection.test", "id", regexp.MustCompile("^[0-9a-fA-F]{8}\\b-[0-9a-fA-F]{4}\\b-[0-9a-fA-F]{4}\\b-[0-9a-fA-F]{4}\\b-[0-9a-fA-F]{12}$")),
					resource.TestCheckResourceAttrPair("airbyte_connection.test", "source_id", "airbyte_source.test", "id"),
					resource.TestCheckResourceAttrPair("airbyte_connection.test", "destination_id", "airbyte_destination.test", "id"),
					resource.TestCheckResourceAttr("airbyte_connection.test", "status", "active"),
					// Check for default API provided values
					resource.TestCheckResourceAttr("airbyte_connection.test", "name", "test_source <> test_destination"), // Default naming convention when name is not supplied
					resource.TestCheckResourceAttr("airbyte_connection.test", "namespace_definition", "source"),
					resource.TestCheckResourceAttr("airbyte_connection.test", "operation_ids.#", "0"),
					resource.TestCheckResourceAttr("airbyte_connection.test", "schedule_type", "manual"),
					resource.TestCheckResourceAttr("airbyte_connection.test", "geography", "auto"),
					// The data source doesn't return this value, but that's ok - the default is true
					resource.TestCheckResourceAttr("airbyte_connection.test", "sync_catalog.0.destination_config.selected", "true"),
				),
			},
		},
	})
}

const testAccResourceConnection = `
resource "airbyte_workspace" "test" {
  name = "test_workspace"
}

resource "airbyte_source_definition" "test" {
  name = "test_source_definition"
  docker_repository = "eabrouwer3/airbyte-test-data-source"
  docker_image_tag = "0.0.1"
  documentation_url = "https://github.com/eabrouwer3/airbyte-test-data-source"
}

resource "airbyte_source" "test" {
  definition_id = airbyte_source_definition.test.id
  workspace_id = airbyte_workspace.test.id
  name = "test_source"
  connection_configuration = jsonencode({})
}

resource "airbyte_destination_definition" "test" {
  name = "test_destination_definition"
  docker_repository = "eabrouwer3/airbyte-test-data-destination"
  docker_image_tag = "0.0.1"
  documentation_url = "https://github.com/eabrouwer3/airbyte-test-data-destination"
}

resource "airbyte_destination" "test" {
  definition_id = airbyte_destination_definition.test.id
  workspace_id = airbyte_workspace.test.id
  name = "test_destination"
  connection_configuration = jsonencode({})
}

data "airbyte_source_schema_catalog" "test" {
  source_id = airbyte_source.test.id
}

resource "airbyte_connection" "test" {
  source_id = airbyte_source.test.id
  destination_id = airbyte_destination.test.id
  status = "active"
  sync_catalog = data.airbyte_source_schema_catalog.test.sync_catalog
}
`
