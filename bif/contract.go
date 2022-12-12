package bif

import (
	"bif-offline-api/proto"
	"encoding/hex"
	"errors"

	protobuf "github.com/golang/protobuf/proto"
)

const (
	// 交易默认值
	GAS_PRICE int64 = 100
	FEE_LIMIT int64 = 1000000

	// CONTRACT_INVOKE Invoke contract by sending
	CONTRACT_INVOKE = 15

	// SUCCESS ...
	SUCCESS int = 0
	FAILED  int = 500
)

func ContractInvoke(r BIFContractInvokeRequest) (*BIFTransactionResponse, error) {

	if r.FeeLimit == 0 {
		r.FeeLimit = FEE_LIMIT
	}
	if r.GasPrice == 0 {
		r.GasPrice = GAS_PRICE
	}
	if r.Input == "" {
		return nil, errors.New("input不能为空")
	}

	// 广播交易
	bifAccountActivateOperation := BIFContractInvokeOperation{
		BIFBaseOperation: BIFBaseOperation{
			OperationType: CONTRACT_INVOKE,
		},
		ContractAddress: r.ContractAddress,
		BifAmount:       r.BIFAmount,
		Input:           r.Input,
	}

	bifRadioTransactionRequest := BIFRadioTransactionRequest{
		SenderAddress:    r.SenderAddress,
		FeeLimit:         r.FeeLimit,
		GasPrice:         r.GasPrice,
		Operation:        bifAccountActivateOperation,
		CeilLedgerSeq:    r.CeilLedgerSeq,
		Remarks:          r.Metadata,
		SenderPrivateKey: r.PrivateKey,
		Nonce:            r.Nonce,
		Seq:              r.Seq,
	}

	BIFTransactionResponse, err := RadioTransaction(bifRadioTransactionRequest)
	if err != nil {
		return nil, err
	}

	return BIFTransactionResponse, nil
}

// RadioTransaction 广播交易
func RadioTransaction(r BIFRadioTransactionRequest) (*BIFTransactionResponse, error) {
	// 二、构建操作、序列化交易
	// 初始化请求参数 BIFTransactionSerializeRequest
	serializeRequest := BIFTransactionSerializeRequest{
		SourceAddress: r.SenderAddress,
		Nonce:         r.Nonce,
		GasPrice:      r.GasPrice,
		FeeLimit:      r.FeeLimit,
		Operation:     r.Operation,
		CeilLedgerSeq: r.CeilLedgerSeq,
		Metadata:      r.Remarks,
	}
	// BIFTransactionSerializeResponse
	serializeResponse := BIFSerializable(serializeRequest)
	if serializeResponse.ErrorCode != SUCCESS {
		return nil, errors.New("BIFSerializable is failed")
	}
	transactionBlob := serializeResponse.Result.TransactionBlob
	blob, err := hex.DecodeString(transactionBlob)
	if err != nil {
		return nil, err
	}

	// 三、签名
	signData, err := Sign([]byte(r.SenderPrivateKey), blob)
	if err != nil {
		return nil, err
	}

	publicKey, err := GetEncPublicKeyByEncPrivateKey([]byte(r.SenderPrivateKey))
	if err != nil {
		return nil, err
	}

	submitRequest := BIFTransactionResponse{
		Serialization: transactionBlob,
		SignData:      hex.EncodeToString(signData),
		PublicKey:     publicKey,
	}

	return &submitRequest, nil
}

// BIFSerializable 交易序列化
func BIFSerializable(r BIFTransactionSerializeRequest) BIFTransactionSerializeResponse {

	if !IsAddressValid(r.SourceAddress) {
		return BIFTransactionSerializeResponse{
			BIFBaseResponse: BIFBaseResponse{
				ErrorCode: FAILED,
				ErrorDesc: "invalid sourceAddress",
			},
		}
	}
	if r.Nonce <= 0 {
		return BIFTransactionSerializeResponse{
			BIFBaseResponse: BIFBaseResponse{
				ErrorCode: FAILED,
				ErrorDesc: "invalid nonce",
			},
		}
	}
	if r.CeilLedgerSeq < 0 {
		return BIFTransactionSerializeResponse{
			BIFBaseResponse: BIFBaseResponse{
				ErrorCode: FAILED,
				ErrorDesc: "invalid ceilLedgerSeq",
			},
		}
	}
	if r.GasPrice < 0 {
		return BIFTransactionSerializeResponse{
			BIFBaseResponse: BIFBaseResponse{
				ErrorCode: FAILED,
				ErrorDesc: "invalid gasPrice",
			},
		}
	}
	if r.FeeLimit < 0 {
		return BIFTransactionSerializeResponse{
			BIFBaseResponse: BIFBaseResponse{
				ErrorCode: FAILED,
				ErrorDesc: "invalid feeLimit",
			},
		}
	}

	var operations []*proto.Operation
	operationResData := GetContractInvokeOperation(r.Operation)
	if operationResData.ErrorCode != 0 {
		return BIFTransactionSerializeResponse{
			BIFBaseResponse: BIFBaseResponse{
				ErrorCode: FAILED,
				ErrorDesc: operationResData.ErrorDesc,
			},
		}
	}
	operations = append(operations, operationResData.Result.Operation)

	if r.CeilLedgerSeq < 0 {
		return BIFTransactionSerializeResponse{
			BIFBaseResponse: BIFBaseResponse{
				ErrorCode: FAILED,
				ErrorDesc: "invalid CeilLedgerSeq",
			},
		}
	}
	var seq int64 = 0
	if r.CeilLedgerSeq > 0 {
		seq = r.CeilLedgerSeq + r.Seq
	}
	transaction := proto.Transaction{
		SourceAddress: r.SourceAddress,
		Nonce:         r.Nonce,
		CeilLedgerSeq: seq,
		FeeLimit:      r.FeeLimit,
		GasPrice:      r.GasPrice,
		Metadata:      []byte(r.Metadata),
		Operations:    operations,
	}
	blobByte, err := protobuf.Marshal(&transaction)
	if err != nil {
		return BIFTransactionSerializeResponse{
			BIFBaseResponse: BIFBaseResponse{
				ErrorCode: FAILED,
				ErrorDesc: err.Error(),
			},
		}
	}
	blob := hex.EncodeToString(blobByte)

	return BIFTransactionSerializeResponse{
		BIFBaseResponse: BIFBaseResponse{
			ErrorCode: SUCCESS,
		},
		Result: BIFTransactionSerializeResult{
			TransactionBlob: blob,
			Hash:            string(GenerateHashHex(blobByte, SHA256)),
		},
	}
}

func GetContractInvokeOperation(r BIFContractInvokeOperation) BIFContractInvokeOperationResponse {

	if r.SourceAddress != "" && !IsAddressValid(r.SourceAddress) {
		return BIFContractInvokeOperationResponse{
			BIFBaseResponse: BIFBaseResponse{
				ErrorCode: FAILED,
				ErrorDesc: "invalid SourceAddress",
			},
		}
	}
	if r.ContractAddress != "" && !IsAddressValid(r.ContractAddress) {
		return BIFContractInvokeOperationResponse{
			BIFBaseResponse: BIFBaseResponse{
				ErrorCode: FAILED,
				ErrorDesc: "invalid ContractAddress",
			},
		}
	}
	if r.SourceAddress == r.ContractAddress {
		return BIFContractInvokeOperationResponse{
			BIFBaseResponse: BIFBaseResponse{
				ErrorCode: FAILED,
				ErrorDesc: "sourceAddress equal contractAddress",
			},
		}
	}

	operation := &proto.Operation{
		SourceAddress: r.SourceAddress,
		Metadata:      []byte(r.Metadata),
		Type:          proto.Operation_PAY_COIN,
		PayCoin: &proto.OperationPayCoin{
			Input:       r.Input,
			Amount:      r.BifAmount,
			DestAddress: r.ContractAddress,
		},
	}

	return BIFContractInvokeOperationResponse{
		BIFBaseResponse: BIFBaseResponse{
			ErrorCode: SUCCESS,
		},
		Result: BIFBaseOperationResult{
			Operation: operation,
		},
	}
}
