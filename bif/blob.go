package bif

import (
	"bif-offline-api/proto"
	"encoding/hex"
	"errors"
	protobuf "github.com/golang/protobuf/proto"
)

// ParseBlobManager
type ParseBlobResponse struct {
	SourceAddress string             `json:"sourceAddress"`
	FeeLimit      int64              `json:"feeLimit"`
	GasPrice      int64              `json:"gasPrice"`
	Nonce         int64              `json:"nonce"`
	Operations    []*proto.Operation `json:"operations"`
	ChainId       int64              `json:"chainId"`
}

// ParseBlobManager
func ParseBlobManager(blob string) (*ParseBlobResponse, error) {
	var blobRequest ParseBlobResponse
	transaction := proto.Transaction{}

	blobHex, err := hex.DecodeString(blob)
	if err != nil {
		return nil, errors.New("Invalid serialization")
	}
	obj := &transaction
	err = protobuf.Unmarshal(blobHex, obj)
	if err != nil {
		return nil, errors.New("Invalid serialization")
	}

	blobRequest.SourceAddress = obj.SourceAddress
	blobRequest.FeeLimit = obj.FeeLimit
	blobRequest.GasPrice = obj.GasPrice
	blobRequest.Nonce = obj.Nonce
	blobRequest.Operations = obj.Operations
	blobRequest.ChainId = obj.ChainId

	return &blobRequest, nil
}
