// Copyright 2022 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package amplitude

import "errors"

var (
	// ErrClosed message.
	ErrClosed = errors.New("the client was already closed")

	// ErrBatchFailed message.
	ErrBatchFailed = errors.New("request failed")
)
