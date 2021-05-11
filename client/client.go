package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/batphonghan/holepunching-go/model"
)

var lPeer = flag.String("lpeer", "111", "peer id")
var rPeer = flag.String("rpeer", "222", "peer id")
var lPort = flag.String("port", "49980", "port")

func main() {
	flag.Parse()
	runUDPClient(*lPeer, *rPeer)
}

const (
	ServerAddr = "localhost:8081"
)

func runUDPClient(lPeerID string, rPeerID string) {
	var buf [512]byte
	raddr, err := net.ResolveUDPAddr("udp4", ServerAddr)
	if err != nil {
		panic(err)
	}

	lAddr, err := net.ResolveUDPAddr("udp4", fmt.Sprintf(":%s", *lPort))
	if err != nil {
		panic(err)
	}
	conn, err := net.ListenUDP("udp4", lAddr)
	if err != nil {
		panic(err)
	}

	fmt.Println("My local IP is", conn.LocalAddr().String())

	p := model.Request{
		Action: "REG",
		PeerID: lPeerID,
	}
	b, err := json.Marshal(p)
	if err != nil {
		panic(err)
	}
	_, err = conn.WriteToUDP(b, raddr)
	if err != nil {
		fmt.Printf("Err send REG MES %v RADD: %+v", err, raddr)
		return
	}

	_, _, err = conn.ReadFromUDP(buf[:])
	if err != nil {
		fmt.Println("Err Getback from REG request", err)
		return
	}

	fmt.Println("Get back ", string(buf[:]))

	tik := time.NewTicker(time.Second * 1)
	for {
		<-tik.C
		p := model.Request{
			Action: "GET",
			PeerID: rPeerID,
		}
		b, err = json.Marshal(p)
		if err != nil {
			fmt.Println("Err marshal ", err, p)
			return
		}
		_, err = conn.WriteToUDP(b, raddr)
		if err != nil {
			fmt.Println("Err send REG MES", err)
			return
		}

		n, _, err := conn.ReadFromUDP(buf[:])
		if err != nil {
			fmt.Println("Err ReadFromUDP", err)
			return
		}
		var resp model.Response
		err = json.Unmarshal(buf[:n], &resp)
		if err != nil {
			fmt.Println("Err Unmarshal", err)
			return
		}

		if resp.RAddr != nil {
			log.Print("Peer address: ", *resp.RAddr)
			peerAddr, err := net.ResolveUDPAddr("udp4", *resp.RAddr)
			if err != nil {
				log.Fatal(err)
			}
			go listen(conn)
			for {
				fmt.Print("Input message >_ ")
				message := make([]byte, 2048)
				fmt.Scanln(&message)
				messageRequest := ChatRequest{
					"Chat",
					*lPeer,
					string(message),
				}
				rq, err := json.Marshal(messageRequest)
				if err != nil {
					log.Print("Error: ", err)
					continue
				}
				fmt.Println("Write to peer add: ", peerAddr.String())
				_, err = conn.WriteToUDP(rq, peerAddr)
				fmt.Println("Result write to UDP ", err)
			}
		}
		fmt.Println("No connect from this peerID: yet")
	}
}

type ChatRequest struct {
	Action   string
	Username string
	Message  string
}

func listen(conn *net.UDPConn) {
	for {
		buf := make([]byte, 2048)
		n, addr, err := conn.ReadFromUDP(buf)
		if err != nil {
			log.Print(err)
			continue
		}
		log.Print("Message from ", addr.IP)
		var message ChatRequest
		err = json.Unmarshal(buf[:n], &message)
		if err != nil {
			log.Print(err)
			continue
		}
		fmt.Println(message.Username, ":", message.Message)
	}
}
