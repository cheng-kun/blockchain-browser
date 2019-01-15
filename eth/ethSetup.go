package eth

import (
	"github.com/nebulaai/nbai-node/ethclient"
	"navigator-api/config"
	"navigator-api/logs"
	"time"
)

/**
 * created on 10/10/18.
 * author: nebula-ai-chengkun
 * Copyright defined in blockchainwebbrowser/LICENSE.txt
 */

type ConnSetup struct {
	ConnWeb *ethclient.Client
}

//setup eth connection
var WebConn = new(ConnSetup)

func ClientInit() {

	for {
		client, err := ethclient.Dial(config.GetConfig().NBAILedgerHost)
		if err != nil {
			logs.GetLogger().Error("Try to reconnect ...")
			time.Sleep(time.Second * 5)
		} else {
			WebConn.ConnWeb = client
			break
		}
	}
}
