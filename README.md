# Api Client Library

Repository in development process, do not use in production environments.

This repository contains the library to communicate with Nuvla API Server.

The module provides multiple clients depending on the type of operations you want to perform. 
The clients are built using the Composition pattern, so that the base client is extended with the specific operations required for the client.

| Client Name       | Description                                                                                                                                                                                                                                                                                                                  | Dev. Status    |
|-------------------|------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|----------------|
| NuvlaClient       | Contains generic HTTP methods implemented.  <br/> - Login  <br/>- Get <br/> - Post  <br/>- Put  <br/>- Remove <br/>- Search <br/>This is the base client which is then imported  in all the other clients using the Composition pattern                                                                                      | Pre-production |
| User Client       | Extends the above mentioned client and provides utilities to execute operations related (and allowed) to the user.  - Add/Remove resource  - Get resource  - Search resource  - LogIn  Then it also has wrappers to access the lower level HTTP operations from NuvlaClient for special operations and non-covered resources | Pre-production |
| NuvlaEdge Client  | Following the same pattern, creates a client with specific operations NuvlaEdge requires:  <br/>- Activate  <br/>- Commission <br/> - Telemetry  <br/>- GetResources (Retrieves all the accessible resources from NuvlaEdge)                                                                                                 | Pre-production |
| Deployment Client | Client to control Deployment resources and related operations                                                                                                                                                                                                                                                                | Pre-production |
| Job Client        | Client to control and retrieve Job resources and related operations                                                                                                                                                                                                                                                          | Pre-production |

## Usage
Basic usage for generic NuvlaClient:

```go
package main

import (
	"fmt"
	nuvla "github.com/nuvla/api-client-go/"
)

func main() {
	// Create a new client
	client := nuvla.NewNuvlaClient("https://nuvla.io", false, false)

	// Login
	client.LoginApiKeys("api-key", "api-secret")

	// Get a resource
	res, _ := client.Get("nuvlabox/nuvlabox-id")
	fmt.Println("Resource: ", res)
	
	// Post a resource
	client.Post("resource-endpoint", "resource-data")

	// Put a resource
	client.Put("resource-endpoint", "resource-data", "data-to-delete")

	// Remove a resource
	client.Remove("resource-id")
}
```

Generic usage for NuvlaClient with detailed configuration:

```go
package main

import (
	"fmt"
	nuvla "github.com/nuvla/api-client-go/"
)

func main() {
	// Create a new client
	// Create options using NewSessionOpts method creates them with default values 
	clientOps := nuvla.NewSessionOpts(&nuvla.SessionOpts{})
	client := nuvla.NewNuvlaClientFromOpts(clientOps)
}
```
