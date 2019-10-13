import (
	"encoding/gob"
	//For Printing On Console
	"fmt"
	//To Get Command Line Arguments
	"os"
	//For Networking
	"log"
	"net"
	//For Configuration
	client "github.com/syedaraiz/Client"
	configuration "github.com/syedaraiz/ConfigurationPackage"
	//To Get Random Clients
	"math/rand"
)

var hasSetup = false
var myBalance = 100

//Connected Users
var connectedClients []client.Client
var connectedClientsConnections []net.Conn

//Function that returns array of input command line Arguments
func getCommandLineArguments() []string {
	argsWithoutProg := os.Args[1:]
	return argsWithoutProg
}
func hasElement(searchArray []int, elementToSearch int) bool {
	for i := 0; i < len(searchArray); i++ {
		if searchArray[i] == elementToSearch {
			return true
		}
	}
	return false
}
func handleConnection(c net.Conn, newClient client.Client) {
	myConn, err := net.Dial(newClient.ListeningIP, newClient.ListeningPort)
	if err != nil {
		println("Invalid IP/Port...")
		return
	}
	connectedClients = append(connectedClients, newClient)
	connectedClientsConnections = append(connectedClientsConnections, myConn)
	log.Println("A client has connected", c.RemoteAddr())
	if len(connectedClients) == configuration.MaxNumberOfUsersToInitiateBlockchain && !hasSetup {
		//Connect Peers
		for i := 0; i < len(connectedClients); i++ {
			searchArray := []int{i}
			for j := 0; j <= i; j++ {
				randomClient := rand.Intn(len(connectedClients))
				if j == configuration.MaxNumberOfUsersToInitiateBlockchain-1 {
					break
				}
				for hasElement(searchArray, randomClient) {
					randomClient = rand.Intn(len(connectedClients))
				}
				searchArray = append(searchArray, randomClient)
				connectedClients[i].ConnectedPeers = append(connectedClients[i].ConnectedPeers, connectedClients[randomClient])
			}
			encoder1 := gob.NewEncoder(connectedClientsConnections[i])
			encodedClient := connectedClients[i]
			err2 := encoder1.Encode(encodedClient)
			if err2 != nil {
				fmt.Println(err2)
			}

		}
		hasSetup = true
		println("Setup Completed.")
	}
}

func acceptAndHandleConnections(ln net.Listener) {
	for true {
		conn, err := ln.Accept()
		if err != nil {
			println("An Error has occured while accepting a connection.")
		} else {
			var thisClient client.Client
			decoder := gob.NewDecoder(conn)
			decoder.Decode(&thisClient)
			fmt.Println(thisClient.Name, " has connected.")

			if len(connectedClients) < configuration.MaxNumberOfUsersToInitiateBlockchain {
				go handleConnection(conn, thisClient)
			} else {
				
				println("Max Number Of Users Connected.")
				conn.Close()
				return
			}
		}
	}
}
