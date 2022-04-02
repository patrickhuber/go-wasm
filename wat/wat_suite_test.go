package wat_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestWat(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Wat Suite")
}
