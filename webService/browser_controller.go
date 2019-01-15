package webService

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"navigator-api/models"
	"net/http"
	"strconv"
)

type RequestSlice struct {
	WalletAddress string  `json:"wallet_address"`
	Balance string        `json:"balance"`
	Percentage string     `json:"percentage"`
}

/**
 * created on 12/10/18.
 * author: nebula-ai-chengkun
 * Copyright defined in blockchainwebbrowser/LICENSE.txt
 */

func BlockManager(router *gin.RouterGroup) {
	router.GET("/:value", GetBlockByValue)
	router.GET("", GetLatestBlocks)
}

func TransactionManager(router *gin.RouterGroup) {
	router.GET("/:txHash",GetTransactionByHash)
	router.GET("", GetLatestTransactions)
}

func WalletManager(router *gin.RouterGroup) {
	router.GET(":walletAddr/transactions", GetWalletTxs) //rank by descending numbers
	router.GET("", GetWalletRank)
}

func UtilityManager(router *gin.RouterGroup) {
	router.GET("blocktime", GetAverageBlockTime)
}

func GetWalletRank(c *gin.Context){
	walletAmount := len(models.GetFirstSlice().WalletWithBalances)
	limit, err := strconv.Atoi(c.Query("limit"))
	if err !=nil {
		c.String(http.StatusBadRequest, errors.New("Invalid limit format").Error())
		return
	}
	var hasOffset bool = len(c.Query(("offset")))>0

	offset, err := strconv.Atoi(c.Query("offset")) //Note: offset should be set as an optional parameter here
		if err !=nil &&hasOffset{
		c.String(http.StatusBadRequest, errors.New("Invalid offset format").Error())
		return
	}
	//offset is zero anyways if there is an error
	if offset>=walletAmount && hasOffset{
		c.String(http.StatusBadRequest, errors.New("Offset too large").Error())
		return
	}
    sortedSlice:=[]*models.WalletWithBalance{}
	if offset+limit<walletAmount{ //as long as offset is valid (even when offset is zero), we should not have an index out of range problem no matter what the limit is
	    sortedSlice=append(sortedSlice,models.GetFirstSlice().WalletWithBalances[offset:offset+limit]...)}else{
		sortedSlice=append(sortedSlice,models.GetFirstSlice().WalletWithBalances[offset:]...)
	}

	requestSlice := []RequestSlice{}
	for _, element := range sortedSlice{
		newElement := RequestSlice{WalletAddress:element.WalletAddress, Balance:element.Balance.String(), Percentage:element.Percentage}
		requestSlice = append(requestSlice,newElement)
	}
	c.JSON(http.StatusOK, gin.H{"result":requestSlice,"wallet_count":len(models.GetFirstSlice().WalletWithBalances)})
}


func GetBlockByValue (c *gin.Context){
	value := c.Param("value")
	valueType := c.Query("value-type")
	switch valueType {
	case "miner-address" :GetBlockByMinerAddr(value, c)
	case "tx-hash": GetBlockByTx(value, c)
	case "block-hash": GetBlockByHash(value, c)
	case "block-number":GetBlockByNumber(value, c)
	default:c.String(http.StatusBadRequest, errors.New("Invalid value type option").Error())
		return
	}
}

func GetBlockByMinerAddr(minerAddress string, c *gin.Context) {
	if (!models.IsValidAddress(minerAddress)){
		c.String(http.StatusBadRequest, errors.New("Invalid miner address format").Error())
		return
	}
	block := models.Block{
		Miner: minerAddress,
	}
	_block, _ := block.FindOneBlock(&block)
	if _block.ID == 0 {
		c.String(http.StatusBadRequest, errors.New("Cannot find a block with this miner address").Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"result": _block})
}

func GetTransactionByHash(c *gin.Context){
	transaction := models.Transaction{
		Hash: c.Param("txHash"),
	}
	_transaction, _ := transaction.FindOneTransaction(&transaction)
	if _transaction.ID == 0 {
		c.String(http.StatusBadRequest, errors.New("Cannot find a transaction with this transaction hash").Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"result": _transaction})
}

func GetBlockByTx(txHash string, c *gin.Context) {
	tx := models.Transaction{
		Hash: txHash,
	}
	transaction, _ := tx.FindOneTransaction(&tx)
	if transaction.ID == 0 {
		c.String(http.StatusBadRequest, errors.New("Invalid transaction hash").Error())
		return
	}

	block := models.Block{
		Number:transaction.BlockNumber,
	}

	_block, _ := block.FindOneBlock(&block)

	if _block.ID == 0 {
		c.String(http.StatusBadRequest, errors.New("Cannot get block information").Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": _block})
}

// get block by block hash
func GetBlockByHash(parameter string, c *gin.Context) {

	block := models.Block{
		Hash: parameter,
	}
	_block, _ := block.FindOneBlockByPara(&block)
	if _block.ID == 0 {
		c.String(http.StatusBadRequest, errors.New("Cannot find block with given hash").Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": _block})
}


// get block by block number
func GetBlockByNumber(number string, c *gin.Context) {
	number_int, err := strconv.ParseInt(number, 10, 64)
	if err != nil {
		c.String(http.StatusBadRequest, errors.New("Invalid block number format").Error())
		return
	}

	block := models.Block{
		Number: number_int,
	}
	_block, _ := block.FindOneBlockByPara(&block)
	if _block.ID == 0 {
		c.String(http.StatusBadRequest, errors.New("Cannot find block with given number").Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": _block})
}


// get most recent blocks by limit, if invalid limit parameter input, return most recent 10 blocks by default
func GetLatestBlocks(c *gin.Context) {
	limit := c.Query("limit")
	offset := c.Query("offset")

	_blocks, err := models.FindManyBlock(limit, offset)

	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": _blocks})
}

func GetAverageBlockTime(c *gin.Context){

	resultSlice, err := models.FindAverageBlockTime()
	if (err!=nil){c.String(http.StatusBadRequest, err.Error())}
	c.JSON(http.StatusOK, gin.H{"latest10": resultSlice[0], "latest100":resultSlice[1], "latest1000": resultSlice[2], "latest10000": resultSlice[3],  "latest100000": resultSlice[4]})

}

func GetLatestTransactions(c *gin.Context) {

	limit := c.Query("limit")
	offset := c.Query("offset")
	//Add new functions that detect if limit and offset have invalid types (can only accept numbers!!!)

	addresses := c.Query("addresses")

	if len(addresses)>0{
		GetTransactionsByMultipleAddress(addresses, limit, offset, c)
		return
	}
	_transactions, err := models.FindManyTransactions(limit, offset)

	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": _transactions})
}

func GetTransactionsByMultipleAddress(addresses string, limit string, offset string, c *gin.Context) {

	transactions, txCount, err, errCode := models.FindTransactionsByMultipleWallet(addresses, limit, offset) //this method finds all transactions in a wallet, therefore we need to return the total transaction count of a wallet
	if err != nil {
		switch errCode {
		case 1: c.String(http.StatusBadRequest, errors.New("Invalid limit format").Error())
			return
		case 2: c.String(http.StatusBadRequest, errors.New("Invalid offset format").Error())
			return
		case 3: c.String(http.StatusBadRequest, errors.New("Offset too large").Error())
			return
		case 4: c.String(http.StatusBadRequest, errors.New("At least one of the addresses input is invalid").Error())
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{"result": transactions, "tx_count": txCount})
}

func GetWalletTxs(c *gin.Context) {
	walletAddr := c.Param("walletAddr")
	if !models.IsValidAddress(walletAddr){
		c.String(http.StatusBadRequest, errors.New("Invalid wallet address format").Error())
		return
	}
	limit := c.Query("limit")
	offset := c.Query("offset")  //note:optional parameter
	transactions, txCount, err, errCode := models.FindWalletTransactionByWalletAddr(walletAddr, limit, offset) //this method finds all transactions in a wallet, therefore we need to return the total transaction count of a wallet
	if err != nil {
		switch errCode {
			case 1: c.String(http.StatusBadRequest, errors.New("Invalid limit format").Error())
			    return
			case 2: c.String(http.StatusBadRequest, errors.New("Invalid offset format").Error())
		        return
			case 3: c.String(http.StatusBadRequest, errors.New("Offset too large").Error())
		        return
		}
	}
	c.JSON(http.StatusOK, gin.H{"result": transactions, "tx_count": txCount})
}


