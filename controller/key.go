package controller

import (
	"bif-offline-api/bif"
	"encoding/hex"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Key(c *gin.Context) {
	var r bif.BIFKeyRequest
	err := c.BindJSON(&r)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusBadRequest,
			"message": "请检查请求参数",
			"data":    "",
		})
		return
	}
	if r.Type != bif.ED25519 && r.Type != bif.SM2 {
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "keyType值必须是1或者2",
			"data":    "",
		})
		return
	}
	keyManager, err := bif.GenerateKeyManager(r.Type)
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
func GetRawPrivateKey(c *gin.Context) {
	var r bif.BIFGetRawPrivateKeyRequest
	err := c.BindJSON(&r)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusBadRequest,
			"message": "请检查请求参数",
			"data":    "",
		})
		return
	}
	if r.EndPrivateKey == "" {
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "endPrivateKey不能为空",
			"data":    "",
		})
		return
	}
	keyType, rawPrivateKey, err := bif.GetRawPrivateKey([]byte(r.EndPrivateKey))
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
			"keyType":       keyType,
			"rawPrivateKey": hex.EncodeToString(rawPrivateKey),
		},
	})
}
func GetRawPublicKey(c *gin.Context) {
	var r bif.BIFGetRawPublicKeyRequest
	err := c.BindJSON(&r)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusBadRequest,
			"message": "请检查请求参数",
			"data":    "",
		})
		return
	}
	if r.EndPublicKey == "" {
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "endPublicKey不能为空",
			"data":    "",
		})
		return
	}
	rawPublicKey := bif.GetRawPublicKey([]byte(r.EndPublicKey))
	if hex.EncodeToString(rawPublicKey) == "" {
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "endPublicKey 无效",
			"data":    "",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "ok",
		"data": gin.H{
			"rawPublicKey": hex.EncodeToString(rawPublicKey),
		},
	})
}

func GetEncPrivateKeyByRaw(c *gin.Context) {
	var r bif.BIFGetEncPrivateKeyRequest
	err := c.BindJSON(&r)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusBadRequest,
			"message": "请检查请求参数",
			"data":    "",
		})
		return
	}
	if r.Type != bif.ED25519 && r.Type != bif.SM2 {
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "keyType值必须是1或者2",
			"data":    "",
		})
		return
	}
	if r.RawPrivateKey == "" {
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "rawPrivateKey不能为空",
			"data":    "",
		})
		return
	}
	rawPrivateKeyHex, err := hex.DecodeString(r.RawPrivateKey)
	encPrivateKey := bif.GetEncPrivateKey(rawPrivateKeyHex[:32], r.Type)
	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "ok",
		"data": gin.H{
			"keyType":       r.Type,
			"encPrivateKey": encPrivateKey,
		},
	})
}
func GetEncPublicKeyByRaw(c *gin.Context) {
	var r bif.BIFGetEncPublicKeyRequest
	err := c.BindJSON(&r)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusBadRequest,
			"message": "请检查请求参数",
			"data":    "",
		})
		return
	}
	if r.Type != bif.ED25519 && r.Type != bif.SM2 {
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "keyType值必须是1或者2",
			"data":    "",
		})
		return
	}
	if r.RawPublicKey == "" {
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "rawPublicKey不能为空",
			"data":    "",
		})
		return
	}
	rawPublicKeyHex, err := hex.DecodeString(r.RawPublicKey)
	encPublicKey := bif.EncPublicKey(rawPublicKeyHex, r.Type)
	if encPublicKey == "" {
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "rawPublicKey 无效",
			"data":    "",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "ok",
		"data": gin.H{
			"keyType":      r.Type,
			"encPublicKey": encPublicKey,
		},
	})
}

func GetEncryptionTypeByPubK(c *gin.Context) {
	var r bif.BIFGetRawPublicKeyRequest
	err := c.BindJSON(&r)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusBadRequest,
			"message": "请检查请求参数",
			"data":    "",
		})
		return
	}
	if r.EndPublicKey == "" {
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "endPublicKey不能为空",
			"data":    "",
		})
		return
	}
	encType, err := bif.GetEncryptionTypeByPubK(r.EndPublicKey)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "endPublicKey 无效",
			"data":    err,
		})
		return
	}
	if encType == "" {
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "endPublicKey 无效",
			"data":    "",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "ok",
		"data": gin.H{
			"keyType": encType,
		},
	})
}
