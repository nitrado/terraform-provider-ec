resource "ec_container_branch_v1" "test" {
  metadata {
    name = "test"
  }

  spec {
    description = "My branch"
    display_name = "Test"
  }
}
