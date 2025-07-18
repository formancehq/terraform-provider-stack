---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "stack_payments_connectors Resource - stack"
subcategory: ""
description: |-
  Resource for managing Formance Payments Connectors. For advanced usage and configuration, see the Payments Connectors documentation https://docs.formance.com/payments/connectors/.
---

# stack_payments_connectors (Resource)

Resource for managing Formance Payments Connectors. For advanced usage and configuration, see the [Payments Connectors documentation](https://docs.formance.com/payments/connectors/).



<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `config` (Dynamic) The configuration for the payment connector. It must not contain sensitive information like API keys or secrets. Advanced usage: See [Payments Connectors documentation](https://docs.formance.com/payments/connectors/) for connector configuration options.
- `credentials` (Dynamic, Sensitive) The credentials for the payment connector. This should include sensitive information like API keys, secrets, certificate, and must be handled securely. Advanced usage: See [Payments Connectors documentation](https://docs.formance.com/payments/connectors/) for connector security best practices.

### Read-Only

- `id` (String) The unique identifier of the payment connector.
