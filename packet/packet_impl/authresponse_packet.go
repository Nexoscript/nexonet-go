package packetimpl

import (
	"fmt"

	"github.com/Nexoscript/nexonet-go/api"
)

type AuthResponsePacket struct {
	api.BasePacket
	Id        string `json:"id"`
	IsSuccess bool   `json:"isSuccess"`
}

func NewAuthResponsePacket(id string, isSuccess bool) *AuthResponsePacket {
	return &AuthResponsePacket{
		BasePacket: api.BasePacket{Type: "AUTH"},
		Id:         id,
		IsSuccess:  isSuccess,
	}
}

func (ap AuthResponsePacket) String() string {
	return fmt.Sprintf("AuthResponsePacket{type=%s, isSuccess=%t, id=%s}", ap.Type, ap.IsSuccess, ap.Id)
}
