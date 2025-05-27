package server

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/Nexoscript/nexonet-go/api"
	"github.com/Nexoscript/nexonet-go/packet"

	packetimpl "github.com/Nexoscript/nexonet-go/packet/packet_impl"
)

var listen net.Listener
var packetManager *packet.PacketManager
var isRunning bool = false

func initilize() {
	isRunning = true
	packetManager = packet.NewPacketManager()
	packetManager.RegisterPacketType("DISCONNECT", func() api.PacketInterface { return &packetimpl.DisconnectPacket{} })
	packetManager.RegisterPacketType("AUTH", func() api.PacketInterface { return &packetimpl.AuthPacket{} })
	packetManager.RegisterPacketType("AUTH_RESPONSE", func() api.PacketInterface { return &packetimpl.AuthResponsePacket{} })
}

func Start(host string, port int64) {
	initilize()
	var err error
	listen, err = net.Listen("tcp", host+":"+strconv.FormatInt(port, 10))
	fmt.Println("Server is listening on " + host + ":" + strconv.FormatInt(port, 10))
	if err != nil {
		fmt.Println("Error while listening:", err.Error())
	}
	go run()
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan

	fmt.Println("Server closing...")
	listen.Close()
}

func run() {
	for isRunning {
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("Error while accepting client connection:", err.Error())
			Close()
			break
		}
		go HandleClientRequest(conn)
	}
}

func Close() {
	if listen != nil {
		listen.Close()
	}
	isRunning = false
}

func SendToClient(id string, packet api.PacketInterface) {

}

func SendToClients(packet api.PacketInterface) {

}

func GetPacketManager() *packet.PacketManager {
	return packetManager
}
