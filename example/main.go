package main

import (
	"os"

	"github.com/euskadi31/go-amplitude"
)

func main() {
	client := amplitude.New(
		os.Getenv("AMPLITUDE_API_KEY"),
		amplitude.WithURL(amplitude.EUResidencyEndpoint),
	)
	defer client.Close()

	evt := &amplitude.Event{
		EventType: "amplitude.client.example",
		UserID:    "user-demo",
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
