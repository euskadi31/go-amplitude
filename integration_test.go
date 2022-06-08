// Copyright 2022 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package amplitude

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestIntegration(t *testing.T) {
	c := New(
		os.Getenv("DEMO_AMPLITUDE_API_KEY"),
		WithURL(StandardEndpoint),
		WithTimeout(time.Second*2),
		WithInterval(time.Second*5),
		WithBatchSize(2),
		WithBufferSize(2),
		WithMaxRetry(2),
		WithRetryInterval(time.Second*5),
	)
	defer c.Close()

	err := c.Enqueue(&Event{
		UserID:      "f892be22-8f8e-445d-83b0-af199b9a5c72",
		DeviceID:    "0a16e988-8f70-4877-bdc6-08997832cfff",
		Timestamp:   1643367217,
		EventType:   "user.created",
		Platform:    "ios",
		OSName:      "iOS",
		OSVersion:   "15.2.1",
		DeviceModel: "iPhone13,3",
		Language:    "fr-FR",
		InsertID:    "a5461410-6b12-4a7a-905d-166cc00af4b2",
	})
	assert.NoError(t, err)

	time.Sleep(15 * time.Second)
}
