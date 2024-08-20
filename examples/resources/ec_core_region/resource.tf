resource "ec_core_region_v1" "test" {
  metadata {
    name        = "test"
    environment = "dev"
  }
  spec {
    description  = "Capacity in test region"
    display_name = "Test Region"

    types {
      name = "cloud"
      sites = [
        "test-cluster"
      ]
    }
  }
}
