package filecoincompat_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestFilecoinCompat(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Filecoin Compat Suite")
}
