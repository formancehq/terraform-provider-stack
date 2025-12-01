<!-- Start SDK Example Usage [usage] -->
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