provider "airbyte" {
  host_url = "http://localhost:8000"
  username = "airbyte"
  password = "password"
  additional_headers = {
    Host = "airbyte.internal"
  }
  timeout = 120
}
