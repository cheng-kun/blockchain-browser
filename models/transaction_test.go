package models

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/nebulaai/nbai-node/common"
	"github.com/nebulaai/nbai-node/ethclient"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"math"
	"math/big"
	"navigator-api/database"
	"navigator-api/logs"
	"net/http"
	"strings"

	//"payment-gateway/payments/nbai"
	"testing"
)

type Datas struct {
	Data Data `json:"data"`
}
type Data struct {
	Quotes Quote `json:"quotes"`
	Rank   int32 `json:"rank"`
}

type Quote struct {
	USD USD `json:"USD"`
}
type USD struct {
	Price     float64 `json:"price"`
	MarketCap float64 `json:"volume_24h"`
}

//const EqualityThreshold = 1e-9

func NearlyEqual(a *big.Float, b *big.Float, epsilon *big.Float) bool {

	absA := new(big.Float).Abs(a);
	absB := new(big.Float).Abs(b);
	absDiff := new(big.Float).Abs(new(big.Float).Sub(a, b));
	if a == b { // shortcut, handles infinities
		return true;
	} else if (a == new(big.Float).SetUint64(0) || b == new(big.Float).SetUint64(0) || absDiff.Cmp(big.NewFloat(math.SmallestNonzeroFloat64)) == -1) {
		// a or b is zero or both are extremely close to it
		// relative error is less meaningful here
		return absDiff.Cmp((new(big.Float).Mul(epsilon, big.NewFloat(math.SmallestNonzeroFloat64)))) == -1;
	} else { // use relative error
		return new(big.Float).Quo(absDiff, new(big.Float).Add(absA, absB)).Cmp(epsilon) == -1;
	}
}

func testMissing(t *testing.T) {

}

func TestGetBalance(t *testing.T) {
	//
	//getConfig := config.GetConfig()
	//rpcHost := getConfig.LedgerHost + getConfig.LedgerPort
	//LedgerHost := "http://18.211.182.244:8051"
	LedgerHost := "http://192.168.88.30:8051"
	rpcHost := LedgerHost
	pubKey := "0x63b6d68e84ae0468c3bb890cb3e24dbc3f8768e8"

	client, err := ethclient.Dial(rpcHost)
	assert.Nil(t, err)

	account := common.HexToAddress(pubKey)
	balance, err := client.BalanceAt(context.Background(), account, nil)
	assert.Nil(t, err)
	balanceFloat := new(big.Float).SetInt(balance)

	fmt.Println("Starting the application...")
	response, err := http.Get("https://api.coinmarketcap.com/v2/ticker/2692/?convert=USD")

	body, _ := ioutil.ReadAll(response.Body)
	datas := Datas{}
	json.Unmarshal(body, &datas)
	price := datas.Data.Quotes.USD.Price
	newBalance := new(big.Float).Quo(balanceFloat, new(big.Float).SetUint64(1000000000000000000))
	total := new(big.Float).Mul(big.NewFloat(price), newBalance)

	fmt.Println("The balance of account ", pubKey, "is: ", newBalance, "NBAI", "Total is: ", total, "USD") // 25893180161173005034
	db := database.Init()
	defer db.Close()
	transactions, _, _, _ := FindWalletTransactionByWalletAddr(pubKey, "100", "0")
	totalSent, _ := new(big.Float).SetString("0")
	totalin, _ := new(big.Float).SetString("0")
	lastBlockBalance, _ := new(big.Float).SetString("0")
	//lastTransfered, _ := new(big.Float).SetString("0")

	for _, tx := range transactions {

		if strings.ToLower(tx.TFrom) == pubKey {
			transerferdFloat, _ := new(big.Float).SetString(tx.Value)
			transferred := new(big.Float).Quo(transerferdFloat, new(big.Float).SetUint64(1000000000000000000))

			blockNumber := big.NewInt(tx.BlockNumber)
			balance, err := client.BalanceAt(context.Background(), account, blockNumber)
			balanceFloat, _ = new(big.Float).SetString(balance.String())
			gasFloat, _ := new(big.Float).SetString(tx.Gas)
			gasPriceFloat, _ := new(big.Float).SetString(tx.GasPrice)

			//1 gwei ==1000000000 wei
			gasCost := gasPriceFloat.Mul(gasFloat, gasPriceFloat)
			gasCost.Quo(gasCost, new(big.Float).SetUint64(1000000000000000000))
			//gasCost.Quo(gasCost, new(big.Float).SetUint64(1000000000))

			currentBlockBalance := new(big.Float).SetPrec(64).Quo(balanceFloat, new(big.Float).SetUint64(1000000000000000000))

			if err != nil {
				logs.GetLogger().Error(err)
			}

			lastBlockBalance.Sub(lastBlockBalance, transferred)
			lastBlockBalance.Sub(lastBlockBalance, gasCost)
			//assert.Equal(t,currentBlockBalance,lastBlockBalance)

			if (
				!assert.True(t, NearlyEqual(currentBlockBalance, lastBlockBalance, big.NewFloat(1e-9)))) {
				fmt.Println("THE PROBLEMATIC BLOCK IS BLOCK", blockNumber, "current balance", currentBlockBalance, "last balance", lastBlockBalance)
			}

			if (currentBlockBalance.Cmp(lastBlockBalance) != 0) {
				logs.GetLogger().Error("Block numebr ", blockNumber, "Balance :", currentBlockBalance, " Transefered ", transferred, "Nonce:", tx.Nonce, "=====Out===>>", lastBlockBalance, " Gas cost: ", gasCost)
			} else {
				fmt.Println("Block numebr ", blockNumber, "Balance :", currentBlockBalance, " Transefered", transferred, "Nounce:", tx.Nonce, "=====Out===>>")
			}

			totalSent = transferred.Add(totalSent, transferred)
			//lastTransfered = transferred
			lastBlockBalance = currentBlockBalance

		}
		if strings.ToLower(tx.TTo) == pubKey {
			balanceFloat, _ = new(big.Float).SetString(tx.Value)
			transerferd := new(big.Float).Quo(balanceFloat, new(big.Float).SetUint64(1000000000000000000))

			blockNumber := big.NewInt(tx.BlockNumber)
			balance, err := client.BalanceAt(context.Background(), account, blockNumber)
			balanceFloat = new(big.Float).SetInt(balance)
			currentBlockBalance := new(big.Float).Quo(balanceFloat, new(big.Float).SetUint64(1000000000000000000))

			if err != nil {
				logs.GetLogger().Error()
			}
			fmt.Println("\t", "==== In ==>>", "block numebr ", blockNumber, "Balance :", currentBlockBalance, " Transefered", transerferd, "Nounce:", tx.Nonce, )
			totalin = transerferd.Add(totalin, transerferd)
			lastBlockBalance = currentBlockBalance

		}

	}
	fmt.Println("Total out: ", totalSent)

	fmt.Println("Total in: ", totalin)

}
