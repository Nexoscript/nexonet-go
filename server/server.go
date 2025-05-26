package server

import (
	"fmt"
	"net"
	"os"
	"strconv"

	"github.com/Nexoscript/nexonet-go/api"
	"github.com/Nexoscript/nexonet-go/packet"

	packetimpl "github.com/Nexoscript/nexonet-go/packet/packet_impl"
)

var packetManager *packet.PacketManager

func Initilize() {
	packetManager = packet.NewPacketManager()
	packetManager.RegisterPacketType("DISCONNECT", func() api.PacketInterface { return &packetimpl.DisconnectPacket{} })
	packetManager.RegisterPacketType("AUTH", func() api.PacketInterface { return &packetimpl.AuthPacket{} })
	packetManager.RegisterPacketType("AUTH_RESPONSE", func() api.PacketInterface { return &packetimpl.AuthResponsePacket{} })
}

func Start(host string, port int64) {
	listen, err := net.Listen("tcp", host+":"+strconv.FormatInt(port, 10))
	fmt.Println("Server is listening on " + host + ":" + strconv.FormatInt(port, 10))
	if err != nil {
		fmt.Println("Error while listening:", err.Error())
	}
	defer listen.Close()
	for {
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("Error while accepting client connection:", err.Error())
			os.Exit(1)
		}
		go HandleClientRequest(conn)
	}
}

func SendToClient(id string, packet api.PacketInterface) {

}
