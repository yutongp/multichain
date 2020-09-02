package iotex

import (
	"context"

	"github.com/renproject/pack"
)

type GasEstimator struct {
	gasPerByte pack.U256
}

func NewGasEstimator(gasPerByte pack.U256) GasEstimator {
	return GasEstimator{
		gasPerByte: gasPerByte,
	}
}

func (gasEstimator GasEstimator) EstimateGasPrice(_ context.Context) (pack.U256, error) {
	return gasEstimator.gasPerByte, nil
}
