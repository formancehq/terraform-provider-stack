# Formance Terraform Provider

## Working locally

Create a `~/.terraformrc` with the following content:

TODO:

- replace {WORKING_DIRECTORY}

```hcl
provider_installation {
  dev_overrides {
    "formancehq/stack" = "${WORKING_DIRECTORY}/build"
  }

  # For all other providers, install them directly from their origin provider
  # registries as normal. If you omit this, Terraform will _only_ use
  # the dev_overrides block, and so no other providers will be available.
  direct {}
}
```


## Docs

- [Publishig new versions](https://developer.hashicorp.com/terraform/registry/providers/publishing#signing-provider-releases).
