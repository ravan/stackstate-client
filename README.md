# StackState Client Library

Enables communication with the StackState rest api and receiver api.

This library is evolving to support all StackState api calls.  Please feel free to add new api calls as needed.


## Installation

```bash
go get github.com/ravan/stackstate-client
```

## Usage

```go 

conf := &sts.StackState{
    ApiUrl:   os.Getenv("STS_URL"),
    ApiKey:   os.Getenv("STS_API_KEY"),
    ApiToken: os.Getenv("STS_TOKEN"),
}

client := NewClient(conf)
status, err := client.Status()

```

## Authorization

The TopologyQuery and TopologyStreamQuery methods require additional authorization on the StackState server.
You will get the error `"The supplied authentication is not authorized to access this resource"`.
Speak to your StackState administrator for more information.