# terraform-provider-stack

Terraform provider for managing Formance Stack resources.

## Features
- Manage Formance Ledger, Payments, Webhooks, Reconciliation Policies, and more
- Supports advanced filtering and metadata queries
- Integrates with Formance Cloud and self-hosted stacks

## Requirements
- [Terraform](https://www.terraform.io/downloads.html) >= 1.0
- Formance Stack (Cloud or self-hosted)
- Go (for building from source)

## Installation

To install the provider, add the following to your Terraform configuration:

```hcl
terraform {
  required_providers {
    stack = {
      source  = "formancehq/stack"
      version = ">= 0.1.0"
    }
  }
}
```

## Provider Configuration

```hcl
provider "stack" {
  organization_id = "your-organization-id"
  stack_id        = "your-stack-id"
  uri             = "https://api.formance.cloud"
  # Optional:
  # cloud {
  #   client_id     = "..."
  #   client_secret = "..."
  #   endpoint      = "..."
  # }
  # retry_config { ... }
  # wait_module_duration = "2m"
}
```

See [docs/index.md](docs/index.md) for the full provider schema.

## Resources

- [Ledger](docs/resources/ledger.md) ([Ledger docs](https://docs.formance.com/ledger/))
- [Payments Pool](docs/resources/payments_pool.md) ([Payments docs](https://docs.formance.com/payments/))
- [Payments Connectors](docs/resources/payments_connectors.md) ([Payments Connectors docs](https://docs.formance.com/payments/connectors/))
- [Reconciliation Policy](docs/resources/reconciliation_policy.md) ([Reconciliation docs](https://docs.formance.com/reconciliation/))
- [Webhooks](docs/resources/webhooks.md) ([Webhooks docs](https://docs.formance.com/webhooks/))

## Advanced Usage & References

- **Ledger Advanced Filtering:** [Formance Ledger Filtering documentation](https://docs.formance.com/ledger/advanced/filtering)
- **Ledger Module Reference:** [Formance Ledger documentation](https://docs.formance.com/ledger/)
- **Payments Reference:** [Formance Payments documentation](https://docs.formance.com/payments/)
- **Payments Connectors Reference:** [Formance Payments Connectors documentation](https://docs.formance.com/payments/connectors/)
- **Reconciliation Reference:** [Formance Reconciliation documentation](https://docs.formance.com/reconciliation/)
- **Webhooks Reference:** [Formance Webhooks documentation](https://docs.formance.com/webhooks/)

## Contributing

Contributions are welcome! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

## License

This project is licensed under the Apache 2.0 License.
