package controller

import (
	"bif-offline-api/bif"
	"github.com/gin-gonic/gin"
	"net/http"
)

func VerifySign(c *gin.Context) {
	var r bif.BIFVerifySignReq
	err := c.BindJSON(&r)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusBadRequest,
			"message": "请检查请求参数",
			"data":    "",
		})
		return
	}
	if r.Serialization == "" {
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusBadRequest,
			"message": "serialization不能为空",
			"data":    "",
		})
		return
	}
	if r.SignData == "" {
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusBadRequest,
			"message": "signData不能为空",
			"data":    "",
		})
		return
	}
	if r.PublicKey == "" {
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusBadRequest,
			"message": "publicKey不能为空",
			"data":    "",
		})
		return
	}

	isOK, err := bif.Verify(r.PublicKey, r.Serialization, r.SignData)
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
			"verifySign": isOK,
		},
	})
}
