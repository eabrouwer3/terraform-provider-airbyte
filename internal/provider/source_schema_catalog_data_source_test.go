package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccDataSourceSourceSchemaCatalog(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceSourceSchemaCatalog,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair("data.airbyte_source_schema_catalog.test", "source_id", "airbyte_source.test", "id"),
					resource.TestCheckResourceAttr("data.airbyte_source_schema_catalog.test", "sync_catalog.#", "1"),
					resource.TestCheckResourceAttr("data.airbyte_source_schema_catalog.test", "sync_catalog.0.source_schema.name", "appliances"),
					resource.TestCheckResourceAttr("data.airbyte_source_schema_catalog.test", "sync_catalog.0.source_schema.json_schema", "{\"type\":\"object\",\"$schema\":\"http://json-schema.org/draft-07/schema#\",\"properties\":{\"id\":{\"type\":\"integer\"},\"uid\":{\"type\":\"string\"},\"brand\":{\"type\":\"string\"},\"equipment\":{\"type\":\"string\"}}}"),
					resource.TestCheckResourceAttr("data.airbyte_source_schema_catalog.test", "sync_catalog.0.source_schema.supported_sync_modes.0", "incremental"),
					resource.TestCheckResourceAttr("data.airbyte_source_schema_catalog.test", "sync_catalog.0.source_schema.supported_sync_modes.1", "full_refresh"),
				),
			},
		},
	})
}

const testAccDataSourceSourceSchemaCatalog = `
resource "airbyte_workspace" "test" {
  name = "test_workspace"
}

resource "airbyte_source_definition" "test" {
	workspace_id = airbyte_workspace.test.id
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

data "airbyte_source_schema_catalog" "test" {
  source_id = airbyte_source.test.id
}
`
