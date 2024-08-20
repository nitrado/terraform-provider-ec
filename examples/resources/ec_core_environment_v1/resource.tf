resource "ec_core_environment_v1" "test" {
  metadata {
    name        = "test"
  }
  spec {
    description  = "Test environment"
    display_name = "Test Environment"
  }
}
