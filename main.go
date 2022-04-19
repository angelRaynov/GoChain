package main

import (
	"github.com/angelRaynov/GoChain/blockchain"
	"log"
	"strconv"
)

func main() {

	chain := blockchain.InitBlockChain()

	log.Println("BLOCKCHAIN CREATED")


	chain.AddBlock("first block after genesis")
	chain.AddBlock("second block after genesis")
	chain.AddBlock("third block after genesis")

	for _, block := range chain.Blocks {
		log.Printf("Prev hash: %x\n", block.PrevHash)
		log.Printf("Data: %s\n", block.Data)
		log.Printf("Hash: %x\n", block.Hash)

		pow := blockchain.NewProofOfWork(block)

		log.Printf("POW: %s\n\n", strconv.FormatBool(pow.Validate()))

	}
}
