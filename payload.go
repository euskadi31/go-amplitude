// Copyright 2022 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package amplitude

import "fmt"

type RequestPayload struct {
	APIKey  string          `json:"api_key"`
	Events  []Event         `json:"events"`
	Options *PayloadOptions `json:"options,omitempty"`
}

type PayloadOptions struct {
	MinIDLength int `json:"min_id_length"`
}

type Payload struct {
	Body     []byte
	Attempts int
	Size     int
}

/*
*

	{
	  "code": 400,
	  "error": "Request missing required field",
	  "missing_field": "api_key",
	  "events_with_invalid_fields": {
	    "time": [
	      3,
	      4,
	      7
	    ]
	  },
	  "events_with_missing_fields": {
	    "event_type": [
	      3,
	      4,
	      7
	    ]
	  }
	}

*
*/
type ErrorResponse struct {
	Code         int    `json:"code"`
	ErrorMessage string `json:"error"`
	MissingField string `json:"missing_field"`
}

func (e *ErrorResponse) Error() string {
	msg := fmt.Sprintf("%d: %s", e.Code, e.ErrorMessage)

	if e.MissingField != "" {
		msg = fmt.Sprintf("%s: missing: %s", msg, e.MissingField)
	}

	return msg
}
