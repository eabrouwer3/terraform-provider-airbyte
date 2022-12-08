package provider

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceSource,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("airbyte_source.test", "id", regexp.MustCompile("^[0-9a-fA-F]{8}\\b-[0-9a-fA-F]{4}\\b-[0-9a-fA-F]{4}\\b-[0-9a-fA-F]{4}\\b-[0-9a-fA-F]{12}$")),
					resource.TestCheckResourceAttrPair("airbyte_source.test", "workspace_id", "airbyte_workspace.test", "id"),
					resource.TestCheckResourceAttrPair("airbyte_source.test", "definition_id", "airbyte_source_definition.test", "id"),
					resource.TestCheckResourceAttrPair("airbyte_source.test", "definition_name", "airbyte_source_definition.test", "name"),
					resource.TestCheckResourceAttr("airbyte_source.test", "name", "test_source"),
					resource.TestCheckResourceAttr("airbyte_source.test", "connection_configuration", "{}"),
				),
			},
		},
	})
}

const testAccResourceSource = `
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
`
