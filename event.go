// Copyright 2022 Axel Etcheverry. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package amplitude

// Event struct.
// see: https://developers.amplitude.com/docs/http-api-v2
type Event struct {
	UserID             string                 `json:"user_id,omitempty"`
	DeviceID           string                 `json:"device_id,omitempty"`
	EventType          string                 `json:"event_type"`
	Timestamp          int64                  `json:"time,omitempty"`
	EventProperties    map[string]interface{} `json:"event_properties,omitempty"`
	UserProperties     map[string]interface{} `json:"user_properties,omitempty"`
	Groups             map[string]interface{} `json:"groups,omitempty"`
	AppVersion         string                 `json:"app_version,omitempty"`
	Platform           string                 `json:"platform,omitempty"`
	OSName             string                 `json:"os_name,omitempty"`
	OSVersion          string                 `json:"os_version,omitempty"`
	DeviceBrand        string                 `json:"device_brand,omitempty"`
	DeviceManufacturer string                 `json:"device_manufacturer1,omitempty"`
	DeviceModel        string                 `json:"device_model,omitempty"`
	Carrier            string                 `json:"carrier,omitempty"`
	Country            string                 `json:"country,omitempty"`
	Region             string                 `json:"region,omitempty"`
	City               string                 `json:"city,omitempty"`
	DMA                string                 `json:"dma,omitempty"`
	Language           string                 `json:"language,omitempty"`
	Price              float64                `json:"price,omitempty"`
	Quantity           int                    `json:"quantity,omitempty"`
	Revenue            float64                `json:"revenue,omitempty"`
	ProductID          string                 `json:"productId,omitempty"`
	RevenueType        string                 `json:"revenueType,omitempty"`
	LocationLat        float64                `json:"location_lat,omitempty"`
	LocationLng        float64                `json:"location_lng,omitempty"`
	IP                 string                 `json:"ip,omitempty"`
	IDFA               string                 `json:"idfa,omitempty"`
	IDFV               string                 `json:"idfv,omitempty"`
	ADID               string                 `json:"adid,omitempty"`
	AndroidID          string                 `json:"android_id,omitempty"`
	EventID            int                    `json:"event_id,omitempty"`
	SessionID          int64                  `json:"session_id,omitempty"`
	InsertID           string                 `json:"insert_id,omitempty"`
	Plan               *Plan                  `json:"plan,omitempty"`
}

type Plan struct {
	Branch  string `json:"branch,omitempty"`
	Source  string `json:"source,omitempty"`
	Version string `json:"version,omitempty"`
}
