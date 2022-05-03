package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"log"

	"golang.org/x/crypto/ripemd160"
)

const (
	ChecksumLength = 4
	//hexadecimal representation of 0
	Version = byte(0x00)
)

type Wallet struct {
	PrivateKey ecdsa.PrivateKey
	PublicKey  []byte
}

//Eliptic curvepub key to BTC address conversion
//1. Take our public key (in bytes)
//2. Run a SHA256 hash on it, then run a RipeMD160 hash on that hash. This is called our PublicHash
//3. take our Publish Hash and append our Version (the global variable from earlier) to it. This is called the Versioned Hash
//4. Run SHA256 on our Versioned Hash twice. Then take the first 4 bytes of that output. This is called the Checksum
//5 .Then we will add our Checksum to the end of our original Versioned Hash. We can call this FinalHash
//6. Lastly, we will base58Encode our FinalHash. This is our wallet address!

func PublicKeyHash(publicKey []byte) []byte {
	hashedPubKey := sha256.Sum256(publicKey)

	hasher := ripemd160.New()
	_, err := hasher.Write(hashedPubKey[:])
	if err != nil {
		log.Panic(err)
	}

	publicRipeMd := hasher.Sum(nil)

	return publicRipeMd
}

func Checksum(ripeMdHash []byte) []byte {
	firstHash := sha256.Sum256(ripeMdHash)
	secondHash := sha256.Sum256(firstHash[:])

	return secondHash[:ChecksumLength]
}

func (w *Wallet) Address() []byte {
	//step 1/2
	pubHash := PublicKeyHash(w.PublicKey)

	//step 3
	versionedHash := append([]byte{Version}, pubHash...)

	//step 4
	checksum := Checksum(versionedHash)

	//step 5
	finalHash := append(versionedHash, checksum...)

	//step 6
	address := base58Encode(finalHash)

	return address
}

func NewKeyPair() (ecdsa.PrivateKey, []byte) {
	curve := elliptic.P256()

	private, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		log.Panic(err)
	}

	pub := append(private.PublicKey.X.Bytes(), private.PublicKey.Y.Bytes()...)

	return *private, pub
}

func MakeWallet() *Wallet {
	privateKey, publicKey := NewKeyPair()
	wallet := Wallet{privateKey, publicKey}
	return &wallet
}
