---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "stack_reconciliation_policy Resource - stack"
subcategory: ""
description: |-
  Resource for managing a Formance Reconciliation Policy. For advanced usage and configuration, see the Reconciliation documentation https://docs.formance.com/reconciliation/.
---

# stack_reconciliation_policy (Resource)

Resource for managing a Formance Reconciliation Policy. For advanced usage and configuration, see the [Reconciliation documentation](https://docs.formance.com/reconciliation/).



<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `ledger_name` (String) The name of the ledger associated with the reconciliation policy.
- `name` (String) The name of the pool.
- `payments_pool_id` (String) The ID of the payments pool associated with the reconciliation policy.

### Optional

- `ledger_query` (Dynamic) The ledger query used to filter transactions for reconciliation. It must be a valid JSON object representing a query Builder. Advanced usage: See [Ledger Advanced Filtering documentation](https://docs.formance.com/ledger/advanced/filtering) for more details.

### Read-Only

- `created_at` (String) The timestamp when the reconciliation policy was created.
- `id` (String) The unique identifier of the reconciliation policy.
