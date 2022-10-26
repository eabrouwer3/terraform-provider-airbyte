package provider

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceSource_basic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("airbyte_source.basic", "id", regexp.MustCompile("^[0-9a-fA-F]{8}\\b-[0-9a-fA-F]{4}\\b-[0-9a-fA-F]{4}\\b-[0-9a-fA-F]{4}\\b-[0-9a-fA-F]{12}$")),
					resource.TestCheckResourceAttrPair("airbyte_source.basic", "workspace_id", "airbyte_workspace.basic", "id"),
					resource.TestCheckResourceAttrPair("airbyte_source.basic", "source_definition_id", "airbyte_source_definition.basic", "id"),
					resource.TestCheckResourceAttr("airbyte_source.basic", "name", "basic_test"),
					resource.TestCheckResourceAttr("airbyte_source.basic", "connection_configuration", "{}"),
				),
			},
		},
	})
}

//func TestAccResourceSource_complex(t *testing.T) {
//	resource.UnitTest(t, resource.TestCase{
//		PreCheck:                 func() { testAccPreCheck(t) },
//		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
//		Steps: []resource.TestStep{
//			{
//				Config: testAccResourceSourceDefinition_complex,
//				Check: resource.ComposeTestCheckFunc(
//					resource.TestMatchResourceAttr("airbyte_source_definition.complex", "id", regexp.MustCompile("^[0-9a-fA-F]{8}\\b-[0-9a-fA-F]{4}\\b-[0-9a-fA-F]{4}\\b-[0-9a-fA-F]{4}\\b-[0-9a-fA-F]{12}$")),
//					resource.TestCheckResourceAttr("airbyte_source_definition.complex", "name", "complex_test"),
//					resource.TestCheckResourceAttr("airbyte_source_definition.complex", "docker_repository", "airbyte/source-github"),
//					resource.TestCheckResourceAttr("airbyte_source_definition.complex", "docker_image_tag", "0.3.7"),
//					resource.TestCheckResourceAttr("airbyte_source_definition.complex", "documentation_url", "https://hub.docker.com/r/airbyte/source-github"),
//					resource.TestCheckResourceAttr("airbyte_source_definition.complex", "default_resource_requirements.cpu_request", "0.25"),
//					resource.TestCheckResourceAttr("airbyte_source_definition.complex", "default_resource_requirements.cpu_limit", "0.25"),
//					resource.TestCheckResourceAttr("airbyte_source_definition.complex", "default_resource_requirements.memory_request", "200Mi"),
//					resource.TestCheckResourceAttr("airbyte_source_definition.complex", "default_resource_requirements.memory_limit", "200Mi"),
//					resource.TestCheckResourceAttr("airbyte_source_definition.complex", "job_specific_resource_requirements.0.cpu_request", "0.5"),
//					resource.TestCheckResourceAttr("airbyte_source_definition.complex", "job_specific_resource_requirements.0.cpu_limit", "0.5"),
//					resource.TestCheckResourceAttr("airbyte_source_definition.complex", "job_specific_resource_requirements.0.memory_request", "500Mi"),
//					resource.TestCheckResourceAttr("airbyte_source_definition.complex", "job_specific_resource_requirements.0.memory_limit", "500Mi"),
//					resource.TestCheckNoResourceAttr("airbyte_source_definition.complex", "job_specific_resource_requirements.1.cpu_request"),
//					resource.TestCheckNoResourceAttr("airbyte_source_definition.complex", "job_specific_resource_requirements.1.cpu_limit"),
//					resource.TestCheckResourceAttr("airbyte_source_definition.complex", "job_specific_resource_requirements.1.memory_request", "50Mi"),
//					resource.TestCheckResourceAttr("airbyte_source_definition.complex", "job_specific_resource_requirements.1.memory_limit", "50Mi"),
//				),
//			},
//		},
//	})
//}

const testAccResourceSource_basic = `
resource "airbyte_workspace" "basic" {
  name = "basic_test"
}

resource "airbyte_source_definition" "basic" {
  name = "basic_test"
  docker_repository = "eabrouwer3/airbyte-test-data-source"
  docker_image_tag = "0.0.1"
  documentation_url = "https://github.com/eabrouwer3/airbyte-test-data-source"
}

resource "airbyte_source" "basic" {
  source_definition_id = airbyte_source_definition.basic.id
  workspace_id = airbyte_workspace.basic.id
  name = "basic_test"
  connection_configuration = jsonencode({})
}
`

//const testAccResourceSource_complex = `
//resource "airbyte_source" "complex" {
//  name = "complex_test"
//  docker_repository = "airbyte/source-github"
//  docker_image_tag = "0.3.7"
//  documentation_url = "https://hub.docker.com/r/airbyte/source-github"
//
//  default_resource_requirements = {
//    cpu_request = "0.25"
//    cpu_limit = "0.25"
//    memory_request = "200Mi"
//    memory_limit = "200Mi"
//  }
//
//  job_specific_resource_requirements = [{
//    job_type = "sync"
//    cpu_request = "0.5"
//    cpu_limit = "0.5"
//    memory_request = "500Mi"
//    memory_limit = "500Mi"
//  }, {
//    job_type = "check_connection"
//    memory_request = "50Mi"
//    memory_limit = "50Mi"
//  }]
//}
//`
