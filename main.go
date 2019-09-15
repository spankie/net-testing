package main

import (
	"errors"
	"fmt"
	"log"
	"net"
	"strings"
)

// Server defines the minimum contract our
// TCP and UDP server implementations must satisfy.
type Server interface {
	handleServer(*net.UDPAddr)
	Run() (*net.UDPAddr, error)
	Close() error
}

// NewServer creates a new Server using given protocol
// and addr.
func NewServer(protocol, addr string) (Server, error) {
	switch strings.ToLower(protocol) {
	case "udp":
		return &UDPServer{
			addr: addr,
		}, nil
	}
	return nil, errors.New("Invalid protocol given")
}

// UDPServer holds the necessary structure for our
// UDP server.
type UDPServer struct {
	addr   string
	server *net.UDPConn
}

// Run starts the UDP server.
func (u *UDPServer) Run() (*net.UDPAddr, error) {
	laddr, err := net.ResolveUDPAddr("udp", u.addr)
	if err != nil {
		return nil, errors.New("could not resolve UDP addr")
	}
	return laddr, nil

	// u.server, err = net.ListenUDP("udp", laddr)
	// if err != nil {
	// return errors.New("could not listen on UDP")
	// }
	//
	// return u.handleConnections()
}

func (u *UDPServer) handleServer(laddr *net.UDPAddr) {
	var err error
	u.server, err = net.ListenUDP("udp", laddr)
	if err != nil {
		log.Println(err)
		return
	}
	go func() {
		err := u.handleConnections()
		if err != nil {
			log.Printf("error occured: %v", err)
		}
	}()

}

func (u *UDPServer) handleConnections() error {
	var err error
	for {
		buf := make([]byte, 2048)
		n, conn, err := u.server.ReadFromUDP(buf)
		if err != nil {
			log.Println(err)
			break
		}
		if conn == nil {
			continue
		}

		go u.handleConnection(conn, buf[:n])
	}
	return err
}

func (u *UDPServer) handleConnection(addr *net.UDPAddr, cmd []byte) {
	r, err := u.server.WriteToUDP([]byte(fmt.Sprintf("Request recieved: %s", cmd)), addr)
	if err != nil {
		log.Printf("error occured writing to the UDP connection: %v", err)
	}
	log.Printf("wrote %v bytes to the connection", r)
}

// Close ensures that the UDPServer is shut down gracefully.
func (u *UDPServer) Close() error {
	return u.server.Close()
}

func main() {
	port := 8081
	// Start the new server
	udp, err := NewServer("udp", fmt.Sprintf(":%v", port))
	if err != nil {
		log.Println("error starting UDP server")
		return
	}

	// Run the servers in goroutines to stop blocking
	laddr, err := udp.Run()
	if err != nil {
		log.Printf("err occured running the server: %#v\n", err)
		udp.Close()
		log.Fatal(err)
	}
	for {
		udp.handleServer(laddr)
	}
}
