package controller

import (
	"bif-offline-api/bif"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ContractInvoke(c *gin.Context) {
	var r bif.BIFContractInvokeReq
	err := c.BindJSON(&r)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusBadRequest,
			"message": err.Error(),
			"data":    "",
		})
		return
	}
	if !bif.IsAddressValid(r.SenderAddress) {
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusBadRequest,
			"message": "非法userBid,请检查",
			"data":    "",
		})
		return
	}
	if !bif.IsAddressValid(r.ContractAddress) {
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusBadRequest,
			"message": "非法的ContractAddress,请检查",
			"data":    "",
		})
		return
	}
	if r.SenderPrivateKey == "" {
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusBadRequest,
			"message": "privateKey不能为空",
			"data":    "",
		})
		return
	}
	if r.FeeLimit < 0 {
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusBadRequest,
			"message": "feeLimit不能小于0",
			"data":    "",
		})
		return
	}
	if r.GasPrice < 0 {
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusBadRequest,
			"message": "gasPrice不能小于0",
			"data":    "",
		})
		return
	}
	if r.Nonce < 0 {
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusBadRequest,
			"message": "nonce不能小于0",
			"data":    "",
		})
		return
	}
	if r.BIFAmount < 0 {
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusBadRequest,
			"message": "bifAmount不能小于0",
			"data":    "",
		})
		return
	}
	if r.Input == "" {
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusBadRequest,
			"message": "input不能为空",
			"data":    "",
		})
		return
	}
	if r.SenderAddress == r.ContractAddress {
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusBadRequest,
			"message": "SenderAddress不能等于ContractAddress",
			"data":    "",
		})
		return
	}
	req := bif.BIFContractInvokeRequest{
		SenderAddress:   r.SenderAddress,
		FeeLimit:        r.FeeLimit,
		PrivateKey:      r.SenderPrivateKey,
		ContractAddress: r.ContractAddress,
		BIFAmount:       r.BIFAmount,
		Input:           r.Input,
		GasPrice:        r.GasPrice,
		Nonce:           r.Nonce,
	}
	transactionResponse, err := bif.ContractInvoke(req)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusInternalServerError,
			"message": err.Error(),
			"data":    "",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "ok",
		"data":    transactionResponse,
	})
}
