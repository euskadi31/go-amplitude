// Copyright 2022 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package amplitude

type RequestPayload struct {
	APIKey  string          `json:"api_key"`
	Events  []*Event        `json:"events"`
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
