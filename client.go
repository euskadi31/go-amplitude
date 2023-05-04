// Copyright 2022 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package amplitude

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
)

const (
	StandardEndpoint    = "https://api2.amplitude.com/2/httpapi"
	EUResidencyEndpoint = "https://api.eu.amplitude.com/2/httpapi"
)

var userAgent = "Amplitude Golang Client (https://github.com/euskadi31/go-amplitude)"

// Client Amplitude interface.
type Client interface {
	Enqueue(event *Event) error
	Close() error
}

type client struct {
	endpoint      string
	key           string
	timeout       time.Duration
	interval      time.Duration
	batchSize     int
	bufferSize    int
	maxRetry      int
	retryInterval time.Duration
	retrySize     int
	httpClient    *http.Client
	msgs          chan Event
	events        []Event
	retries       chan *Payload
	quit          chan struct{}
	shutdown      chan struct{}
	mtx           sync.Mutex
}

// New Amplitude client.
func New(key string, opts ...Option) Client {
	c := &client{
		endpoint:      StandardEndpoint,
		key:           key,
		timeout:       time.Second * 1,
		interval:      time.Second * 10,
		batchSize:     1000,
		bufferSize:    2000,
		maxRetry:      3,
		retryInterval: time.Second * 1,
		retrySize:     1000,
		quit:          make(chan struct{}),
		shutdown:      make(chan struct{}),
	}

	c.httpClient = &http.Client{
		Timeout: c.timeout,
	}
	c.msgs = make(chan Event, c.bufferSize)
	c.events = make([]Event, 0, c.bufferSize)
	c.retries = make(chan *Payload, c.retrySize)

	for _, opt := range opts {
		opt(c)
	}

	go c.loop()

	return c
}

func (c *client) loop() {
	defer close(c.shutdown)

	tick := time.NewTicker(c.interval)
	defer tick.Stop()

	for {
		select {
		case payload := <-c.retries:
			if err := c.sendBatch(payload); err != nil {
				if payload.Attempts > c.maxRetry {
					log.Warn().Msgf("%d messages dropped because they failed to be sent after %d attempts", payload.Size, c.maxRetry)

					continue
				}

				c.retries <- payload
			}
		case event := <-c.msgs:
			c.events = append(c.events, event)

			if len(c.events) == c.bufferSize {
				c.flush()
			}

		case <-tick.C:
			c.flush()

		case <-c.quit:
			log.Debug().Msg("exit requested - draining messages")

			// Drain the msg channel, we have to close it first so no more
			// messages can be pushed and otherwise the loop would never end.
			close(c.msgs)

			for event := range c.msgs {
				c.events = append(c.events, event)

				if len(c.events) == cap(c.events) {
					c.flush()
				}
			}

			c.flush()

			close(c.retries)

			for payload := range c.retries {
				if err := c.sendBatch(payload); err != nil {
					log.Error().Msg("Amplitude send batch failed, events lost !")
				}
			}

			log.Debug().Msg("exit")

			return
		}
	}
}

func (c *client) Close() (err error) {
	defer func() {
		// Always recover, a panic could be raised if `c`.quit was closed which
		// means the method was called more than once.
		if recover() != nil {
			err = ErrClosed
		}
	}()

	close(c.quit)

	<-c.shutdown

	return
}

func (c *client) processErrorResponse(resp *http.Response) error {
	var errorResponse *ErrorResponse

	if err := json.NewDecoder(resp.Body).Decode(&errorResponse); err != nil {
		return fmt.Errorf("json decode failed: %w", err)
	}

	return errorResponse
}

func (c *client) sendBatch(payload *Payload) error {
	ctx := context.Background()

	payload.Attempts++

	r, err := http.NewRequestWithContext(ctx, http.MethodPost, c.endpoint, bytes.NewReader(payload.Body))
	if err != nil {
		return fmt.Errorf("http new request failed: %w", err)
	}

	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("Accept", "application/json")
	r.Header.Set("User-Agent", userAgent)

	// r.Header.Add("Content-Length", strconv.Itoa(len(data)))

	resp, err := c.httpClient.Do(r)
	if err != nil {
		log.Error().Err(err).Msg("")

		return fmt.Errorf("http client send request failed: %w", err)
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Error().Err(err).Msg("http client close response body failed")
		}
	}()

	if resp.StatusCode != http.StatusOK {
		err := c.processErrorResponse(resp)

		log.Error().Err(err).Msgf("Amplitude send batch failed: status code %d", resp.StatusCode)

		return ErrBatchFailed
	}

	log.Debug().Msg("Amplitude sent batch !")

	return nil
}

func (c *client) flush() error {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	if len(c.events) == 0 {
		return nil
	}

	end := c.batchSize
	if length := len(c.events); length < end {
		end = length
	}

	var events []Event

	events, c.events = c.events[0:end], c.events[end:]

	reqPayload := &RequestPayload{
		APIKey: c.key,
		Events: events,
	}

	b, err := json.Marshal(reqPayload)
	if err != nil {
		return fmt.Errorf("json marshal events failed: %w", err)
	}

	payload := &Payload{
		Body: b,
		Size: len(events),
	}

	if err := c.sendBatch(payload); err != nil {
		c.retries <- payload
	}

	return nil
}

func (c *client) Enqueue(event *Event) (err error) {
	if event.Timestamp == 0 {
		event.Timestamp = time.Now().UTC().Unix()
	}

	defer func() {
		// When the `msgs` channel is closed writing to it will trigger a panic.
		// To avoid letting the panic propagate to the caller we recover from it
		// and instead report that the client has been closed and shouldn't be
		// used anymore.
		if recover() != nil {
			err = ErrClosed
		}
	}()

	c.msgs <- *event

	if len(c.msgs) == cap(c.msgs) {
		go func() {
			err = c.flush()
		}()
	}

	return
}
