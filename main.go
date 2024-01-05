package main

import (
	"fmt"
	"github.com/shopspring/decimal"
	"inscription/app"
	"inscription/chain/util"
	"inscription/config"
	"time"
)

func main() {
	var err error
	defer func() {
		if err != nil {
			app.LogErrorf("inscribe failed,reason：%s", err)
			return
		}
		app.LogInfo("all inscription is done ")
		time.Sleep(3 * time.Second)
	}()

	app.LogInfof("Welcome to Use %s ", config.ApplicationConfig)
	mintConfig := inputConfig()

	evmApp := app.NewApp(mintConfig.RpcUrl, 3)

	accuracyEth, err := decimal.NewFromString("1000000000000000000")
	if err != nil {
		return
	}

	accuracyGWei, err := decimal.NewFromString("1000000000")
	if err != nil {
		return
	}
	gasPrice, err := decimal.NewFromString(mintConfig.GasPrice)
	if err != nil {
		return
	}

	app.LogInfof("============executing============")
	app.LogInfof("rpcUrl: %s", mintConfig.RpcUrl)
	app.LogInfof("the number of inscriptions: %d", mintConfig.Times)
	app.LogInfof("gas price: %sGwei", gasPrice.DivRound(accuracyGWei, 18))
	app.LogInfof("gas limit: %s\n\n", mintConfig.GasLimit)

	privateKey := mintConfig.PrivateKey
	//begin
	for i := 1; i <= mintConfig.Times; i++ {
		time.Sleep(time.Duration(mintConfig.Delay) * time.Second)

		balanceStr, er := evmApp.TokenBalanceOf(privateKey)
		err = er
		if err != nil {
			app.LogErrorf("%dth inscription query the balance failed，reason: %s", i, err)
			continue
		}
		balance, er := decimal.NewFromString(balanceStr)
		err = er
		if err != nil {
			app.LogErrorf("%dth inscription convert the balance failed，reason: %s", i, err)
			continue
		}
		app.LogInfof("the balance:%s", balance.DivRound(accuracyEth, 4))

		data := mintConfig.Data
		gasPrice := mintConfig.GasPrice
		gasLimit := mintConfig.GasLimit

		hash, er := evmApp.Inscribe(privateKey, data, gasPrice, gasLimit)
		err = er
		if err != nil {
			app.LogErrorf("%dth inscription failed,reason: %s", i, err)
			continue
		}
		app.LogInfof("%dth inscription suc,hash: %s", i, hash)
	}
}

func inputConfig() *config.Inscription {
	var privateKey, data, rpcUrl, gasPrice, gasLimit string

	fmt.Println("please input privateKey:")
	fmt.Scanln(&privateKey)

	var text string
	fmt.Println("please input text:")
	fmt.Scanln(&text)

	data = util.TextToHex(text)
	fmt.Printf("please comfirm data hex: %s\n", data)

	var confirm string
	fmt.Println("please confirm,input y/n")
	fmt.Scanln(&confirm)

	if confirm == "y" {
		fmt.Printf("data is: %s\n", data)
	} else {
		fmt.Println("please input data hex:")
		fmt.Scanln(&data)
		fmt.Printf("data is: %s\n", data)
	}

	fmt.Println("please input rpcUrl:  (free and stable rpc provider in https://zan.top/home/node-service)")
	fmt.Scanln(&rpcUrl)

	fmt.Println("please input gasPrice:")
	fmt.Scanln(&gasPrice)

	fmt.Println("please input gasLimit:")
	fmt.Scanln(&gasLimit)

	fmt.Println("The default number of inscriptions is 10 and the interval between each inscription is 1 second")
	fmt.Println("input y/n: confirm-y, modify-n")
	fmt.Scanln(&confirm)

	var times, delay int
	if confirm == "y" {
		times = 1
		delay = 1
	} else {
		fmt.Println("please input the number of inscriptions:")
		fmt.Scanln(&times)
		fmt.Println("please input the time interval for each inscription:")
		fmt.Scanln(&delay)
	}

	return &config.Inscription{PrivateKey: privateKey, Data: data, RpcUrl: rpcUrl, GasPrice: gasPrice, GasLimit: gasLimit, Times: times, Delay: delay}

}
