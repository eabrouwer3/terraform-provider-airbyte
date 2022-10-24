package provider

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceWorkspace_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceWorkspace_basic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("data.airbyte_workspace.by_id", "id", regexp.MustCompile("^[0-9a-fA-F]{8}\\b-[0-9a-fA-F]{4}\\b-[0-9a-fA-F]{4}\\b-[0-9a-fA-F]{4}\\b-[0-9a-fA-F]{12}$")),
					resource.TestMatchResourceAttr("data.airbyte_workspace.by_id", "slug", regexp.MustCompile("^basic_test")),
					resource.TestCheckResourceAttr("data.airbyte_workspace.by_id", "name", "basic_test"),
					resource.TestMatchResourceAttr("data.airbyte_workspace.by_slug", "id", regexp.MustCompile("^[0-9a-fA-F]{8}\\b-[0-9a-fA-F]{4}\\b-[0-9a-fA-F]{4}\\b-[0-9a-fA-F]{4}\\b-[0-9a-fA-F]{12}$")),
					resource.TestMatchResourceAttr("data.airbyte_workspace.by_slug", "slug", regexp.MustCompile("^basic_test")),
					resource.TestCheckResourceAttr("data.airbyte_workspace.by_slug", "name", "basic_test"),
				),
			},
		},
	})
}

func TestAccDataSourceWorkspace_complex(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceWorkspace_complex,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("data.airbyte_workspace.by_id", "id", regexp.MustCompile("^[0-9a-fA-F]{8}\\b-[0-9a-fA-F]{4}\\b-[0-9a-fA-F]{4}\\b-[0-9a-fA-F]{4}\\b-[0-9a-fA-F]{12}$")),
					resource.TestMatchResourceAttr("data.airbyte_workspace.by_id", "slug", regexp.MustCompile("^complex_test")),
					resource.TestCheckResourceAttrPair("airbyte_workspace.complex", "name", "data.airbyte_workspace.by_id", "name"),
					resource.TestCheckResourceAttrPair("airbyte_workspace.complex", "email", "data.airbyte_workspace.by_id", "email"),
					resource.TestCheckResourceAttrPair("airbyte_workspace.complex", "display_setup_wizard", "data.airbyte_workspace.by_id", "display_setup_wizard"),
					resource.TestCheckResourceAttrPair("airbyte_workspace.complex", "anonymous_data_collection", "data.airbyte_workspace.by_id", "anonymous_data_collection"),
					resource.TestCheckResourceAttrPair("airbyte_workspace.complex", "news", "data.airbyte_workspace.by_id", "news"),
					resource.TestCheckResourceAttrPair("airbyte_workspace.complex", "security_updates", "data.airbyte_workspace.by_id", "security_updates"),

					resource.TestMatchResourceAttr("data.airbyte_workspace.by_slug", "id", regexp.MustCompile("^[0-9a-fA-F]{8}\\b-[0-9a-fA-F]{4}\\b-[0-9a-fA-F]{4}\\b-[0-9a-fA-F]{4}\\b-[0-9a-fA-F]{12}$")),
					resource.TestMatchResourceAttr("data.airbyte_workspace.by_slug", "slug", regexp.MustCompile("^complex_test")),
					resource.TestCheckResourceAttrPair("airbyte_workspace.complex", "name", "data.airbyte_workspace.by_slug", "name"),
					resource.TestCheckResourceAttrPair("airbyte_workspace.complex", "email", "data.airbyte_workspace.by_slug", "email"),
					resource.TestCheckResourceAttrPair("airbyte_workspace.complex", "display_setup_wizard", "data.airbyte_workspace.by_slug", "display_setup_wizard"),
					resource.TestCheckResourceAttrPair("airbyte_workspace.complex", "anonymous_data_collection", "data.airbyte_workspace.by_slug", "anonymous_data_collection"),
					resource.TestCheckResourceAttrPair("airbyte_workspace.complex", "news", "data.airbyte_workspace.by_slug", "news"),
					resource.TestCheckResourceAttrPair("airbyte_workspace.complex", "security_updates", "data.airbyte_workspace.by_slug", "security_updates"),
				),
			},
		},
	})
}

const testAccDataSourceWorkspace_basic = `
resource "airbyte_workspace" "basic" {
  name = "basic_test"
}

data "airbyte_workspace" "by_id" {
  id = airbyte_workspace.basic.id
}

data "airbyte_workspace" "by_slug" {
  slug = airbyte_workspace.basic.slug
}
`

const testAccDataSourceWorkspace_complex = `
resource "airbyte_workspace" "complex" {
  name = "complex_test"
  email = "test@example.com"
  display_setup_wizard = true
  anonymous_data_collection = false
  news = true
  security_updates = true
  notification_config = [{
    notification_type = "slack"
    send_on_success = true
    slack_webhook = "http://example.com/webhook"
  }, {
    notification_type = "slack"
    send_on_failure = false
    slack_webhook = "https://example2.com/cooler-webhook"
  }]
}

data "airbyte_workspace" "by_id" {
  id = airbyte_workspace.complex.id
}

data "airbyte_workspace" "by_slug" {
  slug = airbyte_workspace.complex.slug
}
`
