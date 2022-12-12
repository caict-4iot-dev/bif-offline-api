package test

import (
	"bif-offline-api/bif"
	"fmt"
	"testing"
)

func TestGenerateKeyManager(t *testing.T) {
	keyManager, err := bif.GenerateKeyManager(bif.ED25519)
	if err != nil {
		t.Error(err)
	}

	fmt.Println("EncPrivateKey: ", keyManager.EncPrivateKey)
	fmt.Println("EncPublicKey: ", keyManager.EncPublicKey)
	fmt.Println("EncAddress: ", keyManager.EncAddress)
}
