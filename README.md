# Amplitude Client for Go

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
