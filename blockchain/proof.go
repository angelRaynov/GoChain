package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"log"
	"math"
	"math/big"
)

const Difficulty = 5

//grab the data from the block
//create a nonce
//create a hash of the data from the block + nonce
//check the proceeding hash to see if it meets the requirements (first few bytes to be 0's)

type ProofOfWork struct {
	Block  *Block
	Target *big.Int
}

func (pow *ProofOfWork) InitNonce(nonce int) []byte {
	data := bytes.Join(
		[][]byte{
			pow.Block.PrevHash,
			pow.Block.Data,
			ToHex(int64(nonce)),
			ToHex(int64(Difficulty)),
		},
		[]byte{},
	)

	return data
}

func (pow *ProofOfWork) Run() (int, []byte) {
	var intHash big.Int
	var hash [32]byte

	nonce := 0

	//essentially an infinite loop due to how big MaxInt64 is
	for nonce < math.MaxInt64 {
		data := pow.InitNonce(nonce)
		hash = sha256.Sum256(data)

		log.Printf("\r%x", hash)

		intHash.SetBytes(hash[:])

		if intHash.Cmp(pow.Target) == -1 {
			break
		} else {
			nonce++
		}
	}

	log.Println()

	return nonce, hash[:]
}

func (pow *ProofOfWork) Validate() bool {
	var intHash big.Int

	data := pow.InitNonce(pow.Block.Nonce)

	hash := sha256.Sum256(data)
	intHash.SetBytes(hash[:])

	return intHash.Cmp(pow.Target) == -1
}

func NewProofOfWork(b *Block) *ProofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-Difficulty))

	return &ProofOfWork{b, target}
}

func ToHex(num int64) []byte {
	buff := new(bytes.Buffer)

	err := binary.Write(buff, binary.BigEndian, num)
	if err != nil {
		log.Panic(err)
	}

	return buff.Bytes()
}
