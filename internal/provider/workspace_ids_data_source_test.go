package provider

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceWorkspaceIds(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceWorkspaceIds,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.airbyte_workspace_ids.all", "id"),
					resource.TestMatchResourceAttr("data.airbyte_workspace_ids.all", "ids.#", regexp.MustCompile("^[1-9][0-9]*$")),
				),
			},
		},
	})
}

const testAccDataSourceWorkspaceIds = `
data "airbyte_workspace_ids" "all" {}
`
