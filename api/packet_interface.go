package api

type BasePacket struct {
	Type string `json:"type"`
}

type PacketInterface interface {
	GetType() string
	String() string
}

func (bp BasePacket) GetType() string {
	return bp.Type
}
