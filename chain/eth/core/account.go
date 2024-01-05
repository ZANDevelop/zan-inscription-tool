package core

import (
	"encoding/hex"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	hdwallet "github.com/miguelmota/go-ethereum-hdwallet"
	"github.com/tyler-smith/go-bip39"
	"inscription/chain/util"
)

type Account struct {
	Address    string `json:"address"`
	PrivateKey string `json:"privateKey"`
	PublicKey  string `json:"publicKey"`
	Mnemonic   string `json:"mnemonic"`
}

func NewAccount() *Account {
	return &Account{}
}

func (a *Account) AccountByMnemonic() (account *Account, err error) {
	entropy, err := bip39.NewEntropy(128)
	if err != nil {
		return nil, err
	}
	mnemonic, err := bip39.NewMnemonic(entropy)
	if err != nil {
		return nil, err
	}
	return a.AccountInfoByMnemonic(mnemonic)
}

func (a *Account) AccountInfoByMnemonic(mnemonic string) (account *Account, err error) {
	wallet, err := hdwallet.NewFromSeed(bip39.NewSeed(mnemonic, ""))
	if err != nil {
		return nil, err
	}
	acc, err := wallet.Derive(hdwallet.MustParseDerivationPath("m/44'/60'/0'/0/0"), false)
	if err != nil {
		return nil, err
	}
	address := acc.Address.Hex()

	// get public key
	// publicKey, _ := wallet.PublicKeyHex(acc),because of this function remove the begin of 04,so i don't use this function
	publicKeyBytes, err := wallet.PublicKeyBytes(acc)
	if err != nil {
		return nil, err
	}
	publicKey := hexutil.Encode(publicKeyBytes)[2:]

	privateKey, err := wallet.PrivateKeyHex(acc)
	if err != nil {
		return nil, err
	}
	account = &Account{
		PrivateKey: privateKey,
		PublicKey:  publicKey,
		Mnemonic:   mnemonic,
		Address:    address,
	}
	return
}

func (a *Account) AccountWithPrivateKey(privateKey string) (account *Account, err error) {
	priData, err := util.HexDecodeString(privateKey)
	if err != nil {
		return nil, err
	}

	privateKeyECDSA, err := crypto.ToECDSA(priData)
	if err != nil {
		return nil, err
	}
	address := crypto.PubkeyToAddress(privateKeyECDSA.PublicKey).Hex()
	publicKey := hex.EncodeToString(crypto.FromECDSAPub(&privateKeyECDSA.PublicKey))
	account = &Account{
		PrivateKey: privateKey,
		PublicKey:  publicKey,
		Mnemonic:   "",
		Address:    address,
	}
	return
}
