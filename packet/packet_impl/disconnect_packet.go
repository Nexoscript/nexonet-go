package packetimpl

import (
	"fmt"

	"github.com/Nexoscript/nexonet-go/api"
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
	return fmt.Sprintf("DisconnectPacket{type=%s, code=%d}", dp.Type, dp.Code)
}
