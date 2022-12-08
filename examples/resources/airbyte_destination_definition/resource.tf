# Basic version of a custom destination definition
resource "airbyte_destination_definition" "basic" {
  workspace_id      = airbyte_workspace.test.id
  name              = "basic_test"
  docker_repository = "eabrouwer3/airbyte-test-data-destination"
  docker_image_tag  = "0.0.1"
  documentation_url = "https://github.com/eabrouwer3/airbyte-test-data-destination"
}

# Much more complex version of a custom destination definition
resource "airbyte_destination_definition" "complex" {
  workspace_id      = airbyte_workspace.test.id
  name              = "complex_test"
  docker_repository = "eabrouwer3/airbyte-test-data-destination"
  docker_image_tag  = "0.0.1"
  documentation_url = "https://github.com/eabrouwer3/airbyte-test-data-destination"

  default_resource_requirements = {
    cpu_request    = "0.25"
    cpu_limit      = "0.25"
    memory_request = "200Mi"
    memory_limit   = "200Mi"
  }

  job_specific_resource_requirements = [{
    job_type       = "sync"
    cpu_request    = "0.5"
    cpu_limit      = "0.5"
    memory_request = "500Mi"
    memory_limit   = "500Mi"
    }, {
    job_type       = "check_connection"
    memory_request = "50Mi"
    memory_limit   = "50Mi"
  }]
}
