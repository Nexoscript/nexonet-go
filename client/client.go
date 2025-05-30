package client

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Nexoscript/nexonet-go/api"
	"github.com/Nexoscript/nexonet-go/packet"
	packetimpl "github.com/Nexoscript/nexonet-go/packet/packet_impl"

	"github.com/google/uuid"
)

const (
	NumIterations = 3
)

var id string
var conn net.Conn
var packetManager *packet.PacketManager
var isAuth bool = false
var isRunning bool = false

func initialize() {
	packetManager = packet.NewPacketManager()
	packetManager.RegisterPacketType("DISCONNECT", func() api.PacketInterface { return &packetimpl.DisconnectPacket{} })
	packetManager.RegisterPacketType("ACCEPT", func() api.PacketInterface { return &packetimpl.AcceptPacket{} })
	packetManager.RegisterPacketType("AUTH", func() api.PacketInterface { return &packetimpl.AuthPacket{} })
	packetManager.RegisterPacketType("AUTH_RESPONSE", func() api.PacketInterface { return &packetimpl.AuthResponsePacket{} })
}

func Connect(host string, port int64) {
	initialize()
	var err error
	conn, err = net.Dial("tcp", host+":"+strconv.FormatInt(port, 10))
	if err != nil {
		fmt.Println("Error while connecting:", err.Error())
		os.Exit(1)
	}
	fmt.Println("Connected to server ", host+":"+strconv.FormatInt(port, 10))
	isRunning = true
	serverReader := bufio.NewReader(conn)
	go run(serverReader)
}

func run(reader *bufio.Reader) {
	for isRunning {
		err := conn.SetReadDeadline(time.Now().Add(500 * time.Millisecond))
		if err != nil {
			return
		}
		var serverResponse string
		serverResponse, err = reader.ReadString('\n')
		if err != nil {
			var netErr net.Error
			if errors.As(err, &netErr) && netErr.Timeout() {
				continue
			}
			if conn == nil {
				Disconnect()
				fmt.Println("Client disconnected")
				break
			}
			if err == io.EOF {
				fmt.Println("Server connection closed.")
			} else {
				fmt.Printf("Error while reading server response: %v\n", err)
			}
			Disconnect()
			return
		}
		err = conn.SetReadDeadline(time.Time{})
		if err != nil {
			return
		}
		serverResponse = strings.TrimSpace(serverResponse)
		if serverResponse != "" {
			if !strings.HasPrefix(serverResponse, "{") {
				serverResponse = "{" + serverResponse
			}
			serializedPacket, err := packetManager.FromJson(serverResponse)
			if err != nil {
				fmt.Printf("Error while deserializing serializedPacket: %v\n", err)
				continue
			}
			switch p := serializedPacket.(type) {
			case *packetimpl.AuthResponsePacket:
				if p.IsSuccess {
					id = p.Id
					isAuth = true
					fmt.Printf("Authentication successfull. Client ID: %s\n", id)
					// clientConnectEvent.onClientConnect(this) - falls vorhanden
				} else {
					fmt.Println("Authentication failed, send new Auth-Packet.")
					authPacket := packetimpl.NewAuthPacket(uuid.New().String())
					SendPacket(conn, authPacket)
				}
			case *packetimpl.DisconnectPacket:
				fmt.Printf("Server has send DISCONNECT-Packet with Code %d. Closing connection.\n", p.Code)
				Disconnect()
			default:
				fmt.Printf("Received serializedPacket of type '%s': %s\n", p.GetType(), serverResponse)
			}
		}
	}
	fmt.Println("Run-Goroutine exited.")
}

func Disconnect() {
	if !isRunning {
		return
	}
	if conn != nil {
		err := conn.Close()
		if err != nil {
			return
		}
	}
	isRunning = false
	isAuth = false
	fmt.Println("Client connection closed.")
}

func SendPacket(conn net.Conn, p api.PacketInterface) {
	if conn == nil {
		fmt.Printf("Error: Can't send packet '%s', connection is nil.\n", p.GetType())
		return
	}
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

func IsRunning() bool {
	return isRunning
}

func IsAuth() bool {
	return isAuth
}

func GetPacketManager() *packet.PacketManager {
	return packetManager
}

func GetConnection() net.Conn {
	return conn
}

func GetId() string {
	return id
}
