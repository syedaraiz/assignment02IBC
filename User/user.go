package user

import (
	"encoding/gob"
	"fmt"
	"net"
	"os"
	"bytes"
	"time"
	structures "github.com/syedaraiz/assignment02IBC/Structures"
	blockchain "github.com/syedaraiz/assignment02IBC/Blockchain"
)
var floodedBlocks [] *blockchain.Block
var thisClient structures.Client
var connectedPeersConnection []net.Conn
var listeningIP string
var listeningPort string
func isTransactionValid(transaction string) bool {
	return true
}
func alreadyFlooded(block *blockchain.Block) bool {
	for i:=0; i<len(floodedBlocks); i++ {
		if (bytes.Equal(blockchain.CalculateHash(floodedBlocks[i]),blockchain.CalculateHash(block))){
			return true
		}
	}
	return false
}
func floodTransactions(){
	fmt.Println("Flooding...")
	if (isTransactionValid("transaction")){
		for i:=0 ; i<len(connectedPeersConnection); i++{
			commandEncoder := gob.NewEncoder(connectedPeersConnection[i])
			commandEncoder.Encode("floodTransaction")
			time.Sleep(time.Second)
			fmt.Println(thisClient.Name)
			commandEncoder.Encode(thisClient.ChainHead)
		}
	}
	fmt.Println("Flooded...")
}

func handleConnection(c net.Conn) {
	for true {
		command := "-1"
		commandDecoder := gob.NewDecoder(c)
		commandDecoder.Decode(&command)
		if command == "floodTransaction"{
			time.Sleep(2*time.Second)
			var tempClient structures.Client
			err3 := commandDecoder.Decode(&tempClient)
			if (alreadyFlooded(tempClient.ChainHead)){
				fmt.Println("Blockchain Flooded...")
			} else {
					floodedBlocks = append(floodedBlocks, tempClient.ChainHead)
					thisClient.ChainHead = tempClient.ChainHead
					//go floodTransactions()
			}
			if err3 != nil {
				fmt.Println(err3)
				println("Error Occured While Decoding Blockchain From Server.")
			}
			blockchain.ListBlocks(tempClient.ChainHead)
		}
	}
}
func handleSatoshiConnection(c net.Conn) {
	//Client Recieved With its connected peers
	decoder := gob.NewDecoder(c)
	err := decoder.Decode(&thisClient)
	if err != nil {
		println("Error Occured While Decoding Client From Server.")
	}
	fmt.Println("length->",len(thisClient.ConnectedPeers))
	for i:=0 ; i<len(thisClient.ConnectedPeers); i++ {
		myConn, err := net.Dial(thisClient.ListeningIP, thisClient.ListeningPort)
		if err != nil {
		println("Invalid IP/Port...")
		return
	  }
		connectedPeersConnection = append(connectedPeersConnection, myConn)
	}
	//Now working on other requests
	for true {
		command := "-1"
		commandDecoder := gob.NewDecoder(c)
		commandDecoder.Decode(&command)
		if command == "floodTransaction"{
			time.Sleep(2*time.Second)
			var tempClient structures.Client
			err3 := commandDecoder.Decode(&tempClient)
			thisClient.ChainHead = tempClient.ChainHead
			blockchain.ListBlocks(thisClient.ChainHead)
			//go floodTransactions()
			if err3 != nil {
				fmt.Println(err3)
				println("Error Occured While Decoding Blockchain From Server.")
			}
			//blockchain.ListBlocks(tempClient.ChainHead)
		}
	}
}

func acceptAndHandleConnections(ln net.Listener) {
	count := 0
	for true {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("An Error has occured while accepting a connection.")
		} else if count == 0{
			go handleSatoshiConnection(conn)
		} else if count > 0{
			fmt.Println("User Connected...")
			go handleConnection(conn)
		}
		count = count + 1
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
	fmt.Println(thisClient.Name," is listening at ","127.0.0.1",listeningPort,"...")
	ln, err := net.Listen(listeningIP, listeningPort)
	if err != nil {
		println("Error Occured")
		return
	}
	go acceptAndHandleConnections(ln)
}

//RunUser ...
func RunUser() {
	startListening()
	conn, err := net.Dial("tcp", ":"+getCommandLineArguments()[2])
	encoder := gob.NewEncoder(conn)
	err2 := encoder.Encode(thisClient)
	if err2 != nil {
		fmt.Println(err2)
	}
	if err != nil {

		// handle error
	}
	for true {
		fmt.Scanln()
	}
}
