package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"strconv"
)

func main() {
	http.HandleFunc("/about", func(rw http.ResponseWriter, r *http.Request) {
		fmt.Println("LALAA 2")
		fmt.Fprintf(rw, "This is from holepunching ")
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
func doServer(port int) {
	msgBuf := make([]byte, 1024)

	// Initiatlize a UDP listener
	ln, err := net.ListenUDP("udp4", &net.UDPAddr{Port: port})
	if err != nil {
		fmt.Printf("Unable to listen on :%d\n", port)
		return
	}

	fmt.Printf("Listening on :%d\n", port)

	for {
		fmt.Println("---")
		// Await incoming packets
		rcvLen, addr, err := ln.ReadFrom(msgBuf)
		if err != nil {
			fmt.Println("Transaction was initiated but encountered an error!")
			continue
		}

		fmt.Printf("Received a packet from: %s\n\tSays: %s\n",
			addr.String(), msgBuf[:rcvLen])

		// Let the client confirm a hole was punched through to us
		reply := "お帰り～"
		copy(msgBuf, []byte(reply))
		_, err = ln.WriteTo(msgBuf[:len(reply)], addr)

		if err != nil {
			fmt.Println("Socket closed unexpectedly!")
			continue
		}

		fmt.Printf("Sent reply to %s\n\tReply: %s\n",
			addr.String(), msgBuf[:len(reply)])
	}
}

func doClient(remote string, port int) {
	msgBuf := make([]byte, 1024)

	// Resolve the passed address as UDP4
	toAddr, err := net.ResolveUDPAddr("udp4", remote+":"+strconv.Itoa(port))
	if err != nil {
		fmt.Printf("Could not resolve %s:%d\n", remote, port)
		return
	}

	fmt.Printf("Trying to punch a hole to %s:%d\n", remote, port)

	// Initiate the transaction (force IPv4 to demo firewall punch)
	conn, err := net.DialUDP("udp4", nil, toAddr)
	defer conn.Close()

	if err != nil {
		fmt.Printf("Unable to connect to %s:%d\n", remote, port)
		return
	}

	// Initiate the transaction, creating the hole
	msg := "ただいま～"
	fmt.Fprintf(conn, msg)
	fmt.Printf("Sent a UDP packet to %s:%d\n\tSent: %s\n", remote, port, msg)

	// Await a response through our firewall hole
	msgLen, fromAddr, err := conn.ReadFromUDP(msgBuf)
	if err != nil {
		fmt.Printf("Error reading UDP response!\n")
		return
	}

	fmt.Printf("Received a UDP packet back from %s:%d\n\tResponse: %s\n",
		fromAddr.IP, fromAddr.Port, msgBuf[:msgLen])

	fmt.Println("Success: NAT traversed! ^-^")
}
