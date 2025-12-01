# V2Schema

Complete schema structure with metadata


## Fields

| Field                                                             | Type                                                              | Required                                                          | Description                                                       | Example                                                           |
| ----------------------------------------------------------------- | ----------------------------------------------------------------- | ----------------------------------------------------------------- | ----------------------------------------------------------------- | ----------------------------------------------------------------- |
| `CreatedAt`                                                       | [time.Time](https://pkg.go.dev/time#Time)                         | :heavy_check_mark:                                                | Schema creation timestamp                                         | 2023-01-01T00:00:00Z                                              |
| `Data`                                                            | [shared.V2SchemaData](../../../pkg/models/shared/v2schemadata.md) | :heavy_check_mark:                                                | Schema data structure for ledger schemas                          |                                                                   |
| `Version`                                                         | *string*                                                          | :heavy_check_mark:                                                | Schema version                                                    | v1.0.0                                                            |