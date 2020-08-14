package filecoincompat

import (
	"context"
	"encoding/hex"

	"github.com/filecoin-project/go-jsonrpc"
	"github.com/filecoin-project/lotus/api"
	filecoinClient "github.com/filecoin-project/lotus/api/client"
	"github.com/filecoin-project/lotus/cli"
	"github.com/multiformats/go-multiaddr"
	"github.com/renproject/pack"
)

const (
	DefaultAddress   = ""
	DefaultAuthToken = ""
)

// ClientOptions are used to parameterise the behaviour of the Client.
type ClientOptions struct {
	Address   string
	AuthToken string
}

// DefaultClientOptions returns ClientOptions with the default settings. These
// settings are valid for use with the default local deployment of the
// multichain. In production, the host, user, and password should be changed.
func DefaultClientOptions() ClientOptions {
	return ClientOptions{
		Address:   DefaultAddress,
		AuthToken: DefaultAuthToken,
	}
}

func (opts ClientOptions) WithAddress(address string) ClientOptions {
	opts.Address = address
	return opts
}

func (opts ClientOptions) WithToken(token string) ClientOptions {
	opts.AuthToken = token
	return opts
}

// A Client interacts with an instance of the Bitcoin network using the RPC
// interface exposed by a Bitcoin node.
type Client interface {
	// Output associated with an outpoint, and its number of confirmations.
	Output(ctx context.Context, outpoint Outpoint) (Output, int64, error)
	// UnspentOutputs spendable by the given address.
	UnspentOutputs(ctx context.Context, minConf, maxConf int64, address Address) ([]Output, error)
	// Confirmations of a transaction in the Bitcoin network.
	Confirmations(ctx context.Context, txHash pack.Bytes32) (int64, error)
	// SubmitTx to the Bitcoin network.
	SubmitTx(ctx context.Context, tx Tx) (pack.Bytes32, error)
}

type client struct {
	apiClient api.FullNode
	closer    jsonrpc.ClientCloser
}

// NewClient returns a new Client.
func NewClient(opts ClientOptions) (Client, error) {
	authToken, err := hex.Decode(opts.AuthToken)
	if err != nil {
		return nil, err
	}
	apiInfo := cli.APIInfo{
		Address: multiaddr.NewMultiaddr(opts.Address),
		Token:   authToken,
	}
	apiAddr, err := apiInfo.DialArgs()
	if err != nil {
		return nil, err
	}
	header := apiInfo.AuthHeader()

	apiClient, closer, err := filecoinClient.NewFullNodeRPC(apiAddr, header)
	if err != nil {
		return nil, err
	}

	return &client{
		apiClient,
		closer,
	}, nil
}

// SubmitTx to the Bitcoin network.
func (client *client) SubmitTx(ctx context.Context, tx Tx) (pack.Bytes32, error) {
	client.api.MpoolPush(ctx, tx.)
}
