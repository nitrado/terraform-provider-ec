resource "ec_container_branch" "test" {
  metadata {
    name = "test"
  }

  spec {
    description = "My branch"
    display_name = "Test"

    retention_policy_rules {
      name = "default"
      keep_count = 10
      keep_days = 30
    }
  }
}
