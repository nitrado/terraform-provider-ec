terraform {
  required_providers {
    ec = {
      source  = "nitrado/ec"
      version = ">=1.0.0"
    }
  }
  required_version = ">= 0.14"
}

provider "ec" {
  host = "<your armada host url>"
}

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

