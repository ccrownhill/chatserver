package main

import (
	"fmt"
	"net"
	"strings"
)

const PORT = "12000"
const BUFSIZE = 1024

func main() {
	clients := make(map[*net.UDPAddr]string) // map client address to client name
	socket, err := net.ResolveUDPAddr("udp4", ":"+PORT)
	if err != nil {
		fmt.Println(err)
		return
	}

	connection, err := net.ListenUDP("udp4", socket)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer connection.Close()
	buffer := make([]byte, BUFSIZE)

	for {
		n, addr, err := connection.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println(err)
			return
		}
		msg := strings.TrimSpace(string(buffer[0:n])) // get rid of newline

		splitMsg := strings.Split(msg, ":")
		if splitMsg[0] == "client" {
			fmt.Println("adding client", splitMsg[1])
			clients[addr] = splitMsg[1]
		} else if len(splitMsg) > 1 && splitMsg[1] == "STOP" {
			fmt.Println("removing user:", splitMsg[0])
			delete(clients, addr)
		} else {
			for a := range clients {
				//fmt.Printf("sending: %s to %q\n", msg, a)
				_, err := connection.WriteToUDP([]byte(msg), a)
				if err != nil {
					fmt.Println(err)
					return
				}
			}
		}
	}
}
