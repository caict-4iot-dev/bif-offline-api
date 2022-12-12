package bif

import "errors"

// getKeyManagerByPrivateKey
func GetKeyManagerByPrivateKey(encPrivateKey string) (*KeyManager, error) {
	priKeyTmp := Base58Decode([]byte(encPrivateKey))
	if len(priKeyTmp) <= 5 {
		return nil, errors.New("private key (" + encPrivateKey + ") is invalid")
	}

	keyType, _, err := GetRawPrivateKey([]byte(encPrivateKey))
	if err != nil {
		return nil, err
	}
	encPublicKey, err := GetEncPublicKey([]byte(encPrivateKey))
	if err != nil {
		return nil, err
	}

	publicKey := GetRawPublicKey([]byte(encPublicKey))
	encAddress := GetEncAddress(publicKey, "", keyType)

	var keyManager KeyManager
	keyManager.EncPrivateKey = encPrivateKey
	keyManager.EncPublicKey = encPublicKey
	keyManager.EncAddress = encAddress

	return &keyManager, nil
}
