resource "ec_armada_branch" "test" {
  metadata {
    name = "test"
  }

  spec {
    description = "My branch"
    display_name = "Test"
  }
}
