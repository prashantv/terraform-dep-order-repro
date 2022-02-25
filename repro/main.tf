terraform {
  required_providers {
    tftest = {
      version = "0.1.0"
      source  = "local/prashantv/tftest"
    }
  }
}

resource "tftest_notifier" "n1" {
  email = "n1@example.com"
}

resource "tftest_notifier" "n2" {
  # Workaround: uncomment this and apply *before* deletion.
  # lifecycle { create_before_destroy = true }
  email = "n2@example.com"
}

resource "tftest_policy" "p1" {
  notifier_ids = [
    tftest_notifier.n1.id,
    tftest_notifier.n2.id,
  ]
}
