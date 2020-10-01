package device

import (
	. "GoServer/model/device"
)

type RequestParam struct {
	Type     AccesswayType `form:"type"`
	ClientID string        `form:"clientID"`
	Version  string        `form:"version"`
}

type RequestData struct {
	DeviceNo string `form:"deviceNo" json:"deviceNo" binding:"required"`
	Token    string `form:"token" json:"token" binding:"required"`
}

type FirmwareInfo struct {
	Size int64  `json:"size"`
	URL  string `json:"url"`
}
