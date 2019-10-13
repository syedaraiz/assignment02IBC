package main

import (
	user "github.com/syedaraiz/assignment02IBC/User"
)

func main() {
	//To run give command line arguments in manner
	//go run client.go [NameOfUser] [ClientListeningPort] [Satoshi Connecting Port]
	user.RunUser()
}
