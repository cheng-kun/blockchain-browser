package main

import (
	"navigator-api/browsersync"
	"navigator-api/config"
	"navigator-api/database"
	"navigator-api/eth"
	"navigator-api/logs"
	"navigator-api/webService"
)

/**
 * created on 12/10/18.
 * author: nebula-ai-chengkun
 * Copyright defined in blockchainwebbrowser/LICENSE.txt
 */

func main() {
	logs.InitLogger()
	config.Init()
	db := database.Init()
	defer db.Close()

	eth.ClientInit();
	//browsersync.BlockBrowserSyncTest();

	go browsersync.BlockBrowserSync()

	go browsersync.ScanWallet()

	webService.StartwebReq()

}
