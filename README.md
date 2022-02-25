# Terraform Dependency Order Repro

## Context

This repo contains a dummy provider that defines 2 resources:
 * Notifier
 * Policy

A policy references a list of notifiers. A single notifier may be referenced
by multiple policies.

Policies must reference a notifier that exists, and notifiers cannot be deleted until
all references to the notifier are removed from policies.

The issue we demonstrate below is that removing a notifier (and its' policy reference)
in a single apply leads to the notifier being removed before the policy is updated
to remove the reference, resulting in an error.

In reality, the foreign key check is handled by the database. We use a simple
file-based store with the check implemented manually for a simple repro.
## Repro

To reproduce, check out this repo, cd into it, and run the following commands:
```
# ensure the repo is in a clean state
$ git clean -fdx

# install the provider to the local plugins directory
$ make install

# apply the initial state
$ make apply

# MANUAL: Update repro/main.tf - comment out "n2" and the reference to "n2" in "p1". e.g.,
```diff
-resource "tftest_notifier" "n2" {
-  email = "n2@example.com"
-}
+# resource "tftest_notifier" "n2" {
+# email = "n2@example.com"
+# }

 resource "tftest_policy" "p1" {
   notifier_ids = [
     tftest_notifier.n1.id,
-    tftest_notifier.n2.id,
+    # tftest_notifier.n2.id,
   ]
```

## apply will now fail, deleting n2 while it's referenced
```
$ make apply
[...]
tftest_notifier.n2: Destroying... [id=notifier-172842.521746]
╷
│ Error: failed to delete Notifier: cannot delete notifier, as policies [policy-172842.567744] still refer to it
```

This is the plan for the update:
```
  # tftest_notifier.n2 will be destroyed
  # (because tftest_notifier.n2 is not in configuration)
  - resource "tftest_notifier" "n2" {
      - email = "n2@example.com" -> null
      - id    = "notifier-172842.521746" -> null
    }

  # tftest_policy.p1 will be updated in-place
  ~ resource "tftest_policy" "p1" {
        id           = "policy-172842.567744"
      ~ notifier_ids = [
            "notifier-172842.545667",
          - "notifier-172842.521746",
        ]
    }
```

The provider tries to destroy before update.

If `lifecycle { create_before_destroy = true }` is set, the ordering destroys after the modification:
```
tftest_policy.p1: Modifying... [id=policy-172842.567744]
tftest_policy.p1: Modifications complete after 0s [id=policy-172842.567744]
tftest_notifier.n2: Destroying... [id=notifier-172842.521746]
tftest_notifier.n2: Destruction complete after 0s
```

However, this workaround has issues:
 * It requires users to manually specify the lifecycle on every single resource manually.
 * It doesn't work if the resource cannot be created first (e.g., due to duplicate checks).
