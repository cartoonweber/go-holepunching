package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"

	"github.com/batphonghan/holepunching-go/model"
)

func main() {
	runUDPServer()
}

var prevAddr string

func runUDPServer() {
	addr, err := net.ResolveUDPAddr("udp", ":8081")
	if err != nil {
		panic(err)
	}
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		panic(err)
	}

	for {
		handleUDPClient(conn)
	}
}

var ips = make(map[string]string, 0)

func handleUDPClient(conn *net.UDPConn) {
	var buff [512]byte
	n, addr, err := conn.ReadFromUDP(buff[:])
	if err != nil {
		fmt.Println("Err readfrom UDP", err)
		return
	}

	var p model.Request
	err = json.Unmarshal(buff[:n], &p)
	if err != nil {
		fmt.Println("Err payload", err)

		return
	}

	switch p.Action {
	case model.Get:
		fmt.Println("Process: GET the IP request: ", p)
		var resp model.Response
		if ip, ok := ips[p.PeerID]; ok {
			resp.RAddr = &ip
		}
		b, err := json.Marshal(resp)
		if err != nil {
			panic(err)
		}
		if resp.RAddr != nil {
			fmt.Printf("\n Write back %+v \n", *resp.RAddr)
		} else {
			fmt.Println("Write back emty")
		}
		_, err = conn.WriteToUDP(b, addr)
		if err != nil {
			fmt.Println("Err write back the ip of peer ", p.PeerID, err)
			return
		}
	case model.Reg:
		fmt.Println("Process: REG the IP request:", p)
		remoteAddr := fmt.Sprintf("%s:%d", addr.IP, addr.Port)
		ips[p.PeerID] = remoteAddr

		messageRequest := ChatRequest{
			"Chat",
			"chatRequest.Username",
			remoteAddr,
		}
		jsonRequest, err := json.Marshal(&messageRequest)
		if err != nil {
			log.Print(err)
			break
		}
		_, err = conn.WriteToUDP(jsonRequest, addr)
		if err != nil {
			log.Print(err)
		}
	}

	fmt.Printf("IP table %+v \n", ips)
}

type ChatRequest struct {
	Action   string
	Username string
	Message  string
}
