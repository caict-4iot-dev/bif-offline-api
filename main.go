package main

import (
	"bif-offline-api/controller"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	g := r.Group("/bifApi/v1")

	g.POST("/contract", controller.ContractInvoke)                         // 封装合约交易
	g.POST("/createAddress", controller.Key)                               // 生成公私钥和地址
	g.POST("/getAddressByPrivateKey", controller.GetKeyByPrivateKey)       // 通过私钥获取key
	g.POST("/getAddressByPublicKey", controller.GetKeyByPublicKey)         // 通过私钥获取key
	g.POST("/deserialization", controller.ParseBlob)                       // blob数据解析
	g.POST("/sign", controller.Sign)                                       // 签名
	g.POST("/verifySign", controller.VerifySign)                           // 签名验签
	g.POST("/getRawPrivateKey", controller.GetRawPrivateKey)               // 星火私钥转原生私钥
	g.POST("/getRawPublicKey", controller.GetRawPublicKey)                 // 星火公钥转原生公钥
	g.POST("/getEncPrivateKeyByRaw", controller.GetEncPrivateKeyByRaw)     // 原生私钥转星火私钥
	g.POST("/getEncPublicKeyByRaw", controller.GetEncPublicKeyByRaw)       // 原生公钥转星火公钥
	g.POST("/getEncryptionTypeByPubK", controller.GetEncryptionTypeByPubK) // 根据星火公钥获取加密类型

	r.Run("0.0.0.0:8888")
}
