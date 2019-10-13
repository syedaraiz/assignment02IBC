package satoshi

import (
	"encoding/gob"
	//For Printing On Console
	"fmt"
	//To Get Command Line Arguments
	"os"
	//For Networking
	"log"
	"net"
	"time"
	//For Configuration
	configuration "github.com/syedaraiz/assignment02IBC/Configuration"
	structures "github.com/syedaraiz/assignment02IBC/Structures"
	blockchain "github.com/syedaraiz/assignment02IBC/Blockchain"
	//To Get Random Clients
	"math/rand"
	 "strconv"
)
var satoshi structures.Client
var hasSetup = false
var myBalance = 100
//Connected Users
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
func handleConnection(c net.Conn, newClient structures.Client) {
	myConn, err := net.Dial(newClient.ListeningIP, newClient.ListeningPort)
	if err != nil {
		println("Invalid IP/Port...")
		return
	}
	satoshi.ConnectedPeers = append(satoshi.ConnectedPeers, newClient)
	connectedClientsConnections = append(connectedClientsConnections, myConn)
	log.Println("A client has connected", c.RemoteAddr())

	if len(satoshi.ConnectedPeers) == configuration.MaxNumberOfUsersToInitiateBlockchain && !hasSetup {
		//Connect Peers
		for i := 0; i < len(satoshi.ConnectedPeers); i++ {
			searchArray := []int{i}
			for j := 0; j <= i; j++ {
				randomClient := rand.Intn(len(satoshi.ConnectedPeers))
				if j == configuration.MaxNumberOfUsersToInitiateBlockchain-1 {
					break
				}
				for hasElement(searchArray, randomClient) {
					randomClient = rand.Intn(len(satoshi.ConnectedPeers))
				}
				searchArray = append(searchArray, randomClient)
				satoshi.ConnectedPeers[i].ConnectedPeers = append(satoshi.ConnectedPeers[i].ConnectedPeers, satoshi.ConnectedPeers[randomClient])
			}
			satoshi.ConnectedPeers[i].ChainHead = satoshi.ChainHead
			encoder1 := gob.NewEncoder(connectedClientsConnections[i])
			encodedClient := satoshi.ConnectedPeers[i]
			err2 := encoder1.Encode(encodedClient)
			if err2 != nil {
				fmt.Println(err2)
			}

		}
		hasSetup = true
		println("Setup Completed.")
		makeATransaction(50, -1, rand.Intn(len(satoshi.ConnectedPeers)))
	}
}
func makeATransaction(numberOfCoins int,from int, to int) {
		randomMiner := rand.Intn(len(satoshi.ConnectedPeers))
		for randomMiner==to || randomMiner==from{
			randomMiner = rand.Intn(len(satoshi.ConnectedPeers))
		}
		fmt.Println("Satoshi Chose ",satoshi.ConnectedPeers[randomMiner].Name," For Transaction.")
		fmt.Println("Transaction Details are as under: ")
		fromUser := "Satoshi"
		if (from!=-1) {
			fromUser = satoshi.ConnectedPeers[from].Name
		}
		fmt.Println("From: ", fromUser)
		fmt.Println("To: ", satoshi.ConnectedPeers[to].Name)
		fmt.Println("Transaction Amount: ", strconv.Itoa(numberOfCoins), configuration.MyCurrencyShortForm,"(",configuration.MyCurrency,")")

		transaction := fromUser+"to"+satoshi.ConnectedPeers[to].Name+"="+strconv.Itoa(numberOfCoins)
		satoshi.ChainHead = blockchain.InsertBlock(transaction, satoshi.ChainHead)

		encoder := gob.NewEncoder(connectedClientsConnections[randomMiner])
		encoder.Encode("floodTransaction")
		time.Sleep(2*time.Second)
		//newEncoder := gob.NewEncoder(connectedClientsConnections[randomMiner])
		encoder.Encode(satoshi)
}
func acceptAndHandleConnections(ln net.Listener) {
	for len(satoshi.ConnectedPeers) < configuration.MaxNumberOfUsersToInitiateBlockchain {
		conn, err := ln.Accept()
		if err != nil {
			println("An Error has occured while accepting a connection.")
		} else {
			var thisClient structures.Client
			decoder := gob.NewDecoder(conn)
			decoder.Decode(&thisClient)
			fmt.Println(thisClient.Name, " has connected.")

			go handleConnection(conn, thisClient)
		}
	}
	for true {
		fmt.Scanln()
	}
}

//RunSatoshi ...
func RunSatoshi() {
	satoshi.AccountBalance = configuration.InitalSatoshiBalance
	satoshi.ListeningIP = "tcp"
	satoshi.ListeningPort = getCommandLineArguments()[0]
	satoshi.ChainHead = blockchain.InsertBlock("GenesisBlock", nil)
	fmt.Println("Satoshi is listening at ",getCommandLineArguments()[0])
	//Creating A Listener At a Port
	ln, err := net.Listen(satoshi.ListeningIP, ":"+satoshi.ListeningPort)

	if err != nil {
		println("Error Occured")
		return
	}
	go acceptAndHandleConnections(ln)
	for true {
		var name string
		fmt.Scanln(&name)
	}
}
