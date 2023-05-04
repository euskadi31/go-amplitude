// Copyright 2022 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package amplitude

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestClient(t *testing.T) {
	var wg sync.WaitGroup

	wg.Add(1)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer wg.Done()

		b, err := io.ReadAll(r.Body)
		assert.NoError(t, err)

		defer r.Body.Close()

		assert.Equal(t, `{"api_key":"foo","events":[{"user_id":"f892be22-8f8e-445d-83b0-af199b9a5c72","device_id":"0a16e988-8f70-4877-bdc6-08997832cfff","event_type":"user.created","time":1643367217,"platform":"ios","os_name":"iOS","os_version":"15.2.1","device_model":"iPhone13,3","language":"fr-FR","insert_id":"a5461410-6b12-4a7a-905d-166cc00af4b2"}]}`, string(b))

		msg := &RequestPayload{}

		err = json.Unmarshal(b, msg)
		assert.NoError(t, err)

		assert.Equal(t, "foo", msg.APIKey)
		assert.Equal(t, 1, len(msg.Events))
	}))
	defer ts.Close()

	c := New(
		"foo",
		WithURL(ts.URL),
		WithTimeout(time.Second*1),
		WithInterval(time.Millisecond*100),
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

	wg.Wait()
}

func TestClientWithRetry(t *testing.T) {
	var wg sync.WaitGroup

	wg.Add(2)

	retry := 0

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer wg.Done()

		if retry == 0 {
			w.WriteHeader(http.StatusInternalServerError)
			retry++

			return
		}

		b, err := io.ReadAll(r.Body)
		assert.NoError(t, err)

		defer r.Body.Close()

		assert.Equal(t, `{"api_key":"foo","events":[{"user_id":"f892be22-8f8e-445d-83b0-af199b9a5c72","device_id":"0a16e988-8f70-4877-bdc6-08997832cfff","event_type":"user.created","time":1643367217,"platform":"ios","os_name":"iOS","os_version":"15.2.1","device_model":"iPhone13,3","language":"fr-FR","insert_id":"a5461410-6b12-4a7a-905d-166cc00af4b2"}]}`, string(b))

		msg := &RequestPayload{}

		err = json.Unmarshal(b, msg)
		assert.NoError(t, err)

		assert.Equal(t, "foo", msg.APIKey)
		assert.Equal(t, 1, len(msg.Events))
	}))
	defer ts.Close()

	c := New(
		"foo",
		WithURL(ts.URL),
		WithTimeout(time.Second*1),
		WithInterval(time.Millisecond*100),
		WithBatchSize(2),
		WithBufferSize(2),
		WithMaxRetry(2),
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

	wg.Wait()
}

func TestClientDroppedMessage(t *testing.T) {
	retry := 0

	var wg sync.WaitGroup

	wg.Add(3)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer wg.Done()

		retry++

		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	c := New(
		"foo",
		WithURL(ts.URL),
		WithTimeout(time.Second*1),
		WithInterval(time.Millisecond*100),
		WithBatchSize(2),
		WithBufferSize(2),
		WithMaxRetry(2),
		WithRetryInterval(time.Second*1),
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

	wg.Wait()

	assert.Equal(t, 3, retry)
}

func TestClientWithMultipleEvents(t *testing.T) {
	var wg sync.WaitGroup

	wg.Add(2)

	hits := 0

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer wg.Done()
		defer func() {
			hits++
		}()

		b, err := io.ReadAll(r.Body)
		assert.NoError(t, err)

		defer r.Body.Close()

		msg := &RequestPayload{}

		err = json.Unmarshal(b, msg)
		assert.NoError(t, err)

		assert.Equal(t, "foo", msg.APIKey)

		switch hits {
		case 0:
			assert.Equal(t, 2, len(msg.Events))
		case 1:
			assert.Equal(t, 2, len(msg.Events))
		}
	}))
	defer ts.Close()

	c := New(
		"foo",
		WithURL(ts.URL),
		WithTimeout(time.Second*1),
		WithInterval(time.Millisecond*500),
		WithBatchSize(2),
		WithBufferSize(3),
		WithMaxRetry(2),
		WithRetryInterval(time.Millisecond*100),
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

	err = c.Enqueue(&Event{
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

	err = c.Enqueue(&Event{
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

	err = c.Enqueue(&Event{
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

	wg.Wait()
}

func TestClientClose(t *testing.T) {
	var wg sync.WaitGroup

	wg.Add(1)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer wg.Done()

		b, err := io.ReadAll(r.Body)
		assert.NoError(t, err)

		defer r.Body.Close()

		assert.Equal(t, `{"api_key":"foo","events":[{"user_id":"f892be22-8f8e-445d-83b0-af199b9a5c72","device_id":"0a16e988-8f70-4877-bdc6-08997832cfff","event_type":"user.created","time":1643367217,"platform":"ios","os_name":"iOS","os_version":"15.2.1","device_model":"iPhone13,3","language":"fr-FR","insert_id":"a5461410-6b12-4a7a-905d-166cc00af4b2"}]}`, string(b))

		msg := &RequestPayload{}

		err = json.Unmarshal(b, msg)
		assert.NoError(t, err)

		assert.Equal(t, "foo", msg.APIKey)
		assert.Equal(t, 1, len(msg.Events))
	}))
	defer ts.Close()

	c := New(
		"foo",
		WithURL(ts.URL),
		WithTimeout(time.Second*1),
		WithInterval(time.Second*2),
		WithBatchSize(2),
		WithBufferSize(2),
		WithMaxRetry(2),
		WithRetryInterval(time.Second*2),
	)

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

	assert.NoError(t, c.Close())

	wg.Wait()
}

func TestClientGetBatchEvents(t *testing.T) {

	c := &client{
		timeout:       time.Second * 1,
		interval:      time.Second * 10,
		batchSize:     2,
		bufferSize:    4,
		maxRetry:      3,
		retryInterval: time.Second * 1,
		retrySize:     1000,
	}

	c.events = []*Event{
		{
			UserID: "f892be22-8f8e-445d-83b0-af199b9a5c71",
		},
		{
			UserID: "f892be22-8f8e-445d-83b0-af199b9a5c72",
		},
		{
			UserID: "f892be22-8f8e-445d-83b0-af199b9a5c73",
		},
	}

	events := c.getBatchEvents()

	assert.Equal(t, 2, len(events))

	events = c.getBatchEvents()

	assert.Equal(t, 1, len(events))
}

/*
func TestClientRace(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	}))
	defer ts.Close()

	c := New(
		"foo",
		WithURL(ts.URL),
		WithTimeout(time.Second*2),
		WithInterval(time.Second*10),
		WithRetryInterval(time.Second*5),
	).(*client)

	for i := 0; i < 2000; i++ {
		now := time.Now().Unix()
		id := uuid.New().String()

		err := c.Enqueue(&Event{
			UserID:      "f892be22-8f8e-445d-83b0-af199b9a5c72",
			DeviceID:    "0a16e988-8f70-4877-bdc6-08997832cfff",
			Timestamp:   now,
			EventType:   "user.created",
			Platform:    "ios",
			OSName:      "iOS",
			OSVersion:   "15.2.1",
			DeviceModel: "iPhone13,3",
			Language:    "fr-FR",
			InsertID:    id,
		})
		assert.NoError(t, err)
	}

	assert.NoError(t, c.Close())

wait:
	for {
		if len(c.msgs) == 0 && len(c.events) == 0 && len(c.retries) == 0 {
			break wait
		}
	}
}
*/
