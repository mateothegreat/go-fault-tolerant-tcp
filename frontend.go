package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

func main() {
	conn, err := net.Dial("tcp", "0.0.0.0:8001")
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()

	//
	// SetReadDeadline sets the deadline for future Read calls and any currently-blocked Read call. A zero value for t means Read will not time out.
	//
	err := conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	if err != nil {
		fmt.Println(err)
		break
	}

	clientReader := bufio.NewReader(os.Stdin)
	serverReader := bufio.NewReader(con)

	//
	//
	//
	go func() {
		pulse := time.Tick(pulseInterval)

		for {
			select {
			case _, ok := <-heartbeat:
				if ok == false {
					return
				}
				fmt.Println("heartbeat")
			case r, ok := <-results:
				if ok == false {
					return
				}
				fmt.Printf("doSomething result %v\n", r)
			case <-time.After(timeout):
				return
			}
		}

	}()
	//
	// Create a goroutine that closes a session after 15 seconds
	//
	go func() {
		<-time.After(time.Duration(15) * time.Second)

		// Use SetLinger to force close the connection
		err := con.(*net.TCPConn).SetLinger(0)
		if err != nil {
			log.Printf("Error when setting linger: %s", err)
		}

		defer con.Close()
	}()

	for {
		// Waiting for the client request
		clientRequest, err := clientReader.ReadString('\n')

		switch err {
		case nil:
			clientRequest := strings.TrimSpace(clientRequest)
			if _, err = con.Write([]byte(clientRequest + "\n")); err != nil {
				log.Printf("failed to send the client request: %v\n", err)
			}
		case io.EOF:
			log.Println("client closed the connection")
			return
		default:
			log.Printf("client error: %v\n", err)
			return
		}

		// Waiting for the server response
		serverResponse, err := serverReader.ReadString('\n')

		switch err {
		case nil:
			log.Println(strings.TrimSpace(serverResponse))
		case io.EOF:
			log.Println("server closed the connection")
			return
		default:
			log.Printf("server error: %v\n", err)
			return
		}
	}
}

func (t *TcpConnector) SendHeartbeat() {

	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for x := range ticker.C {
		messageHeader, err := getMessageHeaderInBytes(129, 0, 0, t.AppState.Config.Username)
		if err != nil {
			t.AppState.Log.Errorf("Error Converting heartbeat request data to bytes: %+v, with seconds: %s", err, x)
			return
		}
		t.DataChannel <- messageHeader
	}
}
