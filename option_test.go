// Copyright 2022 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package amplitude

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestWithURL(t *testing.T) {
	c := &client{}

	WithURL("https://api.amplitude.tld")(c)

	assert.Equal(t, "https://api.amplitude.tld", c.endpoint)
}

func TestWithTimeout(t *testing.T) {
	c := &client{}

	WithTimeout(time.Second * 2)(c)

	assert.Equal(t, time.Second*2, c.timeout)
}

func TestWithInterval(t *testing.T) {
	c := &client{}

	WithInterval(time.Second * 2)(c)

	assert.Equal(t, time.Second*2, c.interval)
}

func TestWithBatchSize(t *testing.T) {
	c := &client{}

	WithBatchSize(2)(c)

	assert.Equal(t, 2, c.batchSize)
}

func TestWithBufferSize(t *testing.T) {
	c := &client{}

	WithBufferSize(2)(c)

	assert.Equal(t, 2, c.bufferSize)
}

func TestWithMaxRetry(t *testing.T) {
	c := &client{}

	WithMaxRetry(2)(c)

	assert.Equal(t, 2, c.maxRetry)
}

func TestWithRetryInterval(t *testing.T) {
	c := &client{}

	WithRetryInterval(time.Second * 2)(c)

	assert.Equal(t, time.Second*2, c.retryInterval)
}

func TestWithRetrySize(t *testing.T) {
	c := &client{}

	WithRetrySize(2)(c)

	assert.Equal(t, 2, c.retrySize)
}

func TestWithHTTPClient(t *testing.T) {
	c := &client{}

	hc := &http.Client{}

	WithHTTPClient(hc)(c)

	assert.Equal(t, hc, c.httpClient)
}
