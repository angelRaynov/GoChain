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
}

func (cli *CommandLine) printUsage() {

	log.Println("Usage: ")

	log.Println("getbalance -address ADDRESS - get balance for ADDRESS")

	log.Println("createblockchain -address ADDRESS creates a blockchain and rewards the mining fee")

	log.Println("printchain - Prints the blocks in the chain")

	log.Println("send -from FROM -to TO -amount AMOUNT - Send amount of coins from one address to another")

}

func (cli *CommandLine) getBalance(address string) {
	chain := blockchain.ContinueBlockchain(address)
	defer chain.DB.Close()

	balance := 0
	UTXOs := chain.FindUnspentTransactionOutputs(address)

	for _, out := range UTXOs {
		balance += out.Value
	}

	log.Printf("Balance of %s: %d\n", address, balance)
}

func (cli *CommandLine) send(from, to string, amount int) {
	chain := blockchain.ContinueBlockchain(from)
	defer chain.DB.Close()

	tx := blockchain.NewTransaction(from, to, amount, chain)

	chain.AddBlock([]*blockchain.Transaction{tx})
	log.Println("Success!")
}

func (cli *CommandLine) validateArgs() {
	if len(os.Args) < 2 {
		cli.printUsage()

		//shut goroutines gracefully
		runtime.Goexit()
	}
}

func (cli *CommandLine) createBlockChain(address string) {
	newChain := blockchain.InitBlockChain(address)
	newChain.DB.Close()
	log.Println("Finished creating chain")
}

func (cli *CommandLine) run() {
	cli.validateArgs()

	getBalanceCmd := flag.NewFlagSet("getbalance", flag.ExitOnError)
	createBlockchainCmd := flag.NewFlagSet("createblockchain", flag.ExitOnError)
	sendCmd := flag.NewFlagSet("send", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)

	getBalanceAddress := getBalanceCmd.String("address", "", "The address to get balance for")
	createBlockchainAddress := createBlockchainCmd.String("address", "", "The address to send genesis block reward to")
	sendFrom := sendCmd.String("from", "", "Source wallet address")
	sendTo := sendCmd.String("to", "", "Destination wallet address")
	sendAmount := sendCmd.Int("amount", 0, "Amount to send")

	switch os.Args[1] {
	case "getbalance":
		err := getBalanceCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "createblockchain":
		err := createBlockchainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "printchain":
		err := printChainCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	case "send":
		err := sendCmd.Parse(os.Args[2:])
		if err != nil {
			log.Panic(err)
		}
	default:
		cli.printUsage()
		runtime.Goexit()
	}

	if getBalanceCmd.Parsed() {
		if *getBalanceAddress == "" {
			getBalanceCmd.Usage()
			runtime.Goexit()
		}
		cli.getBalance(*getBalanceAddress)
	}

	if createBlockchainCmd.Parsed() {
		if *createBlockchainAddress == "" {
			createBlockchainCmd.Usage()
			runtime.Goexit()
		}
		cli.createBlockChain(*createBlockchainAddress)
	}

	if printChainCmd.Parsed() {
		cli.printChain()
	}

	if sendCmd.Parsed() {
		if *sendFrom == "" || *sendTo == "" || *sendAmount <= 0 {
			sendCmd.Usage()
			runtime.Goexit()
		}

		cli.send(*sendFrom, *sendTo, *sendAmount)
	}
}

//printChain will display the entire contents of the blockchain
func (cli *CommandLine) printChain() {
	chain := blockchain.ContinueBlockchain("")
	defer chain.DB.Close()
	iterator := chain.Iterator()

	for {
		block := iterator.Next()
		log.Printf("Previous hash: %x\n", block.PrevHash)
		log.Printf("hash: %x\n", block.Hash)
		pow := blockchain.NewProofOfWork(block)
		log.Printf("Pow: %s\n", strconv.FormatBool(pow.Validate()))
		log.Println()
		// This works because the Genesis block has no PrevHash to point to.
		if len(block.PrevHash) == 0 {
			break
		}
	}
}

func main() {
	defer os.Exit(0)
	cli := CommandLine{}
	cli.run()
}
