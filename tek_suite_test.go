package tek_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestTek(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Tek Suite")
}
