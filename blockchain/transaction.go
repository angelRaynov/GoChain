package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"log"
)

const Reward = 100

type Transaction struct {
	ID      []byte
	Inputs  []TxInput
	Outputs []TxOutput
}

func CoinbaseTx(toAddress, data string) *Transaction {
	if data == "" {
		data = fmt.Sprintf("Coins to %s", toAddress)
	}

	//Since this is the "first" transaction of the block, it has no previous output to reference.
	//This means that we initialize it with no ID, and it's OutputIndex is -1
	txIn := TxInput{[]byte{}, -1, data}

	txOut := TxOutput{Reward, toAddress}

	tx := Transaction{nil, []TxInput{txIn}, []TxOutput{txOut}}

	return &tx
}

func (tx *Transaction) SetID() {
	var encoded bytes.Buffer
	var hash [32]byte

	encoder := gob.NewEncoder(&encoded)
	err := encoder.Encode(tx)
	Handle(err)

	hash = sha256.Sum256(encoded.Bytes())
	tx.ID = hash[:]
}

//IsCoinbase checks a transaction and will only return true if it is a newly minted "coin"
func (tx *Transaction) IsCoinbase() bool {
	return len(tx.Inputs) == 1 && len(tx.Inputs[0].ID) == 0 && tx.Inputs[0].Out == -1
}

func NewTransaction(from, to string, amount int, chain *BlockChain) *Transaction {
	var inputs []TxInput
	var outputs []TxOutput

	//Find Spendable Outputs
	acc, validOutputs := chain.FindSpendableOutputs(from, amount)

	//Check if we have enough money to send the amount that we are asking
	if acc < amount {
		log.Panic("Not enough funds!")
	}

	//If we do, make inputs that point to the outputs we are spending
	for txId, outs := range validOutputs {
		txId, err := hex.DecodeString(txId)
		Handle(err)

		for _, out := range outs {
			input := TxInput{txId, out, from}
			inputs = append(inputs, input)
		}
	}

	outputs = append(outputs, TxOutput{amount, to})

	//If there is any leftover money, make new outputs from the difference.
	if acc > amount {
		outputs = append(outputs, TxOutput{acc - amount, from})
	}

	//Initialize a new transaction with all the new inputs and outputs we made
	tx := Transaction{nil, inputs, outputs}

	//Set a new ID, and return it.
	tx.SetID()

	return &tx
}
