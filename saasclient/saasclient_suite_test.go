package saasclient_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestSaasclient(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Saasclient Suite")
}
