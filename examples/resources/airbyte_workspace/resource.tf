# Most basic version of a workspace you can create
resource "airbyte_workspace" "basic" {
  name = "basic_test"
}

# Much more complex version of a workspace
resource "airbyte_workspace" "complex" {
  name                      = "complex_test"
  email                     = "test@example.com"
  display_setup_wizard      = true
  anonymous_data_collection = false
  news                      = true
  security_updates          = true
  notification_config = [{
    notification_type = "slack"
    send_on_success   = true
    slack_webhook     = "http://example.com/webhook"
    }, {
    notification_type = "slack"
    send_on_failure   = false
    slack_webhook     = "https://example2.com/cooler-webhook"
  }]
}