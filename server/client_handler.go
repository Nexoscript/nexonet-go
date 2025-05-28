package server

import (
	"bufio"
	"fmt"
	"net"

	packetimpl "github.com/Nexoscript/nexonet-go/packet/packet_impl"
)

func HandleClientRequest(conn net.Conn) {
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {

		}
	}(conn)
	clientReader := bufio.NewReader(conn)
	for {
		jsonMessageBytes, err := clientReader.ReadBytes('\n')
		if err != nil {
			if err.Error() == "EOF" {
				fmt.Println("Client disconnected:", conn.RemoteAddr().String())
				return
			}
			fmt.Println("Error while reading:", err.Error())
			return
		}
		jsonString := string(jsonMessageBytes[:len(jsonMessageBytes)-1])
		packet, err := packetManager.FromJson(jsonString)
		if err != nil {
			fmt.Println("Error while deserializing packet:", err.Error())
			continue
		}
		fmt.Printf("Received [%s] from %s: %s\n", packet.GetType(), conn.RemoteAddr().String(), packet.String())
		acceptPacket := packetimpl.NewAcceptPacket(fmt.Sprintf("Received [%s] from %s: %s", packet.GetType(), conn.RemoteAddr().String(), packet.String()))
		jsonString, err = packetManager.ToJson(acceptPacket)
		if err != nil {
			fmt.Printf("Error while serializing packet %s: %s\n", packet.GetType(), err.Error())
			return
		}
		_, err = conn.Write(append([]byte(jsonString), '\n'))
		if err != nil {
			fmt.Println("Error while writing answer:", err.Error())
			continue
		}
	}
}
