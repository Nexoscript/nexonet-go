package packetimpl

import (
	"fmt"

	"github.com/Nexoscript/nexonet-go/api"
)

type AcceptPacket struct {
	api.BasePacket
	Message string `json:"message"`
}

func NewAcceptPacket(message string) *AcceptPacket {
	return &AcceptPacket{
		BasePacket: api.BasePacket{Type: "ACCEPT"},
		Message:    message,
	}
}

func (ap AcceptPacket) String() string {
	return fmt.Sprintf("AcceptPacket{message=%s}", ap.Message)
}
