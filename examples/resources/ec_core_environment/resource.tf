resource "ec_core_environment" "test" {
  metadata {
    name        = "test"
  }
  spec {
    description  = "Test environment"
    display_name = "Test Environment"
  }
}
