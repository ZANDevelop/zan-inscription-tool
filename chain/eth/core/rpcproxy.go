package core

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"inscription/chain/util"
	"inscription/config"
	"math/big"
	"strconv"
	"sync"
	"time"
)

var chainConnections = make(map[string]*Proxy)
var lock sync.RWMutex

type Proxy struct {
	RemoteRpcClient *ethclient.Client
	Timeout         int64
	rpcClient       *rpc.Client
	chainId         *big.Int
	rpcUrl          string
}

// GetProxy
//
//	@Description: get connect from cache
//	@param rpcUrl
//	@param timeout
//	@return *EthChain
//	@return error
func GetProxy(rpcUrl string, timeout int64) (*Proxy, error) {
	if rpcUrl == "" {
		return nil, errors.New("rpc url can't be empty")
	}

	chain, ok := chainConnections[rpcUrl]
	if ok {
		return chain, nil
	}

	// 通过加锁范围
	lock.Lock()
	defer lock.Unlock()

	// 再判断一次
	chain, ok = chainConnections[rpcUrl]
	if ok {
		return chain, nil
	}

	// 创建并存储
	chain, err := newProxy(rpcUrl, timeout)
	if err != nil {
		return nil, err
	}

	chainConnections[rpcUrl] = chain
	return chain, nil
}

// newProxy
//
//	@Description:
//	@param timeout the net connect time, second,default is 60
//	@return *Proxy
func newProxy(rpcUrl string, timeout int64) (chain *Proxy, err error) {
	if timeout <= 0 {
		timeout = 60
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()
	rpcClient, err := rpc.DialContext(ctx, rpcUrl)
	if err != nil {
		return
	}

	remoteRpcClient := ethclient.NewClient(rpcClient)
	chainId, err := remoteRpcClient.ChainID(ctx)
	if err != nil {
		return
	}

	chain = &Proxy{
		chainId:         chainId,
		rpcClient:       rpcClient,
		RemoteRpcClient: remoteRpcClient,
		rpcUrl:          rpcUrl,
		Timeout:         timeout,
	}
	return
}

func (c *Proxy) Close() {
	if c.RemoteRpcClient != nil {
		c.RemoteRpcClient.Close()
	}
	if c.rpcClient != nil {
		c.rpcClient.Close()
	}
}

func (c *Proxy) EstimateGasLimit(msg *CallMsg) (gas string, err error) {

	if len(msg.Msg.Data) > 0 {
		// any contract transaction
		gas = config.DefaultContractGasLimit
	} else {
		// nomal transfer
		gas = config.DefaultEthGasLimit
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(c.Timeout)*time.Second)
	defer cancel()
	gasLimit, err := c.RemoteRpcClient.EstimateGas(ctx, msg.Msg)
	if err != nil {
		return
	}
	gasString := ""
	if len(msg.Msg.Data) > 0 {
		gasFloat := big.NewFloat(0).SetUint64(gasLimit)
		gasFloat = gasFloat.Mul(gasFloat, big.NewFloat(config.GasFactor))
		gasInt, _ := gasFloat.Int(nil)
		gasString = gasInt.String()
	} else {
		gasString = strconv.FormatUint(gasLimit, 10)
	}

	return gasString, nil
}

func (c *Proxy) Nonce(spenderAddressHex string) (uint64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(c.Timeout)*time.Second)
	defer cancel()
	nonce, err := c.RemoteRpcClient.PendingNonceAt(ctx, common.HexToAddress(spenderAddressHex))
	if err != nil {
		return 0, err
	}
	return nonce, nil
}

func (c *Proxy) BuildTxUnSign(address string, transaction *Transaction) (*types.Transaction, error) {
	if transaction.Nonce == "" || transaction.Nonce == "0" {
		if !util.IsValidAddress(address) {
			return nil, errors.New("address format is error")
		}
		nonce, err := c.Nonce(address)
		if err != nil {
			nonce = 0
			err = nil
		}
		transaction.Nonce = strconv.FormatUint(nonce, 10)
	}
	return transaction.GetRawTx()
}

func (c *Proxy) BuildTxSign(privateKey *ecdsa.PrivateKey, txNoSign *types.Transaction) (*BuildTxResult, error) {
	if privateKey == nil || txNoSign == nil {
		return nil, errors.New("param is empty")
	}

	signedTx, err := types.SignTx(txNoSign, types.LatestSignerForChainID(c.chainId), privateKey)
	if err != nil {
		return nil, err
	}
	return &BuildTxResult{
		SignedTx: signedTx,
		TxHex:    signedTx.Hash().String(),
	}, nil
}

func (c *Proxy) SendTx(signedTx *types.Transaction) error {
	if signedTx == nil {
		return errors.New("signed transaction can't be empty")
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(c.Timeout)*time.Second)
	defer cancel()
	err := c.RemoteRpcClient.SendTransaction(ctx, signedTx)
	if err != nil {
		return err
	}
	return nil
}
