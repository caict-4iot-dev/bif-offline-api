package bif

import "bif-offline-api/proto"

// BIFKeyRequest 创建公私钥接口请求体
type BIFKeyRequest struct {
	Type int `json:"keyType"` // 密钥类型 1 ed25519 2 sm2
}
type BIFGetEncPrivateKeyRequest struct {
	RawPrivateKey string `json:"rawPrivateKey"` // 原生私钥
	Type          int    `json:"keyType"`       // 密钥类型 1 ed25519 2 sm2
}
type BIFGetEncPublicKeyRequest struct {
	RawPublicKey string `json:"rawPublicKey"` // 原生公钥
	Type         int    `json:"keyType"`      // 密钥类型 1 ed25519 2 sm2
}
type BIFGetRawPrivateKeyRequest struct {
	EndPrivateKey string `json:"encPrivateKey"` // 星火私钥
}
type BIFGetRawPublicKeyRequest struct {
	EndPublicKey string `json:"encPublicKey"` // 星火公钥
}

// BIFBlobRequest blob数据解析
type BIFBlobRequest struct {
	Blob string `json:"serialization"`
}

// BIFSignReq 签名请求体
type BIFSignReq struct {
	Message    string `json:"message"`    // 必填，签名数据
	PrivateKey string `json:"privateKey"` // 必填，交易源账户私钥
}

// BIFGetKeyByPrivateKeyReq 解析私钥
type BIFGetKeyByPrivateKeyReq struct {
	PrivateKey string `json:"privateKey"` // 必填，交易源账户私钥
}

// BIFGetKeyByPublicKeyReq 解析公钥
type BIFGetKeyByPublicKeyReq struct {
	PublicKey string `json:"publicKey"` // 必填，交易源账户公钥
}

// BIFVerifySignReq 签名请求体
type BIFVerifySignReq struct {
	Serialization string `json:"serialization"` // 必填，签名数据
	SignData      string `json:"signData"`      // 必填，签名数据
	PublicKey     string `json:"publicKey"`     // 必填，交易源账户公钥
}

// BIFContractInvokeRequest 合约调用请求体
type BIFContractInvokeRequest struct {
	SenderAddress   string `json:"senderAddress"`   // 必填，交易源账号，即交易的发起方
	FeeLimit        int64  `json:"feeLimit"`        // 可选，交易花费的手续费，默认1000000L
	PrivateKey      string `json:"privateKey"`      // 必填，交易源账户私钥
	CeilLedgerSeq   int64  `json:"ceilLedgerSeq"`   // 可选，区块高度限制, 如果大于0，则交易只有在该区块高度之前（包括该高度）才有效
	Metadata        string `json:"metadata"`        // 可选，用户自定义给交易的备注，16进制格式
	ContractAddress string `json:"contractAddress"` // 必填，合约账户地址
	BIFAmount       int64  `json:"bifAmount"`       // 必填，转账金额
	Input           string `json:"input"`           // 选填，待触发的合约的main()入参
	GasPrice        int64  `json:"gasPrice"`
	Nonce           int64  `json:"nonce"`
	Seq             int64  `json:"seq"`
}

// BIFContractInvokeReq 合约调用请求体
type BIFContractInvokeReq struct {
	SenderAddress    string `json:"senderAddress"`    // 必填，交易源账号，即交易的发起方
	FeeLimit         int64  `json:"feeLimit"`         // 可选，交易花费的手续费，默认1000000L
	SenderPrivateKey string `json:"senderPrivateKey"` // 必填，交易源账户私钥
	ContractAddress  string `json:"contractAddress"`  // 必填，合约账户地址
	BIFAmount        int64  `json:"bifAmount"`        // 必填，转账金额
	Input            string `json:"input"`            // 选填，待触发的合约的main()入参
	GasPrice         int64  `json:"gasPrice"`
	Nonce            int64  `json:"nonce"`
}

type BIFContractInvokeOperation struct {
	BIFBaseOperation
	ContractAddress string
	BifAmount       int64
	Input           string
}

type BIFBaseOperation struct {
	OperationType int    `json:"operation_type"`
	SourceAddress string `json:"source_address"`
	Metadata      string `json:"metadata"`
}

// BIFRadioTransactionRequest 广播交易请求体
type BIFRadioTransactionRequest struct {
	SenderAddress    string                     `json:"senderAddress"`
	FeeLimit         int64                      `json:"feeLimit"`
	GasPrice         int64                      `json:"gasPrice"`
	Operation        BIFContractInvokeOperation `json:"operation"`
	CeilLedgerSeq    int64                      `json:"ceilLedgerSeq"`
	Remarks          string                     `json:"remarks"`
	SenderPrivateKey string                     `json:"privateKey"`
	Nonce            int64                      `json:"nonce"`
	Seq              int64                      `json:"seq"`
}

// BIFBaseResponse SDK响应返回结构体
type BIFBaseResponse struct {
	ErrorCode int    `json:"error_code"`
	ErrorDesc string `json:"error_desc"`
}

// BIFAccountGetNonceResponse 获取指定账户的nonce值返回体
type BIFAccountGetNonceResponse struct {
	BIFBaseResponse
	Result BIFAccountGetNonceResult `json:"result"`
}

type BIFAccountGetNonceResult struct {
	Nonce int64 `json:"nonce"`
}

// BIFTransactionSerializeRequest ...
type BIFTransactionSerializeRequest struct {
	SourceAddress string                     `json:"source_address"`
	Nonce         int64                      `json:"nonce"`
	GasPrice      int64                      `json:"gas_price"`
	FeeLimit      int64                      `json:"fee_limit"`
	Operation     BIFContractInvokeOperation `json:"operation"`
	CeilLedgerSeq int64                      `json:"ceil_ledger_seq"`
	Metadata      string                     `json:"metadata"`
	Seq           int64                      `json:"seq"`
}

// BIFContractGetInfoResponse 查询合约代码响应体
type BIFContractGetInfoResponse struct {
	BIFBaseResponse
	Result BIFContractGetInfoResult `json:"result"`
}

type BIFContractGetInfoResult struct {
	Contract BIFContractInfo `json:"contract"` // 合约信息
}

// BIFContractInfo 合约信息
type BIFContractInfo struct {
	Type    int    `json:"type"`    // 合约类型，默认0
	Payload string `json:"payload"` // 合约代码
}

type BIFContractInvokeOperationResponse struct {
	BIFBaseResponse
	Result BIFBaseOperationResult `json:"result"`
}

type BIFBaseOperationResult struct {
	Operation *proto.Operation `json:"operation"`
}

// BIFBlockGetNumberResponse ...
type BIFBlockGetNumberResponse struct {
	BIFBaseResponse
	Result BIFBlockGetNumberResult `json:"result"`
}

type BIFBlockGetNumberResult struct {
	Header BIFBlockNumber `json:"header"` // 区块头
}

// BlockHeader 区块头
type BIFBlockNumber struct {
	BlockNumber int64 `json:"seq"` // 最新的区块高度，对应底层字段seq
}

type BIFTransactionSerializeResponse struct {
	BIFBaseResponse
	Result BIFTransactionSerializeResult `json:"result"`
}

type BIFTransactionSerializeResult struct {
	TransactionBlob string `json:"transaction_blob"`
	Hash            string `json:"hash"`
}

type BIFTransactionResponse struct {
	Serialization string `json:"serialization"` // 序列化的交易数据
	SignData      string `json:"signData"`      // 交易签名数据
	PublicKey     string `json:"publicKey"`     // 公钥
}
