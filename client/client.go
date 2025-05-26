package client

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/Nexoscript/nexonet-go/api"
	"github.com/Nexoscript/nexonet-go/packet"
	packetimpl "github.com/Nexoscript/nexonet-go/packet/packet_impl"

	"github.com/google/uuid"
)

const (
	NUM_ITERATIONS = 3
)

var packetManager *packet.PacketManager

func Initilize() {
	packetManager = packet.NewPacketManager()
	packetManager.RegisterPacketType("DISCONNECT", func() api.PacketInterface { return &packetimpl.DisconnectPacket{} })
	packetManager.RegisterPacketType("AUTH", func() api.PacketInterface { return &packetimpl.AuthPacket{} })
	packetManager.RegisterPacketType("AUTH_RESPONSE", func() api.PacketInterface { return &packetimpl.AuthResponsePacket{} })
}

func Start(host string, port int64) {
	conn, err := net.Dial("tcp", host+":"+strconv.FormatInt(port, 10))
	if err != nil {
		fmt.Println("Error while connecting:", err.Error())
		os.Exit(1)
	}
	defer conn.Close()
	fmt.Println("Connected to server ", host+":"+strconv.FormatInt(port, 10))
	serverReader := bufio.NewReader(conn)
	go func() {
		for {
			message, err := serverReader.ReadString('\n')
			if err != nil {
				if err.Error() == "EOF" {
					fmt.Println("Server connection disconnected.")
					return
				}
				fmt.Println("Error while reading:", err.Error())
				return
			}
			fmt.Print("Server: " + message)
		}
	}()
	for i := 1; i <= NUM_ITERATIONS; i++ {
		disconnectPacket := packetimpl.NewDisconnectPacket(1000 + i)
		SendPacket(conn, disconnectPacket)
		time.Sleep(1 * time.Second)

		authPacket := packetimpl.NewAuthPacket(uuid.New().String())
		SendPacket(conn, authPacket)
		time.Sleep(1 * time.Second)
	}

	fmt.Println("Alle Pakete gesendet. Warte auf ausstehende Serverantworten...")
	time.Sleep(2 * time.Second)
	fmt.Println("Client wird beendet.")
}

func SendPacket(conn net.Conn, p api.PacketInterface) {
	jsonString, err := packetManager.ToJson(p)
	if err != nil {
		fmt.Printf("Error while serializing packet %s: %s\n", p.GetType(), err.Error())
		return
	}
	fmt.Printf("Send packet '%s': %s\n", p.GetType(), jsonString)
	_, err = conn.Write(append([]byte(jsonString), '\n'))
	if err != nil {
		fmt.Printf("Error while writing packet %s to server: %s\n", p.GetType(), err.Error())
	}
}
