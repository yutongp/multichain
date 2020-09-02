package iotex

import (
	"context"
	"crypto/tls"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"sync"

	"github.com/golang/protobuf/proto"
	"github.com/iotexproject/go-pkgs/crypto"
	"github.com/iotexproject/go-pkgs/hash"
	iotex "github.com/iotexproject/iotex-address/address"
	"github.com/iotexproject/iotex-core/action"
	"github.com/iotexproject/iotex-proto/golang/iotexapi"
	"github.com/iotexproject/iotex-proto/golang/iotextypes"
	"github.com/renproject/multichain/api/account"
	"github.com/renproject/multichain/api/address"
	"github.com/renproject/multichain/api/contract"
	"github.com/renproject/pack"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials"
)

type tx struct {
	from               address.Address
	to                 address.Address
	value, nonce       pack.U256
	gasLimit, gasPrice pack.U256
	payload            pack.Bytes
	*action.Envelope
	Sig       pack.Bytes65
	PublicKey pack.Bytes
}

func (t *tx) Hash() pack.Bytes {
	sealed, err := t.Serialize()
	if err != nil {
		return nil
	}
	h := hash.Hash256b(sealed)
	return h[:]
}

func (t *tx) From() address.Address {
	return t.from
}

func (t *tx) To() address.Address {
	return t.to
}

func (t *tx) Value() pack.U256 {
	return t.value
}

func (t *tx) Nonce() pack.U256 {
	return t.nonce
}

func (t *tx) Payload() contract.CallData {
	return contract.CallData(t.payload)
}

func (t *tx) Sighash() (pack.Bytes32, error) {
	return pack.Bytes32(t.Envelope.Hash()), nil
}

func (t *tx) Sign(sig pack.Bytes65, publicKey pack.Bytes) error {
	fmt.Println("Sign before", hex.EncodeToString(publicKey))
	copy(t.Sig[:], sig[:])
	t.PublicKey = make([]byte, len(publicKey))
	copy(t.PublicKey[:], publicKey[:])
	fmt.Println("Sign after", hex.EncodeToString(t.Sig[:]))
	fmt.Println("Sign after", hex.EncodeToString(t.PublicKey))
	return nil
}

func (t *tx) Serialize() (pack.Bytes, error) {
	fmt.Println("Serialize", hex.EncodeToString(t.PublicKey))
	pub, err := crypto.BytesToPublicKey(t.PublicKey)
	//fmt.Println("Serialize", hex.EncodeToString(t.PublicKey), err)
	if err != nil {
		return nil, err
	}

	ret := &iotextypes.Action{
		Core:         t.Envelope.Proto(),
		SenderPubKey: pub.Bytes(),
		Signature:    t.Sig[:],
	}
	return proto.Marshal(ret)
}

func (t *tx) Deserialize(ser pack.Bytes) (account.Tx, error) {
	act := &iotextypes.Action{}
	if err := proto.Unmarshal(ser, act); err != nil {
		return nil, err
	}
	pub, err := crypto.BytesToPublicKey(act.GetSenderPubKey())
	if err != nil {
		return nil, err
	}
	from, err := iotex.FromBytes(pub.Hash())
	if err != nil {
		return nil, err
	}
	env := &action.Envelope{}
	err = env.LoadProto(act.GetCore())
	if err != nil {
		return nil, err
	}
	sig := pack.Bytes65{}
	copy(sig[:], act.GetSignature())
	amount, ok := new(big.Int).SetString(act.GetCore().GetTransfer().GetAmount(), 10)
	if !ok {
		return nil, errors.New("amount convert error")
	}
	gasPrice, ok := new(big.Int).SetString(act.GetCore().GetGasPrice(), 10)
	if !ok {
		return nil, errors.New("gas price convert error")
	}
	return &tx{
		address.Address(from.String()),
		address.Address(act.GetCore().GetTransfer().GetRecipient()),
		pack.NewU256FromInt(amount),
		pack.NewU256FromU64(pack.U64(act.GetCore().GetNonce())),
		pack.NewU256FromU64(pack.U64(act.GetCore().GetGasLimit())),
		pack.NewU256FromInt(gasPrice),
		act.GetCore().GetTransfer().GetPayload(),
		env,
		sig,
		act.GetSenderPubKey(),
	}, nil
}

type TxBuilder struct {
}

func (t *TxBuilder) BuildTx(from, to address.Address, value, nonce pack.U256, gasPrice, gasLimit pack.U256, payload pack.Bytes) (account.Tx, error) {
	tsf, _ := action.NewTransfer(
		nonce.Int().Uint64(),
		value.Int(),
		string(to),
		payload,
		gasLimit.Int().Uint64(),
		gasPrice.Int(),
	)
	eb := action.EnvelopeBuilder{}
	evlp := eb.
		SetAction(tsf).
		SetGasLimit(tsf.GasLimit()).
		SetGasPrice(tsf.GasPrice()).
		SetNonce(tsf.Nonce()).
		SetVersion(1).
		Build()
	return &tx{from, to, value, nonce, gasLimit, gasPrice, payload, &evlp, pack.NewBytes65([65]byte{0}), nil}, nil
}

type ClientOptions struct {
	Endpoint string
	Secure   bool
}

type client struct {
	sync.RWMutex
	opts     ClientOptions
	grpcConn *grpc.ClientConn
	client   iotexapi.APIServiceClient
}

func NewClient(opts ClientOptions) account.Client {
	return &client{opts: opts}
}

func (c *client) Tx(ctx context.Context, h pack.Bytes) (account.Tx, pack.U64, error) {
	if err := c.connect(); err != nil {
		return nil, 0, err
	}
	res, err := c.client.GetActions(ctx, &iotexapi.GetActionsRequest{
		Lookup: &iotexapi.GetActionsRequest_ByHash{ByHash: &iotexapi.GetActionByHashRequest{ActionHash: h.String()}}})
	if err != nil {
		return nil, 0, err
	}
	if len(res.ActionInfo) != 1 {
		return nil, 0, errors.New("action number should be one")
	}
	ser, err := proto.Marshal(res.ActionInfo[0].GetAction())
	if err != nil {
		return nil, 0, err
	}
	tx := tx{}
	atx, err := tx.Deserialize(ser)
	return atx, 1, err
}

func (c *client) SubmitTx(ctx context.Context, t account.Tx) error {
	if err := c.connect(); err != nil {
		return err
	}

	iotexTx, ok := t.(*tx)
	if !ok {
		return errors.New("not iotex tx")
	}
	pub, err := crypto.BytesToPublicKey(iotexTx.PublicKey)
	if err != nil {
		return err
	}

	_, err = c.client.SendAction(ctx, &iotexapi.SendActionRequest{Action: &iotextypes.Action{
		Core:         iotexTx.Proto(),
		SenderPubKey: pub.Bytes(),
		Signature:    iotexTx.Sig[:],
	}})
	return err
}

func (c *client) connect() (err error) {
	c.Lock()
	defer c.Unlock()
	// Check if the existing connection is good.
	if c.grpcConn != nil && c.grpcConn.GetState() != connectivity.Shutdown {
		return
	}
	opts := []grpc.DialOption{}
	if c.opts.Secure {
		opts = append(opts, grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{})))
	} else {
		opts = append(opts, grpc.WithInsecure())
	}
	c.grpcConn, err = grpc.Dial(c.opts.Endpoint, opts...)
	c.client = iotexapi.NewAPIServiceClient(c.grpcConn)
	return err
}
