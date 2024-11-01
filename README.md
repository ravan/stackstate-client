# StackState Client Library

Enables communication with the StackState rest api and receiver api.

This library is evolving to support all StackState api calls.  Please feel free to add new api calls as needed.


## Installation

```bash
go get github.com/ravan/stackstate-client
```

## Usage

### Access Rest API Endpoints
```go 

import (
	sts "github.com/ravan/stackstate-client/stackstate"
	"github.com/ravan/stackstate-client/stackstate/api"
	"github.com/ravan/stackstate-client/stackstate/receiver"
)
conf := &sts.StackState{
    ApiUrl:   os.Getenv("STS_URL"),
    ApiKey:   os.Getenv("STS_API_KEY"),
    ApiToken: os.Getenv("STS_TOKEN"),
}

// To access the API rest endpoints
client := api.NewClient(conf)
status, err := client.Status()



```

### Access Receiver API Endpoints

See [StackState k8s extension](https://github.com/ravan/stackstate-k8s-ext/blob/main/cmd/sync/main.go) integration for examples on using the receiver api.

## Authorization

The TopologyQuery and TopologyStreamQuery methods require additional authorization on the StackState server.
You will get the error `"The supplied authentication is not authorized to access this resource"`.
Speak to your StackState administrator for more information.