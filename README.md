# Amplitude Client for Go [![Last release](https://img.shields.io/github/release/euskadi31/go-amplitude.svg)](https://github.com/euskadi31/go-amplitude/releases/latest) [![Documentation](https://godoc.org/github.com/euskadi31/go-amplitude?status.svg)](https://godoc.org/github.com/euskadi31/go-amplitude)

[![Go Report Card](https://goreportcard.com/badge/github.com/euskadi31/go-amplitude)](https://goreportcard.com/report/github.com/euskadi31/go-amplitude)

| Branch | Status                                                                                                                                                    | Coverage                                                                                                                                             |
| ------ | --------------------------------------------------------------------------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------- |
| master | [![Go](https://github.com/euskadi31/go-amplitude/actions/workflows/go.yml/badge.svg)](https://github.com/euskadi31/go-amplitude/actions/workflows/go.yml) | [![Coveralls](https://img.shields.io/coveralls/euskadi31/go-amplitude/master.svg)](https://coveralls.io/github/euskadi31/go-amplitude?branch=master) |

## Example

```go
package main

import (
    "github.com/euskadi31/go-amplitude"
)

func main() {
    client := amplitude.New(
        "my-amplitude-key",
        amplitude.WithURL(amplitude.EUResidencyEndpoint),
    )
    defer client.Close()

    evt := &amplitude.Event{
        EventType: "user.created",
        UserID: "c427ba84-a0c3-48d5-aaef-302734212064",
        EventProperties: map[string]interface{}{
            "from": "mobile",
        },
        UserProperties: map[string]interface{}{
            "birthday_year": "1987",
        },
    }

    if err := client.Enqueue(evt); err != nil {
        panic(err)
    }
}
```
