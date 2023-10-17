resource "ec_armada_site" "test" {
  metadata {
    name = "test"
  }

  spec {
    description = "My test"
    credentials {
      endpoint    = "<your endpoint>"
      certificate = "<your cert>"
      namespace   = "<your ns>"
      token       = "<your token>"
    }
  }
}
