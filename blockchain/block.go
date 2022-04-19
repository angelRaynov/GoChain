package blockchain

type Block struct {
	Hash     []byte
	Data     []byte
	PrevHash []byte
	Nonce    int
}

type BlockChain struct {
	Blocks []*Block
}

func (chain *BlockChain) AddBlock(data string) {
	//get the last block
	prevBlock := chain.Blocks[len(chain.Blocks)-1]

	newBlock := CreateBlock(data, prevBlock.Hash)

	chain.Blocks = append(chain.Blocks, newBlock)
}

func CreateBlock(data string, prevHash []byte) *Block {
	block := &Block{
		[]byte{},
		[]byte(data),
		prevHash,
		0,
	}

	pow := NewProofOfWork(block)
	nonce, hash := pow.Run()

	block.Hash = hash[:]
	block.Nonce = nonce

	return block
}

func Genesis() *Block {
	return CreateBlock("Genesis", []byte{})
}

func InitBlockChain() *BlockChain {
	return &BlockChain{[]*Block{Genesis()}}
}
