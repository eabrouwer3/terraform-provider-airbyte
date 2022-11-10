package provider

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceWorkspace_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		//CheckDestroy:             testAccResourceWorkspaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceWorkspace_basic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("airbyte_workspace.basic", "id", regexp.MustCompile("^[0-9a-fA-F]{8}\\b-[0-9a-fA-F]{4}\\b-[0-9a-fA-F]{4}\\b-[0-9a-fA-F]{4}\\b-[0-9a-fA-F]{12}$")),
					resource.TestMatchResourceAttr("airbyte_workspace.basic", "slug", regexp.MustCompile("^basic_test")),
					resource.TestCheckResourceAttr("airbyte_workspace.basic", "name", "basic_test"),
					resource.TestCheckResourceAttr("airbyte_workspace.basic", "notification_config.#", "0"),
				),
			},
		},
	})
}

func TestAccResourceWorkspace_complex(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceWorkspace_complex,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("airbyte_workspace.complex", "id", regexp.MustCompile("^[0-9a-fA-F]{8}\\b-[0-9a-fA-F]{4}\\b-[0-9a-fA-F]{4}\\b-[0-9a-fA-F]{4}\\b-[0-9a-fA-F]{12}$")),
					resource.TestMatchResourceAttr("airbyte_workspace.complex", "slug", regexp.MustCompile("^complex_test")),
					resource.TestCheckResourceAttr("airbyte_workspace.complex", "name", "complex_test"),
					resource.TestCheckResourceAttr("airbyte_workspace.complex", "email", "test@example.com"),
					resource.TestCheckResourceAttr("airbyte_workspace.complex", "display_setup_wizard", "true"),
					resource.TestCheckResourceAttr("airbyte_workspace.complex", "anonymous_data_collection", "false"),
					resource.TestCheckResourceAttr("airbyte_workspace.complex", "news", "true"),
					resource.TestCheckResourceAttr("airbyte_workspace.complex", "security_updates", "true"),
					resource.TestCheckResourceAttr("airbyte_workspace.complex", "notification_config.#", "2"),
					resource.TestCheckResourceAttr("airbyte_workspace.complex", "notification_config.0.notification_type", "slack"),
					resource.TestCheckResourceAttr("airbyte_workspace.complex", "notification_config.0.send_on_success", "true"),
					resource.TestCheckResourceAttr("airbyte_workspace.complex", "notification_config.0.send_on_failure", "true"),
					resource.TestCheckResourceAttr("airbyte_workspace.complex", "notification_config.0.slack_webhook", "http://example.com/webhook"),
					resource.TestCheckResourceAttr("airbyte_workspace.complex", "notification_config.1.notification_type", "slack"),
					resource.TestCheckResourceAttr("airbyte_workspace.complex", "notification_config.1.send_on_success", "false"),
					resource.TestCheckResourceAttr("airbyte_workspace.complex", "notification_config.1.send_on_failure", "false"),
					resource.TestCheckResourceAttr("airbyte_workspace.complex", "notification_config.1.slack_webhook", "https://example2.com/cooler-webhook"),
				),
			},
			{
				Config: testAccResourceWorkspace_complexChange,
				Check:  resource.TestCheckResourceAttr("airbyte_workspace.complex", "notification_config.#", "1"),
			},
		},
	})
}

// Don't know how to do this yet with the new framework...
//func testAccResourceWorkspaceDestroy(s *terraform.State) error {
//	client := testAccProvider.Meta().(*apiclient.ApiClient)
//
//	for _, rs := range s.RootModule().Resources {
//		if rs.Type == "airbyte_workspace" {
//			_, err := client.GetWorkspaceById(rs.Primary.ID)
//			if err == nil {
//				return fmt.Errorf("Workspace (%s) still exists.", rs.Primary.ID)
//			}
//
//			if !strings.Contains(err.Error(), "Could not find configuration for STANDARD_WORKSPACE") {
//				return err
//			}
//		}
//	}
//
//	return nil
//}

const testAccResourceWorkspace_basic = `
resource "airbyte_workspace" "basic" {
  name = "basic_test"
}
`

const testAccResourceWorkspace_complex = `
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
`

const testAccResourceWorkspace_complexChange = `
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
 }]
}
`
