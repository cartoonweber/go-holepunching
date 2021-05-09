package main

import (
    "fmt"
    "io"
    "net"
    "time"
)

func main() {
  runUDPServer()
}

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

func handleUDPClient(conn *net.UDPConn) {
  var buff [512]byte

  _, addr, err := conn.ReadFromUDP(buff[:])

  if err != nil {
    fmt.Println("Err readfrom UDP", err)
    return
  }

  fmt.Println("ServeIP: ", addr.IP.String(), " port: ", addr.Port)
  daytime := time.Now().String()

  _, err = conn.WriteToUDP([]byte(daytime), addr)

  if err != nil {
    fmt.Println("Err write to UDP", err)
    return
  }
}
func runTCp() {
    service, err := net.ResolveTCPAddr("tcp", ":8080")
    if err != nil {
        panic(err)
    }
    listener, err := net.ListenTCP("tcp", service)
    if err != nil {
        panic(err)
    }

    for {
        conn, err := listener.Accept()
        if err != nil {
            fmt.Println(err)
            continue
        }

        go handleClient(conn)
    }
}

func handleClient(conn net.Conn) {
    defer conn.Close()

    var buf [512]byte
    for {
        n, err := conn.Read(buf[:])
        if err != nil && err == io.EOF {
            return
        }

        if err != nil {
            panic(err)
        }

        fmt.Println("Emit back ", string(buf[:n]))
        _, err = conn.Write(buf[:n])
        if err != nil {
            panic(err)
        }
    }
}

