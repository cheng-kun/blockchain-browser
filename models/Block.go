package models

import (
	"navigator-api/database"
	//"navigator-api/models"
	"strconv"
)

/**
 * created on 10/10/18.
 * author: nebula-ai-chengkun
 * Copyright defined in blockchainwebbrowser/LICENSE.txt
 */

type Block struct {
	ID               int64         `json:"id"`
	Number           int64         `json:"number"`
	Author           string        `json:"author"`
	Difficulty       string        `json:"difficulty"`
	ExtraData        string        `json:"extra_data"`
	GasLimit         string        `json:"gas_limit"`
	GasUsed          string        `json:"gas_used"`
	Hash             string        `json:"hash"`
	LogsBloom        string        `json:"logs_bloom"`
	Miner            string        `json:"miner"`
	MixHash          string        `json:"mix_hash"`
	Nonce            string        `json:"nonce"`
	ParentHash       string        `json:"parent_hash"`
	ReceiptsRoot     string        `json:"receipts_root"`
	Sha3uncles       string        `json:"sha3uncles"`
	Size             string        `json:"size"`
	StateRoot        string        `json:"state_root"`
	Timestamp        string        `json:"timestamp"`
	TotalDifficulty  string        `json:"total_difficulty"`
	TransactionsRoot string        `json:"transactions_root"`
	Transactions     []Transaction `gorm:"ForeignKey:block_number"`
	//DifficultyRaw string `json:"difficulty_raw"`
	//GasLimitRaw string `json:"gas_limit_raw"`
	//GasUsedRaw string `json:"gas_used_raw"`
	//NonceRaw string `json:"nonce_raw"`
	//NumberRaw string `json:"number_raw"`
	//SealFileds string `json:"seal_fileds"`
	//SizeRaw string `json:"size_raw"`
	//TimestampRaw string `json:"timestamp_raw"`
	//TotalDifficultyRaw string `json:"total_difficulty_raw"`
}

func (self *Block) FindOneBlock(condition interface{}) (*Block, error) {
	db := database.GetDB()
	tx := db.Begin()
	tx.Where(condition).First(&self)
	err := tx.Commit().Error
	return self, err
}

func (self *Block) FindOneBlockByPara(condition interface{}) (*Block, error) {
	db := database.GetDB()
	err := db.Or(condition).First(&self).Error
	if err != nil {
		return self, err
	}
	var transactions []Transaction
	db.Where(&Transaction{BlockNumber: self.Number}).Find(&transactions)
	self.Transactions = transactions
	return self, err
}

func FindManyBlock(limit, offset string) ([]Block, error) {
	dbHandler := database.GetDB()
	var models []Block

	offset_int, err := strconv.Atoi(offset)
	if err != nil {
		offset_int = 0
	}

	limit_int, err := strconv.Atoi(limit)
	if err != nil {
		limit_int = 10
	}
	tx := dbHandler.Begin()
	dbHandler.Order("number desc").Offset(offset_int).Limit(limit_int).Find(&models)
	for i, _ := range models {
		tx.Model(&models[i])
		var transactions []Transaction
		dbHandler.Where(&Transaction{BlockNumber: (&models[i]).Number}).Find(&transactions)
		(&models[i]).Transactions = transactions
	}
	err = tx.Commit().Error
	return models, err
}

func FindAverageBlockTime() ([]float64, error) {

	block := Block{}
	_, err := block.FindLatestBlock()
	if (err != nil) {
		return nil, err
	}
	latestBlockNo := block.Number
	time, _ := strconv.ParseFloat(block.Timestamp, 64)

	block2 := Block{
		Number: latestBlockNo - 10,
	}
	_, err2 := block2.FindOneBlockByPara(&block2)
	if (err2 != nil) {
		return nil, err2
	}
	time2, _ := strconv.ParseFloat(block2.Timestamp, 64)

	block3 := Block{
		Number: latestBlockNo - 100,
	}
	_, err3 := block3.FindOneBlockByPara(&block3)
	if (err3 != nil) {
		return nil, err3
	}
	time3, _ := strconv.ParseFloat(block3.Timestamp, 64)

	block4 := Block{
		Number: latestBlockNo - 1000,
	}
	_, err4 := block4.FindOneBlockByPara(&block4)
	if (err4 != nil) {
		return nil, err4
	}
	time4, _ := strconv.ParseFloat(block4.Timestamp, 64)

	block5 := Block{
		Number: latestBlockNo - 10000,
	}
	_, err5 := block5.FindOneBlockByPara(&block5)
	if (err5 != nil) {
		return nil, err5
	}
	time5, _ := strconv.ParseFloat(block5.Timestamp, 64)

	block6 := Block{
		Number: latestBlockNo - 100000,
	}
	_, err6 := block6.FindOneBlockByPara(&block6)
	if (err6 != nil) {
		return nil, err6
	}
	time6, _ := strconv.ParseFloat(block6.Timestamp, 64)

	resultSlice := []float64{
		(time - time2) / 10,
		(time - time3) / 100,
		(time - time4) / 1000,
		(time - time5) / 10000,
		(time - time6) / 100000}

	return resultSlice, err
}

func (self *Block) FindLatestBlock() (*Block, error) {
	db := database.GetDB()
	tx := db.Begin()
	tx.Order("number desc").Last(&self)
	err := tx.Commit().Error
	return self, err
}
