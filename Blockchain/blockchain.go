package blockchain
import (
	"bytes"
	"crypto/sha256"
	"fmt"
)

//Block Of Entire Blockchain
type Block struct {
	Transaction string
	PrevHash    []byte
	PrevPointer *Block
}

func calculateHash(data Block) []byte {
	byteTransaction := []byte(data.Transaction)
	appendedTransaction := append(data.PrevHash, byteTransaction...)
	hashedValue := sha256.Sum256(appendedTransaction)
	return hashedValue[:]
}
//CalculateHash ...
func CalculateHash(data *Block) []byte {
	return calculateHash(*data)
}
//InsertBlock inserts at end of Blockchain
func InsertBlock(Transaction string, chainHead *Block) *Block {
	var toBeInsertedBlock Block
	toBeInsertedBlock.Transaction = Transaction
	if chainHead == nil {
		toBeInsertedBlock.PrevHash = []byte{}
		toBeInsertedBlock.PrevPointer = nil
	}
	toBeInsertedBlock.PrevPointer = chainHead
	if chainHead != nil {
		toBeInsertedBlock.PrevHash = calculateHash(*chainHead)
	}
	return &toBeInsertedBlock
}

//ListBlocks : Printing All Blocks Inside Blockchain
func ListBlocks(chainHead *Block) {
	if chainHead == nil {
		fmt.Println("Blockchain is empty.")
	}
	for currentBlock := chainHead; currentBlock != nil; currentBlock = currentBlock.PrevPointer {
		fmt.Print(currentBlock.Transaction)
		if currentBlock.PrevPointer != nil {
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
	for currentBlock := chainHead; currentBlock.PrevPointer != nil; currentBlock = currentBlock.PrevPointer {
		if currentBlock.Transaction == oldTrans {
			currentBlock.Transaction = newTrans
			return
		}
	}
}

//VerifyChain : Evaluates if Blockchain is altered with or not
func VerifyChain(chainHead *Block) bool {
	if chainHead == nil {
		fmt.Println("Blockchain is empty.")
	}
	for currentBlock := chainHead; currentBlock.PrevPointer != nil; currentBlock = currentBlock.PrevPointer {
		if bytes.Compare(calculateHash(*(currentBlock.PrevPointer)), currentBlock.PrevHash) != 0 {
			fmt.Println("Blockchain is altered with.")
			return false
		}
	}
	fmt.Println("Blockchain Verified.")
	return true
}
