package bif

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"math/big"

	"github.com/ZZMarquis/gm/sm2"
)

const (
	// ED25519算法
	ED25519 = iota + 1
	// SM2算法
	SM2

	// 加密类型关键字
	ED25519_VALUE = 'e'
	SM2_VALUE     = 'z'
	// 字符编码类型关键字
	BASE_58_VALUE = 'f'
)

// KeyManager ...
type KeyManager struct {
	EncAddress    string `json:"encAddress"`    // 地址
	EncPublicKey  string `json:"encPublicKey"`  // 公钥
	EncPrivateKey string `json:"encPrivateKey"` // 私钥
}

// GenerateKeyManager 获取私钥、公钥、地址
func GenerateKeyManager(keyType int) (*KeyManager, error) {
	var rawPrivateKey []byte
	switch keyType {
	case ED25519:
		_, privateKey, err := ed25519.GenerateKey(rand.Reader)
		if err != nil {
			return nil, err
		}
		rawPrivateKey = privateKey[:32]
	case SM2:
		privateKey, _, err := sm2.GenerateKey(rand.Reader)
		if err != nil {
			return nil, err
		}
		rawPrivateKey = privateKey.GetRawBytes()
	default:
		return nil, errors.New("type does not exist")
	}
	encPrivateKey := GetEncPrivateKey(rawPrivateKey, keyType)
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

// GetEncPrivateKey 原生私钥转星火私钥
func GetEncPrivateKey(privateKey []byte, keyType int) string {
	buff := make([]byte, len(privateKey)+5)
	buff[0] = 0x18
	buff[1] = 0x9E
	buff[2] = 0x99

	switch keyType {
	case ED25519:
		buff[3] = ED25519_VALUE
	case SM2:
		buff[3] = SM2_VALUE
	default:
		return ""
	}

	buff[4] = BASE_58_VALUE
	buff = append(buff[:5], privateKey...)

	return Base58Encode(buff)
}

// GetEncPublicKey 星火私钥获取星火公钥
func GetEncPublicKey(encPrivateKey []byte) (string, error) {
	keyType, rawPrivateKey, err := GetRawPrivateKey(encPrivateKey)
	if err != nil {
		return "", err
	}
	var rawPublicKey []byte

	switch keyType {
	case ED25519:
		rawPublicKey = rawPrivateKey[32:]
	case SM2:
		priKey, err := sm2.RawBytesToPrivateKey(rawPrivateKey)
		if err != nil {
			return "", err
		}
		pubKey := sm2.CalculatePubKey(priKey)
		rawPublicKey, err = hex.DecodeString("04" + hex.EncodeToString(pubKey.GetRawBytes()))
		if err != nil {
			return "", err
		}
	default:
		return "", errors.New("type does not exist")
	}

	return EncPublicKey(rawPublicKey, keyType), nil
}

// EncPublicKey 原生公钥转星火公钥
func EncPublicKey(rawPublicKey []byte, keyType int) string {

	buff := make([]byte, len(rawPublicKey)+3)
	buff[0] = 0xB0
	switch keyType {
	case ED25519:
		buff[1] = ED25519_VALUE
	case SM2:
		buff[1] = SM2_VALUE
	default:
		return ""
	}

	buff[2] = BASE_58_VALUE
	buff = append(buff[:3], rawPublicKey...)

	return hex.EncodeToString(buff)
}

// GetRawPublicKey 星火公钥获取原生公钥
func GetRawPublicKey(encPublicKey []byte) []byte {

	rawPublicKey, err := hex.DecodeString(string(encPublicKey))
	if err != nil {
		return nil
	}

	return rawPublicKey[3:]
}

// GetEncAddress 获取地址
func GetEncAddress(publicKey []byte, chainCode string, keyType int) string {
	hashPkey := GenerateHashHex(publicKey, keyType)
	encAddress := Base58Encode(hashPkey[10:])
	if chainCode == "" {
		switch keyType {
		case ED25519:
			return "did:bid:" + "ef" + encAddress

		case SM2:
			return "did:bid:" + "zf" + encAddress

		default:
			return ""
		}
	}

	return "did:bid:" + chainCode + ":" + "ef" + encAddress
}

// GetEncPublicKeyByEncPrivateKey 星火私钥获取星火公钥
func GetEncPublicKeyByEncPrivateKey(encPrivateKey []byte) (string, error) {
	keyType, rawPrivateKey, err := GetRawPrivateKey(encPrivateKey)
	if err != nil {
		return "", err
	}
	var rawPublicKey []byte

	switch keyType {
	case ED25519:
		rawPublicKey = rawPrivateKey[32:]
	case SM2:
		priKey, err := sm2.RawBytesToPrivateKey(rawPrivateKey)
		if err != nil {
			return "", err
		}
		pubKey := sm2.CalculatePubKey(priKey)
		rawPublicKey, err = hex.DecodeString("04" + hex.EncodeToString(pubKey.GetRawBytes()))
		if err != nil {
			return "", err
		}
	default:
		return "", errors.New("type does not exist")
	}

	return EncPublicKey(rawPublicKey, keyType), nil
}

// GetRawPrivateKey 星火私钥转原生私钥
func GetRawPrivateKey(encPrivateKey []byte) (int, []byte, error) {
	priKeyTmp := Base58Decode(encPrivateKey)
	if len(priKeyTmp) <= 5 {
		return 0, nil, errors.New("private key (" + string(encPrivateKey) + ") is invalid")
	}

	if priKeyTmp[3] != ED25519_VALUE && priKeyTmp[3] != SM2_VALUE {
		return 0, nil, errors.New("private key (" + string(encPrivateKey) + ") is invalid")
	}
	var keyType int
	switch priKeyTmp[3] {
	case ED25519_VALUE:
		{
			keyType = ED25519
		}
	case SM2_VALUE:
		{
			keyType = SM2
		}
	default:
		return 0, nil, errors.New("Private key (" + string(encPrivateKey) + ") is invalid")
	}
	if priKeyTmp[4] != BASE_58_VALUE {
		return 0, nil, errors.New("private key (" + string(encPrivateKey) + ") is invalid")
	}

	var rawPrivateKey []byte
	switch keyType {
	case ED25519:
		rawPrivateKey = ed25519.NewKeyFromSeed(priKeyTmp[5:])
	case SM2:
		rawPrivateKey = priKeyTmp[5:]
	}

	return keyType, rawPrivateKey, nil
}

// Sign 签名
func Sign(encPrivate []byte, message []byte) ([]byte, error) {
	keyType, rawPrivateKey, err := GetRawPrivateKey(encPrivate)
	if err != nil {
		return nil, err
	}

	var sign []byte
	switch keyType {
	case ED25519:
		sign25519 := ed25519.Sign(rawPrivateKey, message)
		sign = sign25519
	case SM2:
		priKey, err := sm2.RawBytesToPrivateKey(rawPrivateKey)
		if err != nil {
			return nil, err
		}
		r, s, err := sm2.SignToRS(priKey, []byte("1234567812345678"), message)
		if err != nil {
			return nil, err
		}
		rBytes := r.Bytes()
		sBytes := s.Bytes()
		sig := make([]byte, 64)
		if len(rBytes) == 33 {
			copy(sig[:32], rBytes[1:])
		} else {
			copy(sig[:32], rBytes[:])
		}
		if len(sBytes) == 33 {
			copy(sig[32:], sBytes[1:])
		} else {
			copy(sig[32:], sBytes[:])
		}

		sign = sig
	}

	return sign, nil
}
func Verify(encPublicKey string, msg string, signMsg string) (bool, error) {

	var isOK bool
	pblicKeyHex, err := hex.DecodeString(encPublicKey)
	if err != nil {
		return false, errors.New("public key (" + encPublicKey + ") is invalid, please check")
	}
	signMsgHex, err := hex.DecodeString(signMsg)
	if err != nil {
		return false, errors.New(signMsg + " is invalid, please check")
	}
	msgHex, err := hex.DecodeString(msg)
	if err != nil {
		return false, errors.New(msg + " is invalid, please check")
	}
	var keyType = ED25519
	// 判断算法类型
	if pblicKeyHex[1] == ED25519_VALUE {
		keyType = ED25519
	} else if pblicKeyHex[1] == SM2_VALUE {
		keyType = SM2
	} else {
		return false, errors.New("public key (" + encPublicKey + ") is invalid, please check")
	}
	if pblicKeyHex[2] != BASE_58_VALUE {
		return false, errors.New("public key (" + encPublicKey + ") is invalid, please check")
	}
	rawPublicKey := GetRawPublicKey([]byte(encPublicKey))

	switch keyType {
	case ED25519:
		isOK = ed25519.Verify(rawPublicKey, msgHex, signMsgHex)
	case SM2:
		publicKeyHex := hex.EncodeToString(rawPublicKey)
		publicKey, err := hex.DecodeString(publicKeyHex[2:])
		if err != nil {
			return false, err
		}
		pubKey, _ := sm2.RawBytesToPublicKey(publicKey)
		r := new(big.Int).SetBytes(signMsgHex[:32])
		s := new(big.Int).SetBytes(signMsgHex[32:])
		isOK = sm2.VerifyByRS(pubKey, []byte("1234567812345678"), msgHex, r, s)
	}

	return isOK, nil
}

func GetEncryptionTypeByPubK(encPublicKey string) (string, error) {
	pblicKeyHex, err := hex.DecodeString(encPublicKey)
	if err != nil {
		return "", errors.New("endPublic key (" + encPublicKey + ") is invalid, please check")
	}
	var keyType = "ED25519"
	// 判断算法类型
	if pblicKeyHex[1] == ED25519_VALUE {
		keyType = "ED25519"
	} else if pblicKeyHex[1] == SM2_VALUE {
		keyType = "SM2"
	} else {
		return "", errors.New("endPublic key (" + encPublicKey + ") is invalid, please check")
	}

	return keyType, nil
}
