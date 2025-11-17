package main

import (
	"crypto/rand"
	"fmt"

	"golang.org/x/crypto/curve25519"
)

func main() {
	iPriv, iPub := mustGenerateKeyPair()
	rPriv, rPub := mustGenerateKeyPair()

	var dhInitiator, dhResponder [32]byte
	curve25519.ScalarMult(&dhInitiator, &iPriv, &rPub)
	curve25519.ScalarMult(&dhResponder, &rPriv, &iPub)

	fmt.Printf("DH from Initiator: %x\n", dhInitiator)
	fmt.Printf("DH from Responder: %x\n", dhResponder)
	fmt.Println("Equal:", dhInitiator == dhResponder)
}

func mustGenerateKeyPair() ([32]byte, [32]byte) {
	var priv, pub [32]byte
	mustFillRandom(priv[:])
	curve25519.ScalarBaseMult(&pub, &priv)
	return priv, pub
}

func mustFillRandom(buf []byte) {
	if _, err := rand.Read(buf); err != nil {
		panic(fmt.Errorf("rand read: %w", err))
	}
}
