package controller

import (
	"bif-offline-api/bif"
	"encoding/hex"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Sign(c *gin.Context) {
	var r bif.BIFSignReq
	err := c.BindJSON(&r)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusBadRequest,
			"message": "请检查请求参数",
			"data":    "",
		})
		return
	}
	if r.Message == "" {
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Message不能为空",
			"data":    "",
		})
		return
	}
	if r.PrivateKey == "" {
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusBadRequest,
			"message": "privateKey不能为空",
			"data":    "",
		})
		return
	}
	keyManager, _ := bif.GetKeyManagerByPrivateKey(r.PrivateKey)
	signMsg, err := bif.Sign([]byte(r.PrivateKey), []byte(r.Message))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusInternalServerError,
			"message": err.Error(),
			"data":    "",
		})
		return
	}
	signData := hex.EncodeToString(signMsg)
	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "ok",
		"data": gin.H{
			"serialization": r.Message,
			"signData":      signData,
			"publicKey":     keyManager.EncPublicKey,
		},
	})
}
