package blockchain

import (
	"encoding/hex"
	"github.com/dgraph-io/badger"
	"log"
	"os"
	"runtime"
)

const (
	DBPath = "./tmp/blocks"

	//DBFile can be used to verify that the blockchain exists
	DBFile = "./tmp/blocks/MANIFEST"

	//GenesisData is arbitrary data for our genesis block
	GenesisData = "First Transaction from Genesis"
)

type BlockChain struct {
	LastHash []byte
	DB       *badger.DB
}

type BlockChainIterator struct {
	CurrentHash []byte
	DB          *badger.DB
}

func (chain *BlockChain) Iterator() *BlockChainIterator {
	return &BlockChainIterator{chain.LastHash, chain.DB}
}

func (iterator *BlockChainIterator) Next() *Block {
	var block *Block

	err := iterator.DB.View(func(txn *badger.Txn) error {
		item, err := txn.Get(iterator.CurrentHash)
		Handle(err)

		err = item.Value(func(val []byte) error {
			block = Deserialize(val)
			return nil
		})
		Handle(err)
		return err
	})
	Handle(err)

	iterator.CurrentHash = block.PrevHash

	return block
}

func (chain *BlockChain) AddBlock(transactions []*Transaction) {
	var lastHash []byte

	err := chain.DB.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("last_hash"))
		Handle(err)
		err = item.Value(func(val []byte) error {
			lastHash = val
			return nil
		})
		Handle(err)
		return err
	})
	Handle(err)

	newBlock := CreateBlock(transactions, lastHash)

	err = chain.DB.Update(func(txn *badger.Txn) error {
		err := txn.Set(newBlock.Hash, newBlock.Serialize())
		Handle(err)

		err = txn.Set([]byte("last_hash"), newBlock.Hash)
		Handle(err)

		chain.LastHash = newBlock.Hash
		return err
	})
	Handle(err)
}

func InitBlockChain(address string) *BlockChain {
	//open connection
	var lastHash []byte

	if DBExists(DBFile) {
		log.Println("Blockchain already exists")
		runtime.Goexit()
	}

	opts := badger.DefaultOptions(DBPath)
	db, err := badger.Open(opts)
	Handle(err)

	//check if there is blockchain created
	err = db.Update(func(txn *badger.Txn) error {

		cbtx := CoinbaseTx(address, GenesisData)
		genesis := Genesis(cbtx)

		log.Println("Genesis created")

		err = txn.Set(genesis.Hash, genesis.Serialize())
		Handle(err)

		err = txn.Set([]byte("last_hash"), genesis.Hash)
		Handle(err)

		lastHash = genesis.Hash

		return err

	})
	Handle(err)

	return &BlockChain{lastHash, db}
}

func ContinueBlockchain(address string) *BlockChain {
	if !DBExists(DBFile) {
		log.Println("No blockchain found, please create one first")
		runtime.Goexit()
	}

	var lastHash []byte

	opts := badger.DefaultOptions(DBPath)
	db, err := badger.Open(opts)
	Handle(err)

	err = db.Update(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("last_hash"))
		Handle(err)
		err = item.Value(func(val []byte) error {
			lastHash = val
			return nil
		})

		Handle(err)
		return err
	})
	Handle(err)

	return &BlockChain{lastHash, db}
}

//DBExists checks to see if we've initialized a database
func DBExists(db string) bool {
	if _, err := os.Stat(db); os.IsNotExist(err) {
		return false
	}
	return true
}

//FindUnspentTransactions Transactions that have outputs, but no inputs pointing to them are spendable.
//We will call these Unspent Transactions.
func (chain *BlockChain) FindUnspentTransactions(address string) []Transaction {
	var unspentTxns []Transaction

	spentTxns := make(map[string][]int)

	iter := chain.Iterator()

	for {
		block := iter.Next()

		for _, txn := range block.Transactions {
			txnID := hex.EncodeToString(txn.ID)

		Outputs:
			for outIdx, out := range txn.Outputs {
				if spentTxns[txnID] != nil {
					for _, spentOut := range spentTxns[txnID] {
						if spentOut == outIdx {
							continue Outputs
						}
					}
				}
				if out.CanBeUnlocked(address) {
					unspentTxns = append(unspentTxns, *txn)
				}
			}
			if !txn.IsCoinbase() {
				for _, in := range txn.Inputs {
					if in.CanUnlock(address) {
						inTxnID := hex.EncodeToString(in.ID)
						spentTxns[inTxnID] = append(spentTxns[inTxnID], in.Out)
					}
				}
			}
			if len(block.PrevHash) == 0 {
				break
			}
		}
		return unspentTxns
	}
}

func (chain *BlockChain) FindUnspentTransactionOutputs(address string) []TxOutput {
	var UnspentTxnOutputs []TxOutput
	UnspentTransactions := chain.FindUnspentTransactions(address)

	for _, tx := range UnspentTransactions {
		for _, out := range tx.Outputs {
			if out.CanBeUnlocked(address) {
				UnspentTxnOutputs = append(UnspentTxnOutputs, out)
			}
		}
	}
	return UnspentTxnOutputs
}

func (chain *BlockChain) FindSpendableOutputs(address string, amount int) (int, map[string][]int) {
	unspentOuts := make(map[string][]int)
	unspentTxns := chain.FindUnspentTransactions(address)

	accumulated := 0

Work:
	for _, txn := range unspentTxns {
		txnID := hex.EncodeToString(txn.ID)
		for outIdx, out := range txn.Outputs {
			if out.CanBeUnlocked(address) && accumulated < amount {
				accumulated += out.Value
				unspentOuts[txnID] = append(unspentOuts[txnID], outIdx)

				if accumulated >= amount {
					break Work
				}
			}
		}
	}
	return accumulated, unspentOuts
}
