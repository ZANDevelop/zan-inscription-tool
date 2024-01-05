package config

type Inscription struct {
	Times      int
	Delay      int
	PrivateKey string
	GasPrice   string
	GasLimit   string
	Data       string
	RpcUrl     string
	//MaxPriorityFeePerGas string
}
