package controller

import (
	"bif-offline-api/bif"
	"github.com/gin-gonic/gin"
	"net/http"
)

func ParseBlob(c *gin.Context) {
	var r bif.BIFBlobRequest
	err := c.BindJSON(&r)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusBadRequest,
			"message": "请检查请求参数",
			"data":    "",
		})
		return
	}
	if r.Blob == "" {
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "blob不能为空",
			"data":    "",
		})
		return
	}

	ParseBlobResponse, err := bif.ParseBlobManager(r.Blob)
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
		"data": gin.H{
			"SourceAddress": ParseBlobResponse.SourceAddress,
			"Nonce":         ParseBlobResponse.Nonce,
			"FeeLimit":      ParseBlobResponse.FeeLimit,
			"GasPrice":      ParseBlobResponse.GasPrice,
			"ChainId":       ParseBlobResponse.ChainId,
			"Operations":    ParseBlobResponse.Operations,
		},
	})
}
