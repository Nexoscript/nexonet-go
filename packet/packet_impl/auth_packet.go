package packetimpl

import (
	"fmt"

	"github.com/Nexoscript/nexonet-go/api"
)

type AuthPacket struct {
	api.BasePacket
	Id string `json:"id"`
}

func NewAuthPacket(id string) *AuthPacket {
	return &AuthPacket{
		BasePacket: api.BasePacket{Type: "AUTH"},
		Id:         id,
	}
}

func (ap AuthPacket) String() string {
	return fmt.Sprintf("AuthPacket{id=%s}", ap.Id)
}
