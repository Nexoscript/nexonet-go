package client

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/Nexoscript/nexonet-go/api"
	"github.com/Nexoscript/nexonet-go/packet"
	packetimpl "github.com/Nexoscript/nexonet-go/packet/packet_impl"

	"github.com/google/uuid"
)

const (
	NUM_ITERATIONS = 3
)

var id string
var conn net.Conn
var packetManager *packet.PacketManager
var isAuth bool = false
var isRunning bool = false

func initilize() {
	packetManager = packet.NewPacketManager()
	packetManager.RegisterPacketType("DISCONNECT", func() api.PacketInterface { return &packetimpl.DisconnectPacket{} })
	packetManager.RegisterPacketType("AUTH", func() api.PacketInterface { return &packetimpl.AuthPacket{} })
	packetManager.RegisterPacketType("AUTH_RESPONSE", func() api.PacketInterface { return &packetimpl.AuthResponsePacket{} })
}

func Connect(host string, port int64) {
	initilize()
	var err error
	conn, err = net.Dial("tcp", host+":"+strconv.FormatInt(port, 10))
	if err != nil {
		fmt.Println("Error while connecting:", err.Error())
		os.Exit(1)
	}
	fmt.Println("Connected to server ", host+":"+strconv.FormatInt(port, 10))
	serverReader := bufio.NewReader(conn)
	go run(serverReader)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan

	fmt.Println("Client disconnecting...")
	conn.Close()
}

func run(reader *bufio.Reader) {
	for {
		if isRunning {
			if !isAuth {
				authPacket := packetimpl.NewAuthPacket(uuid.New().String())
				SendPacket(conn, authPacket)
			}
			serverResponse, err := reader.ReadString('\n')
			if err != nil {
				if err == io.EOF {
					// this.logger.log(LoggingType.INFO, "Serververbindung geschlossen.")
					return
				}
				//this.logger.log(LoggingType.ERROR, fmt.Sprintf("Fehler beim Lesen vom Server: %v", err))
				return
			}
			serverResponse = strings.TrimSpace(serverResponse)
			if serverResponse != "" {
				if !strings.HasPrefix(serverResponse, "{") {
					serverResponse = "{" + serverResponse
				}
				packet, err := packetManager.FromJson(serverResponse)
				if err != nil {
					continue
				}
				if isAuth {
					//clientReceivedEvent.OnClientReceived(this, packet)
				}
				if authResponsePacket, ok := packet.(*packetimpl.AuthResponsePacket); ok {
					if authResponsePacket.IsSuccess {
						id = authResponsePacket.Id
						isAuth = true
						//clientConnectEvent.onClientConnect(this)
						continue
					}
					authPacket := packetimpl.NewAuthPacket(uuid.New().String())
					SendPacket(conn, authPacket)
				}
			}
			time.Sleep(1 * time.Millisecond)
			continue
		}
		break
	}
}

func Disconncet() {
	conn.Close()
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
