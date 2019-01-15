package webService

import (
	"github.com/gin-gonic/gin"
	"github.com/itsjamie/gin-cors"
	"time"
)

/**
 * created on 10/10/18.
 * author: nebula-ai-chengkun
 * Copyright defined in blockchainwebbrowser/LICENSE.txt
 */

func StartwebReq() {

	r := gin.Default()

	r.Use(cors.Middleware(cors.Config{
		Origins:         "*",
		Methods:         "GET, PUT, POST, DELETE",
		RequestHeaders:  "Origin, Authorization, Content-Type",
		ExposedHeaders:  "",
		MaxAge:          50 * time.Second,
		Credentials:     true,
		ValidateHeaders: false,
	}))

	v1 := r.Group("")
	BlockManager(v1.Group("blocks"))
	WalletManager(v1.Group("wallets"))
	TransactionManager(v1.Group("transactions"))
	UtilityManager(v1.Group("utility"))
	r.Run(":8093")
}
