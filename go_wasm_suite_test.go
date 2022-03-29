package wasm_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestGoWasm(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "GoWasm Suite")
}
