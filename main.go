package main

import (
	"bytes"
	"crypto/sha256"
	"log"
)

type Block struct {
	Hash     []byte
	Data     []byte
	PrevHash []byte
}

type BlockChain struct {
	blocks []*Block
}

func main() {
	chain := InitBlockChain()
	chain.AddBlock("first block after genesis")
	chain.AddBlock("second block after genesis")
	chain.AddBlock("third block after genesis")

	for _, block := range chain.blocks {
		log.Printf("Prev hash: %x\n",block.PrevHash)
		log.Printf("Data: %s\n",block.Data)
		log.Printf("Hash: %x\n",block.Hash)
	}
}

func (b *Block) DeriveHash() {
	//join previous block info with the current one
	info := bytes.Join([][]byte{b.Data, b.PrevHash}, []byte{})

	//hash the actual info
	hash := sha256.Sum256(info)

	b.Hash = hash[:]
}

func CreateBlock(data string, prevHash []byte) *Block {
	block := &Block{
		[]byte{},
		[]byte(data),
		prevHash,
	}

	block.DeriveHash()

	return block
}

func (chain *BlockChain) AddBlock(data string)  {
	//get the last block
	prevBlock := chain.blocks[len(chain.blocks) - 1]

	newBlock := CreateBlock(data,prevBlock.Hash)

	chain.blocks = append(chain.blocks, newBlock)
}

func Genesis() *Block {
	return CreateBlock("Genesis", []byte{})
}

func InitBlockChain() *BlockChain {
	return &BlockChain{[]*Block{Genesis()}}
}