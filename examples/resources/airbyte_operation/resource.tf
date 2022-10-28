resource "airbyte_workspace" "test" {
  name = "basic_test"
}

resource "airbyte_operation" "normalization" {
  workspace_id         = airbyte_workspace.test.id
  name                 = "normalization_operation"
  operator_type        = "normalization"
  normalization_option = "basic"
}

resource "airbyte_operation" "dbt" {
  workspace_id  = airbyte_workspace.test.id
  name          = "dbt_operation"
  operator_type = "dbt"
  dbt = {
    git_repo_url    = ""
    git_repo_branch = ""
    docker_image    = ""
    dbt_arguments   = ""
  }
}

resource "airbyte_operation" "webhook" {
  workspace_id  = airbyte_workspace.test.id
  name          = "webhook_operation"
  operator_type = "webhook"
  webhook = {
    execution_url     = ""
    execution_body    = ""
    webhook_config_id = ""
  }
}