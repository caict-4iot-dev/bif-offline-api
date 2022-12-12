package controller

import (
	"bif-offline-api/bif"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetKeyByPrivateKey(c *gin.Context) {
	var r bif.BIFGetKeyByPrivateKeyReq
	err := c.BindJSON(&r)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusBadRequest,
			"message": "请检查请求参数",
			"data":    "",
		})
		return
	}
	if r.PrivateKey == "" {
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "privateKey不能为空",
			"data":    "",
		})
		return
	}

	keyManager, err := bif.GetKeyManagerByPrivateKey(r.PrivateKey)
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
			"encAddress":    keyManager.EncAddress,
			"encPublicKey":  keyManager.EncPublicKey,
			"encPrivateKey": keyManager.EncPrivateKey,
		},
	})
}
