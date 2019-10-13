package blockchain
import (
	"bytes"
	"crypto/sha256"
	"fmt"
)

//Block Of Entire Blockchain
type Block struct {
	transaction string
	prevHash    []byte
	prevPointer *Block
}

func calculateHash(data Block) []byte {
	byteTransaction := []byte(data.transaction)
	appendedTransaction := append(data.prevHash, byteTransaction...)
	hashedValue := sha256.Sum256(appendedTransaction)
	return hashedValue[:]
}

//InsertBlock inserts at end of Blockchain
func InsertBlock(transaction string, chainHead *Block) *Block {
	var toBeInsertedBlock Block
	toBeInsertedBlock.transaction = transaction
	if chainHead == nil {
		toBeInsertedBlock.prevHash = []byte{}
		toBeInsertedBlock.prevPointer = nil
	}
	toBeInsertedBlock.prevPointer = chainHead
	if chainHead != nil {
		toBeInsertedBlock.prevHash = calculateHash(*chainHead)
	}
	return &toBeInsertedBlock
}

//ListBlocks : Printing All Blocks Inside Blockchain
func ListBlocks(chainHead *Block) {
	if chainHead == nil {
		fmt.Println("Blockchain is empty.")
	}
	for currentBlock := chainHead; currentBlock != nil; currentBlock = currentBlock.prevPointer {
		fmt.Print(currentBlock.transaction)
		if currentBlock.prevPointer != nil {
			fmt.Print("<-")
		}
	}
	fmt.Println()
}

//ChangeBlock : This Function Changes Block With In Blockchain
func ChangeBlock(oldTrans string, newTrans string, chainHead *Block) {
	if chainHead == nil {
		fmt.Println("Blockchain is empty.")
	}
	for currentBlock := chainHead; currentBlock.prevPointer != nil; currentBlock = currentBlock.prevPointer {
		if currentBlock.transaction == oldTrans {
			currentBlock.transaction = newTrans
			return
		}
	}
}

//VerifyChain : Evaluates if Blockchain is altered with or not
func VerifyChain(chainHead *Block) bool {
	if chainHead == nil {
		fmt.Println("Blockchain is empty.")
	}
	for currentBlock := chainHead; currentBlock.prevPointer != nil; currentBlock = currentBlock.prevPointer {
		if bytes.Compare(calculateHash(*(currentBlock.prevPointer)), currentBlock.prevHash) != 0 {
			fmt.Println("Blockchain is altered with.")
			return false
		}
	}
	fmt.Println("Blockchain Verified.")
	return true
}
