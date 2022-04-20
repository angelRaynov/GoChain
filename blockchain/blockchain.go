package blockchain

import (
	"github.com/dgraph-io/badger"
	"log"
)

const DBPath = "./tmp/blocks"

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

func (chain *BlockChain) AddBlock(data string) {
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

	newBlock := CreateBlock(data, lastHash)

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

func InitBlockChain() *BlockChain {
	//open connection
	var lastHash []byte

	opts := badger.DefaultOptions(DBPath)

	db, err := badger.Open(opts)
	Handle(err)

	//check if there is blockchain created
	err = db.Update(func(txn *badger.Txn) error {
		if _, err := txn.Get([]byte("last_hash")); err == badger.ErrKeyNotFound {
			log.Println("No existing blockchain found")

			genesis := Genesis()

			log.Println("Genesis initialized")

			err = txn.Set(genesis.Hash, genesis.Serialize())
			Handle(err)

			err = txn.Set([]byte("last_hash"), genesis.Hash)
			Handle(err)

			lastHash = genesis.Hash

			return err
		} else {
			item, err := txn.Get([]byte("last_hash"))
			Handle(err)

			err = item.Value(func(val []byte) error {
				lastHash = val
				return nil
			})
			Handle(err)

			return err
		}
	})
	Handle(err)

	return &BlockChain{lastHash, db}
}
