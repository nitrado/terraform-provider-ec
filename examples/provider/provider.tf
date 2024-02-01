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
  host = "<your enterprise console host url>"
  instances {
    name = "my-other-instance"
    host = "<your enterprise console host url>"
  }
}

resource "ec_core_site" "test" {
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

resource "ec_core_site" "test2" {
  instance = "my-other-instance"
  metadata {
    name = "test"
  }

  spec {
    description = "My other test"
    credentials {
      endpoint    = "<your endpoint>"
      certificate = "<your cert>"
      namespace   = "<your ns>"
      token       = "<your token>"
    }
  }
}
