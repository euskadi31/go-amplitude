// Copyright 2026 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package amplitude

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrorResponse(t *testing.T) {

	e := &ErrorResponse{
		Code:         400,
		ErrorMessage: "Bad Request",
	}

	assert.Equal(t, "400: Bad Request", e.Error())

	e.MissingField = "events"

	assert.Equal(t, "400: Bad Request: missing: events", e.Error())
}
