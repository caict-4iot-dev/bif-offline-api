package bif

import (
	"errors"
	"log"
	"strings"
)

// IsAddressValid check address
func IsAddressValid(encAddress string) bool {
	err := encAddressValid(encAddress)
	if err != nil {
		log.Println("IsAddressValid is failed, err: ", err)
		return false
	}

	return true
}

func encAddressValid(encAddress string) error {
	if encAddress == "" {
		return errors.New("invalid address")
	}
	items := strings.Split(encAddress, ":")
	if len(items) != 3 && len(items) != 4 {
		return errors.New("invalid address")
	}
	if len(items) == 3 {
		encAddress = items[2]
	} else {
		encAddress = items[3]
	}

	prefix := string([]byte(encAddress)[:2])
	if !(prefix == "ef") && !(prefix == "zf") {
		return errors.New("invalid address")
	}

	address := []byte(encAddress)[2:]
	base58Address := Base58Decode(address)
	if len(base58Address) != 22 {
		return errors.New("invalid address")
	}

	return nil
}
