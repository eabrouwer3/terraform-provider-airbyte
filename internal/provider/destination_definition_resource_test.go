package provider

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceDestinationDefinition_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		//CheckDestroy:             testAccResourceDestinationDefinitionDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceDestinationDefinition_basic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("airbyte_destination_definition.basic", "id", regexp.MustCompile("^[0-9a-fA-F]{8}\\b-[0-9a-fA-F]{4}\\b-[0-9a-fA-F]{4}\\b-[0-9a-fA-F]{4}\\b-[0-9a-fA-F]{12}$")),
					resource.TestCheckResourceAttr("airbyte_destination_definition.basic", "name", "basic_test"),
					resource.TestCheckResourceAttr("airbyte_destination_definition.basic", "docker_repository", "eabrouwer3/airbyte-test-data-destination"),
					resource.TestCheckResourceAttr("airbyte_destination_definition.basic", "docker_image_tag", "0.0.1"),
					resource.TestCheckResourceAttr("airbyte_destination_definition.basic", "documentation_url", "https://github.com/eabrouwer3/airbyte-test-data-destination"),
					resource.TestCheckNoResourceAttr("airbyte_destination_definition.basic", "default_resource_requirements"),
					resource.TestCheckNoResourceAttr("airbyte_destination_definition.basic", "job_specific_resource_requirements"),
				),
			},
		},
	})
}

func TestAccResourceDestinationDefinition_complex(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceDestinationDefinition_complex,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("airbyte_destination_definition.complex", "id", regexp.MustCompile("^[0-9a-fA-F]{8}\\b-[0-9a-fA-F]{4}\\b-[0-9a-fA-F]{4}\\b-[0-9a-fA-F]{4}\\b-[0-9a-fA-F]{12}$")),
					resource.TestCheckResourceAttr("airbyte_destination_definition.complex", "name", "complex_test"),
					resource.TestCheckResourceAttr("airbyte_destination_definition.complex", "docker_repository", "eabrouwer3/airbyte-test-data-destination"),
					resource.TestCheckResourceAttr("airbyte_destination_definition.complex", "docker_image_tag", "0.0.1"),
					resource.TestCheckResourceAttr("airbyte_destination_definition.complex", "documentation_url", "https://github.com/eabrouwer3/airbyte-test-data-destination"),
					resource.TestCheckResourceAttr("airbyte_destination_definition.complex", "default_resource_requirements.cpu_request", "0.25"),
					resource.TestCheckResourceAttr("airbyte_destination_definition.complex", "default_resource_requirements.cpu_limit", "0.25"),
					resource.TestCheckResourceAttr("airbyte_destination_definition.complex", "default_resource_requirements.memory_request", "200Mi"),
					resource.TestCheckResourceAttr("airbyte_destination_definition.complex", "default_resource_requirements.memory_limit", "200Mi"),
					resource.TestCheckResourceAttr("airbyte_destination_definition.complex", "job_specific_resource_requirements.0.cpu_request", "0.5"),
					resource.TestCheckResourceAttr("airbyte_destination_definition.complex", "job_specific_resource_requirements.0.cpu_limit", "0.5"),
					resource.TestCheckResourceAttr("airbyte_destination_definition.complex", "job_specific_resource_requirements.0.memory_request", "500Mi"),
					resource.TestCheckResourceAttr("airbyte_destination_definition.complex", "job_specific_resource_requirements.0.memory_limit", "500Mi"),
					resource.TestCheckNoResourceAttr("airbyte_destination_definition.complex", "job_specific_resource_requirements.1.cpu_request"),
					resource.TestCheckNoResourceAttr("airbyte_destination_definition.complex", "job_specific_resource_requirements.1.cpu_limit"),
					resource.TestCheckResourceAttr("airbyte_destination_definition.complex", "job_specific_resource_requirements.1.memory_request", "50Mi"),
					resource.TestCheckResourceAttr("airbyte_destination_definition.complex", "job_specific_resource_requirements.1.memory_limit", "50Mi"),
				),
			},
		},
	})
}

const testAccResourceDestinationDefinition_basic = `
resource "airbyte_workspace" "test" {
  name = "test_workspace"
}

resource "airbyte_destination_definition" "basic" {
	workspace_id = airbyte_workspace.test.id
  name = "basic_test"
  docker_repository = "eabrouwer3/airbyte-test-data-destination"
  docker_image_tag = "0.0.1"
  documentation_url = "https://github.com/eabrouwer3/airbyte-test-data-destination"
}
`

const testAccResourceDestinationDefinition_complex = `
resource "airbyte_workspace" "test" {
  name = "test_workspace"
}

resource "airbyte_destination_definition" "complex" {
	workspace_id = airbyte_workspace.test.id
  name = "complex_test"
  docker_repository = "eabrouwer3/airbyte-test-data-destination"
  docker_image_tag = "0.0.1"
  documentation_url = "https://github.com/eabrouwer3/airbyte-test-data-destination"

  default_resource_requirements = {
    cpu_request = "0.25"
    cpu_limit = "0.25"
    memory_request = "200Mi"
    memory_limit = "200Mi"
  }

  job_specific_resource_requirements = [{
    job_type = "sync"
    cpu_request = "0.5"
    cpu_limit = "0.5"
    memory_request = "500Mi"
    memory_limit = "500Mi"
  }, {
    job_type = "check_connection"
    memory_request = "50Mi"
    memory_limit = "50Mi"
  }]
}
`
