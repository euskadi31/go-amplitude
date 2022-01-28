// Copyright 2022 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package amplitude

import (
	"net/http"
	"time"
)

type Option func(*client)

func WithURL(url string) Option {
	return func(c *client) {
		c.endpoint = url
	}
}

func WithTimeout(timeout time.Duration) Option {
	return func(c *client) {
		c.timeout = timeout
	}
}

func WithInterval(interval time.Duration) Option {
	return func(c *client) {
		c.interval = interval
	}
}

func WithBatchSize(size int) Option {
	return func(c *client) {
		c.batchSize = size
	}
}

func WithBufferSize(size int) Option {
	return func(c *client) {
		c.bufferSize = size
	}
}

func WithMaxRetry(retry int) Option {
	return func(c *client) {
		c.maxRetry = retry
	}
}

func WithRetryInterval(interval time.Duration) Option {
	return func(c *client) {
		c.retryInterval = interval
	}
}

func WithRetrySize(size int) Option {
	return func(c *client) {
		c.retrySize = size
	}
}

func WithHTTPClient(httpClient *http.Client) Option {
	return func(c *client) {
		c.httpClient = httpClient
	}
}
