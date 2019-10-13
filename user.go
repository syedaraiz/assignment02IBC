import (
	"encoding/gob"
	"fmt"
	"net"
	"os"

	client "github.com/syedaraiz/assignment02IBC"
)

var thisClient client.Client

var listeningIP string
var listeningPort string

func handleConnection(c net.Conn) {
	//Client Recieved With its connected peers
	var thisClient structures.Client
	decoder := gob.NewDecoder(c)
	err := decoder.Decode(&thisClient)
	if err != nil {
		println("Error Occured While Decoding Client From Server.")
	}
	//Now working on other requests
	for true {
		recvdSlice := make([]byte, 30)
		c.Read(recvdSlice)
		fmt.Println(string(recvdSlice))
	}
}

func acceptAndHandleConnections(ln net.Listener) {
	for true {
		conn, err := ln.Accept()
		if err != nil {
			println("An Error has occured while accepting a connection.")
		} else {
			go handleConnection(conn)
			print("Connection Accepted...")
		}
	}
}
func getCommandLineArguments() []string {
	argsWithoutProg := os.Args[1:]
	return argsWithoutProg
}

//Client ...
func startListening() {
	//Initializations

	listeningIP = "tcp"
	listeningPort = ":" + getCommandLineArguments()[1]
	thisClient.Name = getCommandLineArguments()[0]
	thisClient.AccountBalance = 0
	thisClient.ListeningPort = listeningPort
	thisClient.ListeningIP = listeningIP

	ln, err := net.Listen(listeningIP, listeningPort)
	if err != nil {
		println("Error Occured")
		return
	}
	go acceptAndHandleConnections(ln)
}
