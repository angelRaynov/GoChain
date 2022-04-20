package main

import (
	"flag"
	"github.com/angelRaynov/GoChain/blockchain"
	"log"
	"os"
	"runtime"
	"strconv"
)

type CommandLine struct {
	blockchain *blockchain.BlockChain
}

func (cli *CommandLine) PrintUsage() {
	log.Println("Usage: ")
	log.Println("add -block <BLOCK_DATA> - add a block to the chain")
	log.Println("print - prints the blocks in the chain")
}

func (cli *CommandLine) validateArgs() {
	if len(os.Args) < 2 {
		cli.PrintUsage()

		//shut goroutines gracefully
		runtime.Goexit()
	}
}

func (cli *CommandLine) addBlock(data string) {
	cli.blockchain.AddBlock(data)
	log.Println("Added new block")
}

func (cli *CommandLine) printChain() {
	iterator := cli.blockchain.Iterator()

	for {
		block := iterator.Next()
		log.Printf("Previous hash: %x\n", block.PrevHash)
		log.Printf("Data: %s\n", block.Data)
		log.Printf("Hash: %x\n", block.Hash)

		pow := blockchain.NewProofOfWork(block)

		log.Printf("POW: %s\n\n", strconv.FormatBool(pow.Validate()))

		//genesis block has no prev hash
		if len(block.PrevHash) == 0 {
			break
		}
	}
}

func (cli *CommandLine) run() {
	cli.validateArgs()

	addBlockCmd := flag.NewFlagSet("add", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("print", flag.ExitOnError)
	addBlockData := addBlockCmd.String("block", "", "Block data")

	switch os.Args[1] {
	case "add":
		err := addBlockCmd.Parse(os.Args[2:])
		blockchain.Handle(err)

	case "print":
		err := printChainCmd.Parse(os.Args[2:])
		blockchain.Handle(err)

	default:
		cli.PrintUsage()
		runtime.Goexit()
	}

	// parsed() will return true if the object was used or has been called
	if addBlockCmd.Parsed() {
		if *addBlockData == "" {
			addBlockCmd.Usage()
			runtime.Goexit()
		}
		cli.addBlock(*addBlockData)
	}

	if printChainCmd.Parsed() {
		cli.printChain()
	}
}

func main() {
	defer os.Exit(0)

	chain := blockchain.InitBlockChain()
	defer chain.DB.Close()

	cli := CommandLine{chain}

	cli.run()
}
