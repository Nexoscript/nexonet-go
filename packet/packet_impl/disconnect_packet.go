package packetimpl

import (
	"fmt"
	"nexonet/api"
)

type DisconnectPacket struct {
	api.BasePacket
	Code int `json:"code"`
}

func NewDisconnectPacket(code int) *DisconnectPacket {
	return &DisconnectPacket{
		BasePacket: api.BasePacket{Type: "DISCONNECT"},
		Code:       code,
	}
}

func (dp DisconnectPacket) String() string {
	return fmt.Sprintf("DisconnectPacket{code=%d}", dp.Code)
}
