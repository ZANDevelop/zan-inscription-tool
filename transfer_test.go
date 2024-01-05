package main

import (
	"github.com/shopspring/decimal"
	"inscription/app"
	"inscription/chain/eth/core"
	"inscription/chain/util"
	"testing"
)

func TestMint(t *testing.T) {

	accuracyEth, err := decimal.NewFromString("1000000000000000000")
	rpcUrl := "https://api.zan.top/node/v1/eth/goerli/{apiKey}}"
	evmApp := app.NewApp(rpcUrl, 3)
	account, err := core.NewAccount().AccountWithPrivateKey("{privateKey}")
	if err != nil {
		return
	}

	balanceStr, er := evmApp.TokenBalanceOfAccount(account)
	err = er
	if err != nil {
		return
	}
	balance, er := decimal.NewFromString(balanceStr)
	err = er
	if err != nil {
		return
	}
	t.Logf("the balance: %s", balance.DivRound(accuracyEth, 4))

	text := "{inscription}"
	data := util.TextToHex(text)
	gasLimit := "210000"
	gasPrice := "30000000000" // in wei (30 gwei)
	hash, err := evmApp.Inscribe(account.PrivateKey, data, gasPrice, gasLimit)
	if err != nil {
		t.Log(err)
		return
	}
	t.Log(hash)
}

func TestTransfer(t *testing.T) {

	accuracyEth, err := decimal.NewFromString("1000000000000000000")
	rpcUrl := "https://api.zan.top/node/v1/eth/goerli/{apiKey}}"
	evmApp := app.NewApp(rpcUrl, 3)
	account, err := core.NewAccount().AccountWithPrivateKey("{privateKey}")
	if err != nil {
		return
	}

	balanceStr, er := evmApp.TokenBalanceOfAccount(account)
	err = er
	if err != nil {
		return
	}
	balance, er := decimal.NewFromString(balanceStr)
	err = er
	if err != nil {
		return
	}
	t.Logf("the balance: %s", balance.DivRound(accuracyEth, 4))

	hash, err := evmApp.Transfer(account, "0x0b39fb6bce3381115db85210666585ebb9d32e25", "10000")
	if err != nil {
		t.Log(err)
		return
	}
	t.Log(hash)
}
