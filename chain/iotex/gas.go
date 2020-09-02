package iotex

import (
	"context"

	"github.com/renproject/multichain/api/gas"

	"github.com/renproject/pack"
)

type GasEstimator struct {
	gasPrice pack.U256
}

func NewGasEstimator(gasPrice pack.U256) GasEstimator {
	return GasEstimator{
		gasPrice: gasPrice,
	}
}

func (gasEstimator GasEstimator) EstimateGasPrice(_ context.Context) (pack.U256, error) {
	return gasEstimator.gasPrice, nil
}

func (gasEstimator GasEstimator) EstimateGasLimit(txType gas.TxType) (pack.U256, error) {
	return pack.NewU256FromU64(pack.NewU64(10000)), nil
}
