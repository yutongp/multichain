package iotex_test

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"os"
	"time"

	"github.com/iotexproject/go-pkgs/crypto"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/renproject/multichain/api/address"
	"github.com/renproject/multichain/chain/iotex"
	"github.com/renproject/pack"
)

var _ = Describe("IoTex", func() {
	Context("when decoding address", func() {
		Context("when decoding IoTex address", func() {
			It("should work", func() {
				decoder := iotex.NewAddressDecoder()

				addrStr := "io17ch0jth3dxqa7w9vu05yu86mqh0n6502d92lmp"
				_, err := decoder.DecodeAddress(address.Address(pack.NewString(addrStr)))

				Expect(err).ToNot(HaveOccurred())
			})
		})
	})
})
var _ = Describe("IoTex", func() {
	Context("when submitting transactions", func() {
		Context("when sending IOTX", func() {
			It("should work", func() {
				pkEnv := os.Getenv("pk")
				if pkEnv == "" {
					panic("pk is undefined")
				}
				opts := iotex.ClientOptions{
					Endpoint: "api.testnet.iotex.one:80",
					Secure:   false,
				}
				client := iotex.NewClient(opts)
				builder := iotex.TxBuilder{}
				gasPrice, _ := new(big.Int).SetString("100000000000", 10)
				tx, err := builder.BuildTx("io1vdtfpzkwpyngzvx7u2mauepnzja7kd5rryp0sg", "io1vdtfpzkwpyngzvx7u2mauepnzja7kd5rryp0sg", pack.NewU256FromU64(pack.NewU64(1)), pack.NewU256FromU64(pack.NewU64(132)), pack.NewU256FromInt(gasPrice), pack.NewU256FromU64(pack.NewU64(1000000)), nil)
				Expect(err).ToNot(HaveOccurred())
				sh, err := tx.Sighash()
				Expect(err).ToNot(HaveOccurred())
				sk, err := crypto.HexStringToPrivateKey("0d4d9b248110257c575ef2e8d93dd53471d9178984482817dcbd6edb607f8cc5")
				Expect(err).ToNot(HaveOccurred())
				sig, err := sk.Sign(sh[:])
				Expect(err).ToNot(HaveOccurred())
				var sig65 [65]byte
				copy(sig65[:], sig[:])
				tx.Sign(pack.NewBytes65(sig65), sk.PublicKey().Bytes())
				sigHash, err := tx.Sighash()
				Expect(err).ToNot(HaveOccurred())

				sHash := hex.EncodeToString(sigHash[:])
				Expect(err).ToNot(HaveOccurred())
				fmt.Println("sig hash:", sHash)
				fmt.Println("transaction hash:", hex.EncodeToString(tx.Hash()))

				err = client.SubmitTx(context.Background(), tx)
				Expect(err).ToNot(HaveOccurred())

				// We wait for 10 s before beginning to check transaction.
				time.Sleep(10 * time.Second)
				returnedTx, n, err := client.Tx(context.Background(), tx.Hash())
				Expect(err).ToNot(HaveOccurred())
				Expect(n).To(Equal(1))
				Expect(returnedTx.Nonce()).To(Equal(pack.NewU256FromU64(pack.NewU64(132))))
			})
		})
	})
})
