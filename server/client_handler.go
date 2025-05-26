package server

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
)

func HandleClientRequest(conn net.Conn) {
	defer conn.Close()
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
		responseMap := map[string]string{
			"status":  "OK",
			"message": fmt.Sprintf("Server has received packet '%s': %s", packet.GetType(), packet.String()),
		}
		responseJsonBytes, err := json.Marshal(responseMap)
		if err != nil {
			fmt.Println("Error while serializing answer:", err.Error())
			continue
		}
		_, err = conn.Write(append(responseJsonBytes, '\n'))
		if err != nil {
			fmt.Println("Error while writing answer:", err.Error())
			continue
		}
	}
}
