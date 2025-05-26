package packet

import (
	"encoding/json"
	"fmt"

	"github.com/Nexoscript/nexonet-go/api"
)

type PacketCreator func() api.PacketInterface

type PacketManager struct {
	packetRegistry map[string]PacketCreator
}

func NewPacketManager() *PacketManager {
	return &PacketManager{packetRegistry: make(map[string]PacketCreator)}
}

func (pm *PacketManager) RegisterPacketType(packetType string, creator PacketCreator) {
	if _, exists := pm.packetRegistry[packetType]; exists {
		fmt.Printf("Warn: Packet type '%s' is already registered and will be overwritten.\n", packetType)
	}
	pm.packetRegistry[packetType] = creator
}

func (pm *PacketManager) ToJson(packet api.PacketInterface) (string, error) {
	jsonBytes, err := json.Marshal(packet)
	if err != nil {
		return "", fmt.Errorf("error while seralize packet: %w", err)
	}
	return string(jsonBytes), nil
}

func (pm *PacketManager) FromJson(jsonString string) (api.PacketInterface, error) {
	var base api.BasePacket
	err := json.Unmarshal([]byte(jsonString), &base)
	if err != nil {
		return nil, fmt.Errorf("error while reading type of packet: %w", err)
	}
	creator, found := pm.packetRegistry[base.Type]
	if !found {
		return nil, fmt.Errorf("unknown registered packet type: %s", base.Type)
	}
	packet := creator()
	err = json.Unmarshal([]byte(jsonString), packet)
	if err != nil {
		return nil, fmt.Errorf("error while deseralize packet from type %s: %w", base.Type, err)
	}
	return packet, nil
}
