package provider

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceNormalizationOperation(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceNormalizationOperation,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("airbyte_operation.normalization", "id", regexp.MustCompile("^[0-9a-fA-F]{8}\\b-[0-9a-fA-F]{4}\\b-[0-9a-fA-F]{4}\\b-[0-9a-fA-F]{4}\\b-[0-9a-fA-F]{12}$")),
					resource.TestCheckResourceAttrPair("airbyte_operation.normalization", "workspace_id", "airbyte_workspace.test", "id"),
					resource.TestCheckResourceAttr("airbyte_operation.normalization", "name", "normalization_operation"),
					resource.TestCheckResourceAttr("airbyte_operation.normalization", "operator_type", "normalization"),
					resource.TestCheckResourceAttr("airbyte_operation.normalization", "normalization_option", "basic"),
					resource.TestCheckNoResourceAttr("airbyte_operation.normalization", "dbt"),
					resource.TestCheckNoResourceAttr("airbyte_operation.normalization", "webhook"),
				),
			},
		},
	})
}

const testAccResourceNormalizationOperation = `
resource "airbyte_workspace" "test" {
  name = "basic_test"
}

resource "airbyte_operation" "normalization" {
  workspace_id = airbyte_workspace.test.id
  name = "normalization_operation"
  operator_type = "normalization"
  normalization_option = "basic"
}
`

// I don't know how to configure these myself right now - need an easy to test way to do it
//const testAccResourceDbtOperation = `
//resource "airbyte_workspace" "test" {
//  name = "basic_test"
//}
//
//resource "airbyte_operation" "dbt" {
//  workspace_id = airbyte_workspace.test.id
//  name = "dbt_operation"
//  operator_type = "dbt"
//  dbt = {
//    git_repo_url = ""
//    git_repo_branch = ""
//    docker_image = ""
//    dbt_arguments = ""
//  }
//}
//`
//
//const testAccResourceWebhookOperation = `
//resource "airbyte_workspace" "test" {
//  name = "basic_test"
//}
//
//resource "airbyte_operation" "webhook" {
//  workspace_id = airbyte_workspace.test.id
//  name = "webhook_operation"
//  operator_type = "webhook"
//  webhook = {
//    execution_url = ""
//    execution_body = ""
//    webhook_config_id = ""
//  }
//}
//`
