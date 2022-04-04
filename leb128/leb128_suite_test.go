package leb128_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestLeb128(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Leb128 Suite")
}
