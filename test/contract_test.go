package test

import (
	"bif-offline-api/bif"
	"fmt"
	"testing"
)

func TestContractInvoke(t *testing.T) {
	var r bif.BIFContractInvokeRequest
	senderAddress := "did:bid:efzewQxg38x2Tmb1cpxSC1ZWwMZUxUeV"
	contractAddress := "did:bid:efWVypEKTQoVTunsdBDw8rp4uoG5Lsy5"
	senderPrivateKey := "priSPKhTMRa7SsQLc4wXUDrEZW5wSeKN68Xy5LuCYQmndS75SZ"

	r.SenderAddress = senderAddress
	r.PrivateKey = senderPrivateKey
	r.ContractAddress = contractAddress
	r.BIFAmount = 1
	r.Metadata = "contract invoke"
	// 需要链上获取
	//r.Nonce =
	//r.Seq =

	res, err := bif.ContractInvoke(r)
	if err != nil {
		t.Error(err)
	}

	fmt.Printf("res: %+v", res)
}
