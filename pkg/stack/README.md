# undefined

Developer-friendly & type-safe Go SDK specifically catered to leverage *undefined* API.

<div align="left" style="margin-bottom: 0;">
    <a href="https://www.speakeasy.com/?utm_source=undefined&utm_campaign=go" class="badge-link">
        <span class="badge-container">
            <span class="badge-icon-section">
                <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 30 30" fill="none" style="vertical-align: middle;"><title>Speakeasy Logo</title><path fill="currentColor" d="m20.639 27.548-19.17-2.724L0 26.1l20.639 2.931 8.456-7.336-1.468-.208-6.988 6.062Z"></path><path fill="currentColor" d="m20.639 23.1 8.456-7.336-1.468-.207-6.988 6.06-6.84-.972-9.394-1.333-2.936-.417L0 20.169l2.937.416L0 23.132l20.639 2.931 8.456-7.334-1.468-.208-6.986 6.062-9.78-1.39 1.468-1.273 8.31 1.18Z"></path><path fill="currentColor" d="m20.639 18.65-19.17-2.724L0 17.201l20.639 2.931 8.456-7.334-1.468-.208-6.988 6.06Z"></path><path fill="currentColor" d="M27.627 6.658 24.69 9.205 20.64 12.72l-7.923-1.126L1.469 9.996 0 11.271l11.246 1.596-1.467 1.275-8.311-1.181L0 14.235l20.639 2.932 8.456-7.334-2.937-.418 2.937-2.549-1.468-.208Z"></path><path fill="currentColor" d="M29.095 3.902 8.456.971 0 8.305l20.639 2.934 8.456-7.337Z"></path></svg>
            </span>
            <span class="badge-text badge-text-section">BUILT BY SPEAKEASY</span>
        </span>
    </a>
    <a href="https://opensource.org/licenses/MIT" class="badge-link">
        <span class="badge-container blue">
            <span class="badge-text badge-text-section">LICENSE // MIT</span>
        </span>
    </a>
</div>


<br /><br />
> [!IMPORTANT]
> This SDK is not yet ready for production use. To complete setup please follow the steps outlined in your [workspace](https://app.speakeasy.com/org/formance/formance). Delete this section before > publishing to a package manager.

<!-- Start Summary [summary] -->
## Summary

Formance Stack API: Open, modular foundation for unique payments flows

# Introduction
This API is documented in **OpenAPI format**.

# Authentication
Formance Stack offers one forms of authentication:
  - OAuth2
OAuth2 - an open protocol to allow secure authorization in a simple
and standard method from web, mobile and desktop applications.
<SecurityDefinitions />
<!-- End Summary [summary] -->

<!-- Start Table of Contents [toc] -->
## Table of Contents
<!-- $toc-max-depth=2 -->
* [undefined](#undefined)
* [Introduction](#introduction)
* [Authentication](#authentication)
  * [SDK Installation](#sdk-installation)
  * [SDK Example Usage](#sdk-example-usage)
  * [Authentication](#authentication-1)
  * [Available Resources and Operations](#available-resources-and-operations)
  * [Retries](#retries)
  * [Error Handling](#error-handling)
  * [Server Selection](#server-selection)
  * [Custom HTTP Client](#custom-http-client)
* [Development](#development)
  * [Maturity](#maturity)
  * [Contributions](#contributions)

<!-- End Table of Contents [toc] -->

<!-- Start SDK Installation [installation] -->
## SDK Installation

To add the SDK as a dependency to your project:
```bash
go get github.com/formancehq/formance-sdk-go/v3
```
<!-- End SDK Installation [installation] -->

<!-- Start SDK Example Usage [usage] -->
## SDK Example Usage

### Example

```go
package main

import (
	"context"
	"github.com/formancehq/formance-sdk-go/v3"
	"github.com/formancehq/formance-sdk-go/v3/pkg/models/shared"
	"log"
)

func main() {
	ctx := context.Background()

	s := v3.New(
		v3.WithSecurity(shared.Security{
			ClientID:     "<YOUR_CLIENT_ID_HERE>",
			ClientSecret: "<YOUR_CLIENT_SECRET_HERE>",
			TokenURL:     "/api/auth/oauth/token",
		}),
	)

	res, err := s.GetVersions(ctx)
	if err != nil {
		log.Fatal(err)
	}
	if res.GetVersionsResponse != nil {
		// handle response
	}
}

```
<!-- End SDK Example Usage [usage] -->

<!-- Start Authentication [security] -->
## Authentication

### Per-Client Security Schemes

This SDK supports the following security scheme globally:

| Name                                         | Type   | Scheme                         |
| -------------------------------------------- | ------ | ------------------------------ |
| `ClientID`<br/>`ClientSecret`<br/>`TokenURL` | oauth2 | OAuth2 Client Credentials Flow |

You can configure it using the `WithSecurity` option when initializing the SDK client instance. For example:
```go
package main

import (
	"context"
	"github.com/formancehq/formance-sdk-go/v3"
	"github.com/formancehq/formance-sdk-go/v3/pkg/models/shared"
	"log"
)

func main() {
	ctx := context.Background()

	s := v3.New(
		v3.WithSecurity(shared.Security{
			ClientID:     "<YOUR_CLIENT_ID_HERE>",
			ClientSecret: "<YOUR_CLIENT_SECRET_HERE>",
			TokenURL:     "/api/auth/oauth/token",
		}),
	)

	res, err := s.GetVersions(ctx)
	if err != nil {
		log.Fatal(err)
	}
	if res.GetVersionsResponse != nil {
		// handle response
	}
}

```
<!-- End Authentication [security] -->

<!-- Start Available Resources and Operations [operations] -->
## Available Resources and Operations

<details open>
<summary>Available methods</summary>

#### [Auth.V1](docs/sdks/v1/README.md)

* [CreateClient](docs/sdks/v1/README.md#createclient) - Create client
* [CreateSecret](docs/sdks/v1/README.md#createsecret) - Add a secret to a client
* [DeleteClient](docs/sdks/v1/README.md#deleteclient) - Delete client
* [DeleteSecret](docs/sdks/v1/README.md#deletesecret) - Delete a secret from a client
* [GetOIDCWellKnowns](docs/sdks/v1/README.md#getoidcwellknowns) - Retrieve OpenID connect well-knowns.
* [GetServerInfo](docs/sdks/v1/README.md#getserverinfo) - Get server info
* [ListClients](docs/sdks/v1/README.md#listclients) - List clients
* [ListUsers](docs/sdks/v1/README.md#listusers) - List users
* [ReadClient](docs/sdks/v1/README.md#readclient) - Read client
* [ReadUser](docs/sdks/v1/README.md#readuser) - Read user
* [UpdateClient](docs/sdks/v1/README.md#updateclient) - Update client

### [Formance SDK](docs/sdks/formance/README.md)

* [GetVersions](docs/sdks/formance/README.md#getversions) - Show stack version information

### [Ledger](docs/sdks/ledger/README.md)

* [GetInfo](docs/sdks/ledger/README.md#getinfo) - Show server information
* [GetMetrics](docs/sdks/ledger/README.md#getmetrics) - Read in memory metrics

#### [Ledger.V1](docs/sdks/formancev1/README.md)

* [CreateTransactions](docs/sdks/formancev1/README.md#createtransactions) - Create a new batch of transactions to a ledger
* [AddMetadataOnTransaction](docs/sdks/formancev1/README.md#addmetadataontransaction) - Set the metadata of a transaction by its ID
* [AddMetadataToAccount](docs/sdks/formancev1/README.md#addmetadatatoaccount) - Add metadata to an account
* [CountAccounts](docs/sdks/formancev1/README.md#countaccounts) - Count the accounts from a ledger
* [CountTransactions](docs/sdks/formancev1/README.md#counttransactions) - Count the transactions from a ledger
* [CreateTransaction](docs/sdks/formancev1/README.md#createtransaction) - Create a new transaction to a ledger
* [GetAccount](docs/sdks/formancev1/README.md#getaccount) - Get account by its address
* [GetBalances](docs/sdks/formancev1/README.md#getbalances) - Get the balances from a ledger's account
* [GetBalancesAggregated](docs/sdks/formancev1/README.md#getbalancesaggregated) - Get the aggregated balances from selected accounts
* [GetInfo](docs/sdks/formancev1/README.md#getinfo) - Show server information
* [GetLedgerInfo](docs/sdks/formancev1/README.md#getledgerinfo) - Get information about a ledger
* [GetMapping](docs/sdks/formancev1/README.md#getmapping) - Get the mapping of a ledger
* [GetTransaction](docs/sdks/formancev1/README.md#gettransaction) - Get transaction from a ledger by its ID
* [ListAccounts](docs/sdks/formancev1/README.md#listaccounts) - List accounts from a ledger
* [ListLogs](docs/sdks/formancev1/README.md#listlogs) - List the logs from a ledger
* [ListTransactions](docs/sdks/formancev1/README.md#listtransactions) - List transactions from a ledger
* [ReadStats](docs/sdks/formancev1/README.md#readstats) - Get statistics from a ledger
* [RevertTransaction](docs/sdks/formancev1/README.md#reverttransaction) - Revert a ledger transaction by its ID
* [~~RunScript~~](docs/sdks/formancev1/README.md#runscript) - Execute a Numscript :warning: **Deprecated**
* [UpdateMapping](docs/sdks/formancev1/README.md#updatemapping) - Update the mapping of a ledger

#### [Ledger.V2](docs/sdks/v2/README.md)

* [AddMetadataOnTransaction](docs/sdks/v2/README.md#addmetadataontransaction) - Set the metadata of a transaction by its ID
* [AddMetadataToAccount](docs/sdks/v2/README.md#addmetadatatoaccount) - Add metadata to an account
* [CountAccounts](docs/sdks/v2/README.md#countaccounts) - Count the accounts from a ledger
* [CountTransactions](docs/sdks/v2/README.md#counttransactions) - Count the transactions from a ledger
* [CreateBulk](docs/sdks/v2/README.md#createbulk) - Bulk request
* [CreateExporter](docs/sdks/v2/README.md#createexporter) - Create exporter
* [CreateLedger](docs/sdks/v2/README.md#createledger) - Create a ledger
* [CreatePipeline](docs/sdks/v2/README.md#createpipeline) - Create pipeline
* [CreateTransaction](docs/sdks/v2/README.md#createtransaction) - Create a new transaction to a ledger
* [DeleteAccountMetadata](docs/sdks/v2/README.md#deleteaccountmetadata) - Delete metadata by key
* [DeleteBucket](docs/sdks/v2/README.md#deletebucket) - Delete bucket
* [DeleteExporter](docs/sdks/v2/README.md#deleteexporter) - Delete exporter
* [DeleteLedgerMetadata](docs/sdks/v2/README.md#deleteledgermetadata) - Delete ledger metadata by key
* [DeletePipeline](docs/sdks/v2/README.md#deletepipeline) - Delete pipeline
* [DeleteTransactionMetadata](docs/sdks/v2/README.md#deletetransactionmetadata) - Delete metadata by key
* [ExportLogs](docs/sdks/v2/README.md#exportlogs) - Export logs
* [GetAccount](docs/sdks/v2/README.md#getaccount) - Get account by its address
* [GetBalancesAggregated](docs/sdks/v2/README.md#getbalancesaggregated) - Get the aggregated balances from selected accounts
* [GetExporterState](docs/sdks/v2/README.md#getexporterstate) - Get exporter state
* [GetLedger](docs/sdks/v2/README.md#getledger) - Get a ledger
* [GetLedgerInfo](docs/sdks/v2/README.md#getledgerinfo) - Get information about a ledger
* [GetPipelineState](docs/sdks/v2/README.md#getpipelinestate) - Get pipeline state
* [GetSchema](docs/sdks/v2/README.md#getschema) - Get a schema for a ledger by version
* [GetTransaction](docs/sdks/v2/README.md#gettransaction) - Get transaction from a ledger by its ID
* [GetVolumesWithBalances](docs/sdks/v2/README.md#getvolumeswithbalances) - Get list of volumes with balances for (account/asset)
* [ImportLogs](docs/sdks/v2/README.md#importlogs)
* [InsertSchema](docs/sdks/v2/README.md#insertschema) - Insert or update a schema for a ledger
* [ListAccounts](docs/sdks/v2/README.md#listaccounts) - List accounts from a ledger
* [ListExporters](docs/sdks/v2/README.md#listexporters) - List exporters
* [ListLedgers](docs/sdks/v2/README.md#listledgers) - List ledgers
* [ListLogs](docs/sdks/v2/README.md#listlogs) - List the logs from a ledger
* [ListPipelines](docs/sdks/v2/README.md#listpipelines) - List pipelines
* [ListSchemas](docs/sdks/v2/README.md#listschemas) - List all schemas for a ledger
* [ListTransactions](docs/sdks/v2/README.md#listtransactions) - List transactions from a ledger
* [ReadStats](docs/sdks/v2/README.md#readstats) - Get statistics from a ledger
* [ResetPipeline](docs/sdks/v2/README.md#resetpipeline) - Reset pipeline
* [RestoreBucket](docs/sdks/v2/README.md#restorebucket) - Restore bucket
* [RevertTransaction](docs/sdks/v2/README.md#reverttransaction) - Revert a ledger transaction by its ID
* [StartPipeline](docs/sdks/v2/README.md#startpipeline) - Start pipeline
* [StopPipeline](docs/sdks/v2/README.md#stoppipeline) - Stop pipeline
* [UpdateExporter](docs/sdks/v2/README.md#updateexporter) - Update exporter
* [UpdateLedgerMetadata](docs/sdks/v2/README.md#updateledgermetadata) - Update ledger metadata

#### [Orchestration.V1](docs/sdks/formanceorchestrationv1/README.md)

* [CancelEvent](docs/sdks/formanceorchestrationv1/README.md#cancelevent) - Cancel a running workflow
* [CreateTrigger](docs/sdks/formanceorchestrationv1/README.md#createtrigger) - Create trigger
* [CreateWorkflow](docs/sdks/formanceorchestrationv1/README.md#createworkflow) - Create workflow
* [DeleteTrigger](docs/sdks/formanceorchestrationv1/README.md#deletetrigger) - Delete trigger
* [DeleteWorkflow](docs/sdks/formanceorchestrationv1/README.md#deleteworkflow) - Delete a flow by id
* [FlowsgetServerInfo](docs/sdks/formanceorchestrationv1/README.md#flowsgetserverinfo) - Get server info
* [GetInstance](docs/sdks/formanceorchestrationv1/README.md#getinstance) - Get a workflow instance by id
* [GetInstanceHistory](docs/sdks/formanceorchestrationv1/README.md#getinstancehistory) - Get a workflow instance history by id
* [GetInstanceStageHistory](docs/sdks/formanceorchestrationv1/README.md#getinstancestagehistory) - Get a workflow instance stage history
* [GetWorkflow](docs/sdks/formanceorchestrationv1/README.md#getworkflow) - Get a flow by id
* [ListInstances](docs/sdks/formanceorchestrationv1/README.md#listinstances) - List instances of a workflow
* [ListTriggers](docs/sdks/formanceorchestrationv1/README.md#listtriggers) - List triggers
* [ListTriggersOccurrences](docs/sdks/formanceorchestrationv1/README.md#listtriggersoccurrences) - List triggers occurrences
* [ListWorkflows](docs/sdks/formanceorchestrationv1/README.md#listworkflows) - List registered workflows
* [ReadTrigger](docs/sdks/formanceorchestrationv1/README.md#readtrigger) - Read trigger
* [RunWorkflow](docs/sdks/formanceorchestrationv1/README.md#runworkflow) - Run workflow
* [SendEvent](docs/sdks/formanceorchestrationv1/README.md#sendevent) - Send an event to a running workflow

#### [Orchestration.V2](docs/sdks/formancev2/README.md)

* [CancelEvent](docs/sdks/formancev2/README.md#cancelevent) - Cancel a running workflow
* [CreateTrigger](docs/sdks/formancev2/README.md#createtrigger) - Create trigger
* [CreateWorkflow](docs/sdks/formancev2/README.md#createworkflow) - Create workflow
* [DeleteTrigger](docs/sdks/formancev2/README.md#deletetrigger) - Delete trigger
* [DeleteWorkflow](docs/sdks/formancev2/README.md#deleteworkflow) - Delete a flow by id
* [GetInstance](docs/sdks/formancev2/README.md#getinstance) - Get a workflow instance by id
* [GetInstanceHistory](docs/sdks/formancev2/README.md#getinstancehistory) - Get a workflow instance history by id
* [GetInstanceStageHistory](docs/sdks/formancev2/README.md#getinstancestagehistory) - Get a workflow instance stage history
* [GetServerInfo](docs/sdks/formancev2/README.md#getserverinfo) - Get server info
* [GetWorkflow](docs/sdks/formancev2/README.md#getworkflow) - Get a flow by id
* [ListInstances](docs/sdks/formancev2/README.md#listinstances) - List instances of a workflow
* [ListTriggers](docs/sdks/formancev2/README.md#listtriggers) - List triggers
* [ListTriggersOccurrences](docs/sdks/formancev2/README.md#listtriggersoccurrences) - List triggers occurrences
* [ListWorkflows](docs/sdks/formancev2/README.md#listworkflows) - List registered workflows
* [ReadTrigger](docs/sdks/formancev2/README.md#readtrigger) - Read trigger
* [RunWorkflow](docs/sdks/formancev2/README.md#runworkflow) - Run workflow
* [SendEvent](docs/sdks/formancev2/README.md#sendevent) - Send an event to a running workflow
* [TestTrigger](docs/sdks/formancev2/README.md#testtrigger) - Test trigger

#### [Payments.V1](docs/sdks/formancepaymentsv1/README.md)

* [AddAccountToPool](docs/sdks/formancepaymentsv1/README.md#addaccounttopool) - Add an account to a pool
* [ConnectorsTransfer](docs/sdks/formancepaymentsv1/README.md#connectorstransfer) - Transfer funds between Connector accounts
* [CreateAccount](docs/sdks/formancepaymentsv1/README.md#createaccount) - Create an account
* [CreateBankAccount](docs/sdks/formancepaymentsv1/README.md#createbankaccount) - Create a BankAccount in Payments and on the PSP
* [CreatePayment](docs/sdks/formancepaymentsv1/README.md#createpayment) - Create a payment
* [CreatePool](docs/sdks/formancepaymentsv1/README.md#createpool) - Create a Pool
* [CreateTransferInitiation](docs/sdks/formancepaymentsv1/README.md#createtransferinitiation) - Create a TransferInitiation
* [DeletePool](docs/sdks/formancepaymentsv1/README.md#deletepool) - Delete a Pool
* [DeleteTransferInitiation](docs/sdks/formancepaymentsv1/README.md#deletetransferinitiation) - Delete a transfer initiation
* [ForwardBankAccount](docs/sdks/formancepaymentsv1/README.md#forwardbankaccount) - Forward a bank account to a connector
* [GetAccountBalances](docs/sdks/formancepaymentsv1/README.md#getaccountbalances) - Get account balances
* [GetBankAccount](docs/sdks/formancepaymentsv1/README.md#getbankaccount) - Get a bank account created by user on Formance
* [~~GetConnectorTask~~](docs/sdks/formancepaymentsv1/README.md#getconnectortask) - Read a specific task of the connector :warning: **Deprecated**
* [GetConnectorTaskV1](docs/sdks/formancepaymentsv1/README.md#getconnectortaskv1) - Read a specific task of the connector
* [GetPayment](docs/sdks/formancepaymentsv1/README.md#getpayment) - Get a payment
* [GetPool](docs/sdks/formancepaymentsv1/README.md#getpool) - Get a Pool
* [GetPoolBalances](docs/sdks/formancepaymentsv1/README.md#getpoolbalances) - Get historical pool balances at a particular point in time
* [GetPoolBalancesLatest](docs/sdks/formancepaymentsv1/README.md#getpoolbalanceslatest) - Get latest pool balances
* [GetTransferInitiation](docs/sdks/formancepaymentsv1/README.md#gettransferinitiation) - Get a transfer initiation
* [InstallConnector](docs/sdks/formancepaymentsv1/README.md#installconnector) - Install a connector
* [ListAllConnectors](docs/sdks/formancepaymentsv1/README.md#listallconnectors) - List all installed connectors
* [ListBankAccounts](docs/sdks/formancepaymentsv1/README.md#listbankaccounts) - List bank accounts created by user on Formance
* [ListConfigsAvailableConnectors](docs/sdks/formancepaymentsv1/README.md#listconfigsavailableconnectors) - List the configs of each available connector
* [~~ListConnectorTasks~~](docs/sdks/formancepaymentsv1/README.md#listconnectortasks) - List tasks from a connector :warning: **Deprecated**
* [ListConnectorTasksV1](docs/sdks/formancepaymentsv1/README.md#listconnectortasksv1) - List tasks from a connector
* [ListPayments](docs/sdks/formancepaymentsv1/README.md#listpayments) - List payments
* [ListPools](docs/sdks/formancepaymentsv1/README.md#listpools) - List Pools
* [ListTransferInitiations](docs/sdks/formancepaymentsv1/README.md#listtransferinitiations) - List Transfer Initiations
* [PaymentsgetAccount](docs/sdks/formancepaymentsv1/README.md#paymentsgetaccount) - Get an account
* [PaymentsgetServerInfo](docs/sdks/formancepaymentsv1/README.md#paymentsgetserverinfo) - Get server info
* [PaymentslistAccounts](docs/sdks/formancepaymentsv1/README.md#paymentslistaccounts) - List accounts
* [~~ReadConnectorConfig~~](docs/sdks/formancepaymentsv1/README.md#readconnectorconfig) - Read the config of a connector :warning: **Deprecated**
* [ReadConnectorConfigV1](docs/sdks/formancepaymentsv1/README.md#readconnectorconfigv1) - Read the config of a connector
* [RemoveAccountFromPool](docs/sdks/formancepaymentsv1/README.md#removeaccountfrompool) - Remove an account from a pool
* [~~ResetConnector~~](docs/sdks/formancepaymentsv1/README.md#resetconnector) - Reset a connector :warning: **Deprecated**
* [ResetConnectorV1](docs/sdks/formancepaymentsv1/README.md#resetconnectorv1) - Reset a connector
* [RetryTransferInitiation](docs/sdks/formancepaymentsv1/README.md#retrytransferinitiation) - Retry a failed transfer initiation
* [ReverseTransferInitiation](docs/sdks/formancepaymentsv1/README.md#reversetransferinitiation) - Reverse a transfer initiation
* [~~UninstallConnector~~](docs/sdks/formancepaymentsv1/README.md#uninstallconnector) - Uninstall a connector :warning: **Deprecated**
* [UninstallConnectorV1](docs/sdks/formancepaymentsv1/README.md#uninstallconnectorv1) - Uninstall a connector
* [UpdateBankAccountMetadata](docs/sdks/formancepaymentsv1/README.md#updatebankaccountmetadata) - Update metadata of a bank account
* [UpdateConnectorConfigV1](docs/sdks/formancepaymentsv1/README.md#updateconnectorconfigv1) - Update the config of a connector
* [UpdateMetadata](docs/sdks/formancepaymentsv1/README.md#updatemetadata) - Update metadata
* [UpdatePoolQuery](docs/sdks/formancepaymentsv1/README.md#updatepoolquery) - Update the query of a pool
* [UpdateTransferInitiationStatus](docs/sdks/formancepaymentsv1/README.md#updatetransferinitiationstatus) - Update the status of a transfer initiation

#### [Payments.V3](docs/sdks/v3/README.md)

* [AddAccountToPool](docs/sdks/v3/README.md#addaccounttopool) - Add an account to a pool
* [AddBankAccountToPaymentServiceUser](docs/sdks/v3/README.md#addbankaccounttopaymentserviceuser) - Add a bank account to a payment service user
* [ApprovePaymentInitiation](docs/sdks/v3/README.md#approvepaymentinitiation) - Approve a payment initiation
* [CreateAccount](docs/sdks/v3/README.md#createaccount) - Create a formance account object. This object will not be forwarded to the connector. It is only used for internal purposes.

* [CreateBankAccount](docs/sdks/v3/README.md#createbankaccount) - Create a formance bank account object. This object will not be forwarded to the connector until you called the forwardBankAccount method.

* [CreateLinkForPaymentServiceUser](docs/sdks/v3/README.md#createlinkforpaymentserviceuser) - Create an authentication link for a payment service user on a connector, for oauth flow
* [CreatePayment](docs/sdks/v3/README.md#createpayment) - Create a formance payment object. This object will not be forwarded to the connector. It is only used for internal purposes.

* [CreatePaymentServiceUser](docs/sdks/v3/README.md#createpaymentserviceuser) - Create a formance payment service user object
* [CreatePool](docs/sdks/v3/README.md#createpool) - Create a formance pool object
* [DeletePaymentInitiation](docs/sdks/v3/README.md#deletepaymentinitiation) - Delete a payment initiation by ID
* [DeletePaymentServiceUser](docs/sdks/v3/README.md#deletepaymentserviceuser) - Delete a payment service user by ID
* [DeletePaymentServiceUserConnectionFromConnectorID](docs/sdks/v3/README.md#deletepaymentserviceuserconnectionfromconnectorid) - Delete a connection for a payment service user on a connector
* [DeletePaymentServiceUserConnector](docs/sdks/v3/README.md#deletepaymentserviceuserconnector) - Remove a payment service user from a connector, the PSU will still exist in Formance
* [DeletePool](docs/sdks/v3/README.md#deletepool) - Delete a pool by ID
* [ForwardBankAccount](docs/sdks/v3/README.md#forwardbankaccount) - Forward a Bank Account to a PSP for creation
* [ForwardPaymentServiceUserBankAccount](docs/sdks/v3/README.md#forwardpaymentserviceuserbankaccount) - Forward a payment service user's bank account to a connector
* [ForwardPaymentServiceUserToProvider](docs/sdks/v3/README.md#forwardpaymentserviceusertoprovider) - Register/forward a payment service user on/to a connector
* [GetAccount](docs/sdks/v3/README.md#getaccount) - Get an account by ID
* [GetAccountBalances](docs/sdks/v3/README.md#getaccountbalances) - Get account balances
* [GetBankAccount](docs/sdks/v3/README.md#getbankaccount) - Get a Bank Account by ID
* [GetConnectorConfig](docs/sdks/v3/README.md#getconnectorconfig) - Get a connector configuration by ID
* [GetConnectorSchedule](docs/sdks/v3/README.md#getconnectorschedule) - Get a connector schedule by ID
* [GetPayment](docs/sdks/v3/README.md#getpayment) - Get a payment by ID
* [GetPaymentInitiation](docs/sdks/v3/README.md#getpaymentinitiation) - Get a payment initiation by ID
* [GetPaymentServiceUser](docs/sdks/v3/README.md#getpaymentserviceuser) - Get a payment service user by ID
* [GetPaymentServiceUserLinkAttemptFromConnectorID](docs/sdks/v3/README.md#getpaymentserviceuserlinkattemptfromconnectorid) - Get a link attempt for a payment service user on a connector
* [GetPool](docs/sdks/v3/README.md#getpool) - Get a pool by ID
* [GetPoolBalances](docs/sdks/v3/README.md#getpoolbalances) - Get historical pool balances from a particular point in time
* [GetPoolBalancesLatest](docs/sdks/v3/README.md#getpoolbalanceslatest) - Get latest pool balances
* [GetTask](docs/sdks/v3/README.md#gettask) - Get a task and its result by ID
* [InitiatePayment](docs/sdks/v3/README.md#initiatepayment) - Initiate a payment
* [InstallConnector](docs/sdks/v3/README.md#installconnector) - Install a connector
* [ListAccounts](docs/sdks/v3/README.md#listaccounts) - List all accounts
* [ListBankAccounts](docs/sdks/v3/README.md#listbankaccounts) - List all bank accounts
* [ListConnectorConfigs](docs/sdks/v3/README.md#listconnectorconfigs) - List all connector configurations
* [ListConnectorScheduleInstances](docs/sdks/v3/README.md#listconnectorscheduleinstances) - List all connector schedule instances
* [ListConnectorSchedules](docs/sdks/v3/README.md#listconnectorschedules) - List all connector schedules
* [ListConnectors](docs/sdks/v3/README.md#listconnectors) - List all connectors
* [ListPaymentInitiationAdjustments](docs/sdks/v3/README.md#listpaymentinitiationadjustments) - List all payment initiation adjustments
* [ListPaymentInitiationRelatedPayments](docs/sdks/v3/README.md#listpaymentinitiationrelatedpayments) - List all payments related to a payment initiation
* [ListPaymentInitiations](docs/sdks/v3/README.md#listpaymentinitiations) - List all payment initiations
* [ListPaymentServiceUserConnections](docs/sdks/v3/README.md#listpaymentserviceuserconnections) - List all connections for a payment service user
* [ListPaymentServiceUserConnectionsFromConnectorID](docs/sdks/v3/README.md#listpaymentserviceuserconnectionsfromconnectorid) - List enabled connections for a payment service user on a connector (i.e. the various banks PSUser has enabled on the connector)
* [ListPaymentServiceUserLinkAttemptsFromConnectorID](docs/sdks/v3/README.md#listpaymentserviceuserlinkattemptsfromconnectorid) - List all link attempts for a payment service user on a connector.
Allows to check if users used the link and completed the oauth flow.

* [ListPaymentServiceUsers](docs/sdks/v3/README.md#listpaymentserviceusers) - List all payment service users
* [ListPayments](docs/sdks/v3/README.md#listpayments) - List all payments
* [ListPools](docs/sdks/v3/README.md#listpools) - List all pools
* [RejectPaymentInitiation](docs/sdks/v3/README.md#rejectpaymentinitiation) - Reject a payment initiation
* [RemoveAccountFromPool](docs/sdks/v3/README.md#removeaccountfrompool) - Remove an account from a pool
* [ResetConnector](docs/sdks/v3/README.md#resetconnector) - Reset a connector. Be aware that this will delete all data and stop all existing tasks like payment initiations and bank account creations.
* [RetryPaymentInitiation](docs/sdks/v3/README.md#retrypaymentinitiation) - Retry a payment initiation
* [ReversePaymentInitiation](docs/sdks/v3/README.md#reversepaymentinitiation) - Reverse a payment initiation
* [UninstallConnector](docs/sdks/v3/README.md#uninstallconnector) - Uninstall a connector
* [UpdateBankAccountMetadata](docs/sdks/v3/README.md#updatebankaccountmetadata) - Update a bank account's metadata
* [UpdateLinkForPaymentServiceUserOnConnector](docs/sdks/v3/README.md#updatelinkforpaymentserviceuseronconnector) - Update/Regenerate a link for a payment service user on a connector
* [UpdatePaymentMetadata](docs/sdks/v3/README.md#updatepaymentmetadata) - Update a payment's metadata
* [UpdatePoolQuery](docs/sdks/v3/README.md#updatepoolquery) - Update the query of a pool
* [V3UpdateConnectorConfig](docs/sdks/v3/README.md#v3updateconnectorconfig) - Update the config of a connector

#### [Reconciliation.V1](docs/sdks/formancereconciliationv1/README.md)

* [CreatePolicy](docs/sdks/formancereconciliationv1/README.md#createpolicy) - Create a policy
* [DeletePolicy](docs/sdks/formancereconciliationv1/README.md#deletepolicy) - Delete a policy
* [GetPolicy](docs/sdks/formancereconciliationv1/README.md#getpolicy) - Get a policy
* [GetReconciliation](docs/sdks/formancereconciliationv1/README.md#getreconciliation) - Get a reconciliation
* [ListPolicies](docs/sdks/formancereconciliationv1/README.md#listpolicies) - List policies
* [ListReconciliations](docs/sdks/formancereconciliationv1/README.md#listreconciliations) - List reconciliations
* [Reconcile](docs/sdks/formancereconciliationv1/README.md#reconcile) - Reconcile using a policy
* [ReconciliationgetServerInfo](docs/sdks/formancereconciliationv1/README.md#reconciliationgetserverinfo) - Get server info

#### [Wallets.V1](docs/sdks/formancewalletsv1/README.md)

* [ConfirmHold](docs/sdks/formancewalletsv1/README.md#confirmhold) - Confirm a hold
* [CreateBalance](docs/sdks/formancewalletsv1/README.md#createbalance) - Create a balance
* [CreateWallet](docs/sdks/formancewalletsv1/README.md#createwallet) - Create a new wallet
* [CreditWallet](docs/sdks/formancewalletsv1/README.md#creditwallet) - Credit a wallet
* [DebitWallet](docs/sdks/formancewalletsv1/README.md#debitwallet) - Debit a wallet
* [GetBalance](docs/sdks/formancewalletsv1/README.md#getbalance) - Get detailed balance
* [GetHold](docs/sdks/formancewalletsv1/README.md#gethold) - Get a hold
* [GetHolds](docs/sdks/formancewalletsv1/README.md#getholds) - Get all holds for a wallet
* [GetTransactions](docs/sdks/formancewalletsv1/README.md#gettransactions)
* [GetWallet](docs/sdks/formancewalletsv1/README.md#getwallet) - Get a wallet
* [GetWalletSummary](docs/sdks/formancewalletsv1/README.md#getwalletsummary) - Get wallet summary
* [ListBalances](docs/sdks/formancewalletsv1/README.md#listbalances) - List balances of a wallet
* [ListWallets](docs/sdks/formancewalletsv1/README.md#listwallets) - List all wallets
* [UpdateWallet](docs/sdks/formancewalletsv1/README.md#updatewallet) - Update a wallet
* [VoidHold](docs/sdks/formancewalletsv1/README.md#voidhold) - Cancel a hold
* [WalletsgetServerInfo](docs/sdks/formancewalletsv1/README.md#walletsgetserverinfo) - Get server info

#### [Webhooks.V1](docs/sdks/formancewebhooksv1/README.md)

* [ActivateConfig](docs/sdks/formancewebhooksv1/README.md#activateconfig) - Activate one config
* [ChangeConfigSecret](docs/sdks/formancewebhooksv1/README.md#changeconfigsecret) - Change the signing secret of a config
* [DeactivateConfig](docs/sdks/formancewebhooksv1/README.md#deactivateconfig) - Deactivate one config
* [DeleteConfig](docs/sdks/formancewebhooksv1/README.md#deleteconfig) - Delete one config
* [GetManyConfigs](docs/sdks/formancewebhooksv1/README.md#getmanyconfigs) - Get many configs
* [InsertConfig](docs/sdks/formancewebhooksv1/README.md#insertconfig) - Insert a new config
* [TestConfig](docs/sdks/formancewebhooksv1/README.md#testconfig) - Test one config
* [UpdateConfig](docs/sdks/formancewebhooksv1/README.md#updateconfig) - Update one config

</details>
<!-- End Available Resources and Operations [operations] -->

<!-- Start Retries [retries] -->
## Retries

Some of the endpoints in this SDK support retries. If you use the SDK without any configuration, it will fall back to the default retry strategy provided by the API. However, the default retry strategy can be overridden on a per-operation basis, or across the entire SDK.

To change the default retry strategy for a single API call, simply provide a `retry.Config` object to the call by using the `WithRetries` option:
```go
package main

import (
	"context"
	"github.com/formancehq/formance-sdk-go/v3"
	"github.com/formancehq/formance-sdk-go/v3/pkg/models/shared"
	"github.com/formancehq/formance-sdk-go/v3/pkg/retry"
	"log"
	"pkg/models/operations"
)

func main() {
	ctx := context.Background()

	s := v3.New(
		v3.WithSecurity(shared.Security{
			ClientID:     "<YOUR_CLIENT_ID_HERE>",
			ClientSecret: "<YOUR_CLIENT_SECRET_HERE>",
			TokenURL:     "/api/auth/oauth/token",
		}),
	)

	res, err := s.GetVersions(ctx, operations.WithRetries(
		retry.Config{
			Strategy: "backoff",
			Backoff: &retry.BackoffStrategy{
				InitialInterval: 1,
				MaxInterval:     50,
				Exponent:        1.1,
				MaxElapsedTime:  100,
			},
			RetryConnectionErrors: false,
		}))
	if err != nil {
		log.Fatal(err)
	}
	if res.GetVersionsResponse != nil {
		// handle response
	}
}

```

If you'd like to override the default retry strategy for all operations that support retries, you can use the `WithRetryConfig` option at SDK initialization:
```go
package main

import (
	"context"
	"github.com/formancehq/formance-sdk-go/v3"
	"github.com/formancehq/formance-sdk-go/v3/pkg/models/shared"
	"github.com/formancehq/formance-sdk-go/v3/pkg/retry"
	"log"
)

func main() {
	ctx := context.Background()

	s := v3.New(
		v3.WithRetryConfig(
			retry.Config{
				Strategy: "backoff",
				Backoff: &retry.BackoffStrategy{
					InitialInterval: 1,
					MaxInterval:     50,
					Exponent:        1.1,
					MaxElapsedTime:  100,
				},
				RetryConnectionErrors: false,
			}),
		v3.WithSecurity(shared.Security{
			ClientID:     "<YOUR_CLIENT_ID_HERE>",
			ClientSecret: "<YOUR_CLIENT_SECRET_HERE>",
			TokenURL:     "/api/auth/oauth/token",
		}),
	)

	res, err := s.GetVersions(ctx)
	if err != nil {
		log.Fatal(err)
	}
	if res.GetVersionsResponse != nil {
		// handle response
	}
}

```
<!-- End Retries [retries] -->

<!-- Start Error Handling [errors] -->
## Error Handling

Handling errors in this SDK should largely match your expectations. All operations return a response object or an error, they will never return both.

By Default, an API error will return `sdkerrors.SDKError`. When custom error responses are specified for an operation, the SDK may also return their associated error. You can refer to respective *Errors* tables in SDK docs for more details on possible error types for each operation.

For example, the `GetInfo` function may return the following errors:

| Error Type                | Status Code | Content Type     |
| ------------------------- | ----------- | ---------------- |
| sdkerrors.V2ErrorResponse | 5XX         | application/json |
| sdkerrors.SDKError        | 4XX         | \*/\*            |

### Example

```go
package main

import (
	"context"
	"errors"
	"github.com/formancehq/formance-sdk-go/v3"
	"github.com/formancehq/formance-sdk-go/v3/pkg/models/sdkerrors"
	"github.com/formancehq/formance-sdk-go/v3/pkg/models/shared"
	"log"
)

func main() {
	ctx := context.Background()

	s := v3.New(
		v3.WithSecurity(shared.Security{
			ClientID:     "<YOUR_CLIENT_ID_HERE>",
			ClientSecret: "<YOUR_CLIENT_SECRET_HERE>",
			TokenURL:     "/api/auth/oauth/token",
		}),
	)

	res, err := s.Ledger.GetInfo(ctx)
	if err != nil {

		var e *sdkerrors.V2ErrorResponse
		if errors.As(err, &e) {
			// handle error
			log.Fatal(e.Error())
		}

		var e *sdkerrors.SDKError
		if errors.As(err, &e) {
			// handle error
			log.Fatal(e.Error())
		}
	}
}

```
<!-- End Error Handling [errors] -->

<!-- Start Server Selection [server] -->
## Server Selection

### Select Server by Index

You can override the default server globally using the `WithServerIndex(serverIndex int)` option when initializing the SDK client instance. The selected server will then be used as the default on the operations that use it. This table lists the indexes associated with the available servers:

| #   | Server                                                | Variables                        | Description                                |
| --- | ----------------------------------------------------- | -------------------------------- | ------------------------------------------ |
| 0   | `http://localhost`                                    |                                  | local server                               |
| 1   | `https://{organization}.{environment}.formance.cloud` | `environment`<br/>`organization` | A per-organization and per-environment API |

If the selected server has variables, you may override its default values using the associated option(s):

| Variable       | Option                                           | Supported Values                                         | Default           | Description                                                   |
| -------------- | ------------------------------------------------ | -------------------------------------------------------- | ----------------- | ------------------------------------------------------------- |
| `environment`  | `WithEnvironment(environment ServerEnvironment)` | - `"eu.sandbox"`<br/>- `"eu-west-1"`<br/>- `"us-east-1"` | `"eu.sandbox"`    | The environment name. Defaults to the production environment. |
| `organization` | `WithOrganization(organization string)`          | string                                                   | `"orgID-stackID"` | The organization name. Defaults to a generic organization.    |

#### Example

```go
package main

import (
	"context"
	"github.com/formancehq/formance-sdk-go/v3"
	"github.com/formancehq/formance-sdk-go/v3/pkg/models/shared"
	"log"
)

func main() {
	ctx := context.Background()

	s := v3.New(
		v3.WithServerIndex(1),
		v3.WithEnvironment("us-east-1"),
		v3.WithOrganization("<value>"),
		v3.WithSecurity(shared.Security{
			ClientID:     "<YOUR_CLIENT_ID_HERE>",
			ClientSecret: "<YOUR_CLIENT_SECRET_HERE>",
			TokenURL:     "/api/auth/oauth/token",
		}),
	)

	res, err := s.GetVersions(ctx)
	if err != nil {
		log.Fatal(err)
	}
	if res.GetVersionsResponse != nil {
		// handle response
	}
}

```

### Override Server URL Per-Client

The default server can also be overridden globally using the `WithServerURL(serverURL string)` option when initializing the SDK client instance. For example:
```go
package main

import (
	"context"
	"github.com/formancehq/formance-sdk-go/v3"
	"github.com/formancehq/formance-sdk-go/v3/pkg/models/shared"
	"log"
)

func main() {
	ctx := context.Background()

	s := v3.New(
		v3.WithServerURL("https://orgID-stackID.eu.sandbox.formance.cloud"),
		v3.WithSecurity(shared.Security{
			ClientID:     "<YOUR_CLIENT_ID_HERE>",
			ClientSecret: "<YOUR_CLIENT_SECRET_HERE>",
			TokenURL:     "/api/auth/oauth/token",
		}),
	)

	res, err := s.GetVersions(ctx)
	if err != nil {
		log.Fatal(err)
	}
	if res.GetVersionsResponse != nil {
		// handle response
	}
}

```
<!-- End Server Selection [server] -->

<!-- Start Custom HTTP Client [http-client] -->
## Custom HTTP Client

The Go SDK makes API calls that wrap an internal HTTP client. The requirements for the HTTP client are very simple. It must match this interface:

```go
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}
```

The built-in `net/http` client satisfies this interface and a default client based on the built-in is provided by default. To replace this default with a client of your own, you can implement this interface yourself or provide your own client configured as desired. Here's a simple example, which adds a client with a 30 second timeout.

```go
import (
	"net/http"
	"time"

	"github.com/formancehq/formance-sdk-go/v3"
)

var (
	httpClient = &http.Client{Timeout: 30 * time.Second}
	sdkClient  = v3.New(v3.WithClient(httpClient))
)
```

This can be a convenient way to configure timeouts, cookies, proxies, custom headers, and other low-level configuration.
<!-- End Custom HTTP Client [http-client] -->

<!-- Placeholder for Future Speakeasy SDK Sections -->

# Development

## Maturity

This SDK is in beta, and there may be breaking changes between versions without a major version update. Therefore, we recommend pinning usage
to a specific package version. This way, you can install the same version each time without breaking changes unless you are intentionally
looking for the latest version.

## Contributions

While we value open-source contributions to this SDK, this library is generated programmatically. Any manual changes added to internal files will be overwritten on the next generation. 
We look forward to hearing your feedback. Feel free to open a PR or an issue with a proof of concept and we'll do our best to include it in a future release. 

### SDK Created by [Speakeasy](https://www.speakeasy.com/?utm_source=undefined&utm_campaign=go)

<style>
  :root {
    --badge-gray-bg: #f3f4f6;
    --badge-gray-border: #d1d5db;
    --badge-gray-text: #374151;
    --badge-blue-bg: #eff6ff;
    --badge-blue-border: #3b82f6;
    --badge-blue-text: #3b82f6;
  }

  @media (prefers-color-scheme: dark) {
    :root {
      --badge-gray-bg: #374151;
      --badge-gray-border: #4b5563;
      --badge-gray-text: #f3f4f6;
      --badge-blue-bg: #1e3a8a;
      --badge-blue-border: #3b82f6;
      --badge-blue-text: #93c5fd;
    }
  }
  
  h1 {
    border-bottom: none !important;
    margin-bottom: 4px;
    margin-top: 0;
    letter-spacing: 0.5px;
    font-weight: 600;
  }
  
  .badge-text {
    letter-spacing: 1px;
    font-weight: 300;
  }
  
  .badge-container {
    display: inline-flex;
    align-items: center;
    background: var(--badge-gray-bg);
    border: 1px solid var(--badge-gray-border);
    border-radius: 6px;
    overflow: hidden;
    font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Helvetica, Arial, sans-serif;
    font-size: 11px;
    text-decoration: none;
    vertical-align: middle;
  }

  .badge-container.blue {
    background: var(--badge-blue-bg);
    border-color: var(--badge-blue-border);
  }

  .badge-icon-section {
    padding: 4px 8px;
    border-right: 1px solid var(--badge-gray-border);
    display: flex;
    align-items: center;
  }

  .badge-text-section {
    padding: 4px 10px;
    color: var(--badge-gray-text);
    font-weight: 400;
  }

  .badge-container.blue .badge-text-section {
    color: var(--badge-blue-text);
  }
  
  .badge-link {
    text-decoration: none;
    margin-left: 8px;
    display: inline-flex;
    vertical-align: middle;
  }

  .badge-link:hover {
    text-decoration: none;
  }
  
  .badge-link:first-child {
    margin-left: 0;
  }
  
  .badge-icon-section svg {
    color: var(--badge-gray-text);
  }

  .badge-container.blue .badge-icon-section svg {
    color: var(--badge-blue-text);
  }
</style> 