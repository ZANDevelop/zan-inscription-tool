package app

import (
	"inscription/chain/eth/core"
	"inscription/feature"
)

type App struct {
	token *feature.Token
}

func NewApp(rpcUrl string, timeout int64) *App {
	proxy, err := core.GetProxy(rpcUrl, timeout)
	if err != nil {
		LogErrorf("init app err: %s", err)
		return nil
	}

	return &App{
		token: feature.NewToken(proxy),
	}
}

func (a *App) TokenBalanceOf(privateKey string) (balance string, err error) {
	account, err := core.NewAccount().AccountWithPrivateKey(privateKey)
	if err != nil {
		return "", err
	}
	return a.token.BalanceOf(account.Address)
}

// TokenBalanceOfAccount
//
//	@Description: get balance
//	@receiver a
//	@return balance
//	@return err
func (a *App) TokenBalanceOfAccount(account *core.Account) (balance string, err error) {
	return a.token.BalanceOf(account.Address)
}

func (a *App) Transfer(account *core.Account, toAddress string, value string) (hash string, err error) {
	gasLimit := "21000"
	gasPrice := "30000000000" // in wei (30 gwei)
	return a.token.Transfer(account.PrivateKey, gasPrice, gasLimit, "", value, toAddress, "")
}

func (a *App) Inscribe(privateKey string, data string, gasPrice string, gasLimit string) (hash string, err error) {
	account, err := core.NewAccount().AccountWithPrivateKey(privateKey)
	return a.token.Transfer(privateKey, gasPrice, gasLimit, "", "0", account.Address, data)
}
