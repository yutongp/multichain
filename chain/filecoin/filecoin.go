package filecoin

import (
	"bytes"
	"fmt"

	"github.com/minio/blake2b-simd"
	"github.com/renproject/multichain/compat/filecoincompat"
	"github.com/renproject/pack"

	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/lotus/chain/types"
	"github.com/filecoin-project/specs-actors/actors/abi/big"
	"github.com/filecoin-project/specs-actors/actors/crypto"
)

type txBuilder struct {
}

func NewTxBuilder() filecoincompat.TxBuilder {
	return txBuilder{}
}

// BuildTx returns a simple Bitcoin transaction that consumes the funds from the
// given outputs, and sends the to the given recipients. The difference in the
// sum value of the inputs and the sum value of the recipients is paid as a fee
// to the Bitcoin network.
//
// It is assumed that the required signature scripts require the SIGHASH_ALL
// signatures and the serialized public key:
//
//  builder := txscript.NewScriptBuilder()
//  builder.AddData(append(signature.Serialize(), byte(txscript.SigHashAll)))
//  builder.AddData(serializedPubKey)
//
// Outputs produced for recipients will use P2PKH, P2SH, P2WPKH, or P2WSH
// scripts as the pubkey script, based on the format of the recipient address.
func (txBuilder txBuilder) BuildTx(recipient pack.Bytes, value pack.U256, nonce pack.U256) (filecoincompat.Tx, error) {

	recipientAddress, err := address.NewFromBytes(recipient)
	if err != nil {
		return nil, err
	}

	msgTx := &types.Message{
		From:     fromAddr,
		To:       recipientAddress,
		Value:    types.BigFromBytes(value.Bytes()),
		GasLimit: 10000,
		GasPrice: big.NewInt(0),
		Nonce:    nonce.Int().Uint64(),
	}

	return Tx{
		msgTx:     msgTx,
		signature: nil,
	}, nil
}

// Tx represents a simple Bitcoin transaction that implements the Bitcoin Compat
// API.
type Tx struct {
	msgTx     *types.Message
	signature *crypto.Signature
}

func (tx *Tx) Hash() pack.Bytes32 {
	return pack.NewBytes32(blake2b.Sum256(tx.msgTx.Cid().Bytes()))
}

func (tx *Tx) Sighash() pack.Bytes32 {
	return pack.NewBytes32(blake2b.Sum256(tx.msgTx.Cid().Bytes()))
}

func (tx *Tx) Sign(signature pack.Bytes65, pubKey pack.Bytes) error {
	if tx.signature != nil {
		return fmt.Errorf("already signed")
	}

	tx.signature = &crypto.Signature{
		Type: crypto.SigTypeSecp256k1,
		Data: signature.Bytes(),
	}

	return nil
}

func (tx *Tx) Serialize() (pack.Bytes, error) {
	buf := new(bytes.Buffer)
	// TODO
	// if err := tx.msgTx.Serialize(buf); err != nil {
	// 	return pack.Bytes{}, err
	// }
	return pack.NewBytes(buf.Bytes()), nil
}
