package core

import (
	"errors"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"inscription/chain/util"
	"math/big"
	"strconv"
	"strings"
)

type BuildTxResult struct {
	SignedTx *types.Transaction
	TxHex    string
}

// CallMsg contains parameters for contract calls.
type CallMsg struct {
	Msg ethereum.CallMsg
}

// NewCallMsg creates an empty contract call parameter list.
func NewCallMsg() *CallMsg {
	return new(CallMsg)
}

func (msg *CallMsg) GetFrom() string     { return msg.Msg.From.String() }
func (msg *CallMsg) GetGasLimit() string { return strconv.FormatUint(msg.Msg.Gas, 10) }
func (msg *CallMsg) GetGasPrice() string { return msg.Msg.GasPrice.String() }
func (msg *CallMsg) GetValue() string    { return msg.Msg.Value.String() }
func (msg *CallMsg) GetData() []byte     { return msg.Msg.Data }
func (msg *CallMsg) GetDataHex() string  { return util.HexEncodeToString(msg.Msg.Data) }
func (msg *CallMsg) GetTo() string       { return msg.Msg.To.String() }

func (msg *CallMsg) SetFrom(address string) { msg.Msg.From = common.HexToAddress(address) }
func (msg *CallMsg) SetGasLimit(gas string) {
	i, _ := strconv.ParseUint(gas, 10, 64)
	msg.Msg.Gas = i
}
func (msg *CallMsg) SetGasPrice(price string) {
	i, _ := new(big.Int).SetString(price, 10)
	msg.Msg.GasPrice = i
}

// Set amount with decimal number
func (msg *CallMsg) SetValue(value string) {
	i, _ := new(big.Int).SetString(value, 10)
	msg.Msg.Value = i
}

// Set amount with hexadecimal number
func (msg *CallMsg) SetValueHex(hex string) {
	hex = strings.TrimPrefix(hex, "0x") // must trim 0x !!
	i, _ := new(big.Int).SetString(hex, 16)
	msg.Msg.Value = i
}
func (msg *CallMsg) SetData(data []byte) { msg.Msg.Data = common.CopyBytes(data) }
func (msg *CallMsg) SetDataHex(hex string) {
	data, err := util.HexDecodeString(hex)
	if err != nil {
		return
	}
	msg.Msg.Data = data
}
func (msg *CallMsg) SetTo(address string) {
	if address == "" {
		msg.Msg.To = nil
	} else {
		a := common.HexToAddress(address)
		msg.Msg.To = &a
	}
}

type Transaction struct {
	Nonce    string // nonce of sender account
	GasPrice string // wei per gas
	GasLimit string // gas limit
	To       string // receiver
	Value    string // wei amount
	Data     string // contract invocation input data

	// EIP1559, Default is ""
	MaxPriorityFeePerGas string
}

func NewTransaction(nonce, gasPrice, gasLimit, maxPriorityFeePerGas, to, value, data string) *Transaction {
	return &Transaction{nonce, gasPrice, gasLimit, to, value, data, maxPriorityFeePerGas}
}

func (tx *Transaction) MaxFee() string {
	return tx.GasPrice
}

func (tx *Transaction) SetMaxFee(maxFee string) {
	tx.GasPrice = maxFee
}

func (tx *Transaction) GetRawTx() (*types.Transaction, error) {
	var (
		gasPrice, value, maxFeePerGas *big.Int // default nil

		nonce     uint64 = 0
		gasLimit  uint64 = 90000 // reference https://eth.wiki/json-rpc/API method eth_sendTransaction
		toAddress common.Address
		data      []byte
		valid     bool
		err       error
	)
	if tx.GasPrice != "" {
		if gasPrice, valid = big.NewInt(0).SetString(tx.GasPrice, 10); !valid {
			return nil, errors.New("invalid gasPrice")
		}
	}
	if tx.Value != "" {
		if value, valid = big.NewInt(0).SetString(tx.Value, 10); !valid {
			return nil, errors.New("invalid value")
		}
	}
	if tx.MaxPriorityFeePerGas != "" {
		if maxFeePerGas, valid = big.NewInt(0).SetString(tx.MaxPriorityFeePerGas, 10); !valid {
			return nil, errors.New("invalid max priority fee per gas")
		}
	}
	if tx.Nonce != "" {
		if nonce, err = strconv.ParseUint(tx.Nonce, 10, 64); err != nil {
			return nil, errors.New("invalid Nonce")
		}
	}
	if tx.GasLimit != "" {
		if gasLimit, err = strconv.ParseUint(tx.GasLimit, 10, 64); err != nil {
			return nil, errors.New("invalid gas limit")
		}
	}
	if tx.To != "" && !common.IsHexAddress(tx.To) {
		return nil, errors.New("invalid toAddress")
	}
	toAddress = common.HexToAddress(tx.To)
	if tx.Data != "" {
		if data, err = util.HexDecodeString(tx.Data); err != nil {
			return nil, errors.New("invalid data string")
		}
	}

	if maxFeePerGas == nil || maxFeePerGas.Int64() == 0 {
		// is legacy tx
		return types.NewTx(&types.LegacyTx{
			Nonce:    nonce,
			To:       &toAddress,
			Value:    value,
			Gas:      gasLimit,
			GasPrice: gasPrice,
			Data:     data,
		}), nil
	} else {
		// is dynamic fee tx
		return types.NewTx(&types.DynamicFeeTx{
			Nonce:     nonce,
			To:        &toAddress,
			Value:     value,
			Gas:       gasLimit,
			GasFeeCap: gasPrice,
			GasTipCap: maxFeePerGas,
			Data:      data,
		}), nil
	}
}
