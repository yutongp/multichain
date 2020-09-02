package iotex

import (
	"log"
	"os"
	"time"

	"github.com/renproject/multichain/api/address"
	"github.com/renproject/pack"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("IoTex", func() {
	Context("when decoding address", func() {
		Context("when decoding IoTex address", func() {
			It("should work", func() {
				decoder := NewAddressDecoder()

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
				opts := ClientOptions{
					endpoint: "api.testnet.iotex.one:80",
					secure:   false,
				}
				//client := NewClient(opts)
				//tx := txBuilder{}.BuildTx()
				//
				//// Send the transfer ont to multiple addresses transaction
				//txHash, err := client.MultiTransferOnt(0, 20000, acc, states, acc)
				//Expect(err).ToNot(HaveOccurred())
				//log.Printf("TxHash               %v", txHash.ToHexString())
				//
				//// We wait for 1000 ms before beginning to check transaction.
				//time.Sleep(1 * time.Second)
				//
				//for {
				//	evt, err := client.GetEvent(txHash.ToHexString())
				//	Expect(err).ToNot(HaveOccurred())
				//	if evt != nil {
				//		Expect(evt.State).Should(BeEquivalentTo(1))
				//		log.Printf("Tx success       %v", txHash.ToHexString())
				//		break
				//	}
				//	time.Sleep(1 * time.Second)
				}
			})
		})
	})
})
