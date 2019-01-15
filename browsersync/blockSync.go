package browsersync

import (
	"context"
	"github.com/nebulaai/nbai-node/common"
	"github.com/nebulaai/nbai-node/core/types"
	"math/big"
	"navigator-api/database"
	"navigator-api/eth"
	"navigator-api/logs"
	"navigator-api/models"
	"sort"
	"strconv"
	"time"
)

func ScanWallet() {
	var firstId int64
	var secondId int64
	var startId int64 //this is the query start id. The necessity of this is to distinguish first scan from subsequent scans. This MAY be rendered more concise, just saying.
	var endId int64   //this is the query end id. The necessity of this is to distinguish first scan from subsequent scans.
	var firstScan = true
	for {
		firstId = secondId + 1 //firstId is the end id in last iteration+1
		if firstScan {
			startId = 0;
		} else {
			startId = firstId;
		}
		transaction := models.Transaction{}
		transaction.FindLastTxId()
		secondId = transaction.ID //this selects the latest transaction id at this moment. It's also where we are supposed to begin the query next time
		if firstScan {
			endId = 9223372036854775807;
			firstScan = false;
		} else {
			endId = secondId;
		}
		transactions := []models.Transaction{}
		models.FindTransactionsById(&transactions, startId, endId)
		for _, element := range transactions {
			wallet := models.Wallet{WalletAddress: element.TFrom}
			_wallet, _ := wallet.FindOneWallet(&wallet)
			if _wallet.ID == 0 {
				record := models.Wallet{WalletAddress: element.TFrom, FirstTxID: strconv.FormatInt(element.ID, 10)}
				err := database.GetDB().Save(&record).Error
				if err != nil {
					logs.GetLogger().Error(err)
				}
			}
			wallet2 := models.Wallet{WalletAddress: element.TTo}
			_wallet2, _ := wallet2.FindOneWallet(&wallet2)
			if _wallet2.ID == 0 {
				record := models.Wallet{WalletAddress: element.TTo, FirstTxID: strconv.FormatInt(element.ID, 10)}
				database.GetDB().Save(&record)
				err := database.GetDB().Save(&record).Error
				if err != nil {
					logs.GetLogger().Error(err)
				}
			}
		}
		UpdateAndSortWallet()
		time.Sleep(time.Hour * 24)
	}
}

func UpdateAndSortWallet() {
	var wallets []models.Wallet
	models.FindWalletAddresseses(&wallets)
	totalBalance := new(big.Int).SetUint64(0)
	for _, wallet := range wallets {
		balance, err := getBalance(wallet.WalletAddress)
		if err != nil {
			logs.GetLogger().Error(err)
			return
		}
		totalBalance = new(big.Int).Add(totalBalance, balance)
		models.GetSecondSlice().StoreArray(models.WalletWithBalance{WalletAddress: wallet.WalletAddress, Balance: balance})
	}
	totalBalanceFloat := new(big.Float).Quo(new(big.Float).SetInt(totalBalance), new(big.Float).SetUint64(100))
	sort.Sort(models.ByBalance(models.GetSecondSlice().WalletWithBalances))
	for _, element := range models.GetSecondSlice().WalletWithBalances {
		element.Percentage = new(big.Float).Quo(new(big.Float).SetInt(element.Balance), totalBalanceFloat).String()
	}

	models.GetFirstSlice().WalletWithBalances = make([]*models.WalletWithBalance, len(models.GetSecondSlice().WalletWithBalances))
	copy(models.GetFirstSlice().WalletWithBalances, models.GetSecondSlice().WalletWithBalances)
}

func getBalance(address string) (*big.Int, error) {
	account := common.HexToAddress(address)
	balance, err := eth.WebConn.ConnWeb.BalanceAt(context.Background(), account, nil)
	if err != nil {
		return balance, err
	}
	return balance, err
}

func BlockBrowserSync() {

	for {

		block := models.Block{}

		_blockDB, err := block.FindLatestBlock()

		blockNoDB := _blockDB.Number
		if err != nil {
			logs.GetLogger().Error(err)
		}

		blockNoCurrent, err := eth.WebConn.GetBlockNumber()
		if err != nil {
			eth.ClientInit()
			logs.GetLogger().Error(err)
			continue
		}

		UpdateFromLastToCurrent(blockNoDB, blockNoCurrent.Int64())

		time.Sleep(time.Second * 5)
	}
}

func UpdateFromLastToCurrent(lastBlockNo, currentBlockNo int64) {

	for lastBlockNo < currentBlockNo {
		nextBlock, err := eth.WebConn.GetBlockByNumber(big.NewInt(lastBlockNo + 1))
		if err != nil {
			logs.GetLogger().Error(err)
			break
		}

		err = SaveBlock(nextBlock)
		if err != nil {
			logs.GetLogger().Error(err)
		}

		if nextBlock.Transactions().Len() > 0 {
			err = SaveTransaction(nextBlock)
			if err != nil {
				logs.GetLogger().Error(err)
			}
		}

		lastBlockNo++

	}

}

func SaveBlock(block *types.Block) error {
	browserBlock := &models.Block{}
	browserBlock.Number = block.Number().Int64()
	browserBlock.Hash = block.Hash().String()
	browserBlock.ParentHash = block.ParentHash().String()
	browserBlock.Author = ""
	browserBlock.Miner = block.Coinbase().String()
	browserBlock.Size = block.Size().String()
	browserBlock.GasLimit = strconv.FormatUint(block.GasLimit(), 10)
	browserBlock.GasUsed = strconv.FormatUint(block.GasUsed(), 10)
	browserBlock.Nonce = strconv.FormatUint(block.Nonce(), 10)
	browserBlock.Timestamp = block.Time().String()
	browserBlock.Difficulty = block.Difficulty().String()

	_findblock, _ := browserBlock.FindOneBlock(browserBlock)
	if _findblock.ID == 0 {
		database.GetDB().Save(browserBlock)
	}

	return nil
}

func SaveTransaction(block *types.Block) error {

	chainID, err1 := eth.WebConn.ConnWeb.NetworkID(context.Background())
	if err1 != nil {
		logs.GetLogger().Error(err1)
	}

	for i, tx := range block.Transactions() {
		browserTx := &models.Transaction{}

		txMsg, err2 := tx.AsMessage(types.NewEIP155Signer(chainID));
		if err2 != nil {
			logs.GetLogger().Error(err2)
		}

		browserTx.Hash = tx.Hash().String()
		browserTx.Nonce = strconv.FormatUint(txMsg.Nonce(), 10)
		browserTx.BlockNumber = block.Number().Int64()
		browserTx.TransactionIndex = strconv.Itoa(int(i))
		browserTx.TFrom = txMsg.From().String()

		browserTx.Timestamp = block.Time().String()

		if txMsg.To() != nil {
			browserTx.TTo = txMsg.To().String()
		} else {
			browserTx.TTo = ""
		}

		browserTx.Input = tx.Data()
		browserTx.Value = txMsg.Value().String()
		browserTx.GasPrice = txMsg.GasPrice().String()
		browserTx.Gas = strconv.FormatUint(txMsg.Gas(), 10)
		browserTx.BlockHash = block.Hash().String()
		_browserTx, _ := browserTx.FindOneTransaction(browserTx)
		if _browserTx.ID == 0 {
			database.GetDB().Save(browserTx)
		}

	}

	return nil
}
