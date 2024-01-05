package feature

import (
	"context"
	"errors"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"inscription/chain/eth/core"
	"inscription/chain/util"
	"time"
)

type Token struct {
	proxy *core.Proxy
}

func NewToken(proxy *core.Proxy) *Token {
	return &Token{
		proxy: proxy,
	}
}

func (t *Token) BalanceOf(address string) (balance string, err error) {
	if t.proxy == nil {
		return "", errors.New("the proxy node is empty")
	}
	if !util.IsValidAddress(address) {
		return "", errors.New("invalid hex address")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(t.proxy.Timeout)*time.Second)
	defer cancel()
	balanceResult, err := t.proxy.RemoteRpcClient.BalanceAt(ctx, common.HexToAddress(address), nil)
	if err != nil {
		return "", err
	}
	return balanceResult.String(), nil
}

func (t *Token) Transfer(privateKey, gasPrice, gasLimit, maxPriorityFeePerGas, value, to, data string) (hash string, err error) {
	if gasPrice == "" || gasLimit == "" || to == "" || value == "" {
		return "", errors.New("param is error")
	}
	tx := core.NewTransaction("", gasPrice, gasLimit, maxPriorityFeePerGas, to, value, data)

	priData, err := util.HexDecodeString(privateKey)
	if err != nil {
		return "", err
	}
	privateKeyECDSA, err := crypto.ToECDSA(priData)
	address := crypto.PubkeyToAddress(privateKeyECDSA.PublicKey).Hex()

	//get no sign tx
	txUnSign, err := t.proxy.BuildTxUnSign(address, tx)
	if err != nil {
		return "", err
	}

	//tx sign
	txSign, err := t.proxy.BuildTxSign(privateKeyECDSA, txUnSign)
	if err != nil {
		return "", err
	}

	//send tx
	return txSign.TxHex, t.proxy.SendTx(txSign.SignedTx)
}

func (t *Token) EstimateGasLimit(fromAddress, receiverAddress, gasPrice, amount string, data []byte) (string, error) {
	msg := core.NewCallMsg()
	msg.SetFrom(fromAddress)
	msg.SetTo(receiverAddress)
	msg.SetGasPrice(gasPrice)
	msg.SetValue(amount)
	if data != nil {
		msg.SetData(data)
	}
	return t.proxy.EstimateGasLimit(msg)
}
