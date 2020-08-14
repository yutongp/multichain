package filecoincompat

import (
	"fmt"

	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcutil"
	"github.com/renproject/pack"
)

// GatewayPubKeyScript returns the Filecoin master address.
// This is the pubkey script that is expected to be in the underlying
// transaction output for lock-and-mint cross-chain transactions (when working
// with BTC).
func GatewayPubKeyScript(gpubkey pack.Bytes, ghash pack.Bytes32) ([]byte, error) {
	script, err := GatewayScript(gpubkey, ghash)
	if err != nil {
		return nil, fmt.Errorf("invalid script: %v", err)
	}
	pubKeyScript, err := txscript.NewScriptBuilder().
		AddOp(txscript.OP_HASH160).
		AddData(btcutil.Hash160(script)).
		AddOp(txscript.OP_EQUAL).Script()
	if err != nil {
		return nil, fmt.Errorf("invalid pubkeyscript: %v", err)
	}
	return pubKeyScript, nil
}
