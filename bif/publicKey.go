package bif

import (
	"encoding/hex"
	"errors"
)

// KeyManager ...
type PublicKeyManager struct {
	EncAddress   string `json:"encAddress"`   // 地址
	EncPublicKey string `json:"encPublicKey"` // 公钥
}

// getKeyManagerByPrivateKey
func GetKeyManagerByPublicKey(encPublicKey string) (*PublicKeyManager, error) {
	pblicKeyHex, err := hex.DecodeString(encPublicKey)
	if err != nil {
		return nil, errors.New("public key (" + encPublicKey + ") is invalid, please check")
	}
	var keyType = ED25519
	// 判断算法类型
	if pblicKeyHex[1] == ED25519_VALUE {
		keyType = ED25519
	} else if pblicKeyHex[1] == SM2_VALUE {
		keyType = SM2
	} else {
		return nil, errors.New("public key (" + encPublicKey + ") is invalid, please check")
	}
	if pblicKeyHex[2] != BASE_58_VALUE {
		return nil, errors.New("public key (" + encPublicKey + ") is invalid, please check")
	}

	publicKey := GetRawPublicKey([]byte(encPublicKey))
	encAddress := GetEncAddress(publicKey, "", keyType)

	var keyManager PublicKeyManager
	keyManager.EncPublicKey = encPublicKey
	keyManager.EncAddress = encAddress

	return &keyManager, nil
}
