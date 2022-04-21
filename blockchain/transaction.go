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

//TxOutput From
type TxOutput struct {
	//Value would be representative of the amount of coins in a transaction
	Value int
	//PubKey is needed to "unlock" any coins within an Output. This indicated that YOU are the one that sent it.
	//You are indentified by your PubKey
	PubKey string
}

//TxInput is representative of a reference to a previous TxOutput (To)
type TxInput struct {
	ID []byte
	//ID will find the Transaction that a specific output is inside of
	Out int
	//Out will be the index of the specific output we found within a transaction.
	//For example if a transaction has 4 outputs, we can use this "Out" field to specify which output we are looking for
	Sig string
	//This would be a script that adds data to an outputs' PubKey
	//however for this tutorial the Sig will be indentical to the PubKey.
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

func (in *TxInput) CanUnlock(data string) bool {
	return in.Sig == data
}

func (out *TxOutput) CanBeUnlocked(data string) bool {
	return out.PubKey == data
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
