package blockchain

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

func (in *TxInput) CanUnlock(data string) bool {
	return in.Sig == data
}

func (out *TxOutput) CanBeUnlocked(data string) bool {
	return out.PubKey == data
}
