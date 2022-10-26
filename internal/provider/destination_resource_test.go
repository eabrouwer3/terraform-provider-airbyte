package provider

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceDestination(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceDestination,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("airbyte_destination.test", "id", regexp.MustCompile("^[0-9a-fA-F]{8}\\b-[0-9a-fA-F]{4}\\b-[0-9a-fA-F]{4}\\b-[0-9a-fA-F]{4}\\b-[0-9a-fA-F]{12}$")),
					resource.TestCheckResourceAttrPair("airbyte_destination.test", "workspace_id", "airbyte_workspace.test", "id"),
					resource.TestCheckResourceAttrPair("airbyte_destination.test", "definition_id", "airbyte_destination_definition.test", "id"),
					resource.TestCheckResourceAttrPair("airbyte_destination.test", "definition_name", "airbyte_destination_definition.test", "name"),
					resource.TestCheckResourceAttr("airbyte_destination.test", "name", "test_destination"),
					resource.TestCheckResourceAttr("airbyte_destination.test", "connection_configuration", "{}"),
				),
			},
		},
	})
}

const testAccResourceDestination = `
resource "airbyte_workspace" "test" {
  name = "test_workspace"
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
`
