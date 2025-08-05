package main

import (
	"fmt"
	"strconv"

	"github.com/juniorguerra/golang-blockchain/blockchain"
)

func main() {

	chain := blockchain.InitBlockChain()

	chain.AddBlock("First Block After Genesis")
	chain.AddBlock("second Block After Genesis")
	chain.AddBlock("Third Block After Genesis")

	for _, block := range chain.Blocks {
		fmt.Printf("Previus Hash: %x\n", block.PrevHash)
		fmt.Printf("Data In Block: %s\n", block.Data)
		fmt.Printf("Hash: %x\n", block.Hash)

		pow := blockchain.NewProof(block)
		fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.Validate()))
		fmt.Println()
	}

	fmt.Println("Hello blockchain")

}
